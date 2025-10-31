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

// FilterGroup represents a group of conditions that must all match on the same stream
// Used for parenthesized expressions like (AUDIO_LANGUAGE IS FR AND AUDIO_CODEC IS AAC)
type FilterGroup struct {
	Conditions []FilterCondition
	Operators  []LogicalOperator // Operators within the group (only AND supported for now)
}

// FilterElement can be either a single condition or a group of conditions
type FilterElement struct {
	IsGroup   bool
	Condition *FilterCondition
	Group     *FilterGroup
}

// FilterExpression represents a complete filter with logical operators
type FilterExpression struct {
	Elements  []FilterElement   // Can be conditions or groups
	Operators []LogicalOperator // Operators between elements
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

// ParseFilter parses a filter string with support for grouping via parentheses
// Examples:
//   - "BITRATE > 2000 AND AUDIO_LANGUAGE IS FR" (ungrouped)
//   - "(AUDIO_LANGUAGE IS FR AND AUDIO_CODEC IS AAC)" (grouped - same stream)
//   - "BITRATE > 1000 AND (AUDIO_LANGUAGE IS FR AND AUDIO_CODEC IS AAC)" (mixed)
func (fs *FilterService) ParseFilter(filterStr string) (*FilterExpression, error) {
	if strings.TrimSpace(filterStr) == "" {
		return &FilterExpression{}, nil
	}

	logger.Debugf("Parsing filter: %s", filterStr)

	// Tokenize the filter string
	tokens, err := fs.tokenize(filterStr)
	if err != nil {
		return nil, err
	}

	// Parse tokens into expression
	expr, err := fs.parseTokens(tokens)
	if err != nil {
		return nil, err
	}

	logger.Infof("Parsed filter: %d elements, %d operators", len(expr.Elements), len(expr.Operators))
	return expr, nil
}

// tokenize splits the filter string into tokens, preserving parentheses
func (fs *FilterService) tokenize(filterStr string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	parenDepth := 0

	for _, ch := range filterStr {
		switch ch {
		case '(':
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
			tokens = append(tokens, "(")
			parenDepth++
		case ')':
			if current.Len() > 0 {
				tokens = append(tokens, strings.TrimSpace(current.String()))
				current.Reset()
			}
			tokens = append(tokens, ")")
			parenDepth--
			if parenDepth < 0 {
				return nil, fmt.Errorf("unmatched closing parenthesis")
			}
		case ' ':
			if current.Len() > 0 {
				token := strings.TrimSpace(current.String())
				if token != "" {
					tokens = append(tokens, token)
				}
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(current.String()))
	}

	if parenDepth != 0 {
		return nil, fmt.Errorf("unmatched opening parenthesis")
	}

	return tokens, nil
}

// parseTokens converts tokens into a FilterExpression
func (fs *FilterService) parseTokens(tokens []string) (*FilterExpression, error) {
	expr := &FilterExpression{
		Elements:  []FilterElement{},
		Operators: []LogicalOperator{},
	}

	i := 0
	for i < len(tokens) {
		token := tokens[i]

		// Skip logical operators - they'll be handled separately
		if token == "AND" || token == "OR" {
			expr.Operators = append(expr.Operators, LogicalOperator(token))
			i++
			continue
		}

		// Handle grouped expressions (parentheses)
		if token == "(" {
			groupTokens, endIdx := fs.extractGroup(tokens, i)
			group, err := fs.parseGroup(groupTokens)
			if err != nil {
				return nil, fmt.Errorf("failed to parse group: %w", err)
			}
			expr.Elements = append(expr.Elements, FilterElement{
				IsGroup: true,
				Group:   group,
			})
			i = endIdx + 1
			continue
		}

		// Handle single condition
		condStr := fs.buildConditionString(tokens, &i)
		condition, err := fs.parseCondition(condStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse condition '%s': %w", condStr, err)
		}
		expr.Elements = append(expr.Elements, FilterElement{
			IsGroup:   false,
			Condition: condition,
		})
	}

	// Validate that we have the right number of operators
	if len(expr.Elements) > 0 && len(expr.Elements) != len(expr.Operators)+1 {
		return nil, fmt.Errorf("invalid filter expression: mismatched elements and operators")
	}

	return expr, nil
}

// extractGroup extracts tokens between parentheses
func (fs *FilterService) extractGroup(tokens []string, startIdx int) ([]string, int) {
	depth := 0
	start := startIdx + 1 // Skip opening paren

	for i := startIdx; i < len(tokens); i++ {
		if tokens[i] == "(" {
			depth++
		} else if tokens[i] == ")" {
			depth--
			if depth == 0 {
				return tokens[start:i], i
			}
		}
	}

	return tokens[start:], len(tokens) - 1
}

// parseGroup parses tokens within parentheses into a FilterGroup
func (fs *FilterService) parseGroup(tokens []string) (*FilterGroup, error) {
	group := &FilterGroup{
		Conditions: []FilterCondition{},
		Operators:  []LogicalOperator{},
	}

	i := 0
	for i < len(tokens) {
		token := tokens[i]

		if token == "AND" || token == "OR" {
			group.Operators = append(group.Operators, LogicalOperator(token))
			i++
			continue
		}

		condStr := fs.buildConditionString(tokens, &i)
		condition, err := fs.parseCondition(condStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse group condition '%s': %w", condStr, err)
		}
		group.Conditions = append(group.Conditions, *condition)
	}

	if len(group.Conditions) != len(group.Operators)+1 {
		return nil, fmt.Errorf("invalid group: mismatched conditions and operators")
	}

	return group, nil
}

// buildConditionString builds a condition string from tokens until we hit an operator or end
func (fs *FilterService) buildConditionString(tokens []string, idx *int) string {
	var parts []string

	for *idx < len(tokens) {
		token := tokens[*idx]
		if token == "AND" || token == "OR" || token == "(" || token == ")" {
			break
		}
		parts = append(parts, token)
		*idx++
	}

	return strings.Join(parts, " ")
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
	if len(expr.Elements) == 0 {
		return true // No filter = all pass
	}

	// Evaluate first element (condition or group)
	result := fs.evaluateElement(media, expr.Elements[0])

	// Apply logical operators between elements
	for i, op := range expr.Operators {
		nextResult := fs.evaluateElement(media, expr.Elements[i+1])

		switch op {
		case LogicalAnd:
			result = result && nextResult
		case LogicalOr:
			result = result || nextResult
		}
	}

	return result
}

// evaluateElement evaluates a single element (condition or group)
func (fs *FilterService) evaluateElement(media *medias.FfprobeResult, element FilterElement) bool {
	if element.IsGroup {
		return fs.evaluateGroup(media, element.Group)
	}
	return fs.evaluateCondition(media, *element.Condition)
}

// evaluateGroup evaluates a group of conditions on the same stream
// For a group to match, at least one stream must satisfy ALL conditions in the group
func (fs *FilterService) evaluateGroup(media *medias.FfprobeResult, group *FilterGroup) bool {
	if len(group.Conditions) == 0 {
		return true
	}

	// Determine which type of streams we're filtering
	// All conditions in a group should target the same stream type
	firstField := group.Conditions[0].Field

	switch firstField {
	case FieldAudioCodec, FieldAudioLanguage, FieldAudioChannels, FieldAudioBitrate:
		// Check if at least one audio stream satisfies all conditions
		return fs.evaluateGroupOnAudioStreams(media, group)

	case FieldVideoCodec, FieldWidth, FieldHeight, FieldVideoBitrate:
		// Check if at least one video stream satisfies all conditions
		return fs.evaluateGroupOnVideoStreams(media, group)

	case FieldSubLanguage:
		// Check if at least one subtitle stream satisfies all conditions
		return fs.evaluateGroupOnSubtitleStreams(media, group)

	default:
		// For file-level filters, evaluate conditions normally
		result := fs.evaluateCondition(media, group.Conditions[0])
		for i, op := range group.Operators {
			nextResult := fs.evaluateCondition(media, group.Conditions[i+1])
			switch op {
			case LogicalAnd:
				result = result && nextResult
			case LogicalOr:
				result = result || nextResult
			}
		}
		return result
	}
}

// evaluateGroupOnAudioStreams checks if at least one audio stream satisfies all conditions in the group
func (fs *FilterService) evaluateGroupOnAudioStreams(media *medias.FfprobeResult, group *FilterGroup) bool {
	for _, audio := range media.Audios {
		// Create a temporary FfprobeResult with only this audio stream
		tempMedia := &medias.FfprobeResult{
			Format:    media.Format,
			Videos:    []medias.Video{},
			Audios:    []medias.Audio{audio},
			Subtitles: []medias.Subtitle{},
		}

		// Check if all conditions in the group match this audio stream
		allMatch := true
		for _, cond := range group.Conditions {
			if !fs.evaluateCondition(tempMedia, cond) {
				allMatch = false
				break
			}
		}

		if allMatch {
			return true // Found an audio stream that matches all conditions
		}
	}

	return false // No audio stream matches all conditions
}

// evaluateGroupOnVideoStreams checks if at least one video stream satisfies all conditions in the group
func (fs *FilterService) evaluateGroupOnVideoStreams(media *medias.FfprobeResult, group *FilterGroup) bool {
	for _, video := range media.Videos {
		// Create a temporary FfprobeResult with only this video stream
		tempMedia := &medias.FfprobeResult{
			Format:    media.Format,
			Videos:    []medias.Video{video},
			Audios:    []medias.Audio{},
			Subtitles: []medias.Subtitle{},
		}

		// Check if all conditions in the group match this video stream
		allMatch := true
		for _, cond := range group.Conditions {
			if !fs.evaluateCondition(tempMedia, cond) {
				allMatch = false
				break
			}
		}

		if allMatch {
			return true // Found a video stream that matches all conditions
		}
	}

	return false // No video stream matches all conditions
}

// evaluateGroupOnSubtitleStreams checks if at least one subtitle stream satisfies all conditions in the group
func (fs *FilterService) evaluateGroupOnSubtitleStreams(media *medias.FfprobeResult, group *FilterGroup) bool {
	for _, subtitle := range media.Subtitles {
		// Create a temporary FfprobeResult with only this subtitle stream
		tempMedia := &medias.FfprobeResult{
			Format:    media.Format,
			Videos:    []medias.Video{},
			Audios:    []medias.Audio{},
			Subtitles: []medias.Subtitle{subtitle},
		}

		// Check if all conditions in the group match this subtitle stream
		allMatch := true
		for _, cond := range group.Conditions {
			if !fs.evaluateCondition(tempMedia, cond) {
				allMatch = false
				break
			}
		}

		if allMatch {
			return true // Found a subtitle stream that matches all conditions
		}
	}

	return false // No subtitle stream matches all conditions
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

	if len(expr.Elements) == 0 {
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
