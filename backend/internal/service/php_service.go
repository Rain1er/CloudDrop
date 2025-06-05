package service

import (
	"clouddrop/pkg/util"
	"strings"
)

// PHPShell represents a connection to a PHP shell
type PHPShell struct {
	Name     string
	URL      string
	Password string
	Type     string
	Encode   string
	Note     string
}

// NewPHPShell creates a new PHP shell connection
// 如果多个变量类型相同，只要在函数参数列表中声明一次即可
func NewPHPShell(name, url, password, shellType, encode, note string) *PHPShell {
	return &PHPShell{
		Name:     name,
		URL:      url,
		Password: password,
		Type:     shellType,
		Encode:   encode,
		Note:     note,
	}
}

// GetCurrentDirectory retrieves the current directory from the PHP shell
func (s *PHPShell) GetCurrentDirectory(url string, password string) (string, error) {
	cmd := `echo getcwd();`
	result, err := util.PostRequest(url, password, cmd)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// ListFiles lists all files in the current directory
func (s *PHPShell) ListFiles() ([]string, error) {
	cmd := `
	$files = scandir(getcwd());
	foreach ($files as $file) {
		if ($file != "." && $file != "..") {
			echo $file . "\n";
		}
	}
	`

	result, err := util.PostRequest("", "", cmd)
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}
