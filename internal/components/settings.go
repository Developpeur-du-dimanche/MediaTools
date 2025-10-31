package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

// SettingsDialog represents the settings dialog
type SettingsDialog struct {
	app              fyne.App
	window           fyne.Window
	dialog           *widget.PopUp
	languageSelector *widget.Select
}

// NewSettingsDialog creates a new settings dialog
func NewSettingsDialog(app fyne.App, window fyne.Window) *SettingsDialog {
	sd := &SettingsDialog{
		app:    app,
		window: window,
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
