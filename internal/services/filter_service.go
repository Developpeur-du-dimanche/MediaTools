package services

import (
	"fmt"
	"strings"

	"github.com/Developpeur-du-dimanche/MediaTools/internal/filters"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

// FilterOperator represents comparison operators
type FilterOperator string

const (
	OpEquals      FilterOperator = "IS"
	OpNotEquals   FilterOperator = "IS_NOT"
	OpGreater     FilterOperator = ">"
	OpGreaterEq   FilterOperator = ">="
	OpLess        FilterOperator = "<"
	OpLessEq      FilterOperator = "<="
	OpContains    FilterOperator = "CONTAINS"
	OpNotContains FilterOperator = "NOT_CONTAINS"
)

// LogicalOperator represents logical operators between conditions
type LogicalOperator string

const (
	LogicalAnd LogicalOperator = "AND"
	LogicalOr  LogicalOperator = "OR"
)

// FilterField represents available fields to filter on
type FilterField string

const (
	FieldBitrate        FilterField = "BITRATE"
	FieldVideoBitrate   FilterField = "VIDEO_BITRATE"
	FieldAudioBitrate   FilterField = "AUDIO_BITRATE"
	FieldVideoCodec     FilterField = "VIDEO_CODEC"
	FieldAudioCodec     FilterField = "AUDIO_CODEC"
	FieldAudioLanguage  FilterField = "AUDIO_LANGUAGE"
	FieldSubLanguage    FilterField = "SUBTITLE_LANGUAGE"
	FieldWidth          FilterField = "WIDTH"
	FieldHeight         FilterField = "HEIGHT"
	FieldDuration       FilterField = "DURATION"
	FieldFramerate      FilterField = "FRAMERATE"
	FieldAudioChannels  FilterField = "AUDIO_CHANNELS"
	FieldHasVideo       FilterField = "HAS_VIDEO"
	FieldHasAudio       FilterField = "HAS_AUDIO"
	FieldHasSubtitles   FilterField = "HAS_SUBTITLES"
)

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Field    FilterField
	Operator FilterOperator
	Value    string
}

// FilterExpression represents a complete filter with logical operators
type FilterExpression struct {
	Conditions []FilterCondition
	Operators  []LogicalOperator
}

// FilterService handles media filtering
type FilterService struct {
	filterRegistry map[FilterField]filters.Filter
}

// NewFilterService creates a new filter service
func NewFilterService() *FilterService {
	// Build registry from all registered filters
	registry := make(map[FilterField]filters.Filter)
	for _, filter := range filters.GetAllFilters() {
		config := filter.GetFieldConfig()
		registry[FilterField(config.Key)] = filter
	}

	return &FilterService{
		filterRegistry: registry,
	}
}

// ParseFilter parses a filter string like "BITRATE > 2000 AND AUDIO_LANGUAGE IS FRE"
func (fs *FilterService) ParseFilter(filterStr string) (*FilterExpression, error) {
	if strings.TrimSpace(filterStr) == "" {
		return &FilterExpression{}, nil
	}

	logger.Debugf("Parsing filter: %s", filterStr)

	expr := &FilterExpression{
		Conditions: []FilterCondition{},
		Operators:  []LogicalOperator{},
	}

	// Split by AND/OR while preserving the operators
	parts := []string{}
	current := ""

	tokens := strings.Fields(filterStr)
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if token == "AND" || token == "OR" {
			if current != "" {
				parts = append(parts, strings.TrimSpace(current))
				current = ""
			}
			expr.Operators = append(expr.Operators, LogicalOperator(token))
		} else {
			current += " " + token
		}
	}

	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}

	// Parse each condition
	for _, part := range parts {
		condition, err := fs.parseCondition(part)
		if err != nil {
			return nil, fmt.Errorf("failed to parse condition '%s': %w", part, err)
		}
		expr.Conditions = append(expr.Conditions, *condition)
	}

	if len(expr.Conditions) != len(expr.Operators)+1 {
		return nil, fmt.Errorf("invalid filter expression: mismatched conditions and operators")
	}

	logger.Infof("Parsed filter: %d conditions, %d operators", len(expr.Conditions), len(expr.Operators))
	return expr, nil
}

// parseCondition parses a single condition like "BITRATE > 2000"
func (fs *FilterService) parseCondition(condStr string) (*FilterCondition, error) {
	condStr = strings.TrimSpace(condStr)

	// Try to find operator
	operators := []FilterOperator{
		OpNotContains, OpContains, // Check compound operators first
		OpGreaterEq, OpLessEq, OpNotEquals, OpEquals,
		OpGreater, OpLess,
	}

	for _, op := range operators {
		opStr := string(op)
		if strings.Contains(condStr, " "+opStr+" ") {
			parts := strings.SplitN(condStr, " "+opStr+" ", 2)
			if len(parts) != 2 {
				continue
			}

			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			return &FilterCondition{
				Field:    FilterField(field),
				Operator: op,
				Value:    value,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid condition format: %s", condStr)
}

// ApplyFilter applies the filter expression to a media item
func (fs *FilterService) ApplyFilter(media *medias.FfprobeResult, expr *FilterExpression) bool {
	if len(expr.Conditions) == 0 {
		return true // No filter = all pass
	}

	// Evaluate first condition
	result := fs.evaluateCondition(media, expr.Conditions[0])

	// Apply logical operators
	for i, op := range expr.Operators {
		nextResult := fs.evaluateCondition(media, expr.Conditions[i+1])

		switch op {
		case LogicalAnd:
			result = result && nextResult
		case LogicalOr:
			result = result || nextResult
		}
	}

	return result
}

// evaluateCondition evaluates a single condition against a media item
func (fs *FilterService) evaluateCondition(media *medias.FfprobeResult, cond FilterCondition) bool {
	// Look up the filter from the registry
	filter, exists := fs.filterRegistry[cond.Field]
	if !exists {
		logger.Warnf("Unknown filter field: %s", cond.Field)
		return false
	}

	// Use the filter's Apply method
	return filter.Apply(media, string(cond.Operator), cond.Value)
}

// FilterMediaList filters a list of media items
func (fs *FilterService) FilterMediaList(mediaList []*medias.FfprobeResult, filterStr string) ([]*medias.FfprobeResult, error) {
	expr, err := fs.ParseFilter(filterStr)
	if err != nil {
		return nil, err
	}

	if len(expr.Conditions) == 0 {
		return mediaList, nil
	}

	filtered := make([]*medias.FfprobeResult, 0)
	for _, media := range mediaList {
		if fs.ApplyFilter(media, expr) {
			filtered = append(filtered, media)
		}
	}

	logger.Infof("Filtered %d/%d media items", len(filtered), len(mediaList))
	return filtered, nil
}
