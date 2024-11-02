package info

import (
	"fmt"
	"strings"

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
		label.SetText("Codec name: " + stream.CodecName)
	case strings.HasSuffix(id, "width"):
		label.SetText("Width: " + fmt.Sprint(stream.Width))
	case strings.HasSuffix(id, "height"):
		label.SetText("Height: " + fmt.Sprint(stream.Height))
	case strings.HasSuffix(id, "bit_rate"):
		label.SetText("Bit rate: " + fmt.Sprint(stream.BitRate))
	case strings.HasSuffix(id, "title"):
		label.SetText("Title: " + stream.Tags.Title)
	}
}
