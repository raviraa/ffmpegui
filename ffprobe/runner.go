package ffprobe

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func runCmd(args []string, ignExitCode bool) (string, error) {
	// args = append([]string{"-c"}, args...)
	// cmd := exec.Command("/bin/sh", args...)
	var out bytes.Buffer
	logi.Print(args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil && !ignExitCode {
		fmt.Println(out.String())
		return "", err
	}
	return out.String(), nil
}

func runCmdStr(cmd string, ignExitCode bool) string {
	if str, err := runCmd(strings.Split(cmd, " "), ignExitCode); err == nil {
		return str
	}
	return ""
}

func (mp *ProberCommon) runCmdPipe(cmdstr []string, cmdtype string) error {
	logi.Print(cmdstr)
	cmd := exec.Command(cmdstr[0], cmdstr[1:]...)
	mp.cmd = cmd
	defer func() { mp.cmd = nil }()
	cmd.Stdout = os.Stdout
	stdout, err := cmd.StderrPipe()
	if err != nil {
		fferr(err)
		return err
	}
	outscanner := bufio.NewScanner(stdout)
	if err := cmd.Start(); err != nil {
		fferr(err)
		return err
	}
	readWritepipe(outscanner)
	Ffoutchan <- Ffoutmsg{Ffother, fmt.Sprintf(
		"\n#####\t\t%s\t\t######\n######\t\t\t#####\n", cmdtype)}
	err = cmd.Wait()
	excode := 0
	if exitError, ok := err.(*exec.ExitError); ok {
		ws := exitError.Sys().(syscall.WaitStatus)
		excode = ws.ExitStatus()
	}
	logi.Print("ffmpeg stopped ", excode, err)

	if err == nil || excode == 255 {
		Ffoutchan <- Ffoutmsg{Ffdone, cmdtype}
	} else {
		Ffoutchan <- Ffoutmsg{Fferr, cmdtype}
	}
	return err
}

// Ffouttype is line output type
type Ffouttype int

const (
	//Fferr when scanner return error
	Fferr = iota
	//Ffdone when process exited without err
	Ffdone
	//FframeUpdate is frame, encoding speed update line
	FframeUpdate
	//Ffother anything else
	Ffother
)

// Ffoutmsg ffmpeg output line
type Ffoutmsg struct {
	Typ Ffouttype
	Msg string
}

// Ffoutchan output lines from ffmpeg process
var Ffoutchan chan Ffoutmsg

func readWritepipe(scanner *bufio.Scanner) {
	go func() {
		count := 0
		scanner.Split(scanLines)
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.Contains(txt, "frame=") {
				if count%3 == 0 {
					Ffoutchan <- Ffoutmsg{FframeUpdate, txt}
				}
				count++
			} else {
				Ffoutchan <- Ffoutmsg{Ffother, txt + "\n"}
			}
		}
		if e := scanner.Err(); e != nil {
			fferr(e)
		}
	}()
}

func fferr(e error) {
	Ffoutchan <- Ffoutmsg{Fferr, e.Error()}
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
