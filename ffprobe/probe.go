package ffprobe

import (
	"bufio"
	"errors"
	"strings"

	logging "github.com/op/go-logging"
)

// Prober has logic that changes based on platform
type Prober interface {
	getDevicesCmd() string
	getPrefixCmd() []string
}

// Devices has information about ffmpeg multimedia devices
type Devices struct {
	Audios []string
	Videos []string
}

// Options configures ffmpeg encoding process
type Options struct {
	AudIdx int
	VidIdx int
	// ffmpegOpts []string
}

type proberCommon struct {
	recordCmdPostfix []string
	devicesKey       string
	deviceKey        string //input device
	done             chan bool
	started          bool
	Devices
	Options
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
func GetFfmpegDevices(prober Prober) Devices {
	devs := Devices{}
	devs.Audios = parseFfmpegDeviceType(prober, "audio")
	devs.Videos = parseFfmpegDeviceType(prober, "video")
	return devs
}

func parseFfmpegDeviceType(prober Prober, dtype string) []string {
	devs := map[string][]string{}
	res := runCmdStr(prober.getDevicesCmd(), true)
	resLines := strings.Split(res, "\n")
	filterfn := func(s string) bool { return strings.Contains(s, deviceCommon.deviceKey) }
	resLines = filterList(resLines, filterfn)
	var currDevType, dtypeKey string
	for _, ln := range resLines {
		if strings.Contains(ln, deviceCommon.devicesKey) {
			currDevType = ln
			if strings.Contains(ln, dtype) {
				dtypeKey = ln
			}
		} else if currDevType != "" {
			lnprsds := strings.Split(ln, "] [")
			if len(lnprsds) == 2 {
				lnprsd := lnprsds[1]
				lnprsd = strings.Replace(lnprsd, "]", " ", -1)
				devs[currDevType] = append(devs[currDevType], lnprsd)
			}
		}
	}
	log.Debugf("%s: %+v\n", dtype, devs[dtypeKey])
	return devs[dtypeKey]
}

// SetOptions configures extra encoder options
func SetOptions(opts Options) {
	deviceCommon.Options = opts
}

func getVersion() string {
	return "ffmpeg 1234.22" //TODO
}

// ffmpeg -y  -framerate 5 -video_size 1024x768 -f avfoundation -i 1 TODO.mkv
func getCommand(prober Prober) (cmd []string) {
	cmd = append(cmd, prober.getPrefixCmd()...)
	cmd = append(cmd, "-i", "1:0") //TODO
	cmd = append(cmd, deviceCommon.recordCmdPostfix...)
	// if len(mp.ffmpegOpts) > 1 {
	// 	cmd = append(cmd, mp.ffmpegOpts...)
	// }
	// runCmdPipe(strings.Split("ls -lR ..", " "))
	return cmd
}

// StartEncode starts ffmpeg process with configured options
// and returns stdout scanner
func StartEncode(prober Prober) (*bufio.Scanner, error) {
	if !deviceCommon.started {
		cmd := getCommand(prober)
		scanner, err := deviceCommon.runCmdPipe(cmd)
		if err == nil {
			deviceCommon.started = true
			return scanner, nil
		}
		log.Errorf("StartEncode failed" + err.Error())
		return nil, err
	}
	log.Errorf("already started")
	return nil, errors.New("already started")
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
	GetFfmpegDevices(prober)
	return prober
}
