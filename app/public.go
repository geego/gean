package app

import (
	"bufio"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Publish() {
	command := globalConfig.Build.Publish
	// Prepare exec command
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	cmd := exec.Command(shell, flag, command)
	cmd.Dir = filepath.Join(rootPath, "public")
	// Start print stdout and stderr of process
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	out := bufio.NewScanner(stdout)
	err := bufio.NewScanner(stderr)
	// Print stdout
	go func() {
		for out.Scan() {
			Log(out.Text())
		}
	}()
	// Print stdin
	go func() {
		for err.Scan() {
			Log(err.Text())
		}
	}()
	// Exec command
	cmd.Run()
}
