package globalfilter

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type ContainerFilter struct {
	GlobalFilter
}

func NewContainerFilter() filter.ConditionContract {
	return &ContainerFilter{}
}

func (c *ContainerFilter) CheckGlobal(data *ffprobe.ProbeData) bool {
	if data == nil || data.Format == nil || data.Format.Filename == "" {
		return false
	}

	// get extension
	extension := data.Format.Filename[strings.LastIndex(data.Format.Filename, ".")+1:]

	return c.Filter.CheckString(extension)
}

func (c *ContainerFilter) CheckStream(data *ffprobe.Stream) bool {
	return false
}

func (c *ContainerFilter) GetPossibleConditions() []string {
	return []string{"equals", "not equals"}
}

func (c *ContainerFilter) Name() string {
	return "Container"
}

func (c *ContainerFilter) New() filter.ConditionContract {
	return &ContainerFilter{
		GlobalFilter{
			Filter: filter.Filter{
				Value: c.Value,
			},
		},
	}
}

func (c *ContainerFilter) SetCondition(condition string) {
	c.Condition = filter.FromString(condition)
}

func (c *ContainerFilter) GetEntry() fyne.Widget {
	entry := widget.NewSelect([]string{"mkv", "mp4"}, func(s string) {
		c.Value = s

	})
	return entry
}
