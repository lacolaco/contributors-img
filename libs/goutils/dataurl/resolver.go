package dataurl

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func ResolveImageDataURL(c context.Context, remoteURL string, imageSize int) (string, error) {
	u, _ := url.Parse(remoteURL)
	q := u.Query()
	q.Set("size", fmt.Sprint(imageSize))
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(c, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
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
