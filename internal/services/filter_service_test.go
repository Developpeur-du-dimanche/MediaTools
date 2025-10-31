package services

import (
	"testing"
	"time"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

func TestGroupedFilters(t *testing.T) {
	fs := NewFilterService()

	// Create test data: 2 audio streams
	// Audio 1: FR, AAC, 2 channels
	// Audio 2: ENG, MP3, 6 channels
	testMedia := &medias.FfprobeResult{
		Format: medias.FfprobeData{
			Filename:        "test.mkv",
			DurationSeconds: time.Duration(3600 * time.Second),
			Size:            "1000000",
			BitRate:         "2000000",
		},
		Videos: []medias.Video{
			{StreamIndex: 0, CodecName: "h264", Width: 1920, Height: 1080},
		},
		Audios: []medias.Audio{
			{StreamIndex: 1, CodecName: "aac", Language: "fre", Channels: 2},
			{StreamIndex: 2, CodecName: "mp3", Language: "eng", Channels: 6},
		},
		Subtitles: []medias.Subtitle{},
	}

	tests := []struct {
		name     string
		filter   string
		expected bool
		reason   string
	}{
		{
			name:     "Ungrouped - matches different streams",
			filter:   "AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS mp3",
			expected: true,
			reason:   "Has french audio (stream 1) AND has mp3 audio (stream 2)",
		},
		{
			name:     "Grouped - same stream required",
			filter:   "(AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS mp3)",
			expected: false,
			reason:   "No single audio stream is both french AND mp3",
		},
		{
			name:     "Grouped - matching stream exists",
			filter:   "(AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS aac)",
			expected: true,
			reason:   "Audio stream 1 is french AND aac",
		},
		{
			name:     "Grouped with channels",
			filter:   "(AUDIO_LANGUAGE IS eng AND AUDIO_CHANNELS IS 6)",
			expected: true,
			reason:   "Audio stream 2 is english AND 6 channels",
		},
		{
			name:     "Grouped - no match",
			filter:   "(AUDIO_LANGUAGE IS fre AND AUDIO_CHANNELS IS 6)",
			expected: false,
			reason:   "No audio stream is french AND 6 channels",
		},
		{
			name:     "Mixed grouped and ungrouped",
			filter:   "BITRATE > 1000 AND (AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS aac)",
			expected: true,
			reason:   "Bitrate is 2000000 (>1000) AND has french+aac audio",
		},
		{
			name:     "Multiple groups with OR",
			filter:   "(AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS aac) OR (AUDIO_LANGUAGE IS eng AND AUDIO_CHANNELS IS 6)",
			expected: true,
			reason:   "Has french+aac audio OR english+6ch audio (both exist)",
		},
		{
			name:     "Multiple groups - only one matches",
			filter:   "(AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS mp3) OR (AUDIO_LANGUAGE IS eng AND AUDIO_CHANNELS IS 6)",
			expected: true,
			reason:   "First group fails but second group matches",
		},
		{
			name:     "Multiple groups - none match",
			filter:   "(AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS mp3) OR (AUDIO_LANGUAGE IS eng AND AUDIO_CHANNELS IS 2)",
			expected: false,
			reason:   "Neither group has a matching stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := fs.ParseFilter(tt.filter)
			if err != nil {
				t.Fatalf("Failed to parse filter '%s': %v", tt.filter, err)
			}

			result := fs.ApplyFilter(testMedia, expr)
			if result != tt.expected {
				t.Errorf("Filter '%s' returned %v, expected %v. Reason: %s",
					tt.filter, result, tt.expected, tt.reason)
			}
		})
	}
}

func TestParenthesesParsing(t *testing.T) {
	fs := NewFilterService()

	tests := []struct {
		name        string
		filter      string
		shouldError bool
		description string
	}{
		{
			name:        "Simple grouped filter",
			filter:      "(AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS aac)",
			shouldError: false,
			description: "Should parse simple parentheses",
		},
		{
			name:        "Unmatched opening paren",
			filter:      "(AUDIO_LANGUAGE IS fre",
			shouldError: true,
			description: "Should error on unmatched opening parenthesis",
		},
		{
			name:        "Unmatched closing paren",
			filter:      "AUDIO_LANGUAGE IS fre)",
			shouldError: true,
			description: "Should error on unmatched closing parenthesis",
		},
		{
			name:        "Mixed grouped and ungrouped",
			filter:      "BITRATE > 1000 AND (AUDIO_LANGUAGE IS fre AND AUDIO_CODEC IS aac)",
			shouldError: false,
			description: "Should parse mixed expressions",
		},
		{
			name:        "Multiple groups",
			filter:      "(AUDIO_LANGUAGE IS fre) OR (AUDIO_CODEC IS aac)",
			shouldError: false,
			description: "Should parse multiple groups",
		},
		{
			name:        "Empty filter",
			filter:      "",
			shouldError: false,
			description: "Should handle empty filter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fs.ParseFilter(tt.filter)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for filter '%s' but got none. %s", tt.filter, tt.description)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for filter '%s': %v. %s", tt.filter, err, tt.description)
			}
		})
	}
}
