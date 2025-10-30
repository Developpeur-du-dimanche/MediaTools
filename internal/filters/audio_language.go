package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type AudioLanguageFilter struct{}

func (f AudioLanguageFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any audio streams
	if len(data.Audios) == 0 {
		return false
	}

	// Loop through all audio streams and check if any matches the language
	for _, audio := range data.Audios {
		if compareString(audio.Language, operator, value) {
			return true
		}
	}

	return false
}

func (f AudioLanguageFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "AUDIO_LANGUAGE",
		DisplayName:      "Audio Language",
		Type:             FieldTypeString,
		PredefinedValues: []string{"fre", "eng", "spa", "deu", "ita", "jpn", "kor", "chi", "por", "rus", "ara", "hin"},
	}
}
