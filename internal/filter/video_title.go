package filter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type VideoTitleFilter struct {
	Filter
}

func NewVideoTitleFilter() *VideoTitleFilter {
	return &VideoTitleFilter{}
}

func (c *VideoTitleFilter) Check(data *ffprobe.ProbeData) bool {
	for _, s := range data.Streams {
		// if is video stream and title matches
		if s.CodecType == "video" && c.CheckString(s.Tags.Title) {
			return true
		}
	}
	return false
}

func (c *VideoTitleFilter) Name() string {
	return "Video title"
}

func (c *VideoTitleFilter) GetPossibleConditions() []string {
	return []string{"equals", "contains", "not equals"}
}

func (c *VideoTitleFilter) New() ConditionContract {
	return &VideoTitleFilter{
		Filter{
			Value: c.Value,
		},
	}
}

func (c *VideoTitleFilter) SetCondition(condition string) {
	c.Condition = FromString(condition)
}

func (c *VideoTitleFilter) SetValue(value string) {
	c.Value = value
}

func (c *VideoTitleFilter) GetEntry() fyne.Widget {
	entry := widget.NewEntry()
	entry.OnChanged = func(s string) {
		c.Value = s
	}
	return entry
}

func (c *VideoTitleFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	return true
}

func (c *VideoTitleFilter) CheckStream(data *ffprobe.Stream) bool {
	return false
}
