package filters

import "github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"

type FilterFieldValues []string

// FilterFieldType represents the data type of a filter field
type FilterFieldType string

const (
	FieldTypeNumeric FilterFieldType = "numeric"
	FieldTypeString  FilterFieldType = "string"
	FieldTypeBoolean FilterFieldType = "boolean"
)

// FilterFieldConfig defines all properties of a filter field
type FilterFieldConfig struct {
	Key              string          // Internal field key (e.g., "VIDEO_CODEC")
	DisplayName      string          // Human-readable name (e.g., "Video Codec")
	Type             FilterFieldType // Field data type
	PredefinedValues []string        // Optional list of predefined values for dropdown
	Placeholder      string          // Placeholder text for manual entry
}

// Operator definitions per field type
var OperatorsByType = map[FilterFieldType]FilterFieldValues{
	FieldTypeNumeric: {">", ">=", "<", "<=", "IS", "IS_NOT"},
	FieldTypeString:  {"IS", "IS_NOT", "CONTAINS", "NOT_CONTAINS"},
	FieldTypeBoolean: {"IS", "IS_NOT"},
}

type Filter interface {
	Apply(data *medias.FfprobeResult, operator string, value string) bool
	GetFieldConfig() FilterFieldConfig
}

// GetAllFilters returns all registered filters in the system
func GetAllFilters() []Filter {
	return []Filter{
		TotalBitrateFilter{},
		VideoBitrateFilter{},
		AudioBitrateFilter{},
		VideoCodecFilter{},
		AudioCodecFilter{},
		AudioLanguageFilter{},
		SubtitleLanguageFilter{},
		WidthFilter{},
		HeightFilter{},
		DurationFilter{},
		FramerateFilter{},
		AudioChannelsFilter{},
		HasVideoFilter{},
		HasAudioFilter{},
		HasSubtitlesFilter{},
	}
}
