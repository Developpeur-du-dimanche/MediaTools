package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type HasSubtitlesFilter struct{}

func (f HasSubtitlesFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any subtitle streams
	hasSubtitles := len(data.Subtitles) > 0
	return compareBool(hasSubtitles, operator, value)
}

func (f HasSubtitlesFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "HAS_SUBTITLES",
		DisplayName:      "Has Subtitles",
		Type:             FieldTypeBoolean,
		PredefinedValues: []string{"true", "false"},
	}
}
