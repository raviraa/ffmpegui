package ffprobe

import (
	"fmt"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
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
	Presets        map[string]map[string]interface{}
	InputsDefaults map[string]interface{} `toml:"inputs_defaults"`
	Containers     map[string][]string
	Ffcmdprefix    string
	Framerate      float64
}

// UIInput has config from ui for each input stream
type UIInput struct {
	Devidx    int
	Presetidx int
	Type      Avtype
}

type input struct {
	I string
	F string
}

func avtypestr(typ Avtype) string {
	if typ == Audio {
		return "a"
	}
	return "v"
}

// GetPresets returns available presets
func GetPresets() (cts []string) {
	for k := range config.Presets {
		cts = append(cts, k)
	}
	sort.Strings(cts)
	fmt.Println(cts)
	return
}

func getPresetByIdx(cidx int) string {
	return GetPresets()[cidx]
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
}

var cfgname string

func configFile() string {
	// https://github.com/shibukawa/configdir
	return "ffprobe/common_presets.toml" //TODO
}

// WriteConf saves config file to user conf dir
func WriteConf() {
}
