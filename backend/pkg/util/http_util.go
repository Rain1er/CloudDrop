package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func PostRequest(shellURL string, password string, code string, session string) (string, error) {
	// Convert password to int (assuming it represents a timestamp)
	var timestampValue int
	if _, err := fmt.Sscanf(password, "%d", &timestampValue); err != nil {
		return "", fmt.Errorf("password must be a valid integer: %v", err)
	}

	// Directly construct JSON payload
	jsonStr := fmt.Sprintf(`{"timezone":%d,"sign":"%s"}`, timestampValue, code)

	// Create request
	req, err := http.NewRequest("POST", shellURL, bytes.NewBufferString(jsonStr))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Cookie", session)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
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

func PostRequestWithoutSession(shellURL string, password string, code string) (string, error) {
	// Convert password to int (assuming it represents a timestamp)
	var timestampValue int
	if _, err := fmt.Sscanf(password, "%d", &timestampValue); err != nil {
		return "", fmt.Errorf("password must be a valid integer: %v", err)
	}

	// Directly construct JSON payload
	jsonStr := fmt.Sprintf(`{"timezone":%d,"sign":"%s"}`, timestampValue, code)

	// Create request
	req, err := http.NewRequest("POST", shellURL, bytes.NewBufferString(jsonStr))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}
	defer resp.Body.Close()

	// Extract PHPSESSID from cookies Todo 增加其他类型的session
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "PHPSESSID" {
			return cookie.Value, nil
		}
	}

	return "", nil
}
