package iputil

import (
	"net"
	"net/http"
	"strings"
)

// GuessIPFromRequest returns a guessed client ip. It should not be trusted.
func GuessIPFromRequest(req *http.Request) string {
	xf := req.Header.Get("X-Forwarded-For")
	if xf == "" {
		clientIP, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return req.RemoteAddr
		}
		return clientIP
	}
	return strings.Split(xf, ",")[0]
}
