package sentry

import (
	"context"
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/getsentry/raven-go"
	"github.com/golang/glog"
	"github.com/leclerc04/go-tool/agl/base/mon"
	"github.com/leclerc04/go-tool/agl/base/trace"
	"github.com/leclerc04/go-tool/agl/util/errs"
)

type logger context.Context

var client *raven.Client
var dataClient *raven.Client
var errorLog = func() logger {
	ctx, _ := trace.WithTrace(context.Background(), "log.sentry/error")
	return logger(ctx)
}()
var droppedLog = func() logger {
	ctx, _ := trace.WithTrace(context.Background(), "log.sentry/dropped")
	return logger(ctx)
}()
var dataErrorLog = func() logger {
	ctx, _ := trace.WithTrace(context.Background(), "log.sentry/data_error")
	return logger(ctx)
}()
var notificationLog = func() logger {
	ctx, _ := trace.WithTrace(context.Background(), "log.sentry/notification")
	return logger(ctx)
}()

var (
	metricTotal = mon.NewCounterVec("sentry", "total",
		"Total number of sentry captures.",
		[]string{"channel"})
	metricDroppedTotal = mon.NewCounter("sentry", "dropped_total",
		"Total number of captures that could not be sent to sentry.")
)

// Init intializes the global client. Do not call it twice.
func Init(dsn, dataDSN string) error {
	if client != nil {
		panic("sentry.Init by dsn shouldn't be called twice!")
	}
	if dataDSN != "" && dataClient != nil {
		panic("sentry.Init by data_dsn shouldn't be called twice!")
	}

	if dsn == "" {
		return nil
	}
	var err error
	client, err = NewClient(dsn)
	if err != nil {
		return err
	}
	dataClient, err = NewClient(dataDSN)
	return err
}

func getClientIPPort(req *http.Request) (string, string) {
	userIP := req.Header.Get("X-A2-User-Ip")
	if userIP != "" {
		return userIP, ""
	}
	forwarded := req.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ip := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(ip), ""
	}
	clientIP, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr, ""
	}
	return clientIP, port
}

// bgCtx is the global log variables.
func capture(
	ctx context.Context, bgCtx logger, depth int, err error, level raven.Severity,
	channel string) string {
	if err == nil {
		return ""
	}
	if errs.Unwrap(err) == context.Canceled {
		// ignore all context canceled.
		return ""
	}
	metricTotal.With(prometheus.Labels{
		"channel": channel,
	}).Inc()

	tags := getTags(ctx)
	fingerprint := ""
	extras := func() (extras []raven.Interface) {
		stackTrace := raven.NewStacktrace(3+depth, 3, []string{"bitbucket.org/applysquare/"})
		// if the last frame shouldSkip, send all frames to sentry; or only collect frames not shouldSKip to sentry
		var frames []*raven.StacktraceFrame
		frameSkipper := errs.NewFrameSkipper("" /* always perform the should skip logic */)
		fingerprintHasher := md5.New()
		if stackTrace != nil {
			for _, f := range stackTrace.Frames {
				f.InApp = !frameSkipper.ShouldSkip(f.Filename)
				if f.InApp {
					fmt.Fprint(fingerprintHasher, f.Filename, f.Lineno)
				}
				frames = append(frames, f)
			}
		}
		ex := raven.NewException(err, &raven.Stacktrace{Frames: frames})
		ex.Type = func() string {
			parts := []string{channel, ex.Type}
			if err, ok := err.(*errs.Error); ok {
				parts = append(parts, err.Kind.String())
				fmt.Fprint(fingerprintHasher, err.Stack(true))
				fingerprint = fmt.Sprintf("%x", fingerprintHasher.Sum(nil))
			}
			return strings.Join(parts, "~")
		}()

		extras = append(extras, ex)

		if tags == nil {
			return extras
		}
		if tags.Request != nil {
			httpExtra := raven.NewHttp(tags.Request)
			a2UA := tags.Request.Header.Get("X-A2-Request-User-Agent")
			if a2UA != "" {
				httpExtra.Headers["User-Agent"] = a2UA
			}
			userIP, userPort := getClientIPPort(tags.Request)
			if httpExtra.Env == nil {
				httpExtra.Env = map[string]string{}
			}
			httpExtra.Env["REMOTE_ADDR"] = userIP
			if userPort != "" {
				httpExtra.Env["REMOTE_PORT"] = userPort
			} else {
				delete(httpExtra.Env, "REMOTE_PORT")
			}
			httpExtra.Data = tags.httpData()
			extras = append(extras, httpExtra)
		}
		userExtra := tags.userExtra()
		if userExtra != nil {
			extras = append(extras, userExtra)
		}
		return extras
	}()

	msg := err.Error()
	if tags != nil && tags.creatorStack != "" {
		msg += "\nGo routine creator stack: \n" + tags.creatorStack
	}

	packet := raven.NewPacket(msg, extras...)
	packet.Level = level
	if fingerprint != "" {
		packet.Fingerprint = []string{fingerprint}
	}

	trace.Printf(bgCtx, "error", "err", err)

	if client == nil {
		func() {
			if channel == "data" {
				return
			}
			if level == raven.ERROR {
				if errs.Internal.Is(err) {
					glog.ErrorDepth(2+depth, channel, " : ", errs.ToStringWithFullStack(err))
				} else {
					glog.ErrorDepth(2+depth, channel, " : ", msg)
				}
			} else if level == raven.INFO {
				glog.InfoDepth(2+depth, channel, " : ", msg)
			} else {
				panic("unknwon level: " + level)
			}
		}()
		return "<error not recorded>"
	}

	var currentClient *raven.Client
	if channel == "data" {
		currentClient = dataClient
	}
	if currentClient == nil {
		currentClient = client
	}
	eventID, _ := currentClient.Capture(packet, map[string]string{"channel": channel})
	return eventID
}

