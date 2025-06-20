package service

import (
	"clouddrop/pkg/util"
	"strings"
)

var AspSessions map[int]string

func (s *AspShell) GetShellType() string {
	return "asp"
}

func (s *AspShell) FreshSession(id int, url string, password string) (string, error) {
	return "", nil
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

// ListFiles lists all files in the current directory
func (s *AspShell) ListFiles(id int, url string, password string) ([]string, error) {
	code := ``

	result, err := util.PostRequest(url, password, code, AspSessions[id])
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}

func (s *AspShell) ExecCommand(id int, url string, password string, command string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code, AspSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
