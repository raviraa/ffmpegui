package ffprobe

import (
	"fmt"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
)

var config tomlConfig

type tomlConfig struct {
	Inputs          map[string]*input
	Presets         map[string]map[string]interface{}
	Inputs_defaults map[string]interface{}
	Containers      map[string][]string
}

type input struct {
	Video_size string
	Framerate  string
	F          string
	I          string
}

func inputString(ip *input) (s []string) {
	s = append(s, mapToString(config.Inputs_defaults)...)
	ipmap := make(map[string]interface{})
	for k, v := range structs.Map(ip) {
		if k != "I" && k != "F" {
			ipmap[k] = v
		}
	}
	s = append(s, mapToString(ipmap)...)
	s = append(s, "-i", ip.I)
	s = append(s, "-f", ip.F)
	return
}

func containerCmd(ctname string) []string {
	var s []string
	presetNames, ok := config.Containers[ctname]
	if !ok {
		panic("unknown container key " + ctname)
	}

	for idx, presetName := range presetNames {
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
		s = append(s, "-map", fmt.Sprintf("%d:%s", idx, avtype))
		s = append(s, fmt.Sprintf("-c:%s", avtype), codec.(string))
		for k, v := range prset {
			if k == "fileext" || k == "type" || k == "codec" {
				continue
			}
			s = append(s, "-"+k)
			s = append(s, fmt.Sprintf("%v", v))
		}
		s = append(s, fmt.Sprintf("%d.%s", idx, fileext))
	}

	return s
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

func loadCommonConfig() {
	if _, err := toml.DecodeFile("common_presets.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
}

// WriteConf saves config file to user conf dir
func WriteConf() {
	// https://github.com/shibukawa/configdir
}
