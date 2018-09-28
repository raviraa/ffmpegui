package main

import (
	"reflect"
	"testing"
)

func Test_Mac_getAudios(t *testing.T) {
	macprober := getPlatformProber()
	want := []string{"0  Built-in Microphone"}
	macprober.probeDevices()
	if got := macprober.getDevices().audios; !reflect.DeepEqual(got, want) {
		t.Errorf("proberKeys.getAudios() = %v, want %v", got, want)
	}
}

func Test_Mac_getCmd(t *testing.T) {
	macprober := getPlatformProber()
	want := []string{"ffmpeg", "-y", "-report", "-f", "avfoundation", "-framerate", "24", "-i", "1", "-framerate", "25", "-s", "1920x1080", "TODO.mkv"}
	macprober.probeDevices()
	opts := options{vidIdx: 1}
	macprober.setOptions(opts)
	if got := macprober.getCommand(); !reflect.DeepEqual(got, want) {
		t.Errorf("proberKeys.getAudios() = %#v, want %v", got, want)
	}
}

func Test_getPlatformProber(t *testing.T) {
	tests := []struct {
		name string
		want Prober
	}{
		// TODO: Add test cases HERE!!!
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPlatformProber(); !reflect.DeepEqual(got, tt.want) {
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
		// {"macvideo", "video", []string{"0] Built-in Microphone"}},
	}
	// TODO: mock cmd
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFfmpegDevices(macProber, tt.dtype); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFfmpegDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}
