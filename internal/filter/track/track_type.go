package trackfilter

import "gopkg.in/vansante/go-ffprobe.v2"

type TrackTypeFilter struct {
	TrackFilter
}

func NewTrackTypeFilter() *TrackTypeFilter {
	return &TrackTypeFilter{}
}

func (c *TrackTypeFilter) Check(data *ffprobe.Stream) bool {
	return c.Filter.CheckString(data.CodecType)
}