// ErrorDepth sends the error to sentry asynchronously.
func ErrorDepth(ctx context.Context, depth int, err error) string {
	return Send(ctx, err, SendParams{
		TopFramesToOmit: 1 + depth,
	})
}

// Error a shortcut of ErrorDepth(0, ...).
func Error(ctx context.Context, err error) string {
	return Send(ctx, err, SendParams{
		TopFramesToOmit: 1,
	})
}

// DataErrorDepth sends a data inconsistency error to sentry asynchronously.
func DataErrorDepth(ctx context.Context, depth int, err error) string {
	return Send(ctx, err, SendParams{
		TopFramesToOmit: depth + 1,
		Channel:         "data",
	})
}

// DataError is a shortcut of DataErrorDepth(0, ...).
func DataError(ctx context.Context, err error) string {
	return Send(ctx, err, SendParams{
		TopFramesToOmit: 1,
		Channel:         "data",
	})
}

// DataErrorIfNotFound captures the not found error to data error, and other error to normal error.
func DataErrorIfNotFound(ctx context.Context, err error) string {
	if err == nil {
		return ""
	}
	if errs.NotFound.Is(err) {
		return DataErrorDepth(ctx, 1, err)
	}
	return ErrorDepth(ctx, 1, err)
}

// ErrorIfNotNotFound ignores the not found error, and captures other error to normal error.
func ErrorIfNotNotFound(ctx context.Context, err error) string {
	if err == nil || errs.NotFound.Is(err) {
		return ""
	}
	return ErrorDepth(ctx, 1, err)
}

// WXBotDown sends a message to sentry with it's kind equals "wxbot_down"
func WXBotDown(ctx context.Context, err error) {
	Send(ctx, err, SendParams{
		TopFramesToOmit: 1,
		Channel:         "wxbot_down",
	})
}

// Info sends an info to sentry
func Info(ctx context.Context, msg string, args ...interface{}) {
	Send(ctx, fmt.Errorf(msg, args...), SendParams{
		TopFramesToOmit: 1,
		Severity:        raven.INFO,
		Channel:         "info",
	})
}

// SendParams defines parameters for send.
type SendParams struct {
	Channel         string // e.g.: "", "data".
	TopFramesToOmit int
	Severity        raven.Severity
}

// Send submits an error to sentry.
func Send(ctx context.Context, err error, params SendParams) string {
	if err == nil {
		return ""
	}
	var traceCtx logger
	switch params.Channel {
	case "data":
		traceCtx = dataErrorLog
	case "info":
		traceCtx = notificationLog
	default:
		traceCtx = errorLog
		params.Channel = "error"
	}
	if params.Severity == "" {
		params.Severity = raven.ERROR
	}
	trace.Printf(ctx, "sentry %s: %v", params.Channel, err)
	return capture(
		ctx, traceCtx, params.TopFramesToOmit+1, err, params.Severity,
		params.Channel)
}
