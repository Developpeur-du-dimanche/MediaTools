package filter

import (
	"strings"

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

func (c *ContainerFilter) SetCondition(condition ConditionString) {
	c.condition = condition
}
