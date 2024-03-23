package filter

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type ContainerFilter struct {
	Filter
}

func NewContainerFilter() ConditionContract {
	return &ContainerFilter{}
}

func (c *ContainerFilter) Check(data *ffprobe.ProbeData) bool {
	if data == nil || data.Format == nil || data.Format.Filename == "" {
		return false
	}

	// get extension
	extension := data.Format.Filename[strings.LastIndex(data.Format.Filename, ".")+1:]

	return c.checkString(extension)
}

func (c *ContainerFilter) GetPossibleConditions() []string {
	return []string{"equals", "not equals"}
}

func (c *ContainerFilter) Name() string {
	return "Container"
}

func (c *ContainerFilter) New() ConditionContract {
	return &ContainerFilter{
		Filter{
			value: c.value,
		},
	}
}

func (c *ContainerFilter) SetCondition(condition string) {
	c.condition = FromString(condition)
}

func (c *ContainerFilter) GetEntry() fyne.Widget {
	entry := widget.NewSelect([]string{"mkv", "mp4"}, func(s string) {
		c.value = s

	})
	return entry
}
