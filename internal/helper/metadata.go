package helper

import (
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

type Metadata interface {
	Equals(other Metadata) bool
}

type FileMetadata struct {
	FileName  string
	Directory string
	Path      string
	Format    string
	Codec     string
	Bitrate   string
	Duration  time.Duration
	Size      string
	Video     []*ffprobe.Stream
	Audio     []*ffprobe.Stream
	Subtitle  []*ffprobe.Stream
	Container string
	Extension string
}

func NewFileMetadata(data *ffprobe.ProbeData) (*FileMetadata, error) {

	// get extension
	ext := filepath.Ext(data.Format.Filename)

	filename := filepath.Base(data.Format.Filename)

	folder := filepath.Dir(data.Format.Filename)

	audioStreams := make([]*ffprobe.Stream, 0)
	videoStreams := make([]*ffprobe.Stream, 0)
	subtitleStreams := make([]*ffprobe.Stream, 0)

	for _, stream := range data.Streams {
		switch stream.CodecType {
		case "audio":
			audioStreams = append(audioStreams, stream)
		case "video":
			videoStreams = append(videoStreams, stream)
		case "subtitle":
			subtitleStreams = append(subtitleStreams, stream)
		}
	}

	if ext != "" {
		ext = strings.TrimPrefix(ext, ".")
	} else {
		ext = "unknown"
	}

	return &FileMetadata{
		FileName:  filename,
		Directory: folder,
		Path:      data.Format.Filename,
		Format:    data.Format.FormatName,
		Codec:     data.Format.FormatLongName,
		Bitrate:   data.Format.BitRate,
		Duration:  data.Format.Duration(),
		Size:      data.Format.Size,
		Video:     videoStreams,
		Audio:     audioStreams,
		Subtitle:  subtitleStreams,
		Container: data.Format.FormatName,
		Extension: ext,
	}, nil
}

func (f *FileMetadata) Equals(other *FileMetadata) bool {
	return f.FileName == other.FileName &&
		f.Format == other.Format &&
		f.Codec == other.Codec &&
		f.Bitrate == other.Bitrate &&
		f.Duration == other.Duration &&
		f.Size == other.Size &&
		f.Container == other.Container &&
		f.Extension == other.Extension
}

func (f *FileMetadata) GetVideoStreams() []*ffprobe.Stream {
	return f.Video
}

func (f *FileMetadata) GetAudioStreams() []*ffprobe.Stream {
	return f.Audio
}

func (f *FileMetadata) GetSubtitleStreams() []*ffprobe.Stream {
	return f.Subtitle
}
