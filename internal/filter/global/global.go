package globalfilter

import (
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type GlobalFilterContract interface {
	filter.ConditionContract
	Check(data *ffprobe.ProbeData) bool
}

type GlobalFilter struct {
	filter.Filter
}
