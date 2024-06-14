package globalfilter

import (
	"fyne.io/fyne/v2"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components/widgets"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type BitrateFilter struct {
	GlobalFilter
}

func NewBitrateFilter() *BitrateFilter {
	return &BitrateFilter{}
}

func (c *BitrateFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	return c.Filter.CheckStringToInt(data.Format.BitRate)
}

func (c *BitrateFilter) CheckStream(data *ffprobe.Stream) bool {
	return c.Filter.CheckStringToInt(data.BitRate)
}

func (c *BitrateFilter) Name() string {
	return "Bitrate"
}

func (c *BitrateFilter) GetPossibleConditions() []string {
	return []string{"equals", "greater than", "less than", "greater or equals", "less or equals"}
}

func (c *BitrateFilter) New() filter.ConditionContract {
	return &BitrateFilter{
		GlobalFilter{
			Filter: filter.Filter{
				Value: c.Value,
			},
		},
	}
}

func (c *BitrateFilter) SetCondition(condition string) {
	c.Condition = filter.FromString(condition)
}

func (c *BitrateFilter) GetEntry() fyne.Widget {
	entry := widgets.NewNumericalEntry()
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.OnChanged = func(s string) {
		c.Value = s
	}
	return entry
}
