package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type SubtitleLanguageFilter struct{}

func (f SubtitleLanguageFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	// Check if there are any subtitle streams
	if len(data.Subtitles) == 0 {
		return false
	}

	// Loop through all subtitle streams and check if any matches the language
	for _, subtitle := range data.Subtitles {
		if compareString(subtitle.Language, operator, value) {
			return true
		}
	}

	return false
}

func (f SubtitleLanguageFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:              "SUBTITLE_LANGUAGE",
		DisplayName:      "Subtitle Language",
		Type:             FieldTypeString,
		PredefinedValues: []string{"fre", "eng", "spa", "deu", "ita", "jpn", "kor", "chi", "por", "rus", "ara", "hin"},
	}
}
