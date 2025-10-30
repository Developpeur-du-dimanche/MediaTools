package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type HasAudioFilter struct{}

func (f HasAudioFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any audio streams
	hasAudio := len(data.Audios) > 0
	return compareBool(hasAudio, operator, value)
}

func (f HasAudioFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "HAS_AUDIO",
		DisplayName:      "Has Audio Stream",
		Type:             FieldTypeBoolean,
		PredefinedValues: []string{"true", "false"},
	}
}
