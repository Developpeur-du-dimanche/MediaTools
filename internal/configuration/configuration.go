package configuration

import (
	"strings"

	"fyne.io/fyne/v2"
)

type Configuration interface {
	SaveExtensions(extensions []string)
	Get(key string) string
	GetList(key string) []string
	IsValidExtension(string) bool
}

type configuration struct {
	extensions  []string
	preferences fyne.Preferences
}

func NewConfiguration(app fyne.App) Configuration {

	preferences := app.Preferences()
	extensions := preferences.StringList("extension")

	if extensions == nil || len(extensions) == 0 {
		extensions = []string{
			"mkv",
			"mp4",
			"avi",
			"mov",
			"flv",
			"wmv",
			"webm",
			"mpg",
			"mpeg",
			"m4v",
		}
		preferences.SetStringList("extension", extensions)
	}

	return &configuration{
		extensions:  preferences.StringList("extension"),
		preferences: preferences,
	}
}

func (c *configuration) SaveExtensions(extensions []string) {
	c.preferences.SetStringList("extension", extensions)
	c.extensions = extensions
}

func (c *configuration) Get(key string) string {
	return c.preferences.String(key)
}

func (c *configuration) GetList(key string) []string {
	return c.preferences.StringList(key)
}

func (c *configuration) IsValidExtension(ext string) bool {

	ext = strings.TrimPrefix(ext, ".")

	for _, e := range c.extensions {
		if e == ext {
			return true
		}
	}
	return false
}