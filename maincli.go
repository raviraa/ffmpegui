package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/raviraa/recordscreen/ffprobe"
)

var log = ffprobe.SetLogger()

func beginCli() {
	// log.Info("Starting in CLI")
	prober := ffprobe.GetPlatformProber()
	opts := ffprobe.Options{}
	opts.VidIdx = 3
	opts.Container = 0
	ffprobe.SetOptions(opts)
	stdout, _ := ffprobe.StartEncode(prober)
	readStdout(stdout)
	log.Info("before sleep")
	time.Sleep(20 * time.Second)
	log.Info("sending stop signal")
	ffprobe.StopEncode()
}

func readStdout(scanner *bufio.Scanner) {
	go func() {
		scanner.Split(scanLines)
		// scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			txt := scanner.Text()
			fmt.Println(txt)
			if strings.Contains(txt, "frame=") {
				fmt.Print(txt)
			}
		}
	}()
}

func mainCli() {

	beginCli()
}
