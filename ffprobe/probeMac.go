package ffprobe

type proberMac struct {
	devicesCmd string
}

func newProberMac() Prober {
	var macProber = proberMac{
		//TODO move to conf file
		devicesCmd: "ffmpeg -f avfoundation -list_devices true -i ''",
	}
	return macProber
}

func (pm proberMac) getDevicesCmd() string {
	return pm.devicesCmd
}

func (pm proberMac) getFfmpegCmd(pc ProberCommon) ([]string, error) {
	return pc.getRecCmd("mac", *opts)
}
