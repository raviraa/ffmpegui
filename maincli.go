package main

import (
	"github.com/raviraa/ffmpegui/ffprobe"
)

var log = ffprobe.SetLogger()

/*
func beginCli() {
	// log.Info("Starting in CLI")
	prober := ffprobe.NewProber()
	uiip := ffprobe.UIInput{Type: ffprobe.Video, Devidx: 2, Presetidx: 3}
	ffprobe.SetInputs([]ffprobe.UIInput{uiip})
	stdout, _ := ffprobe.StartEncode(prober, false)
	readStdout(stdout)
	log.Info("before sleep")
	time.Sleep(5 * time.Second)
	log.Info("sending stop signal")
	ffprobe.StopEncode()
}

func readStdout(scanner *bufio.Scanner) {
	go func() {
		scanner.Split(scanLines)
		// scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			txt := scanner.Text()
			fmt.Println(txt)
			if strings.Contains(txt, "frame=") {
				fmt.Print(txt)
			}
		}
	}()
}

func mainCli() {

	beginCli()
}

*/
