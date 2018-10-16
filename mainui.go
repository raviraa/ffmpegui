package main

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/raviraa/ffmpegui/ffprobe"
)

type inputControls struct {
	inpbox   *ui.Box
	ffinputs []*ffprobe.UIInput
}

var (
	lblDesc          *ui.Label
	ctrlStatus       *ui.MultilineEntry
	mwin             *ui.Window
	inps             *inputControls
	prober           ffprobe.Prober
	updateFrameCount = 9 // update frame status every n ffmpeg updates
)

func addInput(idx int, typ ffprobe.Avtype) *ui.Group {
	uiip := ffprobe.UIInput{Devidx: -1, Presetidx: -1, Type: typ}
	inps.ffinputs = append(inps.ffinputs, &uiip)
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
	cboxPresets := ui.NewCombobox()
	entryForm.Append("Presets", cboxPresets, false)
	cboxPresets.OnSelected(func(cb *ui.Combobox) {
		uiip.Presetidx = cb.Selected()
	})
	for _, s := range ffprobe.GetPresets() {
		cboxPresets.Append(s)
	}

	btnfile := ui.NewButton("...")
	entryForm.Append("Save As", btnfile, false)
	btnfile.OnClicked(func(*ui.Button) {
		filename := ui.SaveFile(mwin)
		log.Info("selected file", filename)
	})

	return group
}

func beginUIProbe() {
	log.Info("Starting in GUI mode")
	prober = ffprobe.NewProber()
	lblDesc.SetText(ffprobe.GetVersion())
}

func setupUI() {
	mwin = ui.NewWindow("Record screen, webcam using ffmpeg", 400, 400, false)
	mwin.SetMargined(true)
	mwin.OnClosing(func(mw *ui.Window) bool {
		ffprobe.StopEncode()
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
		addInput(inpidx, ffprobe.Video)
	})
	btnaddaud := ui.NewButton("Add Audio")
	btnaddbox.Append(btnaddaud, false)
	btnaddaud.OnClicked(func(*ui.Button) {
		inpidx := len(inps.ffinputs)
		addInput(inpidx, ffprobe.Audio)
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

	btnstart := ui.NewButton("Start")
	btnhbox.Append(btnstart, false)
	btnstart.OnClicked(onStartClicked)
	btnstop := ui.NewButton("Stop")
	btnstop.OnClicked(onStopClicked)
	btnhbox.Append(btnstop, false)

	mwin.Show()
	go beginUIProbe()
}
func mainUI() {
	if err := ui.Main(setupUI); err != nil {
		log.Panic(err)
	}
	log.Info("Exiting")
}

func main() {
	// mainCli()
	mainUI()
}
