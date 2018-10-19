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
	log.Info(args)
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

func (mp *proberCommon) runCmdPipe(cmdstr []string, wait bool) error {
	log.Info(cmdstr)
	cmd := exec.Command(cmdstr[0], cmdstr[1:]...)
	// cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	outscanner := bufio.NewScanner(stdout)
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		if wait {
			return
		}
		<-mp.done
		log.Infof("Received signal. Sending SIGINT...")
		cmd.Process.Signal(syscall.SIGINT)
		cmd.Wait()
		//indicate signal sent, SIGUSR1 value is ignored
		mp.done <- true
		log.Info("ffmpeg stopped. ", err)
	}()

	readWritepipe(outscanner)
	if wait {
		err := cmd.Wait()
		log.Info("ffmpeg stopped. ", err)
	}
	return nil
}

// Ffoutmsg ffmpeg output line
type Ffoutmsg struct {
	FrameUpdate bool
	Msg         string
}

// Ffoutchan output lines from ffmpeg process
var Ffoutchan chan Ffoutmsg

func readWritepipe(scanner *bufio.Scanner) {
	go func() {
		count := 0
		scanner.Split(scanLines)
		frmtxt := ""
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.Contains(txt, "frame=") {
				if count%3 == 0 {
					frmtxt = txt
					// setStatus(txt)
					Ffoutchan <- Ffoutmsg{true, txt}
				}
				count++
			} else {
				// addInfo(txt + "\n")
				Ffoutchan <- Ffoutmsg{false, txt + "\n"}
			}
		}
		Ffoutchan <- Ffoutmsg{false, frmtxt + "\nDone.\n\n\n\n"}
		Started = false
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
