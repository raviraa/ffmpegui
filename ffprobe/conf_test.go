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
		{"vp9-default", 0, nil, " 0_0.webm"},
		{"vp9-default", 0, nil, "-threads 8"},
		{"opus-default", 0, nil, "-map 0:a -c:a libopus 0_0.opus"},
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

func testInputs(ps ...string) []UIInput {
	var ips []UIInput
	for _, p := range ps {
		ips = append(ips, config.newInput(p))
	}
	return ips
}

type td struct {
	fname string
	want  string
}

func testOptions() map[*Options][]td {
	return map[*Options][]td{
		&Options{UIInputs: testInputs("opus-default")}: []td{
			{"getRecCmd", "-map 0:a -c:a aac"},
			{"getMuxCommand", "-i 0.opus -map 0:a -c copy 20"},
			{"getConcatCmd", "-i 0_0.aac -i 1_0.aac"},
			{"getConcatCmd", "libopus 0.opus"},
		},
		&Options{UIInputs: testInputs("vp9-default")}: []td{
			{"getRecCmd", "-c:v libx264"},
			{"getRecCmd", "-crf 0"},
			{"getMuxCommand", "-i 0.webm -map 0:v -c copy "},
			{"getConcatCmd", "[0:v:0][1:v:0]concat=n=2:v=1[out]"},
		},
		&Options{UIInputs: testInputs("opus-default", "vp9-default")}: []td{
			{"getRecCmd", "-c:a aac 1_0.aac"},
			{"getRecCmd", "1_1.mkv"},
			{"getMuxCommand", "-i 0.opus -i 1.webm -map 0:a -map 1:v -c copy "},
			{"getConcatCmd", "libopus 0.opus"},
		},
	}
}

func TestCmds(t *testing.T) {
	loadCommonConfig(cfgname)
	opts := testOptions()
	/* := []struct {
		plt string
		opts *Options
		want string
	}{
		{"mac", &Options{UIInputs: []UIInput{UIInput{Type: Audio, Presetidx: 1}}},
			"-f avfoundation -i none:0 -map 0:a -c:a libopus"},
		{"mac", &Options{UIInputs: []UIInput{UIInput{Type: Video, Presetidx: 3}}},
			"-f avfoundation -i 0:none -map 0:v -c:v vp9"},
		{"mac", &Options{UIInputs: []UIInput{UIInput{Type: Audio, Presetidx: 2}, UIInput{Type: Video, Presetidx: 3}}},
			"-f avfoundation -i 0:none -map 1:v -c:v vp9"},
	}*/
	for opt, tt := range opts {
		for _, d := range tt {
			t.Run(d.fname, func(t *testing.T) {
				var cmds []string
				var err error
				switch d.fname {
				case "getRecCmd":
					cmds, err = getRecCmd("mac", *opt)
				case "getMuxCommand":
					cmds, err = getMuxCommand(*opt)
				case "getConcatCmd":
					config.resumeCount = 1
					cmds, err = getConcatCmd(*opt, 0)
				}
				if err != nil {
					t.Error(err)
				}
				if got := strings.Join(cmds, " "); !strings.Contains(got, d.want) {
					t.Errorf("%s(%v) got %#v, should contain %#v", d.fname, *opt, got, d.want)
				}
			})
		}
	}
}
