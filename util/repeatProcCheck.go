package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// RepeatProcCheck 检查进程是否重复开启 当前只支持windows平台
func RepeatProcCheck(pidFileName string) (isRepeat bool) {
	pid := strconv.Itoa(os.Getpid())
	ppid := strconv.Itoa(os.Getppid())
	fmt.Println(pid, ppid)
	execPath := getExecPathByPid(pid)
	if execPath == "" {
		isRepeat = true
		return isRepeat
	}
	homePath, err := Home()
	if err != nil {
		fmt.Println(err)
	}
	pidFilePath := path.Join(homePath, pidFileName)
	pidFileR, _ := os.OpenFile(pidFilePath, os.O_RDWR|os.O_CREATE, 0777)
	defer pidFileR.Close()
	lastPidBytes, _ := ioutil.ReadAll(pidFileR)
	lastPid := string(lastPidBytes)
	pidFileR.Close()
	lastExecPath := getExecPathByPid(lastPid)
	if execPath == lastExecPath {
		isRepeat = true
		return isRepeat
	}
	if lastExecPath != "" {
		fmt.Println(lastExecPath)
	}
	pidFileW, _ := os.OpenFile(pidFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	defer pidFileW.Close()
	log.Println("记录当前pid:", pid, " 到"+pidFilePath)
	pidFileW.WriteString(pid)
	pidFileW.Close()
	isRepeat = false
	return isRepeat
}
func getExecPathByPid(pid string) string {
	//get executablepath,name
	order := "wmic process where processid=" + pid + " get executablepath"
	fmt.Println(order)
	cmd := exec.Command("cmd", "/c", order)
	output, _ := cmd.Output()
	outputArr := strings.Split(string(output), "\n")
	return strings.TrimSpace(outputArr[1])
}
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
