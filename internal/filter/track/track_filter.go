package trackfilter

import (
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type TrackFilterContract interface {
	Name() string
	Check(data *ffprobe.Stream) bool
	GetPossibleConditions() []string
	New() TrackFilterContract
	SetCondition(condition string)
	GetEntry() interface{}
}

type TrackFilter struct {
	filter.Filter
}
