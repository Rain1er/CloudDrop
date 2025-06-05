package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func PostRequest(shellURL string, password string, cmd string) (string, error) {
	data := url.Values{}
	data.Set(password, cmd)

	resp, err := http.PostForm(shellURL, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(body), nil
}
