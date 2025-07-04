package sentry

import (
	"context"
	"net/http"

	"github.com/leclecr04/go-tool/agl/util/errs"

	raven "github.com/getsentry/raven-go"
)

type contextKey struct{}

// Tags contains annotation of current context to be reported
// to sentry.
type Tags struct {
	Request      *http.Request
	graphqlQuery string
	graphqlVars  string
	creatorStack string
}

// CreateDetachedContext creates an independent context with annontation of
// the passed in context of the creator. Used for context of a spawned go routine,
// that could outlive the original context.
func CreateDetachedContext(ctx context.Context) context.Context {
	newTags := &Tags{}
	tags := getTags(ctx)
	if tags != nil {
		*newTags = *tags
	}
	newTags.creatorStack = errs.GetStack(1)
	return WithTags(context.Background(), newTags)
}

// WithTags returns a context attached with some sentry tags.
func WithTags(ctx context.Context, tags *Tags) context.Context {
	exTags := getTags(ctx)
	if exTags == nil {
		return context.WithValue(ctx, contextKey{}, tags)
	}
	if tags.Request != nil {
		exTags.Request = tags.Request
	}
	return ctx
}

// AttachGraphQLInfo attaches graphql information to the sentry context.
func AttachGraphQLInfo(ctx context.Context, query string, vars string) {
	tags := getTags(ctx)
	if tags == nil {
		return
	}
	tags.graphqlQuery = query
	tags.graphqlVars = vars
}

// getTags returns the attached tags, never make this exported.
func getTags(ctx context.Context) *Tags {
	v := ctx.Value(contextKey{})
	if v == nil {
		return nil
	}
	return v.(*Tags)
}

func (t *Tags) httpData() map[string]string {
	if t.graphqlQuery == "" {
		return nil
	}
	ret := map[string]string{}
	ret["gql_query"] = t.graphqlQuery
	if t.graphqlVars != "" {
		ret["gql_vars"] = t.graphqlVars
	}
	return ret
}

func (t *Tags) userExtra() *raven.User {
	if t.Request == nil {
		return nil
	}
	userID := t.Request.Header.Get("X-A2-User-Uuid")
	if userID == "" {
		return nil
	}
	u := &raven.User{}
	u.ID = userID
	u.Email = t.Request.Header.Get("X-A2-User-Email")
	return u
}
