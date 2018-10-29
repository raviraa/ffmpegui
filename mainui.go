package main

import (
	"fmt"
	"log"

	"github.com/andlabs/ui"
	"github.com/raviraa/ffmpegui/ffprobe"
)

type inputControls struct {
	inpbox   *ui.Box
	ffinputs []*ffprobe.UIInput
}

var (
	lblDesc    *ui.Label
	ctrlStatus *ui.MultilineEntry
	mwin       *ui.Window
	btnpause   *ui.Button
	btnstop    *ui.Button
	btnstart   *ui.Button
	inps       *inputControls
	prober     ffprobe.ProberCommon
)

func addInput(idx int, typ ffprobe.Avtype, cnfuiip *ffprobe.UIInput) {
	var uiip ffprobe.UIInput
	if cnfuiip == nil {
		uiip = ffprobe.UIInput{Devidx: -1, Presetidx: -1, Type: typ} //TODO use newInput and move ffinputs to ffprobe
	} else {
		uiip = *cnfuiip
	}
	inps.ffinputs = append(inps.ffinputs, &uiip)
	fmt.Println(uiip)
	group := ui.NewGroup(fmt.Sprintf("Input %d", idx))
	inps.inpbox.Append(group, false)
	group.SetMargined(true)
	vbox := ui.NewVerticalBox()
	group.SetChild(vbox)

	entryForm := ui.NewForm()
	vbox.Append(entryForm, false)
	entryForm.SetPadded(true)

	cboxDevs := ui.NewCombobox()
	entryForm.Append("Devices", cboxDevs, false)
	devs := ffprobe.GetFfmpegDevices(prober)
	switch typ {
	case ffprobe.Audio:
		for _, s := range devs.Audios {
			cboxDevs.Append(s)
		}
	case ffprobe.Video:
		for _, s := range devs.Videos {
			cboxDevs.Append(s)
		}
	}
	cboxDevs.OnSelected(func(cb *ui.Combobox) {
		uiip.Devidx = cb.Selected()
	})
	if uiip.Devidx != -1 {
		cboxDevs.SetSelected(uiip.Devidx)
	}
	cboxPresets := ui.NewCombobox()
	entryForm.Append("Presets", cboxPresets, false)
	cboxPresets.OnSelected(func(cb *ui.Combobox) {
		uiip.Presetidx = cb.Selected()
	})
	for _, s := range prober.GetPresets(uiip.Type) {
		cboxPresets.Append(s)
	}
	if uiip.Presetidx != -1 {
		cboxPresets.SetSelected(uiip.Presetidx)
	}

	btnfile := ui.NewButton("Remove")
	entryForm.Append("", btnfile, false)

}

func beginUIProbe() {
	logi.Print("Starting in GUI mode")
	prober = ffprobe.NewProber()
	lblDesc.SetText(ffprobe.GetVersion())
	startFfoutReader()
	for idx, uiip := range ffprobe.GetInputs() {
		addInput(idx, uiip.Type, &uiip)
	}
}

func setupUI() {
	mwin = ui.NewWindow("Record screen, webcam using ffmpeg", 400, 400, false)
	mwin.SetMargined(true)
	mwin.OnClosing(func(mw *ui.Window) bool {
		prober.KillEncode()
		mwin.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		mwin.Destroy()
		return true
	})

	mvbox := ui.NewVerticalBox()
	mwin.SetChild(mvbox)

	btnaddbox := ui.NewHorizontalBox()
	mvbox.Append(btnaddbox, false)
	btnaddinp := ui.NewButton("Add Video")
	btnaddbox.Append(btnaddinp, false)
	btnaddinp.OnClicked(func(*ui.Button) {
		inpidx := len(inps.ffinputs)
		addInput(inpidx, ffprobe.Video, nil)
	})
	btnaddaud := ui.NewButton("Add Audio")
	btnaddbox.Append(btnaddaud, false)
	btnaddaud.OnClicked(func(*ui.Button) {
		inpidx := len(inps.ffinputs)
		addInput(inpidx, ffprobe.Audio, nil)
	})
	inps = &inputControls{}
	inps.inpbox = ui.NewHorizontalBox()
	mvbox.Append(inps.inpbox, false)
	// mvbox.Append(makeInputForm(), false)
	mvbox.Append(ui.NewHorizontalSeparator(), false)
	ctrlStatus = ui.NewMultilineEntry()
	mvbox.Append(ctrlStatus, true)
	ctrlStatus.SetReadOnly(true)

	lblDesc = ui.NewLabel("")
	mvbox.Append(lblDesc, false)
	mvbox.Append(ui.NewHorizontalSeparator(), false)
	btnhbox := ui.NewHorizontalBox()
	mvbox.Append(btnhbox, false)

	btnstart = ui.NewButton("Start")
	btnhbox.Append(btnstart, false)
	btnstart.OnClicked(onStartClicked)
	btnpause = ui.NewButton("Pause")
	btnpause.Disable()
	btnhbox.Append(btnpause, false)
	btnpause.OnClicked(onPauseClicked)
	btnstop = ui.NewButton("Stop")
	btnstop.Disable()
	btnstop.OnClicked(onStopClicked)
	btnhbox.Append(btnstop, false)

	mwin.Show()
	ui.QueueMain(func() {
		beginUIProbe()
	})
}

func mainUI() {
	if err := ui.Main(setupUI); err != nil {
		log.Panic(err)
	}
	logi.Print("Exiting")
}

func main() {
	// mainCli()
	mainUI()
}

func setStatus(s string) {
	ui.QueueMain(func() {
		lblDesc.SetText(s)
	})
}

func addInfo(s string) {
	ui.QueueMain(func() {
		ctrlStatus.Append(s)
	})
}
