package filters

import (
	"strconv"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type TotalBitrateFilter struct{}

func (f TotalBitrateFilter) Apply(data *medias.FfprobeResult, operator string, value string) bool {
	actualBitrate, err := strconv.ParseInt(data.Format.Bitrate, 10, 64)
	if err != nil {
		return false
	}

	targetBitrate := parseBitrateValue(value)
	if targetBitrate == 0 {
		return false
	}

	return compareNumeric(actualBitrate, operator, targetBitrate)
}

func (f TotalBitrateFilter) GetFieldConfig() FilterFieldConfig {
	return FilterFieldConfig{
		Key:         "BITRATE",
		DisplayName: "Bitrate (Total)",
		Type:        FieldTypeNumeric,
		Placeholder: "e.g., 2000kbps or 2mbps",
	}
}
