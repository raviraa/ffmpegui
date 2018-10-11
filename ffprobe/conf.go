package ffprobe

import (
	"fmt"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
)

var config *tomlConfig

type tomlConfig struct {
	Inputs         map[string]map[string]*input
	Presets        map[string]map[string]interface{}
	InputsDefaults map[string]interface{} `toml:"inputs_defaults"`
	Containers     map[string][]string
	options        Options
	Ffconf         ffconf
}

type ffconf struct {
	Ffcmdprefix string
}

type input struct {
	I string
	F string
}

func inputCmd(pltip map[string]*input, avtype string) (s []string) {
	ip, ok := pltip[avtype]
	if !ok {
		panic("unknown platform device type " + avtype)
	}
	s = append(s, mapToString(config.InputsDefaults)...)
	ipmap := make(map[string]interface{})
	for k, v := range structs.Map(ip) {
		if k != "I" && k != "F" {
			ipmap[k] = v
		}
	}
	s = append(s, mapToString(ipmap)...)
	devidx := config.options.VidIdx
	if avtype == "a" {
		devidx = config.options.AudIdx
	}
	s = append(s, "-f", ip.F)
	s = append(s, "-i", fmt.Sprintf(ip.I, devidx))
	return
}

func containerCmd(ctname string, avidx int) []string {
	var s []string
	presetNames, ok := config.Containers[ctname]
	if !ok {
		panic("unknown container key " + ctname)
	}

	presetName := presetNames[avidx]
	prset, ok := config.Presets[presetName]
	if !ok {
		panic(fmt.Sprintf("unknown preset %s in container %s", presetName, ctname))
	}
	fileext, ok1 := prset["fileext"]
	avtype, ok2 := prset["type"]
	codec, ok3 := prset["codec"]
	if !ok1 || !ok2 || !ok3 {
		panic("fileext, codec or type missing for " + presetName)
	}

	//  -map 0:a  -c:a libopus ... fname.opus
	s = append(s, "-map", fmt.Sprintf("%d:%s", avidx, avtype))
	s = append(s, fmt.Sprintf("-c:%s", avtype), codec.(string))
	for k, v := range prset {
		if k == "fileext" || k == "type" || k == "codec" {
			continue
		}
		s = append(s, "-"+k)
		s = append(s, fmt.Sprintf("%v", v))
	}
	s = append(s, fmt.Sprintf("%d.%s", avidx, fileext))

	return s
}

func getConfCmd(plt string) (s []string) {
	avidx := 0
	s = append(s, strings.Split(config.Ffconf.Ffcmdprefix, " ")...)
	if config.options.VidIdx != -1 {
		s = append(s, inputCmd(config.Inputs[plt], "v")...)
		s = append(s, containerCmd(getContainerIdx(config.options.Container), avidx)...)
		avidx++
	}
	if config.options.AudIdx != -1 {
		s = append(s, inputCmd(config.Inputs[plt], "a")...)
		s = append(s, containerCmd(getContainerIdx(config.options.Container), avidx)...)
		avidx++
		avidx++
	}
	fmt.Println(strings.Join(s, " "))
	return
}

// GetContainers returns available presets
func GetContainers() (cts []string) {
	for k := range config.Containers {
		cts = append(cts, k)
	}
	sort.Strings(cts)
	fmt.Println(cts)
	return
}

func getContainerIdx(cidx int) string {
	return GetContainers()[cidx]
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

// SetOptions configures extra encoder options
func SetOptions(opts Options) {
	config.options = opts
}

var cfgname string

func configFile() string {
	cfgname = "common_presets.toml"
	// https://github.com/shibukawa/configdir
	return "ffprobe/common_presets.toml" //TODO
}

// WriteConf saves config file to user conf dir
func WriteConf() {
}
