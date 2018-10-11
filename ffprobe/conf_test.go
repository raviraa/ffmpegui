package ffprobe

import (
	"reflect"
	"testing"
)

func TestConfInputs(t *testing.T) {
	loadCommonConfig("common_presets.toml")
	SetOptions(Options{})
	tests := []struct {
		plt, name string
		want      string
	}{
		{"mac", "a", "avfoundation"},
		{"mac", "v", "avfoundation"},
	}
	for _, tt := range tests {
		if got := config.Inputs[tt.plt][tt.name]; !reflect.DeepEqual(got.F, tt.want) {
			t.Errorf("conf input got %#v, want %#v", got, tt.want)
		}
	}
}

func TestConfInputCmds(t *testing.T) {
	loadCommonConfig("common_presets.toml")
	SetOptions(Options{})
	tests := []struct {
		plt, name string
		want      []string
	}{
		{"mac", "a", []string{"-thread_queue_size", "512", "-f", "avfoundation", "-i", "none:0"}},
		// {"mac", "v", []string{"-thread_queue_size", "512", "-framerate", "25", "-video_size", "1920x1080", "-i", "0:none", "-f", "avfoundation"}},
	}
	for _, tt := range tests {
		if got := inputCmd(config.Inputs[tt.plt], tt.name); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("conf input got %#v, want %#v", got, tt.want)
		}
	}
}

func TestContainer(t *testing.T) {
	loadCommonConfig("common_presets.toml")
	SetOptions(Options{})
	tests := []struct {
		name  string
		avidx int
		want  []string
	}{
		{"webm - vp9 default with no audio", 0, []string{"-map", "0:v", "-c:v", "vp9", "0.webm"}},
		{"webm - vp9 default with opus default", 0, []string{"-map", "0:v", "-c:v", "vp9", "0.webm"}},
		{"webm - vp9 default with opus default", 1, []string{"-map", "1:a", "-c:a", "libopus", "1.opus"}},
	}
	for _, tt := range tests {
		if got := containerCmd(tt.name, tt.avidx); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("container cmd %s, got %#v, want %#v", tt.name, got, tt.want)
		}
	}
}
