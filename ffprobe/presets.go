package ffprobe

var defaultPresetStr = `#default presets
Ffcmdprefix = "ffmpeg -benchmark -y -loglevel verbose -hide_banner"

[inputs_defaults]
thread_queue_size = 512

[inputs]
# -video_size 1920x1080 -framerate 25 -thread_queue_size 512  -f avfoundation -i 3:none \
[inputs.mac.v]
f = 'avfoundation'
i = '%d:none'
[inputs.mac.a]
f = 'avfoundation'
i = 'none:%d'


[presets]

[presets.opus-default]
fileext = 'opus'
avtype = 'a'
codec = 'libopus'

[presets.vp9-default]
fileext = "webm"
avtype = "v"
codec = "vp9"
[presets.vp9-default.options]
threads = 8
[presets.h264-default]
fileext = "mkv"
avtype = "v"
codec = "libx264"

[presets.opus-voice]
fileext = 'opus'
avtype = 'a'
codec = 'libopus'
[presets.opus-voice.options]
"b:a" = '32k'
application = 'voip'
af = "highpass=f=200, lowpass=f=1000"
# TODO: -vf/-af/-filter and -filter_complex cannot be used together for the same stream.



# [presets.vp9-realtime]
# pix_fmt='yuv420p' # TODO output defaults
# fileext = "webm"
# avtype = "v"
# codec = "vp9"
# [presets.vp9-realtime.options]
# speed=5
# tile-columns=6
# frame-parallel=1
# threads=8
# static-thresh=0
# max-intra-rate=300
# deadline="realtime"
# crf=9
# lag-in-frames=9
# # TODO quality min max
# row-mt=1
# error-resilient=1

[presets.capture-v]
fileext = "mkv"
avtype = "v"
codec = "libx264"
[presets.capture-v.options]
preset = "ultrafast"
crf = 0
[presets.capture-a]
fileext = 'aac'
avtype = 'a'
codec = 'aac'
`
