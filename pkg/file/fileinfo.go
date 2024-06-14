package file

import (
	"context"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

type FileInfoContract interface {
	Equals(other FileInfoContract) bool
}

type FileInfo struct {
	Path            string
	VideoStreams    *[]ffprobe.Stream
	AudioStreams    *[]ffprobe.Stream
	SubtitleStreams *[]ffprobe.Stream
	Info            *ffprobe.ProbeData
}

func NewFileInfo(path string) (*FileInfo, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, path)

	if err != nil {
		return nil, err
	}

	var videoStreams []ffprobe.Stream
	var audioStreams []ffprobe.Stream
	var subtitleStreams []ffprobe.Stream

	for _, stream := range data.Streams {
		switch stream.CodecType {
		case "video":
			videoStreams = append(videoStreams, *stream)
		case "audio":
			audioStreams = append(audioStreams, *stream)
		case "subtitle":
			subtitleStreams = append(subtitleStreams, *stream)
		}
	}

	return &FileInfo{
		Path:            path,
		Info:            data,
		VideoStreams:    &videoStreams,
		AudioStreams:    &audioStreams,
		SubtitleStreams: &subtitleStreams,
	}, nil
}

func (f FileInfo) Equals(other FileInfoContract) bool {
	otherFileInfo, ok := other.(*FileInfo)
	if !ok {
		return false
	}

	if f.Path != otherFileInfo.Path {
		return false
	}

	if len(f.Info.Streams) != len(otherFileInfo.Info.Streams) {
		return false
	}

	return true
}
