package filter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type AudioLanguageFilter struct {
	Filter
}

func NewAudioLanguageFilter() *AudioLanguageFilter {
	return &AudioLanguageFilter{}
}

func (c *AudioLanguageFilter) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is audio stream and language matches
		if s.CodecType == "audio" && c.checkString(s.Tags.Language) {
			return true
		}
	}
	return false
}

func (c *AudioLanguageFilter) Name() string {
	return "Audio Language"
}

func (c *AudioLanguageFilter) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *AudioLanguageFilter) New() ConditionContract {
	return &AudioLanguageFilter{
		Filter{
			value: c.value,
		},
	}

}

func (c *AudioLanguageFilter) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *AudioLanguageFilter) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.value = s
	}
	return entry
}
