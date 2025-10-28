package services

import (
	"fmt"
	"strconv"
	"strings"

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
type FilterService struct{}

// NewFilterService creates a new filter service
func NewFilterService() *FilterService {
	return &FilterService{}
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
	switch cond.Field {
	case FieldBitrate:
		return fs.compareBitrateStr(media.Format.BitRate, cond.Operator, cond.Value)

	case FieldVideoBitrate:
		// Video bitrate not available in simplified structure
		return false

	case FieldAudioBitrate:
		// Audio bitrate not available in simplified structure
		return false

	case FieldVideoCodec:
		if len(media.Videos) == 0 {
			return false
		}
		return fs.compareString(media.Videos[0].CodecName, cond.Operator, cond.Value)

	case FieldAudioCodec:
		if len(media.Audios) == 0 {
			return false
		}
		return fs.compareString(media.Audios[0].CodecName, cond.Operator, cond.Value)

	case FieldAudioLanguage:
		return fs.hasAudioLanguage(media, cond.Operator, cond.Value)

	case FieldSubLanguage:
		return fs.hasSubtitleLanguage(media, cond.Operator, cond.Value)

	case FieldWidth:
		if len(media.Videos) == 0 {
			return false
		}
		return fs.compareInt(media.Videos[0].Width, cond.Operator, cond.Value)

	case FieldHeight:
		if len(media.Videos) == 0 {
			return false
		}
		return fs.compareInt(media.Videos[0].Height, cond.Operator, cond.Value)

	case FieldDuration:
		return fs.compareDurationTime(media.Format.DurationSeconds, cond.Operator, cond.Value)

	case FieldFramerate:
		// Framerate not available in simplified structure
		return false

	case FieldAudioChannels:
		if len(media.Audios) == 0 {
			return false
		}
		return fs.compareInt(media.Audios[0].Channels, cond.Operator, cond.Value)

	case FieldHasVideo:
		hasVideo := len(media.Videos) > 0
		return fs.compareBool(hasVideo, cond.Operator, cond.Value)

	case FieldHasAudio:
		hasAudio := len(media.Audios) > 0
		return fs.compareBool(hasAudio, cond.Operator, cond.Value)

	case FieldHasSubtitles:
		hasSubs := len(media.Subtitles) > 0
		return fs.compareBool(hasSubs, cond.Operator, cond.Value)

	default:
		logger.Warnf("Unknown filter field: %s", cond.Field)
		return false
	}
}

// compareBitrateStr compares bitrate values from string format (supports kb/s, mb/s)
func (fs *FilterService) compareBitrateStr(actualBitrateStr string, op FilterOperator, valueStr string) bool {
	// Parse actual bitrate string to int64
	actualBitrate, err := strconv.ParseInt(actualBitrateStr, 10, 64)
	if err != nil {
		return false
	}

	// Parse value with unit (e.g., "2000kbps", "2mbps")
	valueStr = strings.ToLower(strings.TrimSpace(valueStr))
	valueStr = strings.ReplaceAll(valueStr, " ", "")

	var targetBitrate int64

	if strings.HasSuffix(valueStr, "kbps") || strings.HasSuffix(valueStr, "kb/s") {
		numStr := strings.TrimSuffix(strings.TrimSuffix(valueStr, "kbps"), "kb/s")
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return false
		}
		targetBitrate = int64(num * 1000)
	} else if strings.HasSuffix(valueStr, "mbps") || strings.HasSuffix(valueStr, "mb/s") {
		numStr := strings.TrimSuffix(strings.TrimSuffix(valueStr, "mbps"), "mb/s")
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return false
		}
		targetBitrate = int64(num * 1000000)
	} else {
		// Assume bits per second
		num, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return false
		}
		targetBitrate = num
	}

	return fs.compareInt64(actualBitrate, op, targetBitrate)
}

// compareDurationTime compares duration values (time.Duration type)
func (fs *FilterService) compareDurationTime(actualDuration interface{}, op FilterOperator, valueStr string) bool {
	// Convert to seconds
	var actualSeconds float64
	switch v := actualDuration.(type) {
	case float64:
		actualSeconds = v
	case int64: // time.Duration in nanoseconds
		actualSeconds = float64(v) / float64(1000000000)
	default:
		return false
	}

	target, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return false
	}
	return fs.compareFloat(actualSeconds, op, target)
}

// hasAudioLanguage checks if any audio stream matches the language
func (fs *FilterService) hasAudioLanguage(media *medias.FfprobeResult, op FilterOperator, lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))

	for _, audio := range media.Audios {
		audioLang := strings.ToLower(audio.Language)
		if fs.compareString(audioLang, op, lang) {
			return true
		}
	}
	return false
}

// hasSubtitleLanguage checks if any subtitle stream matches the language
func (fs *FilterService) hasSubtitleLanguage(media *medias.FfprobeResult, op FilterOperator, lang string) bool {
	lang = strings.ToLower(strings.TrimSpace(lang))

	for _, sub := range media.Subtitles {
		subLang := strings.ToLower(sub.Language)
		if fs.compareString(subLang, op, lang) {
			return true
		}
	}
	return false
}

// Comparison helpers
func (fs *FilterService) compareInt(actual int, op FilterOperator, valueStr string) bool {
	target, err := strconv.Atoi(valueStr)
	if err != nil {
		return false
	}
	return fs.compareInt64(int64(actual), op, int64(target))
}

func (fs *FilterService) compareInt64(actual int64, op FilterOperator, target int64) bool {
	switch op {
	case OpEquals:
		return actual == target
	case OpNotEquals:
		return actual != target
	case OpGreater:
		return actual > target
	case OpGreaterEq:
		return actual >= target
	case OpLess:
		return actual < target
	case OpLessEq:
		return actual <= target
	default:
		return false
	}
}

func (fs *FilterService) compareFloat(actual float64, op FilterOperator, target float64) bool {
	switch op {
	case OpEquals:
		return actual == target
	case OpNotEquals:
		return actual != target
	case OpGreater:
		return actual > target
	case OpGreaterEq:
		return actual >= target
	case OpLess:
		return actual < target
	case OpLessEq:
		return actual <= target
	default:
		return false
	}
}

func (fs *FilterService) compareString(actual string, op FilterOperator, target string) bool {
	actual = strings.ToLower(strings.TrimSpace(actual))
	target = strings.ToLower(strings.TrimSpace(target))

	switch op {
	case OpEquals:
		return actual == target
	case OpNotEquals:
		return actual != target
	case OpContains:
		return strings.Contains(actual, target)
	case OpNotContains:
		return !strings.Contains(actual, target)
	default:
		return false
	}
}

func (fs *FilterService) compareBool(actual bool, op FilterOperator, valueStr string) bool {
	target := strings.ToLower(valueStr) == "true" || valueStr == "1"

	switch op {
	case OpEquals:
		return actual == target
	case OpNotEquals:
		return actual != target
	default:
		return false
	}
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
