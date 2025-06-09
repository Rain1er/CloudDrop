package util

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func PostRequest(shellURL string, password string, code string) (string, error) {
	// TODO 需要保存会话状态，保存第一次请求返回的cookie,如何做呢？感觉可以直接用map保存在内存中，反正不是很多
	data := url.Values{}
	data.Set(password, code)

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
