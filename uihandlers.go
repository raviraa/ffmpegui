package main

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/raviraa/ffmpegui/ffprobe"
)

var paused bool
var resumeCount int
var reqMux bool

func onStopClicked(btn *ui.Button) {
	reqMux = true
	prober.KillEncode()
	btnpause.Disable()
	btnstop.Disable()
}

func onPauseClicked(btn *ui.Button) {
	if paused {
		// resume logic
		resumeCount++
		onStartClicked(btn)
		paused = false
		btnpause.SetText("Pause")
		btnstop.Enable()
	} else {
		prober.KillEncode()
		paused = true
		btnpause.SetText("Resume")
		btnstop.Disable()
	}
}

func onStartClicked(btn *ui.Button) {
	if len(inps.ffinputs) == 0 {
		lblDesc.SetText("No input streams to start")
		return
	}

	var inpsar []ffprobe.UIInput
	for _, i := range inps.ffinputs {
		inpsar = append(inpsar, *i)
	}
	logi.Print(fmt.Sprintf("%#v", inpsar))
	prober.SetInputs(inpsar, resumeCount)
	ffprobe.WriteUIOpts()
	prober.StartEncode()
	btn.Disable()
	btnpause.Enable()
	btnstop.Enable()
}

func startFfoutReader() {
	go func() {
		for {
			out := <-ffprobe.Ffoutchan
			switch out.Typ {
			case ffprobe.FframeUpdate:
				setStatus(out.Msg)
			case ffprobe.Ffother:
				addInfo(out.Msg)
			case ffprobe.Ffdone:
				fmt.Println(out)
				if reqMux {
					reqMux = false
					prober.StartMux()
				} else if out.Msg == "mux" {
					resumeCount = 0
					btnstart.Enable()
					prober.RmTmpFiles()
				}
			default:
				loge.Print("unknown msg", out)
			}
		}
	}()
}
