package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type AudioBitrateFilter struct{}

func (f AudioBitrateFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	if len(data.Audios) == 0 {
		return false
	}

	for _, audio := range data.Audios {
		actualBitrate := parseBitrateValue(audio.Bitrate)

		targetBitrate := parseBitrateValue(value)
		if targetBitrate == 0 {
			continue
		}

		if compareNumeric(actualBitrate, operator, targetBitrate) {
			return true
		}
	}
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
