package ffprobe

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Prober has logic that changes based on platform
type Prober interface {
	getDevicesCmd() string
	getFfmpegCmd(ProberCommon) ([]string, error)
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

var (
	cfgname     = "common_presets.toml"
	uiOptsFname = "uiopts.toml"
	presetFname = "presets.toml"
	logi        = log.New(os.Stdout, "INFO: ", log.Lshortfile|log.Ltime)
	loge        = log.New(os.Stderr, "ERROR: ", log.Lshortfile|log.Ltime)
)

// ProberCommon has all options; common and  platform prober
type ProberCommon struct {
	devicesKey string
	deviceKey  string //input device
	opts       *Options
	cmd        *exec.Cmd
	prober     Prober
	config     *tomlConfig
	Devices
}

var opts *Options = &Options{} // TODO1 move to devcommon

//SetInputs to set configure input streams
func (pc ProberCommon) SetInputs(uiips []UIInput, resumeCount int) {
	opts.UIInputs = uiips
	pc.config.resumeCount = resumeCount
}

// GetInputs gets
func GetInputs() []UIInput {
	return opts.UIInputs
}

// GetLoggers returns info and error logger
func GetLoggers() (*log.Logger, *log.Logger) {
	return logi, loge
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
func GetFfmpegDevices(p ProberCommon) Devices {
	devs := Devices{}
	devs.Audios = parseFfmpegDeviceType(p, "audio")
	devs.Videos = parseFfmpegDeviceType(p, "video")
	return devs
}

func parseFfmpegDeviceType(p ProberCommon, dtype string) []string {
	devs := map[string][]string{}
	res := runCmdStr(p.prober.getDevicesCmd(), true)
	resLines := strings.Split(res, "\n")
	filterfn := func(s string) bool { return strings.Contains(s, p.deviceKey) }
	resLines = filterList(resLines, filterfn)
	var currDevType, dtypeKey string
	for _, ln := range resLines {
		if strings.Contains(ln, p.devicesKey) {
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
	logi.Printf("%s: %+v\n", dtype, devs[dtypeKey])
	return devs[dtypeKey]
}

// GetVersion returns ffmpeg version
func GetVersion() string {
	return "ffmpeg 1234.22" //TODO
}

func getCommand(prober Prober, pc ProberCommon) ([]string, error) {
	return prober.getFfmpegCmd(pc)
}

// NewProber returns prober for correct platform
func NewProber() ProberCommon {
	var dc = ProberCommon{}
	switch runtime.GOOS {
	case "darwin":
		dc.prober = newProberMac()
	default:
		panic("OS not supported")
	}
	dc.opts = opts
	Ffoutchan = make(chan Ffoutmsg)
	dc.probeDefaults()
	GetFfmpegDevices(dc)
	dc.config = loadCommonConfig(presetTomlStr())
	return dc
}
