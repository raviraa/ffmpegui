package ffprobe

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
)

func (pc ProberCommon) presetOpts(presetName string, avidx int) ([]string, *preset, error) {
	var s []string
	p, ok := pc.config.Presets[presetName]
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

func (pc ProberCommon) presetCmd(presetName string, avidx int) ([]string, error) {
	var s []string
	sopts, p, err := pc.presetOpts(presetName, avidx)
	if err != nil {
		return nil, err
	}

	//  -map 0:a  -c:a libopus ... fname.opus
	s = append(s, "-map", fmt.Sprintf("%d:%s", avidx, p.Avtype))
	s = append(s, sopts...)
	s = append(s, pc.config.mkFile(p.Fileext, pc.config.resumeCount, avidx))
	return s, nil
}

func (pc ProberCommon) inputCmd(pltip map[string]*input, uiip UIInput) ([]string, error) {
	var s []string
	avtype := uiip.Type.str()
	ip, ok := pltip[avtype]
	if !ok {
		return nil, fmt.Errorf("unknown platform device type " + avtype)
	}
	s = append(s, mapToString(pc.config.InputsDefaults)...)
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
func (pc ProberCommon) getRecCmd(plt string, opts Options) ([]string, error) {
	var s []string
	s = append(s, strings.Split(pc.config.Ffcmdprefix, " ")...)
	for avidx, uiip := range opts.UIInputs {
		ipcmd, err := pc.inputCmd(pc.config.Inputs[plt], uiip)
		if err != nil {
			return nil, err
		}
		pc, err := pc.presetCmd(capturePreset(&uiip), avidx)
		if err != nil {
			return nil, err
		}
		s = append(s, ipcmd...)
		s = append(s, pc...)
		avidx++
	}
	logi.Print(strings.Join(s, " "))
	return s, nil
}

func (pc ProberCommon) getConcatCmd(opts Options, avidx int) ([]string, error) {
	var s []string
	s = append(s, strings.Split(pc.config.Ffcmdprefix, " ")...)
	uiip := opts.UIInputs[avidx]
	typ, ext, err := pc.ctnrInfo(uiip, true)
	if err != nil {
		return nil, err
	}
	var filter string
	for i := 0; i <= pc.config.resumeCount; i++ {
		s = append(s, "-i")
		s = append(s, pc.config.mkFile(ext, i, avidx))
		filter += fmt.Sprintf("[%d:%s:0]", i, typ)
	}
	streams := "v=1:a=0"
	if typ == "a" {
		streams = "v=0:a=1"
	}

	s = append(s, "-filter_complex", filter+
		fmt.Sprintf("concat=n=%d:%s[out]", pc.config.resumeCount+1, streams))
	s = append(s, "-map", "[out]")
	psidx, err := pc.getPresetByIdx(uiip.Presetidx, uiip.Type)
	if err != nil {
		return nil, err
	}
	pccmd, pset, err := pc.presetOpts(psidx, avidx)
	if err != nil {
		return nil, err
	}
	s = append(s, pccmd...)
	s = append(s, pc.config.mkFile(pset.Fileext, avidx))
	return s, nil
}

// -i 0.webm  -i 1.webm -map 0:v -map 1:a -c copy  200601021504.webm
func (pc ProberCommon) getMuxCommand(opts Options) ([]string, error) {
	var s []string
	var fileext string
	s = append(s, strings.Split(pc.config.Ffcmdprefix, " ")...)
	for avidx, uiip := range opts.UIInputs {
		// -i 1.webm
		typ, ext, err := pc.ctnrInfo(uiip, false)
		if err != nil {
			return nil, err
		}
		// fname := fmt.Sprintf("%v.%v", avidx, ext)
		fname := pc.config.mkFile(ext, avidx)
		s = append(s, "-i", fname)
		if typ == "v" || fileext == "" {
			fileext = ext
		}
	}
	for avidx, uiip := range opts.UIInputs {
		typ, _, _ := pc.ctnrInfo(uiip, false)
		//-map 0:v
		s = append(s, "-map", fmt.Sprintf("%d:%s", avidx, typ))
	}
	// -c copy
	s = append(s, "-c", "copy")

	tnow := time.Now().Format("200601021504")
	s = append(s, pc.config.mkFile(fileext, tnow))

	return s, nil
}

func (pc ProberCommon) ctnrInfo(uiip UIInput, capture bool) (typ string, ext string, err error) {
	var p preset
	var ok bool
	if capture {
		p, ok = pc.config.Presets[capturePreset(&uiip)]
	} else {
		psidx, err := pc.getPresetByIdx(uiip.Presetidx, uiip.Type)
		if err != nil {
			return "", "", err
		}
		p, ok = pc.config.Presets[psidx]
	}
	if !ok {
		return "", "", fmt.Errorf("unknown preset")
	}
	return p.Avtype, p.Fileext, nil
}

func capturePreset(i *UIInput) string {
	return "capture-" + i.Type.str()
}
