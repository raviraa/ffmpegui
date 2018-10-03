package ffprobe

import (
	"reflect"
	"testing"
	"time"
)

func Test_Mac_getDevices(t *testing.T) {
	want := []string{"0  Built-in Microphone"}
	got := GetFfmpegDevices(&macProber).audios
	if !reflect.DeepEqual(got, want) {
		t.Errorf("proberKeys.getAudios() = %v, want %v", got, want)
	}
}

func Test_Mac_getCmd(t *testing.T) {
	macprober := GetPlatformProber()
	want := []string{"ffmpeg", "-y"}
	opts := Options{VidIdx: 1, AudIdx: 0}
	SetOptions(opts)
	if got := getCommand(macprober); !reflect.DeepEqual(got[0:2], want) {
		t.Errorf("proberKeys.getAudios() = %#v, want %v", got[0:2], want)
	}
}

func Test_getPlatformProber(t *testing.T) {
	tests := []struct {
		name string
		want Prober
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPlatformProber(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPlatformProber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFfmpegDevices(t *testing.T) {
	tests := []struct {
		name  string
		dtype string
		want  []string
	}{
		{"mac", "audio", []string{"0  Built-in Microphone"}},
		{"macvideo", "video", []string{"0  FaceTime HD Camera", "1  Capture screen 0", "2  Capture screen 1", "3  Capture screen 2"}},
	}
	// TODO: mock cmd run
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFfmpegDeviceType(&macProber, tt.dtype)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFfmpegDevices() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StartProcessFail(t *testing.T) {
	prober := &proberMac{recordCmdPrefix: []string{"pytho"}}
	scanner, _ := StartEncode(prober)
	if scanner != nil {
		t.Errorf("expected process fail")
	}
}

func Test_ProcessInterrupt(t *testing.T) {
	prober := &proberMac{recordCmdPrefix: []string{"sleep", "10"}} //ls", "-l"}}
	tbeg := time.Now().UnixNano()
	StartEncode(prober)
	if !StopEncode() || (time.Now().UnixNano()-tbeg > 1e9) {
		t.Error("process interrupt failed or too slow")
	}
}

func Test_StartProcessOutput(t *testing.T) {
	prober := &proberMac{recordCmdPrefix: []string{"echo", "success\n"}}
	scanner, _ := StartEncode(prober)
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
	if ffout != "success" {
		t.Error("wrong process output" + ffout)
	}
}
