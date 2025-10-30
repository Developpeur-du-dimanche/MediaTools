package filters

import (
	"strconv"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type AudioChannelsFilter struct{}

func (f AudioChannelsFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any audio streams
	if len(data.Audios) == 0 {
		return false
	}

	// Parse the target channels value
	targetChannels, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}

	// Compare the channels of the first audio stream
	actualChannels := int64(data.Audios[0].Channels)
	return compareNumeric(actualChannels, operator, targetChannels)
}

func (f AudioChannelsFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "AUDIO_CHANNELS",
		DisplayName: "Audio Channels",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 6",
	}
}
