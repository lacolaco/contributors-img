package httptrace

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewTransport(base http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(base)
}

func NewClient(base http.RoundTripper) *http.Client {
	return &http.Client{Transport: NewTransport(base)}
}
