package configuration

import (
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
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

func NewConfiguration(app fyne.App, window *fyne.Window) Configuration {

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

	ffmpeg := preferences.String("ffmpeg")

	if ffmpeg == "" {

		// get ffmpeg from path
		path, err := exec.LookPath("ffmpeg")

		if err != nil {
			dialog.ShowInformation("MediaTools", lang.L("ffmpeg_not_found"), *window)
		} else {
			preferences.SetString("ffmpeg", path)
		}

		preferences.SetString("ffmpeg", path)
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
