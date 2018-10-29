package ffprobe

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func init() {
	devnull, e := os.Create(os.DevNull)
	if e != nil {
		panic(e)
	}
	logi = log.New(devnull, "", log.Lshortfile|log.Ltime)
}

func TestConfInputs(t *testing.T) {
	pc := NewProber()
	tests := []struct {
		plt, name string
		want      string
	}{
		{"mac", "a", "avfoundation"},
		{"mac", "v", "avfoundation"},
	}
	for _, tt := range tests {
		if got := pc.config.Inputs[tt.plt][tt.name]; !reflect.DeepEqual(got.F, tt.want) {
			t.Errorf("conf input got %#v, want %#v", got, tt.want)
		}
	}
}

func TestConfInputCmds(t *testing.T) {
	pc := NewProber()
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
		inps, err := pc.inputCmd(pc.config.Inputs[tt.plt], tt.uiip)
		if got := strings.Join(inps, " "); !strings.Contains(got, tt.want) || !reflect.DeepEqual(err, tt.err) {
			t.Errorf("conf input got %#v, should contain %#v", got, tt.want)
		}
	}
}

func TestContainer(t *testing.T) {
	pc := NewProber()
	opts.VidPath = ""
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
		cmds, err := pc.presetCmd(tt.name, tt.avidx)
		got := strings.Join(cmds, " ")
		if !strings.Contains(got, tt.want) || !reflect.DeepEqual(err, tt.err) {
			t.Errorf("container cmd %s, got (%#v, %#v) want %#v", tt.name, got, err, tt.want)
		}
	}
}

func testInputs(pc ProberCommon, ps ...string) []UIInput {
	var ips []UIInput
	for _, p := range ps {
		ips = append(ips, pc.config.newInput(p, pc))
	}
	return ips
}

type td struct {
	fname string
	want  string
}

func testOptions(pc ProberCommon) map[*Options][]td {
	return map[*Options][]td{
		&Options{UIInputs: testInputs(pc, "opus-default")}: []td{
			{"getRecCmd", "-map 0:a -c:a aac"},
			{"getMuxCommand", "-i 0.opus -map 0:a -c copy 20"},
			{"getConcatCmd", "[0:a:0][1:a:0]concat=n=2:v=0:a=1[out]"},
			{"getConcatCmd", "-i 0_0.aac -i 1_0.aac"},
			{"getConcatCmd", "libopus 0.opus"},
		},
		&Options{UIInputs: testInputs(pc, "vp9-default")}: []td{
			{"getRecCmd", "-c:v libx264"},
			{"getRecCmd", "-crf 0"},
			{"getMuxCommand", "-i 0.webm -map 0:v -c copy "},
			{"getConcatCmd", "-i 0_0.mkv"},
			{"getConcatCmd", "-c:v vp9"},
			{"getConcatCmd", "[0:v:0][1:v:0]concat=n=2:v=1:a=0[out]"},
		},
		&Options{UIInputs: testInputs(pc, "opus-default", "vp9-default")}: []td{
			{"getRecCmd", "-c:a aac 1_0.aac"},
			{"getRecCmd", "1_1.mkv"},
			{"getMuxCommand", "-i 0.opus -i 1.webm -map 0:a -map 1:v -c copy "},
			{"getConcatCmd", "libopus 0.opus"},
			{"mkFiles", "0_0.aac 1_0.aac 0.opus 0_1.mkv 1_1.mkv 1.webm"},
		},
	}
}

func TestCmds(t *testing.T) {
	pc := NewProber()
	opts.VidPath = ""
	pc.config.resumeCount = 1
	topts := testOptions(pc)
	for opt, tt := range topts {
		for _, d := range tt {
			t.Run(d.fname, func(t *testing.T) {
				var cmds []string
				var err error
				switch d.fname {
				case "getRecCmd":
					cmds, err = pc.getRecCmd("mac", *opt)
				case "getMuxCommand":
					cmds, err = pc.getMuxCommand(*opt)
				case "getConcatCmd":
					cmds, err = pc.getConcatCmd(*opt, 0)
				case "mkFiles":
					pc.config.tmpFiles = nil
					_, err = pc.getConcatCmd(*opt, 0)
					_, err = pc.getConcatCmd(*opt, 1)
					cmds = pc.config.tmpFiles
				}
				if err != nil {
					t.Error(err)
				}
				if got := strings.Join(cmds, " "); !strings.Contains(got, d.want) {
					t.Errorf("%s(%v)\n got %#v\n should contain %#v", d.fname, *opt, got, d.want)
				}
			})
		}
	}
}
