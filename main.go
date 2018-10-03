package main

import (
	"time"

	"github.com/raviraa/recordscreen/ffprobe"
)

var log = ffprobe.Log

func beginCli() {
	// log.Info("Starting in CLI")
	prober := ffprobe.GetPlatformProber()
	// devs := ffprobe.GetFfmpegDevices(prober)
	opts := ffprobe.Options{}
	opts.VidIdx = 1
	ffprobe.SetOptions(opts)
	stdout := ffprobe.StartEncode(prober)
	time.Sleep(5 * time.Second)
	slurpstdout
	log.Info("sending stop signal")
	ffprobe.StopEncode()
}

func mainCli() {

	beginCli()
}
