package main

import (
	"github.com/andlabs/ui"
	"github.com/raviraa/recordscreen/ffprobe"
)

var (
	lblDesc    *ui.Label
	ctrlStatus *ui.MultilineEntry
	cboxAud    *ui.Combobox
	cboxVid    *ui.Combobox
	cboxCtnr   *ui.Combobox
	mwin       *ui.Window
	// txtStatus  strings.Builder
	prober           ffprobe.Prober
	updateFrameCount = 9 // update frame status every n ffmpeg updates
)

func makeInputForm() *ui.Group {
	group := ui.NewGroup("Options")
	group.SetMargined(true)
	vbox := ui.NewVerticalBox()
	group.SetChild(vbox)

	entryForm := ui.NewForm()
	vbox.Append(entryForm, false)
	entryForm.SetPadded(true)

	cboxCtnr = ui.NewCombobox()
	entryForm.Append("Profiles", cboxCtnr, false)
	cboxVid = ui.NewCombobox()
	entryForm.Append("Video Device", cboxVid, false)
	cboxAud = ui.NewCombobox()
	entryForm.Append("Audio Device", cboxAud, false)

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
	prober = ffprobe.GetPlatformProber()
	ui.QueueMain(func() {
		// lblDesc.SetText(ffprobe.getversio)TODO
		devs := ffprobe.GetFfmpegDevices(prober)
		for _, s := range devs.Audios {
			cboxAud.Append(s)
		}
		for _, s := range devs.Videos {
			cboxVid.Append(s)
		}
		for _, s := range ffprobe.GetContainers() {
			cboxCtnr.Append(s)
		}
	})
	// cmd := prober.getCommand()
	// log.Info("Using cmd:", cmd)
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

	lblDesc = ui.NewLabel("")
	mvbox.Append(lblDesc, false)

	mvbox.Append(makeInputForm(), false)
	mvbox.Append(ui.NewHorizontalSeparator(), false)
	ctrlStatus = ui.NewMultilineEntry()
	mvbox.Append(ctrlStatus, true)
	ctrlStatus.SetReadOnly(true)

	mvbox.Append(ui.NewHorizontalSeparator(), false)
	btnhbox := ui.NewHorizontalBox()
	mvbox.Append(btnhbox, false)

	btnstart := ui.NewButton("Start")
	btnhbox.Append(btnstart, false)
	btnstart.OnClicked(onStartClicked)
	btnstop := ui.NewButton("Stop")
	btnstop.OnClicked(func(btn *ui.Button) {
		ffprobe.StopEncode() //TODO chekc err
	})
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
