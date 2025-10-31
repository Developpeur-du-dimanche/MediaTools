package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type AudioCodecFilter struct{}

func (f AudioCodecFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any audio streams
	if len(data.Audios) == 0 {
		return false
	}

	// Loop through all audio streams and check if any matches the codec
	for _, audio := range data.Audios {
		if compareString(audio.CodecName, operator, value) {
			return true
		}
	}

	return false
}

func (f AudioCodecFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "AUDIO_CODEC",
		DisplayName:      "Audio Codec",
		Type:             FieldTypeString,
		PredefinedValues: []string{"aac", "mp3", "ac3", "eac3", "dts", "flac", "opus", "vorbis", "pcm"},
	}
}
