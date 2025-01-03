package jsonfilter

import (
	"testing"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func TestCheck(t *testing.T) {
	// Données JSON simulées
	data := ffprobe.ProbeData{
		Format: &ffprobe.Format{
			BitRate: "128000",
		},
		Streams: []*ffprobe.Stream{
			{
				TagList: ffprobe.Tags{
					"Language": "eng",
				},
			},
		},
	}

	tests := []struct {
		name     string
		jsonPath string
		value    string
		expected bool
	}{
		{
			name:     "Valid BitRate Check",
			jsonPath: "$.format.bit_rate",
			value:    "128000",
			expected: true,
		},
		{
			name:     "Invalid BitRate Check",
			jsonPath: "$.format.BitRate",
			value:    "256000",
			expected: false,
		},
		{
			name:     "Valid Language Check",
			jsonPath: "$.streams[0].tags.Language",
			value:    "eng",
			expected: true,
		},
		{
			name:     "Invalid Language Check",
			jsonPath: "$.streams[0].tags.Language",
			value:    "fr",
			expected: false,
		},
		{
			name:     "Non-existent Path",
			jsonPath: "$.format.NonExistentKey",
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
			result := f.Check(&data, test.value)

			// Vérification du résultat
			if result != test.expected {
				t.Errorf("Check(%q, %q) = %v; want %v", test.jsonPath, test.value, result, test.expected)
			}
		})
	}
}
