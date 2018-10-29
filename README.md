`ffmpegui` is a GUI frontend for ffmpeg which can record screen and webcam with audio. Encoding and source options can be configured from the UI.


* `ffmpeg` needs to installed for separately. Use `homebrew` for Mac and package manager for your linux distribution to install ffmpeg.
* Recordings are stored in Desktop and configuration files are stored in [user level configuration](https://github.com/shibukawa/configdir) directory.
* Recording is made using intermediate fast capture-{a,v} profile and seperate file for each pause/resume. And converted to seleced profile and joined together.

## TODO

* linux support
* better input and output resolution, framerate configuration `ffprobe -f avfoundation -i 0`
* capture selected application window(s)
* support for input videos to transcode
* tray icon with status and stop, pause support
* capture to webp/gif animations