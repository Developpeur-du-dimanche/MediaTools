package filter

import (
	"fyne.io/fyne/v2"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type BitrateFilter struct {
	Filter
}

func NewBitrateFilter() *BitrateFilter {
	return &BitrateFilter{}
}

func (c *BitrateFilter) Check(data *ffprobe.ProbeData) bool {
	return c.CheckStringToInt(data.Format.BitRate)
}

func (c *BitrateFilter) Name() string {
	return "Bitrate"
}

func (c *BitrateFilter) GetPossibleConditions() []string {
	return []string{"equals", "greater than", "less than", "greater or equals", "less or equals"}
}

func (c *BitrateFilter) New() ConditionContract {
	return &BitrateFilter{
		Filter{
			Value: c.Value,
		},
	}
}

func (c *BitrateFilter) SetCondition(condition string) {
	c.Condition = FromString(condition)
}

func (c *BitrateFilter) GetEntry() fyne.Widget {
	entry := NewNumericalEntry()
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.OnChanged = func(s string) {
		c.Value = s
	}
	return entry
}

func (c *BitrateFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	return true
}

func (c *BitrateFilter) CheckStream(data *ffprobe.Stream) bool {
	return false
}
