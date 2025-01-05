package info

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const (
	video = "v"
)

type VideoInfo struct {
	id string
}

func NewVideoInfo(id string) *VideoInfo {
	return &VideoInfo{
		id: id,
	}
}

func (v *VideoInfo) GetNodes() []widget.TreeNodeID {
	i := strings.Split(v.id, " ")[1]
	return []widget.TreeNodeID{
		video + i + "_codec_name",
		video + i + "_width",
		video + i + "_height",
		video + i + "_bit_rate",
	}
}

func (v *VideoInfo) From(label *widget.Label, stream *ffprobe.Stream, id widget.TreeNodeID) {
	switch true {
	case strings.HasSuffix(id, "codec_name"):
		label.SetText(lang.L("codec") + ": " + stream.CodecName)
	case strings.HasSuffix(id, "width"):
		label.SetText(lang.L("width") + ": " + fmt.Sprint(stream.Width))
	case strings.HasSuffix(id, "height"):
		label.SetText(lang.L("height") + ": " + fmt.Sprint(stream.Height))
	case strings.HasSuffix(id, "bit_rate"):
		label.SetText(lang.L("bitrate") + ": " + fmt.Sprint(stream.BitRate))
	case strings.HasSuffix(id, "title"):
		label.SetText(lang.L("title") + ": " + stream.Tags.Title)
	}
}
