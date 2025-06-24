package service

import (
	"clouddrop/pkg/util"
	"log"
	"os"
	"strings"
)

var JavaSessions map[int]string

func (s *JavaShell) GetShellType() string {
	return "java"
}

func (s *JavaShell) FreshSession(id int, url string, password string) (string, error) {
	if JavaSessions == nil {
		NetSessions = make(map[int]string)
	}
	code, err := os.ReadFile("./pkg/api/java/Check.class")
	if err != nil {
		return "", err
	}
	encode := util.Encrypt(string(code), password)
	NetSessions[id], err = util.PostRequestWithoutSession(url, password, encode)
	if err != nil {
		return "", err
	}
	session := JavaSessions[id]
	log.Println("当前JSESSIONID " + session)

	enResult, err := util.PostRequest(url, password, encode, session)
	if err != nil {
		return "", err
	}
	res := util.Decrypt(enResult, password)
	return res, nil

}

// BaseInfo
func (s *JavaShell) BaseInfo(id int, url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code, JavaSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (s *JavaShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	return "", nil
}

func (s *JavaShell) ExecCode(id int, code string, url string, password string) (string, error) {
	return "", nil
}

func (s *JavaShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error) {

	return "", nil
}

func (s *JavaShell) FileZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	return "", nil
}

func (s *JavaShell) FileUnZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	return "", nil
}

// FileList lists all files in the current directory
func (s *JavaShell) FileList(id int, path string, url string, password string) (string, error) {
	return "", nil
}

func (s *JavaShell) FileShow(id int, path string, url string, password string) (string, error) {
	return "", nil
}
