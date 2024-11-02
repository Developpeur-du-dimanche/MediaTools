package globalfilter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type AudioLanguageFilter struct {
	GlobalFilter
}

func NewAudioLanguageFilter() *AudioLanguageFilter {
	return &AudioLanguageFilter{}
}

func (c *AudioLanguageFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is audio stream and language matches
		if s.CodecType == "audio" && c.Filter.CheckString(s.Tags.Language) {
			return true
		}
	}
	return false
}

func (c *AudioLanguageFilter) CheckStream(data *ffprobe.Stream) bool {
	return data.CodecType == "audio" && data.Tags.Language != "" && c.CheckString(data.Tags.Language)
}

func (c *AudioLanguageFilter) Name() string {
	return "Audio Language"
}

func (c *AudioLanguageFilter) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *AudioLanguageFilter) New() filter.ConditionContract {
	return &AudioLanguageFilter{
		GlobalFilter{
			Filter: filter.Filter{
				Value: c.Value,
			},
		},
	}
}

func (c *AudioLanguageFilter) SetCondition(condition string) {
	c.Condition = filter.FromString(condition)
}

func (c *AudioLanguageFilter) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.Value = s
	}
	return entry
}
