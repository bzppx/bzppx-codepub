package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"time"
	"os"
)

func NewCommandX() *CommandX {
	return &CommandX{}
}

type CommandX struct {
}

const Command_ExecType_SyncErrorStop = 1   // 同步执行，遇到错误停止
const Command_ExecType_SyncErrorAccess = 2 // 同步执行，遇到错误继续
const Command_ExecType_Asy = 3             // 异步执行

type CommandXParams struct {
	Path               string
	Command            string
	CommandExecType    int
	CommandExecTimeout int
}

// 执行命令
func (c *CommandX) Exec(commandXParams CommandXParams) (err error) {
	if commandXParams.Command == "" {
		return nil
	}
	if commandXParams.CommandExecType == Command_ExecType_Asy {
		err = c.asyExec(commandXParams)
	} else {
		err = c.syncExec(commandXParams)
	}
	return
}

// 同步执行, 等待结果返回
func (c *CommandX) syncExec(commandXParams CommandXParams) (err error) {

	fileName, err := c.createTmpShellFile(commandXParams.Command)
	if err != nil {
		return
	}
	outChan := make(chan error, 1)
	var out bytes.Buffer
	go func() {
		defer func() {
			err = os.Remove(fileName)
			if err != nil {
				log.Println("remove file error:"+err.Error())
			}
			rec := recover()
			if rec != nil {
				outChan <- fmt.Errorf("%v", rec)
			}
		}()
		cmd := c.command(fileName, commandXParams.Path)
		cmd.Stderr = &out
		select {
		case outChan <- cmd.Run():
			return
		case <-time.After(time.Duration(commandXParams.CommandExecTimeout) * time.Second):
			cmd.Process.Kill()
			outChan <- errors.New("sync command exec timeout")
			return
		}
	}()
	err = <-outChan
	if (err != nil) && (err.Error() == "sync command exec timeout") {
		return err
	}
	if (err != nil) && (out.String() != "") {
		return errors.New(out.String() + "," + err.Error())
	}
	return
}

// 异步执行，不返回结果
func (c *CommandX) asyExec(commandXParams CommandXParams) (err error) {

	fileName, err := c.createTmpShellFile(commandXParams.Command)
	if err != nil {
		return
	}
	go func() {
		defer func() {
			err = os.Remove(fileName)
			if err != nil {
				log.Println("remove file error:"+err.Error())
			}
			rec := recover()
			if rec != nil {
				log.Panicf("%v", rec)
			}
		}()
		cmd := c.command(fileName, commandXParams.Path)
		outChan := make(chan error, 1)
		var out bytes.Buffer
		cmd.Stderr = &out
		select {
		case outChan <- cmd.Run():
		case <-time.After(time.Duration(commandXParams.CommandExecTimeout) * time.Second):
			cmd.Process.Kill()
			outChan <- errors.New("asy exec command exec timeout")
		}
		err = <-outChan
		if (err != nil) && (err.Error() == "asy exec command exec timeout") {
			log.Println(err.Error())
			return
		}
		if (err != nil) && (out.String() != "") {
			log.Println("asy exec command error:"+out.String() + "," + err.Error())
			return
		}
	}()
	return nil
}

// 获取 command
// filename 文件名
// path 执行命令的目录
func (c *CommandX) command(fileName string, path string) (cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd = exec.Command(fileName)
	} else {
		cmd = exec.Command("/bin/bash", fileName)
	}
	cmd.Dir = path
	return
}

// 创建临时的 shell 脚本文件
// path 脚本执行目录
// content 创建的脚本内容
func (c *CommandX) createTmpShellFile(content string) (tmpFile string, err error) {

	file, err := ioutil.TempFile("", "codepub_tmp")
	if err != nil {
		return
	}
	defer func() {
		file.Close()
		if runtime.GOOS == "windows" {
			tmpFile = file.Name() + ".cmd"
		}else {
			tmpFile = file.Name() + ".sh"
		}
		err = os.Rename(file.Name(), tmpFile)
	}()

	file.Chmod(0777)
	if runtime.GOOS == "windows" {

	}else {
		file.WriteString("#!/bin/bash\n")
		file.WriteString("set -e\n")
	}
	_, err = file.WriteString(content)

	return
}
