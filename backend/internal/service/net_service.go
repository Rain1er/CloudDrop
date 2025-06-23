package service

import (
	"clouddrop/pkg/util"
	"strings"
)

var NetSessions map[int]string

func (s *CSharpShell) GetShellType() string {
	return "csharp"
}

func (s *CSharpShell) FreshSession(id int, url string, password string) (string, error) {
	return "", nil
}

// BaseInfo
func (s *CSharpShell) BaseInfo(id int, url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code, NetSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (s *CSharpShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	return "", nil
}

func (s *CSharpShell) ExecCode(id int, code string, url string, password string) (string, error) {
	return "", nil
}

func (s *CSharpShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error) {

	return "", nil
}

func (s *CSharpShell) FileZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	return "", nil
}
func (s *CSharpShell) FileUnZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	return "", nil
}

// FileList lists all files in the current directory
func (s *CSharpShell) FileList(id int, path string, url string, password string) (string, error) {
	return "", nil
}

func (s *CSharpShell) FileShow(id int, path string, url string, password string) (string, error) {
	return "", nil
}
