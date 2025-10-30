package components

import (
	"context"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/services"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

// MergeVideosComponent provides UI for merging multiple videos
type MergeVideosComponent struct {
	widget.BaseWidget

	window        fyne.Window
	ffmpegService *services.FFmpegService
	selectedFiles []*medias.FfprobeResult

	filesList     *widget.List
	mergeButton   *widget.Button
	outputEntry   *widget.Entry
	outputRow     *fyne.Container
	progressBar   *widget.ProgressBar
	statusLabel   *widget.Label
	refreshButton *widget.Button

	refreshList func() []*medias.FfprobeResult
	onComplete  func(outputPath string)
}

// NewMergeVideosComponent creates a new merge videos component
func NewMergeVideosComponent(window fyne.Window, ffmpegService *services.FFmpegService, refreshList func() []*medias.FfprobeResult) *MergeVideosComponent {

	mvc := &MergeVideosComponent{
		window:        window,
		ffmpegService: ffmpegService,
		selectedFiles: refreshList(),
	}
	mvc.refreshList = refreshList
	mvc.initUI()
	mvc.ExtendBaseWidget(mvc)
	return mvc
}

func (mvc *MergeVideosComponent) initUI() {
	// Files list with reordering capability
	mvc.filesList = widget.NewList(
		func() int {
			return len(mvc.selectedFiles)
		},
		func() fyne.CanvasObject {
			return NewMergeItem()
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			mi := obj.(*MergeItem)
			file := mvc.selectedFiles[id]
			mi.SetLabel(filepath.Base(file.Format.Filename))

			mi.SetUpButtonOnTapped(func() {
				mvc.moveFileUp(int(id))
			})
			mi.SetDownButtonOnTapped(func() {
				mvc.moveFileDown(int(id))
			})
			mi.SetRemoveButtonOnTapped(func() {
				mvc.removeFile(int(id))
			})

			mi.DisableUpButton()
			mi.DisableDownButton()
			if id > 0 {
				mi.DisableUpButton()
			}
			if id < len(mvc.selectedFiles)-1 {
				mi.EnableDownButton()
			}
		},
	)

	// Output file entry
	mvc.outputEntry = widget.NewEntry()
	mvc.outputEntry.SetPlaceHolder("Output file path (e.g., merged_output.mp4)")
	mvc.outputEntry.Text = "merged_output.mp4"

	browseButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			mvc.outputEntry.SetText(writer.URI().Path())
			writer.Close()
		}, mvc.window)
	})

	mvc.outputRow = container.New(layout.NewFormLayout(), browseButton, mvc.outputEntry)

	// Progress bar
	mvc.progressBar = widget.NewProgressBar()
	mvc.progressBar.Hide()

	// Status label
	mvc.statusLabel = widget.NewLabel("")
	mvc.statusLabel.Hide()

	// Merge button
	mvc.mergeButton = widget.NewButtonWithIcon("Merge Videos", theme.MediaPlayIcon(), func() {
		mvc.startMerge()
	})
	mvc.mergeButton.Importance = widget.HighImportance

	mvc.refreshButton = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		mvc.selectedFiles = mvc.refreshList()
		mvc.filesList.Refresh()
	})

	if len(mvc.selectedFiles) < 2 {
		mvc.mergeButton.Disable()
	}
}

func (mvc *MergeVideosComponent) CreateRenderer() fyne.WidgetRenderer {
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("Merge %d Videos", len(mvc.selectedFiles)),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	instructions := widget.NewLabel("Reorder files using up/down arrows. Files will be merged in the order shown.")
	instructions.Wrapping = fyne.TextWrapWord

	content := container.NewBorder(
		container.NewVBox(
			header,
			widget.NewSeparator(),
			container.New(layout.NewFormLayout(), mvc.refreshButton, instructions),
			widget.NewLabel("Files to merge:"),
		),
		container.NewVBox(
			widget.NewLabel("Output file:"),
			mvc.outputRow,
			widget.NewLabel(""),
			mvc.progressBar,
			mvc.statusLabel,
			widget.NewLabel(""),
			mvc.mergeButton,
		),
		nil,
		nil,
		container.NewStack(container.NewScroll(mvc.filesList)),
	)

	return widget.NewSimpleRenderer(content)
}

func (mvc *MergeVideosComponent) moveFileUp(index int) {
	if index <= 0 || index >= len(mvc.selectedFiles) {
		return
	}
	mvc.selectedFiles[index], mvc.selectedFiles[index-1] = mvc.selectedFiles[index-1], mvc.selectedFiles[index]
	mvc.filesList.Refresh()
}

func (mvc *MergeVideosComponent) moveFileDown(index int) {
	if index < 0 || index >= len(mvc.selectedFiles)-1 {
		return
	}
	mvc.selectedFiles[index], mvc.selectedFiles[index+1] = mvc.selectedFiles[index+1], mvc.selectedFiles[index]
	mvc.filesList.Refresh()
}

func (mvc *MergeVideosComponent) removeFile(index int) {
	if index < 0 || index >= len(mvc.selectedFiles) {
		return
	}
	mvc.selectedFiles = append(mvc.selectedFiles[:index], mvc.selectedFiles[index+1:]...)
	mvc.filesList.Refresh()

	if len(mvc.selectedFiles) < 2 {
		mvc.mergeButton.Disable()
	}
}

func (mvc *MergeVideosComponent) startMerge() {
	outputPath := mvc.outputEntry.Text
	if outputPath == "" {
		dialog.ShowError(fmt.Errorf("please specify an output file"), mvc.window)
		return
	}

	// Disable UI during merge
	mvc.mergeButton.Disable()
	mvc.outputEntry.Disable()
	mvc.progressBar.Show()
	mvc.progressBar.SetValue(0)
	mvc.statusLabel.SetText("Merging videos...")
	mvc.statusLabel.Show()

	// Extract file paths
	inputPaths := make([]string, len(mvc.selectedFiles))
	for i, file := range mvc.selectedFiles {
		inputPaths[i] = file.Format.Filename
	}

	// Start merge in background
	go func() {
		ctx := context.Background()
		err := mvc.ffmpegService.MergeVideos(ctx, inputPaths, outputPath, func(progress float64, message string) {
			mvc.progressBar.SetValue(progress)
			mvc.statusLabel.SetText(message)
		})

		// Re-enable UI
		mvc.mergeButton.Enable()
		mvc.outputEntry.Enable()

		if err != nil {
			logger.Errorf("Merge failed: %v", err)
			mvc.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			dialog.ShowError(err, mvc.window)
		} else {
			mvc.statusLabel.SetText(fmt.Sprintf("Successfully merged to: %s", outputPath))
			dialog.ShowInformation("Success", fmt.Sprintf("Videos merged successfully!\n\nOutput: %s", outputPath), mvc.window)

			if mvc.onComplete != nil {
				mvc.onComplete(outputPath)
			}
		}
	}()
}
