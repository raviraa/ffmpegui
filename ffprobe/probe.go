package ffprobe

import (
	"errors"
	"runtime"
	"strings"

	logging "github.com/op/go-logging"
)

// Prober has logic that changes based on platform
type Prober interface {
	getDevicesCmd() string
	getFfmpegCmd() ([]string, error)
}

// Devices has information about ffmpeg multimedia devices
type Devices struct {
	Audios []string
	Videos []string
}

type plt int

const (
	lin plt = iota
	mac     //TODO use this
)

type proberCommon struct {
	devicesKey string
	deviceKey  string //input device
	done       chan bool
	opts       *Options
	Devices
}

// Started indicates whether ffmpeg process is running
var Started bool //TODO needs mutex?

var opts *Options // TODO move to devcommon

func init() {
	opts = &Options{}
	cfgname = "common_presets.toml"
	uiOptsFname = "uiopts.toml"
}

// SetInputs to set configure input streams
func SetInputs(uiips []UIInput, resumeCount int) {
	opts.UIInputs = uiips
	config.resumeCount = resumeCount
}

// GetInputs gets
func GetInputs() []UIInput {
	return opts.UIInputs
}

var log *logging.Logger

// SetLogger starts logger, TODO memory logger
func SetLogger() *logging.Logger {
	log = logging.MustGetLogger("probe")
	return log
}

// TODO move to NewProber
var deviceCommon = proberCommon{
	deviceKey:  "input device",
	devicesKey: "devices:",
	done:       make(chan bool),
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

// GetVersion returns ffmpeg version
func GetVersion() string {
	return "ffmpeg 1234.22" //TODO
}

func getCommand(prober Prober) ([]string, error) {
	return prober.getFfmpegCmd()
}

// StartMux concats all resume split streams,
// and avstreams to final container
func StartMux(prober Prober) error {
	// TODO start stop status
	for avidx := range opts.UIInputs {
		cmd, err := getConcatCmd(*opts, avidx)
		if err != nil {
			return err
		}
		if err = deviceCommon.runCmdPipe(cmd, true); err != nil {
			return err
		}
	}
	cmd, err := getMuxCommand(*opts)
	if err != nil {
		return err
	}
	if err = deviceCommon.runCmdPipe(cmd, true); err != nil {
		return err
	}
	return nil
}

// StartEncode starts ffmpeg process with configured options
// and returns stdout scanner
func StartEncode(prober Prober, startmux bool) error {
	if !Started {
		var cmd []string
		var err error
		if startmux {
			cmd, err = getMuxCommand(*opts)
		} else {
			cmd, err = getCommand(prober)
		}
		if err != nil {
			log.Errorf("StartEncode failed" + err.Error())
			return err
		}
		if err = deviceCommon.runCmdPipe(cmd, false); err == nil {
			Started = true
			return nil
		}
		log.Errorf("StartEncode failed" + err.Error())
		return err
	}
	log.Errorf("already started")
	return errors.New("already started")
}

// StopEncode stop ffmpeg process
func StopEncode() bool {
	if Started {
		deviceCommon.done <- true //send done signal and..
		<-deviceCommon.done       //wait for process done
		Started = false
		return true
	}
	log.Errorf("already stopped")
	return false
}

// NewProber returns prober for correct platform
func NewProber() Prober {
	var prober Prober
	switch runtime.GOOS {
	case "darwin":
		prober = &macProber
	default:
		panic("OS not supported")
	}
	deviceCommon.opts = opts
	Ffoutchan = make(chan Ffoutmsg)
	deviceCommon.probeDefaults()
	GetFfmpegDevices(prober)
	loadCommonConfig(configPath())
	return prober
}
