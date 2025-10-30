package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type AudioBitrateFilter struct{}

func (f AudioBitrateFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Audio bitrate is not available in the FfprobeResult struct
	// The struct only contains codec information, not bitrate per stream
	return false
}

func (f AudioBitrateFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "AUDIO_BITRATE",
		DisplayName: "Bitrate (Audio)",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 320kbps",
	}
}
