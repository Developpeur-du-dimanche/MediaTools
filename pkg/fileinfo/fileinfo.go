package fileinfo

import (
	"context"
	"path/filepath"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

type FileInfo interface {
	GetFilename() string
	GetFolder() string
	GetPath() string
	GetInfo() *ffprobe.ProbeData
	GetVideoStreams() *[]ffprobe.Stream
	GetAudioStreams() *[]ffprobe.Stream
	GetSubtitleStreams() *[]ffprobe.Stream
	Equals(other FileInfo) bool
}

type fileInfo struct {
	Path            string
	Folder          string
	Filename        string
	VideoStreams    *[]ffprobe.Stream
	AudioStreams    *[]ffprobe.Stream
	SubtitleStreams *[]ffprobe.Stream
	Info            *ffprobe.ProbeData
}

func NewFileInfo(path string) (FileInfo, error) {
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

	return &fileInfo{
		Path:            path,
		Filename:        getFilename(path),
		Folder:          getFolder(path),
		Info:            data,
		VideoStreams:    &videoStreams,
		AudioStreams:    &audioStreams,
		SubtitleStreams: &subtitleStreams,
	}, nil
}

func (f *fileInfo) GetFilename() string {
	return f.Filename
}

func getFilename(path string) string {
	// extract filename from path
	return filepath.Base(path)
}

func getFolder(path string) string {
	// extract folder from path
	return filepath.Dir(path)
}

func (f *fileInfo) GetFolder() string {
	return f.Folder
}

func (f *fileInfo) GetPath() string {
	return f.Path
}

func (f *fileInfo) GetInfo() *ffprobe.ProbeData {
	return f.Info
}

func (f *fileInfo) GetVideoStreams() *[]ffprobe.Stream {
	return f.VideoStreams
}

func (f *fileInfo) GetAudioStreams() *[]ffprobe.Stream {
	return f.AudioStreams
}

func (f *fileInfo) GetSubtitleStreams() *[]ffprobe.Stream {
	return f.SubtitleStreams
}

func (f *fileInfo) Equals(other FileInfo) bool {
	otherFileInfo, ok := other.(*fileInfo)
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
