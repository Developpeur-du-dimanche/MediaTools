package info

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const (
	audio = "a"
)

type AudioInfo struct {
	id string
}

func NewAudioInfo(id string) *AudioInfo {
	return &AudioInfo{
		id: id,
	}
}

func (a *AudioInfo) GetNodes() []widget.TreeNodeID {
	i := strings.Split(a.id, " ")[1]
	return []widget.TreeNodeID{
		audio + i + "title",
		audio + i + "language",
		audio + i + "channels",
		audio + i + "codec_name",
		audio + i + "bit_rate",
	}
}

func (a *AudioInfo) From(label *widget.Label, stream *ffprobe.Stream, id widget.TreeNodeID) {
	switch true {
	case strings.HasSuffix(id, "codec_name"):
		label.SetText(lang.L("codec_name") + ": " + stream.CodecName)
	case strings.HasSuffix(id, "bit_rate"):
		label.SetText(lang.L("bitrate") + ": " + fmt.Sprint(stream.BitRate))
	case strings.HasSuffix(id, "channels"):
		label.SetText("Channels: " + fmt.Sprint(stream.Channels))
	case strings.HasSuffix(id, "language"):
		if stream.Tags.Language != "" {
			label.SetText(lang.L("language") + ": " + fmt.Sprint(stream.Tags.Language))
		} else {
			label.SetText(lang.L("language") + ": " + lang.L("unknown"))
		}
	case strings.HasSuffix(id, "title"):
		label.SetText(lang.L("title") + ": " + stream.Tags.Title)
	}
}
