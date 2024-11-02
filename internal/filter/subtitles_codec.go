package filter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitleCodecFilter struct {
	Filter
}

func NewSubtitleCodecFilter() *SubtitleCodecFilter {
	return &SubtitleCodecFilter{}
}

func (c *SubtitleCodecFilter) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is subtitle stream and codec matches
		if s.CodecType == "subtitle" && c.CheckString(s.CodecName) {
			return true
		}
	}
	return false
}

func (c *SubtitleCodecFilter) Name() string {
	return "Subtitle codec"
}

func (c *SubtitleCodecFilter) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *SubtitleCodecFilter) New() ConditionContract {
	return &SubtitleCodecFilter{
		Filter{
			Value: c.Value,
		},
	}
}

func (c *SubtitleCodecFilter) SetCondition(condition string) {
	c.Condition = FromString(condition)
}

func (sc *SubtitleCodecFilter) GetEntry() fyne.Widget {
	selectWidget := widget.NewSelect([]string{"subrip", "ass"}, func(s string) {
		sc.Value = s
	})
	return selectWidget
}

func (c *SubtitleCodecFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	return true
}

func (c *SubtitleCodecFilter) CheckStream(data *ffprobe.Stream) bool {
	return false
}
