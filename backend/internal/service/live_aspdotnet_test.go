//go:build live

package service

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type liveTarget struct {
	name       string
	url        string
	sh         Shell
	codeSample string
}

func TestLiveAspDotNet(t *testing.T) {
	targets := []liveTarget{
		{name: "java-jsp", url: envOr("CLOUDDROP_TEST_JSP_URL", "http://192.168.0.241:55001/api.jsp"), sh: &JavaShell{}, codeSample: `System.out.print("hello world");`},
		{name: "java-jspx", url: envOr("CLOUDDROP_TEST_JSPX_URL", "http://192.168.0.241:55001/api.jspx"), sh: &JavaShell{}, codeSample: `System.out.print("hello world");`},
		{name: "asp", url: envOr("CLOUDDROP_TEST_ASP_URL", "http://10.1.1.139/api.asp"), sh: &AspShell{}, codeSample: `GlobalResult = "hello world"`},
		{name: "aspx", url: envOr("CLOUDDROP_TEST_ASPX_URL", "http://10.1.1.139/api.aspx"), sh: &CSharpShell{}, codeSample: `public class RunCode { public override string ToString() { return "hello world"; } }`},
		{name: "ashx", url: envOr("CLOUDDROP_TEST_ASHX_URL", "http://10.1.1.139/api.ashx"), sh: &CSharpShell{}, codeSample: `public class RunCode { public override string ToString() { return "hello world"; } }`},
		{name: "asmx", url: envOr("CLOUDDROP_TEST_ASMX_URL", "http://10.1.1.139/api.asmx"), sh: &CSharpShell{}, codeSample: `public class RunCode { public override string ToString() { return "hello world"; } }`},
	}
	for i, target := range targets {
		target := target
		t.Run(target.name, func(t *testing.T) {
			id := 9000 + i
			runLiveTarget(t, id, target)
		})
	}
}

func runLiveTarget(t *testing.T, id int, target liveTarget) {
	check, err := target.sh.FreshSession(id, target.url, "")
	if err != nil {
		t.Fatalf("FreshSession failed: %v", err)
	}
	if !containsSuccess(check) {
		t.Fatalf("FreshSession returned %q", check)
	}

	baseInfo, err := target.sh.BaseInfo(id, target.url, "")
	if err != nil {
		t.Fatalf("BaseInfo failed: %v", err)
	}
	if strings.TrimSpace(baseInfo) == "" {
		t.Fatalf("BaseInfo empty")
	}

	echo, err := target.sh.ExecCommand(id, "echo CloudDropLiveTest", target.url, "")
	if err != nil {
		t.Fatalf("ExecCommand failed: %v", err)
	}
	if !strings.Contains(echo, "CloudDropLiveTest") {
		t.Fatalf("ExecCommand returned %q", echo)
	}

	codeOutput, err := target.sh.ExecCode(id, target.codeSample, target.url, "")
	if err != nil {
		t.Fatalf("ExecCode failed: %v", err)
	}
	if !strings.Contains(strings.ToLower(codeOutput), "hello world") {
		t.Fatalf("ExecCode returned %q", codeOutput)
	}

	windows := isWindowsTarget(baseInfo)
	root := remoteTempRoot(t, id, target, windows)
	dir := remoteJoin(windows, root, "src")
	outDir := remoteJoin(windows, root, "out")
	zipPath := remoteJoin(windows, root, "src.zip")
	inputPath := remoteJoin(windows, dir, "input.txt")

	setup := setupCommand(windows, root, dir, outDir, inputPath)
	if _, err := target.sh.ExecCommand(id, setup, target.url, ""); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer target.sh.ExecCommand(id, cleanupCommand(windows, root), target.url, "")

	list, err := target.sh.FileList(id, dir, target.url, "")
	if err != nil {
		t.Fatalf("FileList failed: %v", err)
	}
	if !strings.Contains(strings.ToLower(list), "input.txt") {
		t.Fatalf("FileList missing input.txt: %q", list)
	}

	content, err := target.sh.FileShow(id, inputPath, target.url, "")
	if err != nil {
		t.Fatalf("FileShow failed: %v", err)
	}
	if !strings.Contains(content, "CloudDropLiveTest") {
		t.Fatalf("FileShow returned %q", content)
	}

	zipRes, err := target.sh.FileZip(id, dir, zipPath, target.url, "")
	if err != nil {
		t.Fatalf("FileZip failed: %v", err)
	}
	if !containsSuccess(zipRes) {
		t.Fatalf("FileZip returned %q", zipRes)
	}

	unzipRes, err := target.sh.FileUnZip(id, zipPath, outDir, target.url, "")
	if err != nil {
		t.Fatalf("FileUnZip failed: %v", err)
	}
	if !containsSuccess(unzipRes) && strings.TrimSpace(unzipRes) != "" {
		t.Fatalf("FileUnZip returned %q", unzipRes)
	}

	candidates := []string{
		remoteJoin(windows, outDir, "input.txt"),
		remoteJoin(windows, outDir, "src", "input.txt"),
	}
	var lastErr error
	for _, candidate := range candidates {
		content, lastErr = target.sh.FileShow(id, candidate, target.url, "")
		if lastErr == nil && strings.Contains(content, "CloudDropLiveTest") {
			return
		}
	}
	t.Fatalf("unzipped file verification failed, lastErr=%v content=%q", lastErr, content)
}

func remoteTempRoot(t *testing.T, id int, target liveTarget, windows bool) string {
	command := `printf %s "${TMPDIR:-/tmp}"`
	if windows {
		command = "echo %TEMP%"
	}
	out, err := target.sh.ExecCommand(id, command, target.url, "")
	if err != nil {
		t.Fatalf("temp query failed: %v", err)
	}
	temp := strings.TrimSpace(strings.ReplaceAll(out, "\r", ""))
	if windows && (temp == "" || strings.Contains(temp, "%TEMP%")) {
		temp = "C:\\Windows\\Temp"
	}
	if !windows && temp == "" {
		temp = "/tmp"
	}
	return remoteJoin(windows, strings.TrimRight(temp, "\\/"), fmt.Sprintf("clouddrop_live_%d_%d", id, time.Now().UnixNano()))
}

func isWindowsTarget(baseInfo string) bool {
	lower := strings.ToLower(baseInfo)
	return strings.Contains(lower, "windows") || strings.Contains(lower, "win32")
}

func remoteJoin(windows bool, parts ...string) string {
	sep := "/"
	if windows {
		sep = "\\"
	}
	cleaned := make([]string, 0, len(parts))
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			cleaned = append(cleaned, strings.TrimRight(part, "\\/"))
			continue
		}
		cleaned = append(cleaned, strings.Trim(part, "\\/"))
	}
	return strings.Join(cleaned, sep)
}

func setupCommand(windows bool, root, dir, outDir, inputPath string) string {
	if windows {
		return fmt.Sprintf(`mkdir "%s" & mkdir "%s" & echo CloudDropLiveTest>"%s" & echo Success`, dir, outDir, inputPath)
	}
	return fmt.Sprintf(`rm -rf %s; mkdir -p %s %s; printf 'CloudDropLiveTest\n' > %s; echo Success`, shQuote(root), shQuote(dir), shQuote(outDir), shQuote(inputPath))
}

func cleanupCommand(windows bool, root string) string {
	if windows {
		return fmt.Sprintf(`if exist "%s" rmdir /s /q "%s"`, root, root)
	}
	return fmt.Sprintf(`rm -rf %s`, shQuote(root))
}

func shQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
