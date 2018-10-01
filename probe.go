package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Prober finds devices based on platform
type Prober interface {
	probeDevices()
	getDevices() devices
	setOptions(opts options)
	getCommand() []string
	getVersion() string
	start()
	stop()
	runCmdPipe([]string) //error //TODO private not expose?
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
	done             chan bool
	started          bool
	devices
	options
}

var macProber = proberKeys{
	deviceKey:       "input device",
	devicesCmd:      "ffmpeg -f avfoundation -list_devices true -i ''",
	recordCmdPrefix: strings.Split("ffmpeg -y -f avfoundation -framerate 24", " "),
	// TODO: -s scale, resolution input oputpu
	recordCmdPostfix: strings.Split("-framerate 25 -s 1920x1080 TODO.mkv", " "), //-preset ultrafast aaa.mkv
	devicesKey:       "devices:",
	done:             make(chan bool),
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
	// TODO: without audio
	cmd = append(cmd, fmt.Sprintf("%d:%d", mp.options.vidIdx, mp.options.audIdx))
	cmd = append(cmd, mp.recordCmdPostfix...)
	if len(mp.ffmpegOpts) > 1 {
		cmd = append(cmd, mp.ffmpegOpts...)
	}
	// runCmdPipe(strings.Split("ls -lR ..", " "))
	return cmd
}

func (mp *proberKeys) runCmdPipe(cmdstr []string) { //TODO ret error
	log.Info(cmdstr)
	cmd := exec.Command(cmdstr[0], cmdstr[1:]...)
	cmd.Stdout = os.Stdout //TODO use logger and txtarea
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go func() {
		<-mp.done
		log.Info("Received stop signal. Sending quit")
		defer stdin.Close()
		io.WriteString(stdin, "q")
		// io.WriteString(stdin, "+") TODO send keys for more/less status
		cmd.Wait()
		mp.done <- true //signal process done
	}()
}

func (mp *proberKeys) start() {
	if !mp.started {
		cmd := mp.getCommand()
		mp.started = true
		mp.runCmdPipe(cmd)
	} else {
		log.Errorf("already started")
	}
}

func (mp *proberKeys) stop() {
	if mp.started {
		mp.done <- true //send done signal and..
		<-mp.done       //wait for process done
		mp.started = false
	} else {
		log.Errorf("already stopped")
	}
}

func getPlatformProber() Prober {
	var prober Prober = &macProber
	return prober
}
