package ffprobe

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
)

func presetCmd(presetName string, avidx int) ([]string, error) {
	var s []string
	prset, ok := config.Presets[presetName]
	if !ok {
		return nil, fmt.Errorf("unknown preset %s", presetName)
	}
	fileext, ok1 := prset["fileext"]
	avtype, ok2 := prset["type"]
	codec, ok3 := prset["codec"]
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("fileext, codec or type missing for " + presetName)
	}

	//  -map 0:a  -c:a libopus ... fname.opus
	s = append(s, "-map", fmt.Sprintf("%d:%s", avidx, avtype))
	s = append(s, fmt.Sprintf("-c:%s", avtype), codec.(string))
	if avtype == "v" {
		s = append(s, "-framerate", fmt.Sprintf("%v", opts.Framerate))
	}
	for k, v := range prset {
		if k == "fileext" || k == "type" || k == "codec" {
			continue
		}
		s = append(s, "-"+k)
		s = append(s, fmt.Sprintf("%v", v))
	}
	s = append(s, fmt.Sprintf("%d.%s", avidx, fileext))

	return s, nil
}

func inputCmd(pltip map[string]*input, uiip UIInput) ([]string, error) {
	var s []string
	avtype := avtypestr(uiip.Type)
	ip, ok := pltip[avtype]
	if !ok {
		return nil, fmt.Errorf("unknown platform device type " + avtype)
	}
	s = append(s, mapToString(config.InputsDefaults)...)
	s = append(s, "-framerate", fmt.Sprintf("%v", opts.Framerate))
	ipmap := make(map[string]interface{})
	for k, v := range structs.Map(ip) {
		if k != "I" && k != "F" {
			ipmap[k] = v
		}
	}
	s = append(s, mapToString(ipmap)...)
	s = append(s, "-f", ip.F)
	s = append(s, "-i", fmt.Sprintf(ip.I, uiip.Devidx))
	return s, nil
}

func getConfCmd(plt string, opts Options) ([]string, error) {
	var s []string
	s = append(s, strings.Split(opts.Ffcmdprefix, " ")...)
	for avidx, uiip := range opts.UIInputs {
		ipcmd, err := inputCmd(config.Inputs[plt], uiip)
		if err != nil {
			return nil, err
		}
		psidx, err := getPresetByIdx(uiip.Presetidx)
		if err != nil {
			return nil, err
		}
		pc, err := presetCmd(psidx, avidx)
		if err != nil {
			return nil, err
		}
		s = append(s, ipcmd...)
		s = append(s, pc...)
		avidx++
	}
	log.Info(strings.Join(s, " "))
	return s, nil
}

// -i 0.webm  -i 1.webm -map 0:v -map 1:a -c copy  200601021504.webm
func getMuxCommand(opts Options) ([]string, error) {
	var s []string
	var fileext string
	s = append(s, strings.Split(opts.Ffcmdprefix, " ")...)
	for avidx, uiip := range opts.UIInputs {
		psidx, err := getPresetByIdx(uiip.Presetidx)
		if err != nil {
			return nil, err
		}
		ps, ok := config.Presets[psidx]
		if !ok {
			return nil, fmt.Errorf("wrong preset index %v", psidx)
		}
		fname := fmt.Sprintf("%v.%v", avidx, ps["fileext"])
		// -i 1.webm
		s = append(s, "-i", fname)
		if ps["type"] == "v" || fileext == "" {
			fileext = ps["fileext"].(string)
		}
	}
	for avidx, uiip := range opts.UIInputs {
		psidx, _ := getPresetByIdx(uiip.Presetidx)
		ps, _ := config.Presets[psidx]
		//-map 0:v
		s = append(s, "-map", fmt.Sprintf("%d:%s", avidx, ps["type"]))
	}
	// -c copy
	s = append(s, "-c", "copy")

	tnow := time.Now()
	s = append(s, fmt.Sprintf("%s.%s", tnow.Format("200601021504"), fileext))

	return s, nil
}
