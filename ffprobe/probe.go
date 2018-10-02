package ffprobe

import (
	"bufio"
	"strings"

	logging "github.com/op/go-logging"
)

// Prober contains logic that changes based on platform
type Prober interface {
	getDevicesCmd() string
	getPrefixCmd() []string
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

type proberCommon struct {
	recordCmdPostfix []string
	devicesKey       string
	deviceKey        string //input device
	done             chan bool
	started          bool
	devices
	options
}

var log = logging.MustGetLogger("probe")

// Log TODO, how to initialize once?
var Log = log

var deviceCommon = proberCommon{
	deviceKey: "input device",
	// TODO: -s scale, resolution input oputpu
	recordCmdPostfix: strings.Split("-framerate 25 -s 1920x1080 TODO.mkv", " "), //-preset ultrafast aaa.mkv
	devicesKey:       "devices:",
	done:             make(chan bool),
}

func filterList(ss []string, f func(string) bool) (res []string) {
	for _, s := range ss {
		if f(s) {
			res = append(res, s)
		}
	}
	return res
}

// GetFfmpegDevices returns audio and video devices available
func GetFfmpegDevices(prober Prober, pk proberCommon, dtype string) devices {
	devs := devices{}
	devs.audios = parseFfmpegDeviceType(prober, pk, "audio")
	devs.videos = parseFfmpegDeviceType(prober, pk, "video")
	return devs
}

func parseFfmpegDeviceType(prober Prober, pk proberCommon, dtype string) []string {
	devices := map[string][]string{} //0: "Built-in Microphone"}
	res := runCmdStr(prober.getDevicesCmd(), true)
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

// SetOptions configures extra encoder options
func SetOptions(opts options) {
	deviceCommon.options = opts
}

func getVersion() string {
	return "ffmpeg 1234.22" //TODO
}

func getCommand(prober Prober) (cmd []string) {
	mp := deviceCommon
	cmd = append(cmd, prober.getPrefixCmd()...)
	if len(mp.ffmpegOpts) > 1 {
		cmd = append(cmd, mp.ffmpegOpts...)
	}
	// runCmdPipe(strings.Split("ls -lR ..", " "))
	return cmd
}

// StartEncode starts ffmpeg process with configured options
func StartEncode(prober Prober) *bufio.Scanner {
	if !deviceCommon.started {
		cmd := getCommand(prober)
		scanner, err := deviceCommon.runCmdPipe(cmd)
		if err == nil {
			deviceCommon.started = true
			return scanner
		}
		log.Errorf("StartEncode failed" + err.Error())
	} else {
		log.Errorf("already started")
	}
	return nil
}

// StopEncode stop ffmpeg process
func StopEncode() bool {
	if deviceCommon.started {
		deviceCommon.done <- true //send done signal and..
		<-deviceCommon.done       //wait for process done
		log.Info("process stopped")
		deviceCommon.started = false
		return true
	}
	log.Errorf("already stopped")
	return false
}

// GetPlatformProber returns prober for correct platform
func GetPlatformProber() Prober {
	var prober Prober = &macProber
	return prober
}
