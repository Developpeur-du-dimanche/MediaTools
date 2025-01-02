package filter

import (
	"gopkg.in/vansante/go-ffprobe.v2"
)

type GlobalFilterContract interface {
	ConditionContract
	Check(data *ffprobe.ProbeData) bool
}

type GlobalFilter struct {
	Filter
}
