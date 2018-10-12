package ffprobe

import (
	"reflect"
	"strings"
	"testing"
)

func TestConfInputs(t *testing.T) {
	loadCommonConfig("common_presets.toml")
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
	loadCommonConfig(cfgname)
	tests := []struct {
		plt  string
		uiip UIInput
		want []string
	}{
		{"mac", UIInput{Type: Audio}, []string{"-thread_queue_size", "512", "-f", "avfoundation", "-i", "none:0"}},
		{"mac", UIInput{Type: Video}, []string{"-thread_queue_size", "512", "-f", "avfoundation", "-i", "0:none"}},
	}
	for _, tt := range tests {
		if got := inputCmd(config.Inputs[tt.plt], &tt.uiip); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("conf input got %#v, want %#v", got, tt.want)
		}
	}
}

func TestContainer(t *testing.T) {
	loadCommonConfig(cfgname)
	// SetOptions(Options{})
	tests := []struct {
		name  string
		avidx int
		want  []string
	}{
		{"vp9-default", 0, []string{"-map", "0:v", "-c:v", "vp9", "0.webm"}},
		{"opus-default", 1, []string{"-map", "1:a", "-c:a", "libopus", "1.opus"}},
	}
	for _, tt := range tests {
		if got := containerCmd(tt.name, tt.avidx); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("container cmd %s, got %#v, want %#v", tt.name, got, tt.want)
		}
	}
}

func TestGetConfCmd(t *testing.T) {
	loadCommonConfig(cfgname)
	tests := []struct {
		plt  string
		opts *Options
		want string
	}{
		{"mac", &Options{UIInputs: []*UIInput{&UIInput{Type: Audio, Presetidx: 0}}},
			"-f avfoundation -i none:0 -map 0:a -c:a libopus 0.opus"},
		{"mac", &Options{UIInputs: []*UIInput{&UIInput{Type: Video, Presetidx: 2}}},
			"-f avfoundation -i 0:none -map 0:v -c:v vp9 0.webm"},
		{"mac", &Options{UIInputs: []*UIInput{&UIInput{Type: Audio, Presetidx: 1}, &UIInput{Type: Video, Presetidx: 2}}},
			"-f avfoundation -i 0:none -map 1:v -c:v vp9 1.webm"},
	}
	for _, tt := range tests {
		if got := strings.Join(getConfCmd(tt.plt, *tt.opts), " "); !strings.Contains(got, tt.want) {
			t.Errorf("conf input got %#v, should contain %#v", got, tt.want)
		}
	}
}
