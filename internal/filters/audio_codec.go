package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type AudioCodecFilter struct{}

func (f AudioCodecFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any audio streams
	if len(data.Audios) == 0 {
		return false
	}

	// Check the codec of the first audio stream
	return compareString(data.Audios[0].CodecName, operator, value)
}

func (f AudioCodecFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "AUDIO_CODEC",
		DisplayName:      "Audio Codec",
		Type:             FieldTypeString,
		PredefinedValues: []string{"aac", "mp3", "ac3", "eac3", "dts", "flac", "opus", "vorbis", "pcm"},
	}
}
