package ffprobe

import (
	"os"
	"os/user"
	"path"
	"syscall"
)

func (p *ProberCommon) probeDefaults() {
	u, _ := user.Current()
	p.opts.VidPath = path.Join(u.HomeDir, "Desktop")
	p.opts.Framerate = 24.0
	p.deviceKey = "input device"
	p.devicesKey = "devices:"
}

// StartMux concats all resume split streams to single stream,
// and then avstreams to final container
func (p *ProberCommon) StartMux() {
	go func() {
		for avidx := range opts.UIInputs {
			cmd, err := p.getConcatCmd(*opts, avidx)
			if err != nil {
				fferr(err)
				return
			}
			if err = p.runCmdPipe(cmd, "concat"); err != nil {
				fferr(err)
				return
			}
		}
		cmd, err := p.getMuxCommand(*opts)
		if err != nil {
			fferr(err)
			return
		}
		if err = p.runCmdPipe(cmd, "mux"); err != nil {
			fferr(err)
			return
		}
	}()
}

// StartEncode starts ffmpeg process with configured options
func (p *ProberCommon) StartEncode() {
	go func() {
		cmd, err := getCommand(p.prober, *p)
		if err != nil {
			loge.Printf("StartEncode failed" + err.Error())
			fferr(err)
			return
		}
		p.runCmdPipe(cmd, "rec")
	}()
	return
}

// KillEncode stop ffmpeg process
func (p *ProberCommon) KillEncode() {
	if p.cmd != nil {
		p.cmd.Process.Signal(syscall.SIGINT)
	} else {
		loge.Println("no process running")
	}
}

//RmTmpFiles removes intermediate stage avfiles
func (p *ProberCommon) RmTmpFiles() {
	toremove := make(map[string]bool)
	for i := 0; i < len(p.config.tmpFiles)-1; i++ {
		toremove[p.config.tmpFiles[i]] = true
	}
	logi.Print("Removing files ", toremove)
	for f := range toremove {
		if e := os.Remove(f); e != nil {
			loge.Print(e)
		}
	}
}
