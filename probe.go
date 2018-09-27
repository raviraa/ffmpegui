package main

import (
	"strings"
)

// Prober finds devices based on platform
type Prober interface {
	getAudios() []string
}

type proberKeys struct {
	devicesCmd string
	devicesKey string
	deviceKey  string //input device
}

var macProber = proberKeys{
	deviceKey:  "input device",
	devicesCmd: "ffmpeg -f avfoundation -list_devices true -i ''",
	devicesKey: "devices:",
}

func parseFfmpegDevices(pk proberKeys, dtype string) []string {
	devices := map[string][]string{} //0: "Built-in Microphone"}
	res := runCmdStr(pk.devicesCmd, true)
	resLines := strings.Split(res, "\n")
	filterfn := func(s string) bool { return strings.Contains(s, pk.deviceKey) }
	resLines = filterList(resLines, filterfn)
	var currDevType, dtypeKey string
	for _, ln := range resLines {
		if strings.Contains(ln, pk.devicesKey) {
			currDevType = ln
			if strings.Contains(ln, dtype) {
				dtypeKey = ln
			}
		} else if currDevType != "" {
			lnStripped := strings.Split(ln, "] [")
			devices[currDevType] = append(devices[currDevType], lnStripped[1])
		}
	}
	// log.Debugf("%+v", devices)
	return devices[dtypeKey]
}

func (mp proberKeys) getAudios() []string {
	devices := []string{}
	return devices
}

func getPlatformProber() Prober {
	var prober Prober = macProber
	return prober
}
