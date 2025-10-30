package filters

import (
	"strconv"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type HeightFilter struct{}

func (f HeightFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any video streams
	if len(data.Videos) == 0 {
		return false
	}

	// Parse the target height value
	targetHeight, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}

	// Compare the height of the first video stream
	actualHeight := int64(data.Videos[0].Height)
	return compareNumeric(actualHeight, operator, targetHeight)
}

func (f HeightFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "HEIGHT",
		DisplayName: "Height (px)",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 1080",
	}
}
