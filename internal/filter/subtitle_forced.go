package filter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitleForcedFilter struct {
	Filter
}

func NewSubtitleForcedFilter() *SubtitleForcedFilter {
	return &SubtitleForcedFilter{}
}

func (c *SubtitleForcedFilter) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is subtitle stream and forced matches
		if s.CodecType == "subtitle" && c.checkInt(s.Disposition.Forced) {
			return true
		}
	}
	return false
}

func (c *SubtitleForcedFilter) Name() string {
	return "Subtitle forced"
}

func (c *SubtitleForcedFilter) GetPossibleConditions() []string {
	return []string{"equals", "not equals"}
}

func (c *SubtitleForcedFilter) New() ConditionContract {
	return &SubtitleForcedFilter{
		Filter{
			value: c.value,
		},
	}
}

func (c *SubtitleForcedFilter) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *SubtitleForcedFilter) GetEntry() fyne.Widget {
	entry := widget.NewSelect([]string{"true", "false"}, func(s string) {
		switch s {
		case "true":
			c.value = "1"
		case "false":
			c.value = "0"
		}
	})
	return entry
}

func (c *SubtitleForcedFilter) SetValue(value string) {
	c.value = value
}
