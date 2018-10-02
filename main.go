package main

import (
	"time"

	"github.com/raviraa/recordscreen/ffprobe"
)

var log = ffprobe.Log

func beginCli() {
	// log.Info("Starting in CLI")
	prober := ffprobe.GetPlatformProber()
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
