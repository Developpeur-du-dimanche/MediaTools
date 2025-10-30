package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type FramerateFilter struct{}

func (f FramerateFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Framerate is not available in the FfprobeResult struct
	// The Video struct does not contain framerate information
	return false
}

func (f FramerateFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "FRAMERATE",
		DisplayName: "Framerate (fps)",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 30",
	}
}
