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

func runCmdPipe(cmdstr []string) {
	log.Info(cmdstr)
	// cmd := exec.Command("sh", "-c", "echo stdout; echo 1>&2 stderr")
	cmd := exec.Command(cmdstr[0], cmdstr[1:]...)
	// out, err := cmd.StderrPipe()
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("before start")
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	log.Info("after start")
	// go io.Copy(os.Stdout, stderr)
	go func() {
		// defer stdin.Close()
		buf := new(bytes.Buffer)
		n, _ := buf.ReadFrom(out)
		// if e != nil {
		fmt.Println(n)
		fmt.Println(buf.String())
		// }
	}()
	// slurp, _ := ioutil.ReadAll(stderr)
	// for ln, err := stderr.Read(buf); err != nil; {
	// 	log.Info(ln)
	// }
	log.Info("after buf for")

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
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
}

func mainCli() {

	beginCli()
}
