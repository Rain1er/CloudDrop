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

// FileList lists all files in the current directory
func (s *AspShell) FileList(id int, path string, url string, password string) (string, error) {
	code := ``

	result, err := util.PostRequest(url, password, code, AspSessions[id])
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *AspShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code, AspSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
