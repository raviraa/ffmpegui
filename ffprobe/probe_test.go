package ffprobe

import (
	"reflect"
	"strings"
	"testing"
)

func TestMac_getDevices(t *testing.T) {
	want := []string{"0  Built-in Microphone"}
	got := GetFfmpegDevices(NewProber()).Audios
	if !reflect.DeepEqual(got, want) {
		t.Errorf("proberKeys.getAudios() = %v, want %v", got, want)
	}
}

func TestMac_getCmd(t *testing.T) {
	macprober := NewProber()
	opts.Framerate = 24
	macprober.config.Ffcmdprefix = "ffmpeg -benchmark -y -loglevel verbose"
	macprober.SetInputs([]UIInput{UIInput{Type: Video}}, 0)
	defer func() { opts = &Options{} }()
	want := "ffmpeg -benchmark -y -loglevel verbose -thread_queue_size 512 -framerate 24 -f avfoundation -i 0:none -map 0:v -c:v libx264"
	cmds, _ := getCommand(macprober.prober, macprober)
	if got := strings.Join(cmds, " "); !strings.Contains(got, want) {
		t.Errorf("getCommand(mac) = %#v, should contain %#v", got, want)
	}
}

func TestGetPlatformProber(t *testing.T) {
	tests := []struct {
		name string
		want Prober
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProber(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPlatformProber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFfmpegDevices(t *testing.T) {
	tests := []struct {
		name  string
		dtype string
		want  string
	}{
		{"mac", "audio", "0  Built-in Microphone"},
		{"macvideo", "video", "0  FaceTime HD Camera"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFfmpegDeviceType(NewProber(), tt.dtype)
			if !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("parseFfmpegDevices() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartProcessFail(t *testing.T) {
	prober := NewProber()
	prober.config.Ffcmdprefix = "pytho"
	prober.StartEncode()
	if msg := <-Ffoutchan; msg.Typ != Fferr {
		t.Errorf("expected process fail")
	}
}

func TestProcessInterrupt(t *testing.T) {
	prober := NewProber()
	prober.config.Ffcmdprefix = "sleep 10"
	prober.StartEncode()
	<-Ffoutchan
	prober.KillEncode()
	if msg := <-Ffoutchan; msg.Typ != Fferr {
		t.Error("process interrupt failed", msg)
	}
}

func TestStartProcessOutput(t *testing.T) {
	prober := NewProber()
	prober.config.Ffcmdprefix = "ls asdf1234"
	prober.StartEncode()
	<-Ffoutchan
	ffout := <-Ffoutchan
	if ffout.Msg != "ls: -c:v: No such file or directory\n" {
		t.Error("wrong process output", ffout)
	}
}

func ffoutreadUntil(count int, typ Ffouttype) int {
	i := 0
	for {
		m := <-Ffoutchan
		if m.Typ == typ {
			i++
			if i == count {
				return count
			}
		}
	}
}

func TestStartMux(t *testing.T) {
	pc := NewProber()
	pc.SetInputs(testInputs(pc, "opus-default"), 0)
	pc.config.Ffcmdprefix = "echo asdf1234"
	pc.StartMux()
	if dones := ffoutreadUntil(2, Ffdone); dones != 2 {
		t.Error("expected 2 ffmpeg process ", dones)
	}
}
