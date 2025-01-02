package filter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitleTitleFilter struct {
	GlobalFilter
}

func NewSubtitleTitleFilter() *SubtitleTitleFilter {
	return &SubtitleTitleFilter{}
}

func (c *SubtitleTitleFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is subtitle stream and title matches
		if s.CodecType == "subtitle" && c.CheckString(s.Tags.Title) {
			return true
		}
	}
	return false
}

func (c *SubtitleTitleFilter) CheckStream(data *ffprobe.Stream) bool {
	return data.CodecType == "subtitle" && c.CheckString(data.Tags.Title)
}

func (c *SubtitleTitleFilter) Name() string {
	return "Subtitle title"
}

func (c *SubtitleTitleFilter) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *SubtitleTitleFilter) SetCondition(condition string) {
	c.Condition = FromString(condition)
}

func (c *SubtitleTitleFilter) SetValue(value string) {
	c.Value = value
}

func (c *SubtitleTitleFilter) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.Value = s
	}
	return entry
}

func (c *SubtitleTitleFilter) New() ConditionContract {
	return &SubtitleLanguageFilter{
		GlobalFilter{
			Filter: Filter{
				Value: c.Value,
			},
		},
	}
}
