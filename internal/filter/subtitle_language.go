package filter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitleLanguageFilter struct {
	Filter
}

func NewSubtitleLanguageFilter() *SubtitleLanguageFilter {
	return &SubtitleLanguageFilter{}
}

func (c *SubtitleLanguageFilter) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is subtitle stream and language matches
		if s.CodecType == "subtitle" && c.CheckString(s.Tags.Language) {
			return true
		}
	}
	return false
}

func (c *SubtitleLanguageFilter) Name() string {
	return "Subtitle language"
}

func (c *SubtitleLanguageFilter) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *SubtitleLanguageFilter) New() ConditionContract {
	return &SubtitleLanguageFilter{
		Filter{
			Value: c.Value,
		},
	}
}

func (c *SubtitleLanguageFilter) SetCondition(condition string) {
	c.Condition = FromString(condition)
}

func (c *SubtitleLanguageFilter) SetValue(value string) {
	c.Value = value
}

func (c *SubtitleLanguageFilter) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.Value = s
	}
	return entry
}

func (c *SubtitleLanguageFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	return true
}

func (c *SubtitleLanguageFilter) CheckStream(data *ffprobe.Stream) bool {
	return false
}
