#!/bin/bash

set -x

# realtime vp9
#ffmpeg -y -loglevel verbose -thread_queue_size 512  -f avfoundation -i none:0 -map 0:a  -c:a libopus TODO.opus -framerate 25 -thread_queue_size 512  -f avfoundation -i 3:none  -map 1:v -c:v vp9 -s 1920x1080  -benchmark -speed 6 -tile-columns 4 -frame-parallel 1 -threads 8 -static-thresh 0 -max-intra-rate 300 -deadline realtime -lag-in-frames 0 -error-resilient 1 TODO.webm
ffmpeg -y -loglevel verbose \
 -thread_queue_size 512  -f avfoundation -i none:0 \
 -map 0:a  -c:a libopus TODO.opus  \
 -video_size 1920x1080 -framerate 25 -thread_queue_size 512  -f avfoundation -i 3:none \
 -map 1:v -c:v vp9 -s 1920x1080 -framerate 25  -benchmark -speed 5 -tile-columns 6 -frame-parallel 1 -threads 8 -static-thresh 0 -max-intra-rate 300 -deadline realtime  -crf 9 -lag-in-frames 9 -row-mt 1 -error-resilient 1 TODO.webm

#ffmpeg -video_size 1024x768 -framerate 25 -f x11grab -i :0.0+100,200
#-f pulse -ac 2 -i default output.mkv



mpv TODO.webm