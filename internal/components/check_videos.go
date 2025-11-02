package components

import (
	"context"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/services"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

// CheckVideosComponent provides UI for checking video integrity
type CheckVideosComponent struct {
	widget.BaseWidget

	window        fyne.Window
	ffmpegService *services.FFmpegService
	selectedFiles []*medias.FfprobeResult

	// UI elements
	resultsList *widget.List
	progressBar *widget.ProgressBar
	statusLabel *widget.Label
	checkButton *widget.Button
	filesList   *widget.List

	// Data
	checkResults []*services.VideoCheckResult

	onComplete func(results []*services.VideoCheckResult)
}

// NewCheckVideosComponent creates a new component for checking videos
func NewCheckVideosComponent(window fyne.Window, files []*medias.FfprobeResult, ffmpegService *services.FFmpegService) *CheckVideosComponent {
	cvc := &CheckVideosComponent{
		window:        window,
		ffmpegService: ffmpegService,
		selectedFiles: files,
		checkResults:  make([]*services.VideoCheckResult, 0),
	}

	cvc.initUI()
	cvc.ExtendBaseWidget(cvc)
	return cvc
}

func (cvc *CheckVideosComponent) initUI() {
	// Files list
	cvc.filesList = widget.NewList(
		func() int {
			return len(cvc.selectedFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			file := cvc.selectedFiles[id]
			label.SetText(filepath.Base(file.Format.Filename))
		},
	)

	// Results list
	cvc.resultsList = widget.NewList(
		func() int {
			return len(cvc.checkResults)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel(""), // Status icon
				widget.NewLabel(""), // Filename
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			hbox := obj.(*fyne.Container)
			if id < len(cvc.checkResults) {
				result := cvc.checkResults[id]
				statusLabel := hbox.Objects[0].(*widget.Label)
				fileLabel := hbox.Objects[1].(*widget.Label)

				if result.IsValid {
					statusLabel.SetText("✓ OK")
				} else {
					statusLabel.SetText("✗ CORRUPTED")
				}
				fileLabel.SetText(filepath.Base(result.FilePath))

				// Add click to show details
				if !result.IsValid && result.Error != "" {
					fileLabel.SetText(fileLabel.Text + " (click for details)")
				}
			}
		},
	)

	cvc.resultsList.OnSelected = func(id widget.ListItemID) {
		if id < len(cvc.checkResults) {
			result := cvc.checkResults[id]
			if !result.IsValid && result.Error != "" {
				dialog.ShowInformation(
					"Error Details",
					fmt.Sprintf("File: %s\n\nErrors:\n%s", filepath.Base(result.FilePath), result.Error),
					cvc.window,
				)
			}
		}
		cvc.resultsList.UnselectAll()
	}

	// Progress bar
	cvc.progressBar = widget.NewProgressBar()
	cvc.progressBar.Hide()

	// Status label
	cvc.statusLabel = widget.NewLabel("")
	cvc.statusLabel.Hide()

	// Check button
	cvc.checkButton = widget.NewButtonWithIcon("Start Checking", theme.MediaPlayIcon(), func() {
		cvc.startChecking()
	})
	cvc.checkButton.Importance = widget.HighImportance
}

func (cvc *CheckVideosComponent) CreateRenderer() fyne.WidgetRenderer {
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("Check Video Integrity - %d Files", len(cvc.selectedFiles)),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	instructions := widget.NewLabel("Click 'Start Checking' to verify all selected videos for corruption.")
	instructions.Wrapping = fyne.TextWrapWord

	filesSection := container.NewBorder(
		widget.NewLabel("Files to check:"),
		nil,
		nil,
		nil,
		cvc.filesList,
	)

	resultsSection := container.NewVBox(
		widget.NewLabel("Results:"),
		cvc.resultsList,
	)

	content := container.NewBorder(
		container.NewVBox(
			header,
			widget.NewSeparator(),
			instructions,
			widget.NewLabel(""),
		),
		container.NewVBox(
			widget.NewLabel(""),
			cvc.progressBar,
			cvc.statusLabel,
			widget.NewLabel(""),
			cvc.checkButton,
		),
		nil,
		nil,
		container.NewHSplit(filesSection, resultsSection),
	)

	return widget.NewSimpleRenderer(content)
}

func (cvc *CheckVideosComponent) startChecking() {
	// Reset results
	cvc.checkResults = make([]*services.VideoCheckResult, 0)
	cvc.resultsList.Refresh()

	// Disable UI during check
	cvc.checkButton.Disable()
	cvc.progressBar.Show()
	cvc.progressBar.SetValue(0)
	cvc.statusLabel.SetText("Checking videos...")
	cvc.statusLabel.Show()

	// Start checking in background
	go func() {
		ctx := context.Background()
		results, err := cvc.ffmpegService.BatchCheckVideos(ctx, cvc.selectedFiles, func(progress float64, message string) {
			cvc.progressBar.SetValue(progress)
			cvc.statusLabel.SetText(message)
		})

		// Re-enable UI
		cvc.checkButton.Enable()

		if err != nil {
			logger.Errorf("Check failed: %v", err)
			cvc.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			dialog.ShowError(err, cvc.window)
			return
		}

		cvc.checkResults = results
		cvc.resultsList.Refresh()

		// Count corrupted files
		corruptedCount := 0
		for _, result := range results {
			if !result.IsValid {
				corruptedCount++
			}
		}

		cvc.statusLabel.SetText(fmt.Sprintf("Complete: %d OK, %d corrupted", len(results)-corruptedCount, corruptedCount))

		if corruptedCount > 0 {
			dialog.ShowInformation(
				"Check Complete",
				fmt.Sprintf("Found %d corrupted file(s) out of %d.\n\nClick on a corrupted file in the results to see details.", corruptedCount, len(results)),
				cvc.window,
			)
		} else {
			dialog.ShowInformation(
				"Check Complete",
				fmt.Sprintf("All %d files are valid!", len(results)),
				cvc.window,
			)
		}

		if cvc.onComplete != nil {
			cvc.onComplete(results)
		}
	}()
}
