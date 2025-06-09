package service

import (
	"clouddrop/pkg/util"
	"strings"
)

func (s *JavaShell) GetShellType() string {
	return "java"
}

// GetCurrentDirectory retrieves the current directory from the PHP shell
func (s *JavaShell) GetCurrentDirectory(url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// ListFiles lists all files in the current directory
func (s *JavaShell) ListFiles(url string, password string) ([]string, error) {
	code := ``

	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}

func (s *JavaShell) ExecCommand(url string, password string, command string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
