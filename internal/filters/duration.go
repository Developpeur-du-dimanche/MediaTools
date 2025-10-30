package filters

import (
	"strconv"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type DurationFilter struct{}

func (f DurationFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Parse the target duration value (in seconds)
	targetDuration, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}

	// Get the actual duration in seconds from Format.DurationSeconds
	actualDuration := int64(data.Format.DurationSeconds)
	return compareNumeric(actualDuration, operator, targetDuration)
}

func (f DurationFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "DURATION",
		DisplayName: "Duration (seconds)",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 3600 (seconds)",
	}
}
