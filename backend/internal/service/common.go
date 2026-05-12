package service

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const paramHeaderName = "X-Token"

func readPayload(pathParts ...string) ([]byte, error) {
	parts := append([]string{"pkg", "api"}, pathParts...)
	path := filepath.Join(parts...)
	candidates := []string{
		path,
		filepath.Join(append([]string{"..", ".."}, parts...)...),
		filepath.Join(append([]string{"..", "..", ".."}, parts...)...),
	}
	var lastErr error
	for _, candidate := range candidates {
		data, err := os.ReadFile(candidate)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("read payload %s failed: %w", path, lastErr)
}

func payloadWithSuffix(payload []byte, suffix string) []byte {
	out := make([]byte, 0, len(payload)+len(suffix))
	out = append(out, payload...)
	out = append(out, suffix...)
	return out
}

func phpQuote(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\n", "\\n")
	return "\"" + value + "\""
}

func aspQuote(value string) string {
	value = strings.ReplaceAll(value, "\"", "\"\"")
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")
	return "\"" + value + "\""
}

func encodeParamsHeader(params map[string]string) map[string]string {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	encoded := base64.RawURLEncoding.EncodeToString([]byte(values.Encode()))
	return map[string]string{
		paramHeaderName: encoded,
	}
}

func decodeMaybeBase64(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	data, err := base64.StdEncoding.DecodeString(trimmed)
	if err == nil {
		return string(data)
	}
	if data, err = base64.RawStdEncoding.DecodeString(trimmed); err == nil {
		return string(data)
	}
	if pad := len(trimmed) % 4; pad != 0 {
		if data, err = base64.StdEncoding.DecodeString(trimmed + strings.Repeat("=", 4-pad)); err == nil {
			return string(data)
		}
	}
	return value
}

func normalizeASMXURL(shellURL string) string {
	lower := strings.ToLower(shellURL)
	if idx := strings.Index(lower, ".asmx"); idx >= 0 {
		return shellURL[:idx+len(".asmx")]
	}
	return shellURL
}

func isASMXURL(shellURL string) bool {
	return strings.Contains(strings.ToLower(shellURL), ".asmx")
}

func containsSuccess(value string) bool {
	return strings.Contains(strings.ToLower(value), "success")
}

func chooseWindowsCmd(osType string) string {
	if strings.Contains(strings.ToLower(osType), "linux") || strings.Contains(strings.ToLower(osType), "unix") {
		return "/bin/bash"
	}
	return "C:/Windows/System32/cmd.exe"
}

func stripHTML(value string) string {
	re := regexp.MustCompile(`<[^>]+>`)
	return re.ReplaceAllString(value, "")
}

func extractCSharpClassName(code string) string {
	re := regexp.MustCompile(`(?m)\bclass\s+([$_a-zA-Z][$_a-zA-Z0-9]*)\b`)
	if matches := re.FindStringSubmatch(code); len(matches) > 1 {
		return matches[1]
	}
	return "RunCode"
}
