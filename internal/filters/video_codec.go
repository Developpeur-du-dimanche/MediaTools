package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type VideoCodecFilter struct{}

func (f VideoCodecFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any video streams
	if len(data.Videos) == 0 {
		return false
	}

	// Check the codec of the first video stream
	return compareString(data.Videos[0].CodecName, operator, value)
}

func (f VideoCodecFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "VIDEO_CODEC",
		DisplayName:      "Video Codec",
		Type:             FieldTypeString,
		PredefinedValues: []string{"h264", "h265", "hevc", "vp9", "av1", "mpeg4", "mpeg2video", "xvid"},
	}
}
