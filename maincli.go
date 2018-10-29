package main

import (
	"github.com/raviraa/ffmpegui/ffprobe"
)

var logi, loge = ffprobe.GetLoggers()

/*
func beginCli() {
	// logi.Print("Starting in CLI")
	prober := ffprobe.NewProber()
	uiip := ffprobe.UIInput{Type: ffprobe.Video, Devidx: 2, Presetidx: 3}
	ffprobe.SetInputs([]ffprobe.UIInput{uiip})
	stdout, _ := ffprobe.StartEncode(prober, false)
	readStdout(stdout)
	logi.Print("before sleep")
	time.Sleep(5 * time.Second)
	logi.Print("sending stop signal")
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
