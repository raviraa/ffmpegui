package ffprobe

import (
	"os/user"
	"path"
)

func (p proberCommon) probeDefaults() {
	u, _ := user.Current()
	p.opts.VidPath = path.Join(u.HomeDir, "Desktop")
	p.opts.Framerate = 24.0
}
