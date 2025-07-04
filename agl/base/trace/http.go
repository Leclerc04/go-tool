package trace

import (
	"io"
	"net/http"

	"golang.org/x/net/trace"
)

// Render calls trace.Render.
func Render(w io.Writer, req *http.Request, sensitive bool) {
	trace.Render(w, req, sensitive)
}

// RenderEvents calls trace.RenderEvents.
func RenderEvents(w http.ResponseWriter, req *http.Request, sensitive bool) {
	trace.RenderEvents(w, req, sensitive)
}
