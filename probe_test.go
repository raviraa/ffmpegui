package main

import (
	"reflect"
	"testing"
)

func Test_proberKeys_getAudios(t *testing.T) {
	tests := []struct {
		name string
		mp   proberKeys
		want map[int]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := proberKeys{}
			if got := mp.getAudios(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("proberKeys.getAudios() = %v, want %v", got, tt.want)
			}
		})
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
		{"mac", "audio", []string{"0] Built-in Microphone"}},
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
