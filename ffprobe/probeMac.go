package ffprobe

type proberMac struct {
	devicesCmd string
	*proberCommon
}

var macProber = proberMac{
	devicesCmd:   "ffmpeg -f avfoundation -list_devices true -i ''", //TODO move to conf file
	proberCommon: &deviceCommon,
}

func (pm *proberMac) getDevicesCmd() string {
	return pm.devicesCmd
}

func (pm *proberMac) getFfmpegCmd() ([]string, error) {
	return getRecCmd("mac", *opts)
}
