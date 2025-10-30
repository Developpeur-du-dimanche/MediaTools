package filters

import (
	"strconv"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type WidthFilter struct{}

func (f WidthFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any video streams
	if len(data.Videos) == 0 {
		return false
	}

	// Parse the target width value
	targetWidth, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}

	// Compare the width of the first video stream
	actualWidth := int64(data.Videos[0].Width)
	return compareNumeric(actualWidth, operator, targetWidth)
}

func (f WidthFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "WIDTH",
		DisplayName: "Width (px)",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 1920",
	}
}
