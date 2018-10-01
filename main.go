package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

func runCmdStr(cmd string, ignExitCode bool) string {
	if str, err := runCmd(strings.Split(cmd, " "), ignExitCode); err == nil {
		return str
	}
	return ""
}

func runCmd(args []string, ignExitCode bool) (string, error) {
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
	return out.String(), nil
}

func checkRequirements() {
	runCmdStr("ffmpeg -version", false)
}

func beginCli() {
	log.Info("Starting in CLI")
	prober := getPlatformProber()
	prober.probeDevices()
	opts := options{vidIdx: 1}
	prober.setOptions(opts)
	cmd := prober.getCommand()
	log.Info("Using cmd:", cmd)
	prober.start()
	time.Sleep(3 * time.Second)
	log.Info("sending stop signal")
	prober.stop()
}

func mainCli() {

	beginCli()
}
