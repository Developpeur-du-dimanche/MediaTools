package components

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/ncruces/zenity"
)

// SettingsDialog represents the settings dialog
type SettingsDialog struct {
	app                 fyne.App
	window              fyne.Window
	dialog              *widget.PopUp
	languageSelector    *widget.Select
	ffmpegPathEntry     *widget.Entry
	ffmpegChangedLabel  *widget.Label
	onFFmpegPathChanged func(string)
}

// NewSettingsDialog creates a new settings dialog
func NewSettingsDialog(app fyne.App, window fyne.Window, onFFmpegPathChanged func(string)) *SettingsDialog {
	sd := &SettingsDialog{
		app:                 app,
		window:              window,
		onFFmpegPathChanged: onFFmpegPathChanged,
	}

	return sd
}

// Show displays the settings dialog
func (sd *SettingsDialog) Show() {
	// Language selector
	sd.languageSelector = widget.NewSelect([]string{"English", "Français"}, sd.onLanguageChanged)
	sd.languageSelector.PlaceHolder = lang.L("Language")

	// Set initial value based on current locale
	currentLang := sd.app.Preferences().StringWithFallback("language", "en")
	if currentLang == "fr" {
		sd.languageSelector.SetSelected("Français")
	} else {
		sd.languageSelector.SetSelected("English")
	}

	// Language section
	languageLabel := widget.NewLabelWithStyle(lang.L("Language"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	languageSection := container.NewVBox(
		languageLabel,
		sd.languageSelector,
		widget.NewSeparator(),
	)

	// FFmpeg path section
	ffmpegLabel := widget.NewLabelWithStyle(lang.L("FFmpegPath"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sd.ffmpegPathEntry = widget.NewEntry()
	sd.ffmpegPathEntry.SetPlaceHolder(lang.L("FFmpegPathPlaceholder"))

	sd.ffmpegChangedLabel = widget.NewLabel("")

	// Load saved FFmpeg path
	savedPath := sd.app.Preferences().StringWithFallback("ffmpeg_path", "ffmpeg")
	sd.ffmpegPathEntry.SetText(savedPath)

	browseButton := widget.NewButton(lang.L("Browse"), sd.onBrowseFFmpegPath)

	savePathButton := widget.NewButton(lang.L("Save"), func() {
		newPath := sd.ffmpegPathEntry.Text
		if newPath == "" {
			newPath = "ffmpeg"
		}
		sd.app.Preferences().SetString("ffmpeg_path", newPath)
		if sd.onFFmpegPathChanged != nil {
			sd.onFFmpegPathChanged(newPath)
		}

		// Show confirmation
		sd.ffmpegChangedLabel.SetText(lang.L("FFmpegPathSaved"))
		go func() {
			time.Sleep(10 * time.Second)
			sd.ffmpegChangedLabel.SetText("")
		}()
	})

	ffmpegPathRow := container.NewBorder(nil, nil, nil, container.NewHBox(browseButton, savePathButton), sd.ffmpegPathEntry)

	ffmpegSection := container.NewVBox(
		ffmpegLabel,
		ffmpegPathRow,
		sd.ffmpegChangedLabel,
		widget.NewSeparator(),
	)

	// Close button
	closeButton := widget.NewButton(lang.L("Close"), func() {
		if sd.dialog != nil {
			sd.dialog.Hide()
		}
	})

	// Main content
	content := container.NewVBox(
		widget.NewLabelWithStyle(lang.L("Settings"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		languageSection,
		ffmpegSection,
		container.NewPadded(),
		container.NewCenter(closeButton),
	)

	// Create popup
	sd.dialog = widget.NewModalPopUp(
		container.NewPadded(content),
		sd.window.Canvas(),
	)
	sd.dialog.Resize(fyne.NewSize(400, 250))
	sd.dialog.Show()
}

func (sd *SettingsDialog) onLanguageChanged(selected string) {
	var langCode string
	switch selected {
	case "Français":
		langCode = "fr"
	default:
		langCode = "en"
	}

	currentLanguage := sd.app.Preferences().StringWithFallback("language", "en")
	if langCode == currentLanguage {
		return
	}

	// Save preference
	sd.app.Preferences().SetString("language", langCode)

	// Show restart notification
	dialogLabel := widget.NewLabel("Please restart the application to apply the language change.\nVeuillez redémarrer l'application pour appliquer le changement de langue.")
	dialogLabel.Wrapping = fyne.TextWrapWord

	var popup *widget.PopUp

	okButton := widget.NewButton("OK", func() {
		if popup != nil {
			popup.Hide()
		}
	})

	dialogContainer := container.NewVBox(
		dialogLabel,
		widget.NewSeparator(),
		container.NewCenter(okButton),
	)

	popup = widget.NewModalPopUp(dialogContainer, sd.window.Canvas())
	popup.Resize(fyne.NewSize(400, 150))
	popup.Show()
}

func (sd *SettingsDialog) onBrowseFFmpegPath() {
	// Use native file explorer dialog
	filePath, err := zenity.SelectFile(
		zenity.Title("Select FFmpeg Executable"),
		zenity.FileFilters{
			{Name: "Executable Files", Patterns: []string{"*.exe", "ffmpeg", "ffmpeg.exe"}},
			{Name: "All Files", Patterns: []string{"*.*"}},
		},
	)
	if err != nil {
		// User cancelled or error occurred
		logger.Debugf("FFmpeg file selection cancelled or error: %v", err)
		return
	}

	if filePath == "" {
		return
	}

	sd.ffmpegPathEntry.SetText(filePath)
}
