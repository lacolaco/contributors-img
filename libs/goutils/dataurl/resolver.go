package dataurl

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func ResolveImageDataURL(remoteURL string, imageSize int) (string, error) {
	u, _ := url.Parse(remoteURL)
	q := u.Query()
	q.Set("size", fmt.Sprint(imageSize))
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
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
