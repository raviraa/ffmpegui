package main

/*
import (
	"strings"

	"github.com/andlabs/ui"
	"github.com/raviraa/recordscreen/ffprobe"
)

var (
	lblDesc    *ui.Label
	ctrlStatus *ui.MultilineEntry
	cboxAud    *ui.Combobox
	cboxVid    *ui.Combobox
	mwin       *ui.Window
	// txtStatus  strings.Builder
	prober ffprobe.Prober
)


func onStartClicked(btn *ui.Button) {
	log.Info("start clicked..")
	opts := ffprobe.Options{
		vidIdx: cboxVid.Selected(),
		audIdx: cboxAud.Selected(),
	}
	prober.setOptions(opts)
	cmd := prober.getCommand()
	prober.start()
	ctrlStatus.Append(strings.Join(cmd, " ") + "\n")
}

func makeInputForm() *ui.Group {
	group := ui.NewGroup("Options")
	group.SetMargined(true)
	vbox := ui.NewVerticalBox()
	group.SetChild(vbox)

	entryForm := ui.NewForm()
	vbox.Append(entryForm, false)
	entryForm.SetPadded(true)

	cboxVid = ui.NewCombobox()
	entryForm.Append("Video Device", cboxVid, false)
	entryForm.Append("Record Audio", ui.NewCheckbox(""), false)
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
	prober = getPlatformProber()
	prober.probeDevices()
	ui.QueueMain(func() {
		lblDesc.SetText(prober.getVersion())
		devs := prober.getDevices()
		for _, s := range devs.audios {
			cboxAud.Append(s)
		}
		for _, s := range devs.videos {
			cboxVid.Append(s)
		}
	})
	// cmd := prober.getCommand()
	// log.Info("Using cmd:", cmd)
}

func setupUI() {
	mwin = ui.NewWindow("Record screen, webcam using ffmpeg", 400, 400, false)
	mwin.SetMargined(true)
	mwin.OnClosing(func(mw *ui.Window) bool {
		prober.stop()
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
		prober.stop()
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

*/
func main() {
	mainCli()
	// mainUI()
}
