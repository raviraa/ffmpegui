package ffprobe

import (
	"errors"
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
		err  error
		want string
	}{
		{"mac", UIInput{Type: Audio}, nil, "-f avfoundation -i none:0"},
		{"mac", UIInput{Type: Video}, nil, "-f avfoundation -i 0:none"},
	}
	for _, tt := range tests {
		inps, err := inputCmd(config.Inputs[tt.plt], tt.uiip)
		if got := strings.Join(inps, " "); !strings.Contains(got, tt.want) || !reflect.DeepEqual(err, tt.err) {
			t.Errorf("conf input got %#v, should contain %#v", got, tt.want)
		}
	}
}

func TestContainer(t *testing.T) {
	loadCommonConfig(cfgname)
	tests := []struct {
		name  string
		avidx int
		err   error
		want  string
	}{
		{"vp9-default", 0, nil, "-map 0:v -c:v vp9"},
		{"vp9-default", 0, nil, " 0.webm"},
		{"opus-default", 0, nil, "-map 0:a -c:a libopus 0.opus"},
		{"wrongpreset", 1, errors.New("unknown preset wrongpreset"), ""},
	}
	for _, tt := range tests {
		cmds, err := presetCmd(tt.name, tt.avidx)
		got := strings.Join(cmds, " ")
		if !strings.Contains(got, tt.want) || !reflect.DeepEqual(err, tt.err) {
			t.Errorf("container cmd %s, got (%#v, %#v) want %#v", tt.name, got, err, tt.want)
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
		{"mac", &Options{UIInputs: []UIInput{UIInput{Type: Audio, Presetidx: 1}}},
			"-f avfoundation -i none:0 -map 0:a -c:a libopus"},
		{"mac", &Options{UIInputs: []UIInput{UIInput{Type: Video, Presetidx: 3}}},
			"-f avfoundation -i 0:none -map 0:v -c:v vp9"},
		{"mac", &Options{UIInputs: []UIInput{UIInput{Type: Audio, Presetidx: 2}, UIInput{Type: Video, Presetidx: 3}}},
			"-f avfoundation -i 0:none -map 1:v -c:v vp9"},
	}
	for _, tt := range tests {
		cmds, _ := getConfCmd(tt.plt, *tt.opts)
		// TODO check err
		if got := strings.Join(cmds, " "); !strings.Contains(got, tt.want) {
			t.Errorf("conf input got %#v, should contain %#v", got, tt.want)
		}
	}
}

func TestGetMuxCommand(t *testing.T) {
	loadCommonConfig(cfgname)
	tests := []struct {
		opts Options
		want string
	}{
		{Options{UIInputs: []UIInput{UIInput{Type: Audio, Presetidx: 1}}},
			"-i 0.opus -map 0:a -c copy 20"},
		{Options{UIInputs: []UIInput{UIInput{Type: Video, Presetidx: 3}}},
			"-i 0.webm -map 0:v -c copy "},
		{Options{UIInputs: []UIInput{UIInput{Type: Audio, Presetidx: 2}, UIInput{Type: Video, Presetidx: 3}}},
			"-i 0.opus -map 0:a -i 1.webm -map 1:v -c copy "},
	}
	for _, tt := range tests {
		cmds, _ := getMuxCommand(tt.opts)
		if got := strings.Join(cmds, " "); !strings.Contains(got, tt.want) {
			t.Errorf("conf input got %#v, should contain %#v", got, tt.want)
		}
	}
}
