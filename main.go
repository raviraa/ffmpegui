package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

func runCmdStr(cmd string, ignExitCode bool) string {
	return runCmd(strings.Split(cmd, " "), ignExitCode)
}

func runCmd(args []string, ignExitCode bool) string {
	// args = append([]string{"-c"}, args...)
	// cmd := exec.Command("/bin/sh", args...)
	var out bytes.Buffer
	log.Info(args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil && !ignExitCode {
		fmt.Println(out.String())
		panic(err)
	}
	return out.String()
}

func checkRequirements() {
	runCmdStr("ffmpeg -version", false)
}

func getDevices() {
	//ffmpeg -f avfoundation -list_devices true -i ''
}

func main() {

	checkRequirements()
}
