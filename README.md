`ffmpegui` is a GUI frontend for ffmpeg which can record screen and webcam with audio. Encoding and source options can be configured from the UI.


* `ffmpeg` needs to installed for separately. Use `homebrew` for Mac and package manager for your linux distribution to install ffmpeg.
* Make sure to use appropriate presets like *vp9-realtime* to capture screen. ffmpeg record speed should be approximately 1x.
* Container format(webm, mp4) is determined based on the first preset video codec selected.

## TODO

* linux support
* better input and output resolution, framerate configuration `ffprobe -f avfoundation -i 0`
* capture selected application window(s), robotgo?
* support for input videos to transcode
* tray icon with status and stop, pause support
* capture to webp/gif animations