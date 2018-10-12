package ffprobe

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
)

func containerCmd(presetName string, avidx int) []string {
	var s []string
	prset, ok := config.Presets[presetName]
	if !ok {
		panic(fmt.Sprintf("unknown preset %s", presetName))
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

func inputCmd(pltip map[string]*input, uiip *UIInput) (s []string) {
	avtype := avtypestr(uiip.Type)
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
	s = append(s, "-f", ip.F)
	s = append(s, "-i", fmt.Sprintf(ip.I, uiip.Devidx))
	return
}

func getConfCmd(plt string, opts Options) (s []string) {
	s = append(s, strings.Split(config.Ffcmdprefix, " ")...)
	for avidx, uiip := range opts.UIInputs {
		s = append(s, inputCmd(config.Inputs[plt], uiip)...)
		s = append(s, containerCmd(getPresetByIdx(uiip.Presetidx), avidx)...)
		avidx++
	}
	log.Info(strings.Join(s, " "))
	return
}
