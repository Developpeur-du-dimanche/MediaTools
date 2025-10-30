package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type VideoBitrateFilter struct{}

func (f VideoBitrateFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Video bitrate not available in simplified FfprobeResult structure
	// Would need to enhance FfprobeResult to include stream bitrates
	return false
}

func (f VideoBitrateFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "VIDEO_BITRATE",
		DisplayName: "Bitrate (Video)",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 1500kbps",
	}
}
