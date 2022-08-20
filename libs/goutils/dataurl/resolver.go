package dataurl

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var DefaultHTTPClient = &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}

func Convert(c context.Context, remoteURL string, extraParams map[string]string) (string, error) {
	u, _ := url.Parse(remoteURL)
	q := u.Query()
	for k, v := range extraParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(c, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	contentType := resp.Header.Get("Content-Type")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(body)), nil
}
