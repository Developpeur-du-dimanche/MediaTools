package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/ncruces/zenity"
)

type OpenFile struct {
	widget.BaseWidget
	button *widget.Button
	window fyne.Window

	OnFileOpen       func(path string)
	OnScanTerminated func()
}

func NewOpenFile(parent fyne.Window, onFileOpened func(path string)) *OpenFile {
	of := &OpenFile{
		window:     parent,
		OnFileOpen: onFileOpened,
	}
	of.button = widget.NewButtonWithIcon("Open File", theme.FileIcon(), of.openFileDialog)
	of.ExtendBaseWidget(of)
	return of
}

func (of *OpenFile) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(of.button)
}

func (of *OpenFile) openFileDialog() {
	// Use native Windows File Explorer dialog
	filePath, err := zenity.SelectFile(
		zenity.Title("Select Media File"),
		zenity.FileFilters{
			{Name: "Video Files", Patterns: []string{"*.mp4", "*.mkv", "*.avi", "*.mov", "*.wmv", "*.flv", "*.webm", "*.m4v"}},
			{Name: "Audio Files", Patterns: []string{"*.mp3", "*.flac", "*.wav", "*.aac", "*.ogg", "*.m4a", "*.wma"}},
			{Name: "All Files", Patterns: []string{"*.*"}},
		},
	)
	if err != nil {
		// User cancelled or error occurred
		logger.Debugf("File selection cancelled or error: %v", err)
		return
	}

	if filePath == "" {
		return
	}

	if of.OnFileOpen != nil {
		of.OnFileOpen(filePath)
	}

	if of.OnScanTerminated != nil {
		of.OnScanTerminated()
	}
}
