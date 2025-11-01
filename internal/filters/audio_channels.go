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
	for _, audio := range data.Audios {
		actualChannels := int64(audio.Channels)
		if compareNumeric(actualChannels, operator, targetChannels) {
			return true
		}
	}
	return false
}

func (f AudioChannelsFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "AUDIO_CHANNELS",
		DisplayName: "Audio Channels",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 6",
	}
}
