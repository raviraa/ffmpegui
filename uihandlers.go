package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/andlabs/ui"
	"github.com/raviraa/ffmpegui/ffprobe"
)

func onStopClicked(btn *ui.Button) {
	stopped := ffprobe.StopEncode() //TODO enable/disable start/stop buttons
	if stopped {
		ffstderr, err := ffprobe.StartEncode(prober, true)
		if err == nil {
			ctrlStatus.Append("==============\nMuxing streams\n")
			readWritepipe(ffstderr)
		} else {
			lblDesc.SetText("ffmpeg mux start failed: " + err.Error())
		}

	}
}

func onStartClicked(btn *ui.Button) {
	var inpsar []ffprobe.UIInput
	for _, i := range inps.ffinputs {
		inpsar = append(inpsar, *i)
	}
	log.Info(fmt.Sprintf("%#v", inpsar))
	ffprobe.SetInputs(inpsar)
	ffstderr, err := ffprobe.StartEncode(prober, false)
	if err == nil {
		ctrlStatus.Append("==============\n")
		lblDesc.SetText("ffmpeg started succesfully")
		readWritepipe(ffstderr)
	} else {
		lblDesc.SetText("ffmpeg start failed: " + err.Error())
	}
}

func readWritepipe(scanner *bufio.Scanner) {
	go func() {
		count := 0
		scanner.Split(scanLines)
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.Contains(txt, "frame=") {
				if count%updateFrameCount == 0 {
					ui.QueueMain(func() {
						lblDesc.SetText(txt)
					})
				}
				count++
			} else {
				ui.QueueMain(func() {
					ctrlStatus.Append(txt + "\n")
				})
			}
		}
		panic("todo")
		// ffprobe.Started = false
	}()
}

// scanLines splits into lines for either \r or \n
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
