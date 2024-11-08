package globalfilter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type SubtitleLanguageFilter struct {
	GlobalFilter
}

func NewSubtitleLanguageFilter() *SubtitleLanguageFilter {
	return &SubtitleLanguageFilter{}
}

func (c *SubtitleLanguageFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is subtitle stream and language matches
		if s.CodecType == "subtitle" && c.CheckString(s.Tags.Language) {
			return true
		}
	}
	return false
}

func (c *SubtitleLanguageFilter) CheckStream(data *ffprobe.Stream) bool {
	return data.CodecType == "subtitle" && c.CheckString(data.Tags.Language)
}

func (c *SubtitleLanguageFilter) Name() string {
	return "Subtitle language"
}

func (c *SubtitleLanguageFilter) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *SubtitleLanguageFilter) New() filter.ConditionContract {
	return &SubtitleLanguageFilter{
		GlobalFilter{
			Filter: filter.Filter{
				Value: c.Value,
			},
		},
	}
}

func (c *SubtitleLanguageFilter) SetCondition(condition string) {
	c.Condition = filter.FromString(condition)
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
