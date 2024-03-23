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
		if s.CodecType == "subtitle" && c.checkString(s.Tags.Language) {
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
			value: c.value,
		},
	}
}

func (c *SubtitleLanguageFilter) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *SubtitleLanguageFilter) SetValue(value string) {
	c.value = value
}

func (c *SubtitleLanguageFilter) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.value = s
	}
	return entry
}
