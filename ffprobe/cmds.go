package ffprobe

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
)

func presetOpts(presetName string, avidx int) ([]string, *preset, error) {
	var s []string
	p, ok := config.Presets[presetName]
	if !ok {
		return nil, nil, fmt.Errorf("unknown preset %s", presetName)
	}
	s = append(s, fmt.Sprintf("-c:%s", p.Avtype), p.Codec)
	for k, v := range p.Options {
		s = append(s, "-"+k)
		s = append(s, fmt.Sprintf("%v", v))
	}
	return s, &p, nil
}

func presetCmd(presetName string, avidx int) ([]string, error) {
	var s []string
	sopts, p, err := presetOpts(presetName, avidx)
	if err != nil {
		return nil, err
	}

	//  -map 0:a  -c:a libopus ... fname.opus
	s = append(s, "-map", fmt.Sprintf("%d:%s", avidx, p.Avtype))
	s = append(s, sopts...)
	s = append(s, config.mkFile(p.Fileext, config.resumeCount, avidx))
	return s, nil
}

func inputCmd(pltip map[string]*input, uiip UIInput) ([]string, error) {
	var s []string
	avtype := uiip.Type.str()
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

// recording from capture devices, always uses capture-a or capture-v preset
func getRecCmd(plt string, opts Options) ([]string, error) {
	var s []string
	s = append(s, strings.Split(config.Ffcmdprefix, " ")...)
	for avidx, uiip := range opts.UIInputs {
		ipcmd, err := inputCmd(config.Inputs[plt], uiip)
		if err != nil {
			return nil, err
		}
		pc, err := presetCmd(capturePreset(&uiip), avidx)
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

func getConcatCmd(opts Options, avidx int) ([]string, error) {
	var s []string
	s = append(s, strings.Split(config.Ffcmdprefix, " ")...)
	uiip := opts.UIInputs[avidx]
	typ, ext, err := ctnrInfo(uiip, true)
	if err != nil {
		return nil, err
	}
	var filter string
	for i := 0; i <= config.resumeCount; i++ {
		s = append(s, "-i")
		s = append(s, config.mkFile(ext, i, avidx))
		filter += fmt.Sprintf("[%d:%s:0]", i, typ)
	}
	streams := "v=1:a=0"
	if typ == "a" {
		streams = "v=0:a=1"
	}

	s = append(s, "-filter_complex", filter+
		fmt.Sprintf("concat=n=%d:%s[out]", config.resumeCount+1, streams))
	s = append(s, "-map", "[out]")
	psidx, err := getPresetByIdx(uiip.Presetidx, uiip.Type)
	if err != nil {
		return nil, err
	}
	pccmd, pc, err := presetOpts(psidx, avidx)
	if err != nil {
		return nil, err
	}
	s = append(s, pccmd...)
	s = append(s, config.mkFile(pc.Fileext, avidx))
	return s, nil
}

// -i 0.webm  -i 1.webm -map 0:v -map 1:a -c copy  200601021504.webm
func getMuxCommand(opts Options) ([]string, error) {
	var s []string
	var fileext string
	s = append(s, strings.Split(config.Ffcmdprefix, " ")...)
	for avidx, uiip := range opts.UIInputs {
		// -i 1.webm
		typ, ext, err := ctnrInfo(uiip, false)
		if err != nil {
			return nil, err
		}
		fname := fmt.Sprintf("%v.%v", avidx, ext)
		s = append(s, "-i", fname)
		if typ == "v" || fileext == "" {
			fileext = ext
		}
	}
	for avidx, uiip := range opts.UIInputs {
		typ, _, _ := ctnrInfo(uiip, false)
		//-map 0:v
		s = append(s, "-map", fmt.Sprintf("%d:%s", avidx, typ))
	}
	// -c copy
	s = append(s, "-c", "copy")

	tnow := time.Now().Format("200601021504")
	s = append(s, config.mkFile(fileext, tnow))

	return s, nil
}

func ctnrInfo(uiip UIInput, capture bool) (typ string, ext string, err error) {
	var p preset
	var ok bool
	if capture {
		p, ok = config.Presets[capturePreset(&uiip)]
	} else {
		psidx, err := getPresetByIdx(uiip.Presetidx, uiip.Type)
		if err != nil {
			return "", "", err
		}
		p, ok = config.Presets[psidx]
	}
	if !ok {
		return "", "", fmt.Errorf("unknown preset")
	}
	return p.Avtype, p.Fileext, nil
}

func capturePreset(i *UIInput) string {
	return "capture-" + i.Type.str()
}
