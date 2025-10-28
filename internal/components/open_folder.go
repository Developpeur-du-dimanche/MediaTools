package components

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	fynedialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/services"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/ncruces/zenity"
)

type OpenFolder struct {
	widget.BaseWidget
	button *widget.Button
	window fyne.Window

	progressDialog *fynedialog.CustomDialog
	progressBar    *widget.ProgressBar
	progressLabel  *widget.Label

	onFolderOpen     func(path string)
	onScanProgress   func(progress services.ScanProgress)
	OnScanTerminated func()
	cancelFunc       context.CancelFunc
}

func NewOpenFolder(parent fyne.Window, onFolderOpened func(path string), onScanProgress func(progress services.ScanProgress)) *OpenFolder {
	of := &OpenFolder{
		window:         parent,
		onFolderOpen:   onFolderOpened,
		onScanProgress: onScanProgress,
		progressBar:    widget.NewProgressBar(),
		progressLabel:  widget.NewLabel("Preparing scan..."),
	}

	of.progressDialog = fynedialog.NewCustomWithoutButtons("Scanning Folder",
		widget.NewCard("", "",
			container.NewVBox(
				of.progressLabel,
				of.progressBar,
			),
		),
		of.window,
	)
	of.progressDialog.Hide()
	of.button = widget.NewButtonWithIcon("Open Folder", theme.FolderIcon(), of.openFolderDialog)
	of.ExtendBaseWidget(of)
	return of
}

func (of *OpenFolder) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(of.button)
}

func (of *OpenFolder) openFolderDialog() {
	// Use native Windows File Explorer dialog
	folderPath, err := zenity.SelectFile(
		zenity.Title("Select Folder"),
		zenity.Directory(),
	)
	if err != nil {
		// User cancelled or error occurred
		logger.Debugf("Folder selection cancelled or error: %v", err)
		return
	}

	if folderPath == "" {
		return
	}

	if of.onFolderOpen != nil {
		of.onFolderOpen(folderPath)
	}

	of.scanFolderAsync(folderPath)
}

func (of *OpenFolder) scanFolderAsync(folderPath string) {
	// Show progress dialog
	of.progressBar.SetValue(0)
	of.progressLabel.SetText("Starting scan...")
	of.progressDialog.Show()

	// Notify parent that scan started
	if of.onFolderOpen != nil {
		of.onFolderOpen(folderPath)
	}
}

// UpdateProgress updates the progress dialog with scan progress
func (of *OpenFolder) UpdateProgress(progress services.ScanProgress) {
	if progress.IsComplete {
		of.progressLabel.SetText("Scan complete!")
		of.progressBar.SetValue(1.0)

		// Hide dialog after completion
		of.progressDialog.Hide()

		if of.OnScanTerminated != nil {
			of.OnScanTerminated()
		}
	} else {
		of.progressLabel.SetText(progress.CurrentFile)
		if progress.TotalFiles > 0 {
			of.progressBar.SetValue(float64(progress.FilesScanned) / float64(progress.TotalFiles))
		}
	}

	// Notify parent about progress
	if of.onScanProgress != nil {
		of.onScanProgress(progress)
	}
}

// HideProgress hides the progress dialog
func (of *OpenFolder) HideProgress() {
	of.progressDialog.Hide()
}

// CancelScan cancels an ongoing scan operation
func (of *OpenFolder) CancelScan() {
	if of.cancelFunc != nil {
		of.cancelFunc()
		logger.Info("Folder scan cancelled by user")
	}
}

// SetCancelFunc sets the cancel function for the scan
func (of *OpenFolder) SetCancelFunc(cancel context.CancelFunc) {
	of.cancelFunc = cancel
}
