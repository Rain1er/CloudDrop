package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"os"
)

var AspSessions map[int]string

func (s *AspShell) GetShellType() string {
	return "asp"
}

func (s *AspShell) FreshSession(id int, url string, password string) (string, error) {
	if AspSessions == nil {
		AspSessions = make(map[int]string)
	}
	code, err := os.ReadFile("./pkg/api/asp/Check.asp")
	code = fmt.Append(code, "\ncall main()")

	if err != nil {
		return "", err
	}
	encode := util.Encrypt(string(code), password)
	AspSessions[id], err = util.PostRequestWithoutSession(url, password, encode)
	if err != nil {
		return "", err
	}

	session := AspSessions[id] // if key not exist, it returns "" , bcz type is string
	log.Println("当前ASPSESSIONID " + session)

	enResult, err := util.PostRequest(url, password, encode, session)
	if err != nil {
		return "", err
	}
	res := util.Decrypt(enResult, password)
	return res, nil

}

// BaseInfo
func (s *AspShell) BaseInfo(id int, url string, password string) (string, error) {
	code, err := os.ReadFile("./pkg/api/asp/BaseInfo.asp")
	if err != nil {
		return "", nil
	}
	code = fmt.Append(code, "\ncall main()")
	res, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *AspShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	code, err := os.ReadFile("./pkg/api/asp/OS.asp")
	if err != nil {
		return "", err
	}
	code = fmt.Append(code, "\ncall main()")

	osType, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	var cmdPath string // in if-block-level scope, it must define at out
	if osType == "Linux" {
		cmdPath = "/bin/bash"
	} else {
		cmdPath = "C:/Windows/System32/cmd.exe"
	}

	code, err = os.ReadFile("./pkg/api/asp/BaseInfo.asp")
	if err != nil {
		return "", nil
	}
	code = fmt.Appendf(code, "\ncall main(\"%s\", \"true\", \"%s\")", cmdPath, command)
	res, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *AspShell) ExecCode(id int, code string, url string, password string) (string, error) {
	return "", nil
}

func (s *AspShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/asp/Database.asp")

	// step 1, database is "", get all dbname
	if database == "" {
		code = fmt.Appendf(code,
			"\ncall listDatabases(\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")", driver, host, port, user, pass, encoding)
	} else {
		// step 2 , user choses dbname
		code = fmt.Appendf(code,
			"\ncall main(\"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", %s, \"%s\",)", driver, host, port, user, pass, database, sql, option, encoding)
	}

	res, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *AspShell) FileZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/asp/FileZip.asp")
	code = fmt.Appendf(code, "\ncall main(\"%s\", \"%s\")", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *AspShell) FileUnZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/asp/FileUnZip.asp")
	code = fmt.Appendf(code, "\ncall main(\"%s\", \"%s\")", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	return res, nil
}

// FileList lists all files in the current directory
func (s *AspShell) FileList(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/asp/FileList.asp")
	code = fmt.Appendf(code, "\ncall main(\"%s\")", path)

	res, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s *AspShell) FileShow(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/asp/FileShow.asp")
	code = fmt.Appendf(code, "\ncall main(\"%s\")", path)

	res, err := util.HookPost(url, password, string(code), AspSessions[id], s.GetShellType())
	if err != nil {
		return "", err
	}
	return res, nil
}
