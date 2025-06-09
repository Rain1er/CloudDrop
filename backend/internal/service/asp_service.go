package service

import (
	"clouddrop/pkg/util"
	"strings"
)

func (s *AspShell) GetShellType() string {
	return "asp"
}

// GetCurrentDirectory retrieves the current directory from the PHP shell
func (s *AspShell) GetCurrentDirectory(url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// ListFiles lists all files in the current directory
func (s *AspShell) ListFiles(url string, password string) ([]string, error) {
	code := ``

	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}
