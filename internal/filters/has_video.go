package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type HasVideoFilter struct{}

func (f HasVideoFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any video streams
	hasVideo := len(data.Videos) > 0
	return compareBool(hasVideo, operator, value)
}

func (f HasVideoFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "HAS_VIDEO",
		DisplayName:      "Has Video Stream",
		Type:             FieldTypeBoolean,
		PredefinedValues: []string{"true", "false"},
	}
}
