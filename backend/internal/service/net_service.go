package service

import (
	"clouddrop/pkg/util"
	"strings"
)

func (s *CSharpShell) GetShellType() string {
	return "csharp"
}

// GetCurrentDirectory retrieves the current directory from the PHP shell
func (s *CSharpShell) GetCurrentDirectory(url string, password string) (string, error) {
	code := ``
	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// ListFiles lists all files in the current directory
func (s *CSharpShell) ListFiles(url string, password string) ([]string, error) {
	code := ``

	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}
