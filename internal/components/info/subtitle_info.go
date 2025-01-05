package info

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

const (
	subtitle = "s"
)

type SubtitleInfo struct {
	id string
}

func NewSubtitleInfo(id string) *SubtitleInfo {
	return &SubtitleInfo{
		id: id,
	}
}

func (v *SubtitleInfo) GetNodes() []widget.TreeNodeID {
	i := strings.Split(v.id, " ")[1]
	return []widget.TreeNodeID{
		subtitle + i + "_title",
		subtitle + i + "_codec_name",
		subtitle + i + "_bit_rate",
		subtitle + i + "_language",
		subtitle + i + "_forced",
		subtitle + i + "_default",
		subtitle + i + "hearing_impaired",
	}
}

func (v *SubtitleInfo) From(label *widget.Label, stream *ffprobe.Stream, id widget.TreeNodeID) {
	switch true {
	case strings.HasSuffix(id, "codec_name"):
		label.SetText("Codec name: " + stream.CodecName)
	case strings.HasSuffix(id, "bit_rate"):
		label.SetText("Bit rate: " + fmt.Sprint(stream.BitRate))
	case strings.HasSuffix(id, "title"):
		label.SetText("Title: " + stream.Tags.Title)
	case strings.HasSuffix(id, "language"):
		if stream.Tags.Language != "" {
			label.SetText(lang.L("language") + ": " + fmt.Sprint(stream.Tags.Language))
		} else {
			label.SetText(lang.L("language") + ": " + lang.L("unknown"))
		}
	case strings.HasSuffix(id, "forced"):
		if stream.Disposition.Forced == 1 {
			label.SetText(lang.L("forced") + ": " + lang.L("yes"))
		} else {
			label.SetText(lang.L("forced") + ": " + lang.L("no"))
		}
	case strings.HasSuffix(id, "default"):
		if stream.Disposition.Default == 1 {
			label.SetText("Default: Yes")
		} else {
			label.SetText("Default: No")
		}
	case strings.HasSuffix(id, "hearing_impaired"):
		if stream.Disposition.HearingImpaired == 1 {
			label.SetText(lang.L("hearing_impaired") + ": " + lang.L("yes"))
		} else {
			label.SetText(lang.L("hearing_impaired") + ": " + lang.L("no"))
		}
	}
}
