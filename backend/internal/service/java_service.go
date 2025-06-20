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

// ListFiles lists all files in the current directory
func (s *JavaShell) ListFiles(id int, url string, password string) ([]string, error) {
	code := ``

	result, err := util.PostRequest(url, password, code, JavaSessions[id])
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}

func (s *JavaShell) ExecCommand(id int, url string, password string, command string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code, JavaSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
