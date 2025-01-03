package jsonfilter

import (
	"testing"

	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func TestCheck(t *testing.T) {
	// Données JSON simulées

	fileMetadata := helper.FileMetadata{
		FileName:  "test.mp4",
		Directory: "/tmp",
		Bitrate:   "128000",
		Audio: []*ffprobe.Stream{
			{
				TagList: ffprobe.Tags{
					"Language": "eng",
				},
			},
			{
				TagList: ffprobe.Tags{
					"Language": "ita",
				},
			},
		},

		Format:    "/path/to/test.mp4",
		Codec:     "h264",
		Duration:  0,
		Size:      "1024",
		Container: "mp4",
		Extension: "mp4",
	}

	tests := []struct {
		name     string
		jsonPath string
		value    string
		expected bool
	}{
		{
			name:     "Valid BitRate Check",
			jsonPath: "$.Bitrate",
			value:    "128000",
			expected: true,
		},
		{
			name:     "Invalid BitRate Check",
			jsonPath: "$.BitRate",
			value:    "256000",
			expected: false,
		},
		{
			name:     "Valid Language Check",
			jsonPath: "$.Audio[*].TagList.Language",
			value:    "eng",
			expected: true,
		},
		{
			name:     "Valid Language Check",
			jsonPath: "$.Audio[*].TagList.Language",
			value:    "ita",
			expected: true,
		},
		{
			name:     "Invalid Language Check",
			jsonPath: "$.Audio[0].TagList.Language",
			value:    "fr",
			expected: false,
		},
		{
			name:     "Non-existent Path",
			jsonPath: "$.NonExistentKey",
			value:    "value",
			expected: false,
		},
		{
			name:     "Invalid JsonPath Syntax",
			jsonPath: "$.format[BitRate",
			value:    "128000",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Création d'un filtre
			f := filter{JsonPath: test.jsonPath}

			// Appel de la méthode Check
			result := f.Check(&fileMetadata, Equals, test.value)

			// Vérification du résultat
			if result != test.expected {
				t.Errorf("Check(%q, %q) = %v; want %v", test.jsonPath, test.value, result, test.expected)
			}
		})
	}
}
