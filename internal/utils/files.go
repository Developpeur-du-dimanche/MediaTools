package utils

import (
	"path/filepath"

	"fyne.io/fyne/v2"
)

func IsValidExtensions(filename string, validExtensions []string) bool {
	ext := filepath.Ext(filename)
	for _, validExt := range validExtensions {

		if validExt[0] != '.' {
			validExt = "." + validExt
		}

		if ext == validExt {
			return true
		}
	}

	return false
}

func GetValidExtensions() []string {
	// This function can be used to retrieve the valid extensions from preferences or a config file
	return fyne.CurrentApp().Preferences().StringListWithFallback("extensions", []string{"mkv", "mp4", "avi"})
}
