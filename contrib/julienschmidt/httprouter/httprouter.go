// Package httprouter provides functions to trace the julienschmidt/httprouter package (https://github.com/julienschmidt/httprouter).
package httprouter // import "github.com/meetcircle/dd-trace-go/contrib/julienschmidt/httprouter"

import (
	"net/http"
	"strings"

	"github.com/meetcircle/dd-trace-go/contrib/internal/httputil"
	"github.com/meetcircle/dd-trace-go/ddtrace/ext"
	"github.com/meetcircle/dd-trace-go/ddtrace/tracer"

	"github.com/julienschmidt/httprouter"
)

// Router is a traced version of httprouter.Router.
type Router struct {
	*httprouter.Router
	config *routerConfig
}

// New returns a new router augmented with tracing.
func New(opts ...RouterOption) *Router {
	cfg := new(routerConfig)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	if cfg.analyticsRate > 0 {
		cfg.spanOpts = append(cfg.spanOpts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
	}
	return &Router{httprouter.New(), cfg}
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// get the resource associated to this request
	route := req.URL.Path
	_, ps, _ := r.Router.Lookup(req.Method, route)
	for _, param := range ps {
		route = strings.Replace(route, param.Value, ":"+param.Key, 1)
	}
	resource := req.Method + " " + route
	httputil.TraceAndServe(r.Router, w, req, r.config.serviceName, resource, r.config.spanOpts...)
}
