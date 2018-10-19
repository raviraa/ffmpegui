package ffprobe

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/shibukawa/configdir"
)

var config *tomlConfig

// Avtype is audio video type for streams
type Avtype int

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
}

// UIInput has config from ui for each input stream
type UIInput struct {
	Devidx    int
	Presetidx int
	Type      Avtype
	// Video_size string TODO
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
	return s
}

func (c *tomlConfig) newInput(pset string) UIInput {
	ip := UIInput{Devidx: -1}
	p, ok := c.Presets[pset]
	if ok {
		ip.Type = strType(p.Avtype)
		for i, v := range GetPresets(ip.Type) {
			if v == pset {
				ip.Presetidx = i
			}
		}
	}
	return ip
}

// GetPresets returns available presets
func GetPresets(typ Avtype) (ps []string) {
	for k, p := range config.Presets {
		if p.Avtype == typ.str() {
			ps = append(ps, k)
		}
	}
	sort.Strings(ps)
	return
}

func getPresetByIdx(cidx int, typ Avtype) (string, error) {
	ps := GetPresets(typ)
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

func loadCommonConfig(fname string) {
	config = &tomlConfig{}
	if _, err := toml.DecodeFile(fname, config); err != nil {
		fmt.Println(err)
		return
	}
	readUIOpts()
}

var cfgname string
var uiOptsFname string

func configPath() string {
	// https://github.com/shibukawa/configdir
	return "ffprobe/common_presets.toml" //TODO
}

func configDir() *configdir.Config {
	cdir := configdir.New("", "ffmpegui")
	folders := cdir.QueryFolders(configdir.Global)
	return folders[0]
}

// WriteUIOpts saves config file to user conf dir
func WriteUIOpts() error {
	dir := configDir()
	log.Infof("writing conf to %s", dir.Path)
	var b bytes.Buffer
	enc := toml.NewEncoder(&b)
	err := enc.Encode(opts)
	if err != nil {
		log.Errorf("writing ui options failed: %s", err)
		return err
	}
	err = dir.WriteFile(uiOptsFname, b.Bytes())
	if err != nil {
		log.Errorf("writing ui options failed: %s", err)
		return err
	}
	return nil
}

//go:generate go get -u github.com/jteeuwen/go-bindata/...
//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -prefix assets/ assets/
func readUIOpts() error {
	_, err := toml.DecodeFile(filepath.Join(configDir().Path, uiOptsFname), opts)
	if err != nil {
		log.Errorf("reading UI opts failed %s", err)
		return err
	}
	log.Infof("read ui options: ", opts)
	return nil
}
