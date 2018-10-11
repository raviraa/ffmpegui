package main

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/andlabs/ui"
	"github.com/raviraa/recordscreen/ffprobe"
)

func onStartClicked(btn *ui.Button) {
	log.Info("start clicked..")
	opts := ffprobe.Options{
		VidIdx:    cboxVid.Selected(),
		AudIdx:    cboxAud.Selected(),
		Container: cboxCtnr.Selected(),
	}
	log.Info(opts)
	ffprobe.SetOptions(opts)
	ffstderr, err := ffprobe.StartEncode(prober)
	readWritepipe(ffstderr)
	if err == nil {
		lblDesc.SetText("ffmpeg started succesfully")
	} else {
		lblDesc.SetText("ffmpeg start failed: " + err.Error())
	}
	// ctrlStatus.Append(strings.Join(cmd, " ") + "\n")
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
