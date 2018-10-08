package ffprobe

import "strings"

type proberMac struct {
	devicesCmd      string
	recordCmdPrefix []string
	*proberCommon
}

var macProber = proberMac{
	devicesCmd:      "ffmpeg -f avfoundation -list_devices true -i ''",
	recordCmdPrefix: strings.Split("ffmpeg -y -loglevel verbose -f avfoundation -framerate 24", " "),
	proberCommon:    &deviceCommon,
}

func (pm *proberMac) getDevicesCmd() string {
	return pm.devicesCmd
}

func (pm *proberMac) getPrefixCmd() []string {
	// conf
	return pm.recordCmdPrefix
}
