package service

import (
	"clouddrop/pkg/util"
	"strings"
)

var JavaSessions map[int]string

func (s *JavaShell) GetShellType() string {
	return "java"
}

func (s *JavaShell) FreshSession(id int, url string, password string) (string, error) {
	return "", nil
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

// FileList lists all files in the current directory
func (s *JavaShell) FileList(id int, path string, url string, password string) (string, error) {
	code := ``

	result, err := util.PostRequest(url, password, code, JavaSessions[id])
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *JavaShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code, JavaSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
