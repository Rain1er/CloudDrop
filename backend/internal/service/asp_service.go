package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"os"
	"strings"
)

var AspSessions map[int]string

func (s *AspShell) GetShellType() string {
	return "asp"
}

func (s *AspShell) FreshSession(id int, url string, password string) (string, error) {
	if AspSessions == nil {
		AspSessions = make(map[int]string)
	}
	code, err := os.ReadFile("./pkg/api/asp/Check.asp")
	code = fmt.Append(code, "\nmain()")

	if err != nil {
		return "", err
	}
	encode := util.Encrypt(string(code), password)
	AspSessions[id], err = util.PostRequestWithoutSession(url, password, encode)
	if err != nil {
		return "", err
	}

	session := NetSessions[id] // if key not exist, it returns "" , bcz type is string
	log.Println("当前ASPSESSIONID " + session)

	enResult, err := util.PostRequest(url, password, encode, session)
	if err != nil {
		return "", err
	}
	res := util.Decrypt(enResult, password)
	return res, nil

}

// BaseInfo
func (s *AspShell) BaseInfo(id int, url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code, AspSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (s *AspShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	return "", nil
}

func (s *AspShell) ExecCode(id int, code string, url string, password string) (string, error) {
	return "", nil
}

func (s *AspShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error) {
	return "", nil
}

func (s *AspShell) FileZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	return "", nil
}

func (s *AspShell) FileUnZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	return "", nil
}

// FileList lists all files in the current directory
func (s *AspShell) FileList(id int, path string, url string, password string) (string, error) {
	return "", nil
}

func (s *AspShell) FileShow(id int, path string, url string, password string) (string, error) {
	return "", nil
}
