package main

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/raviraa/ffmpegui/ffprobe"
)

func onStopClicked(btn *ui.Button) {
	stopped := ffprobe.StopEncode() //TODO enable/disable start/stop buttons
	if stopped {
		// err := ffprobe.StartEncode(prober, true)
		err := ffprobe.StartMux(prober)
		if err == nil {
			ctrlStatus.Append("==============\nMuxing streams\n")
		} else {
			lblDesc.SetText("ffmpeg mux start failed: " + err.Error())
		}
	}
}

var paused bool
var resumeCount int

func onPauseClicked(btn *ui.Button) {
	if paused {
		// resume logic
		resumeCount++
		onStartClicked(btn)
		paused = false
		btnpause.SetText("Pause")
	} else {
		if ffprobe.StopEncode() {
			paused = true
			btnpause.SetText("Resume")
		}
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
	log.Info(fmt.Sprintf("%#v", inpsar))
	ffprobe.SetInputs(inpsar, resumeCount)
	ffprobe.WriteUIOpts()
	// return
	err := ffprobe.StartEncode(prober, false)
	if err == nil {
		ctrlStatus.Append("==============\n")
		lblDesc.SetText("ffmpeg started succesfully")
	} else {
		lblDesc.SetText("ffmpeg start failed: " + err.Error())
	}
}

func startFfoutReader() {
	go func() {
		for {
			out := <-ffprobe.Ffoutchan
			if out.FrameUpdate {
				setStatus(out.Msg)
			} else {
				addInfo(out.Msg)
			}
		}
	}()
}
