package ffprobe

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
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
			"-i 0.opus -i 1.webm -map 0:a -map 1:v -c copy "},
	}
	for _, tt := range tests {
		cmds, _ := getMuxCommand(tt.opts)
		if got := strings.Join(cmds, " "); !strings.Contains(got, tt.want) {
			t.Errorf("conf input got %#v, should contain %#v", got, tt.want)
		}
	}
}

func encodeExpected(
	t *testing.T, label string, val interface{}, wantStr string,
) {
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		t.Errorf("encode failed: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, wantStr) {
		t.Errorf("%s: should contain\n-----\n%q\n-----\nbut got\n-----\n%q\n-----\n",
			label, wantStr, got)
	}
}

func TestUIInputEncode(t *testing.T) {
	tests := []struct {
		opts Options
		want []string
	}{
		{Options{Framerate: 24.0, UIInputs: []UIInput{UIInput{Type: Audio, Presetidx: 1}}},
			[]string{"[[UIInputs]]\n  Devidx = 0\n  Presetidx = 1\n  Type = 1\n", "Framerate = 24.0"}},
		{Options{UIInputs: []UIInput{UIInput{Type: Video, Presetidx: 3}}},
			[]string{"[[UIInputs]]\n  Devidx = 0\n  Presetidx = 3\n  Type = 0\n"}},
	}
	for _, tt := range tests {
		for _, want := range tt.want {
			encodeExpected(t, "encode UIInput", tt.opts, want)
		}
	}
}
