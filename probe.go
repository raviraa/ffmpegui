package main

import (
	"strings"
)

// Prober finds devices based on platform
type Prober interface {
	probeDevices()
	getDevices() devices
	setOptions(opts options)
	getCommand() []string
	getVersion() string
}

type devices struct {
	audios []string
	videos []string
}

type options struct {
	audIdx     int
	vidIdx     int
	ffmpegOpts []string
}

type proberKeys struct {
	devicesCmd       string
	recordCmdPrefix  []string
	recordCmdPostfix []string
	devicesKey       string
	deviceKey        string //input device
	devices
	options
}

var macProber = proberKeys{
	deviceKey:  "input device",
	devicesCmd: "ffmpeg -f avfoundation -list_devices true -i ''",
	//TODO env FFREPORT=file=ffreport.log:level=32
	recordCmdPrefix: strings.Split("ffmpeg -y -report -f avfoundation -framerate 24", " "),
	//ffmpeg -y -video_size 1024x768 -framerate 5 -f avfoundation -i "3"  TODO.mkv
	recordCmdPostfix: strings.Split("-framerate 25 -s 1920x1080 TODO.mkv", " "), //-preset ultrafast aaa.mkv
	devicesKey:       "devices:",
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
			lnprsds := strings.Split(ln, "] [")
			if len(lnprsds) == 2 {
				lnprsd := lnprsds[1]
				lnprsd = strings.Replace(lnprsd, "]", " ", -1)
				devices[currDevType] = append(devices[currDevType], lnprsd)
			}
		}
	}
	log.Debugf("%s: %+v\n", dtype, devices[dtypeKey])
	return devices[dtypeKey]
}

func (mp proberKeys) getDevices() devices {
	return mp.devices
}
func (mp *proberKeys) probeDevices() {
	devs := devices{}
	devs.audios = parseFfmpegDevices(*mp, "audio")
	devs.videos = parseFfmpegDevices(*mp, "video")
	mp.devices = devs
}

func (mp *proberKeys) setOptions(opts options) {
	mp.options = opts
}

func (mp proberKeys) getVersion() string {
	return "ffmpeg 1234.22"
}

func (mp proberKeys) getCommand() (cmd []string) {
	cmd = append(cmd, mp.recordCmdPrefix...)
	cmd = append(cmd, "-i")
	cmd = append(cmd, "1")
	// cmd = append(cmd, fmt.Sprintf("%d:%d", mp.options.vidIdx, mp.options.audIdx))
	cmd = append(cmd, mp.recordCmdPostfix...)
	if len(mp.ffmpegOpts) > 1 {
		cmd = append(cmd, mp.ffmpegOpts...)
	}
	// runCmdPipe(cmd)
	// runCmdPipe(strings.Split("ls -lR ..", " "))
	return cmd
}

func getPlatformProber() Prober {
	var prober Prober = &macProber
	return prober
}
