package ffprobe

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/shibukawa/configdir"
)

// Avtype is audio video type for streams
type Avtype int

/*
type Enum string TODO: from config to enum

const (
	EnumAbc Enum = "abc"
	EnumDef Enum = "def"
)*/

const (
	// Video type
	Video Avtype = iota
	// Audio type
	Audio
)

type tomlConfig struct {
	Inputs         map[string]map[string]*input
	Presets        map[string]preset
	InputsDefaults map[string]interface{} `toml:"inputs_defaults"`
	Containers     map[string][]string
	resumeCount    int
	Ffcmdprefix    string
	tmpFiles       []string
}

// UIInput has config from ui for each input stream
type UIInput struct {
	Devidx    int
	Presetidx int
	Type      Avtype
	// Video_size string TODO2
}

type preset struct {
	Fileext string
	Avtype  string //TODO string change to Avtype
	Codec   string
	Options map[string]interface{}
}

// Options for ffmpeg, saved/restored on encode start from WriteUIOpts()
type Options struct {
	UIInputs  []UIInput
	Framerate float64
	VidPath   string
}

// input config for platform capture device
type input struct {
	I string
	F string
}

func strType(typ string) Avtype {
	if typ == "a" {
		return Audio
	}
	return Video
}
func (typ Avtype) str() string {
	if typ == Audio {
		return "a"
	}
	return "v"
}

func (c *tomlConfig) mkFile(fext string, fparts ...interface{}) string {
	var s string
	for i, p := range fparts {
		if i != 0 {
			s += "_"
		}
		s += fmt.Sprintf("%v", p)
	}
	s += "." + fext
	s = filepath.Join(opts.VidPath, s)
	c.tmpFiles = append(c.tmpFiles, s)
	return s
}

func (c *tomlConfig) newInput(pset string, pc ProberCommon) UIInput {
	ip := UIInput{Devidx: -1}
	p, ok := c.Presets[pset]
	if ok {
		ip.Type = strType(p.Avtype)
		for i, v := range pc.GetPresets(ip.Type) {
			if v == pset {
				ip.Presetidx = i
			}
		}
	}
	return ip
}

// GetPresets returns available presets
func (pc ProberCommon) GetPresets(typ Avtype) (ps []string) {
	for k, p := range pc.config.Presets {
		if p.Avtype == typ.str() && !strings.HasPrefix(k, "capture-") {
			ps = append(ps, k)
		}
	}
	sort.Strings(ps)
	return
}

func (pc ProberCommon) getPresetByIdx(cidx int, typ Avtype) (string, error) {
	ps := pc.GetPresets(typ)
	if cidx >= len(ps) || cidx < 0 {
		return "", errors.New("invalid preset")
	}
	return ps[cidx], nil
}

func mapToString(m map[string]interface{}) (s []string) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := fmt.Sprintf("%v", m[k])
		if len(vs) != 0 {
			s = append(s, "-"+strings.ToLower(k), vs)
		}
	}
	return s
}

func loadCommonConfig(confstr string) *tomlConfig {
	config := &tomlConfig{}
	if _, err := toml.Decode(confstr, config); err != nil {
		panic(err)
	}
	readUIOpts()
	return config
}

func presetTomlStr() string {
	fname := configDir(presetFname)
	b, e := ioutil.ReadFile(fname)
	if e != nil {
		ioutil.WriteFile(fname, []byte(defaultPresetStr), 0644)
		return defaultPresetStr
	}
	return string(b)
}

func configDir(fname string) string {
	cdir := configdir.New("", "ffmpegui")
	folders := cdir.QueryFolders(configdir.Global)
	return filepath.Join(folders[0].Path, fname)
}

// WriteUIOpts saves ui selection options
func WriteUIOpts() error {
	fname := configDir(uiOptsFname)
	logi.Printf("writing conf to %s", fname)
	var b bytes.Buffer
	enc := toml.NewEncoder(&b)
	err := enc.Encode(opts)
	if err != nil {
		loge.Printf("writing ui options failed: %s", err)
		return err
	}
	err = ioutil.WriteFile(fname, b.Bytes(), 0644)
	if err != nil {
		loge.Printf("writing ui options failed: %s", err)
		return err
	}
	return nil
}

func readUIOpts() error {
	_, err := toml.DecodeFile(configDir(uiOptsFname), opts)
	if err != nil {
		loge.Printf("reading UI opts failed %s", err)
		return err
	}
	logi.Print("read ui options: ", opts)
	return nil
}
