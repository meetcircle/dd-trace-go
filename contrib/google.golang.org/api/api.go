// Package api provides functions to trace the google.golang.org/api package.
package api // import "github.com/meetcircle/dd-trace-go/contrib/google.golang.org/api"

//go:generate go run make_endpoints.go

import (
	"net/http"

	"golang.org/x/oauth2/google"
	"github.com/meetcircle/dd-trace-go/contrib/google.golang.org/api/internal"
	httptrace "github.com/meetcircle/dd-trace-go/contrib/net/http"
	"github.com/meetcircle/dd-trace-go/ddtrace"
	"github.com/meetcircle/dd-trace-go/ddtrace/ext"
)

// apiEndpoints are all of the defined endpoints for the Google API; it is populated
// by "go generate".
var apiEndpoints *internal.Tree

// NewClient creates a new oauth http client suitable for use with the google
// APIs with all requests traced automatically.
func NewClient(options ...Option) (*http.Client, error) {
	cfg := newConfig(options...)
	client, err := google.DefaultClient(cfg.ctx, cfg.scopes...)
	if err != nil {
		return nil, err
	}
	client.Transport = WrapRoundTripper(client.Transport, options...)
	return client, nil
}

// WrapRoundTripper wraps a RoundTripper intended for interfacing with
// Google APIs and traces all requests.
func WrapRoundTripper(transport http.RoundTripper, options ...Option) http.RoundTripper {
	cfg := newConfig(options...)
	rtOpts := []httptrace.RoundTripperOption{
		httptrace.WithBefore(func(req *http.Request, span ddtrace.Span) {
			e, ok := apiEndpoints.Get(req.URL.Hostname(), req.Method, req.URL.Path)
			if ok {
				span.SetTag(ext.ServiceName, e.ServiceName)
				span.SetTag(ext.ResourceName, e.ResourceName)
			} else {
				span.SetTag(ext.ServiceName, "google")
				span.SetTag(ext.ResourceName, req.Method+" "+req.URL.Hostname())
			}
			if cfg.serviceName != "" {
				span.SetTag(ext.ServiceName, cfg.serviceName)
			}
		}),
	}
	if cfg.analyticsRate > 0 {
		rtOpts = append(rtOpts, httptrace.RTWithAnalyticsRate(cfg.analyticsRate))
	}
	return httptrace.WrapRoundTripper(transport, rtOpts...)
}
