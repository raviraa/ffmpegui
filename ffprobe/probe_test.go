package ffprobe

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func init() {
	SetLogger()
}

func Test_Mac_getDevices(t *testing.T) {
	want := []string{"0  Built-in Microphone"}
	got := GetFfmpegDevices(&macProber).Audios
	if !reflect.DeepEqual(got, want) {
		t.Errorf("proberKeys.getAudios() = %v, want %v", got, want)
	}
}

func Test_Mac_getCmd(t *testing.T) {
	macprober := NewProber()
	loadCommonConfig(cfgname)
	SetInputs([]UIInput{UIInput{Type: Audio}})
	defer func() { opts = &Options{} }()
	want := "ffmpeg -benchmark -y -loglevel verbose -thread_queue_size 512 -framerate 24 -f avfoundation -i none:0 -map 0:v -c:v libx264 -framerate 24 -preset faster 0.mkv"
	cmds, _ := getCommand(macprober)
	if got := strings.Join(cmds, " "); !reflect.DeepEqual(got, want) {
		t.Errorf("getCommand = %#v, want %v", got, want)
	}
}

func Test_getPlatformProber(t *testing.T) {
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

func Test_parseFfmpegDevices(t *testing.T) {
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
			got := parseFfmpegDeviceType(&macProber, tt.dtype)
			if !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("parseFfmpegDevices() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StartProcessFail(t *testing.T) {
	prober := NewProber()
	config.Ffcmdprefix = "pytho"
	scanner, _ := StartEncode(prober, false)
	if scanner != nil {
		t.Errorf("expected process fail")
	}
}

func Test_ProcessInterrupt(t *testing.T) {
	prober := NewProber()
	config.Ffcmdprefix = "sleep 10"
	tbeg := time.Now().UnixNano()
	StartEncode(prober, false)
	if !StopEncode() || (time.Now().UnixNano()-tbeg > 1e9) {
		t.Error("process interrupt failed or too slow")
	}
}

func Test_StartProcessOutput(t *testing.T) {
	prober := NewProber()
	config.Ffcmdprefix = "ls asdf1234"
	scanner, _ := StartEncode(prober, false)
	var ffout string
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			txt := scanner.Text()
			ffout += txt
		}
		done <- true
	}()
	<-done
	if ffout != "ls: asdf1234: No such file or directory" {
		t.Error("wrong process output" + ffout)
	}
}
