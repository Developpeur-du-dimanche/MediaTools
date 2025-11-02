package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"


type SubtitleCodecFilter struct{}

func (f SubtitleCodecFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	if len(data.Subtitles) == 0 {
		return false
	}

	for _, subtitle := range data.Subtitles {
		if compareString(subtitle.CodecName, operator, value) {
			return true
		}
	}

	return false
}

func (f SubtitleCodecFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "SUBTITLE_CODEC",
		DisplayName:      "Subtitle Codec",
		Type:             FieldTypeString,
		PredefinedValues: []string{"srt", "ass", "subrip", "pgs", "dvd_subtitle", "webvtt"},
		TargetType:      TargetSubtitle,
	}
}