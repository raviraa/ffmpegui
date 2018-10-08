package ffprobe

import (
	"reflect"
	"testing"
)

func TestCommonConf(t *testing.T) {
	loadCommonConfig()
	SetOptions(Options{})
	tests := []struct {
		name string
		want []string
	}{
		{"audio", []string{"-thread_queue_size", "512", "-i", "none:0", "-f", "avfoundation"}},
		{"video", []string{"-thread_queue_size", "512", "-framerate", "25", "-video_size", "1920x1080", "-i", "0:none", "-f", "avfoundation"}},
	}
	for _, tt := range tests {
		if got := inputString(config.Inputs[tt.name]); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("conf input got %#v, want %#v", got, tt.want)
		}
	}
}

func TestContainer(t *testing.T) {
	loadCommonConfig()
	SetOptions(Options{})
	tests := []struct {
		name string
		want []string
	}{
		{"webm - vp9 default with no audio", []string{"-map", "0:v", "-c:v", "vp9", "0.webm"}},
		{"webm - vp9 default with opus default", []string{"-map", "0:v", "-c:v", "vp9", "0.webm", "-map", "1:a", "-c:a", "libopus", "1.opus"}},
	}
	for _, tt := range tests {
		if got := containerCmd(tt.name); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("container cmd %s, got %#v, want %#v", tt.name, got, tt.want)
		}
	}
}
