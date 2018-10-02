package ffprobe

import "strings"

type proberMac struct {
	devicesCmd      string
	recordCmdPrefix []string
	*proberCommon
}

var macProber = proberMac{
	devicesCmd:      "ffmpeg -f avfoundation -list_devices true -i ''",
	recordCmdPrefix: strings.Split("ffmpeg -y -f avfoundation -framerate 24", " "),
	proberCommon:    &deviceCommon,
}

func (pm *proberMac) getDevicesCmd() string {
	return pm.devicesCmd
}

func (pm *proberMac) getPrefixCmd() []string {
	// cmd = append(cmd, "-i")
	// // TODO: without audio
	// cmd = append(cmd, fmt.Sprintf("%d:%d", mp.options.vidIdx, mp.options.audIdx))
	// cmd = append(cmd, mp.recordCmdPostfix...)
	return pm.recordCmdPrefix
}
