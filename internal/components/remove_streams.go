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

// RemoveStreamsComponent provides UI for removing streams based on criteria
type RemoveStreamsComponent struct {
	widget.BaseWidget

	window        fyne.Window
	ffmpegService *services.FFmpegService
	selectedFiles []*medias.FfprobeResult

	// UI elements
	operationSelect  *widget.Select
	streamTypeSelect *widget.Select
	criteriaEntry    *widget.Entry
	criteriaSelect   *widget.Select
	outputDirEntry   *widget.Entry
	outputDirRow     *fyne.Container
	progressBar      *widget.ProgressBar
	statusLabel      *widget.Label
	processButton    *widget.Button
	filesList        *widget.List

	onComplete func(results []string)
}

// NewRemoveStreamsComponent creates a new component for stream removal
func NewRemoveStreamsComponent(window fyne.Window, files []*medias.FfprobeResult, ffmpegService *services.FFmpegService) *RemoveStreamsComponent {
	rsc := &RemoveStreamsComponent{
		window:        window,
		ffmpegService: ffmpegService,
		selectedFiles: files,
	}

	rsc.initUI()
	rsc.ExtendBaseWidget(rsc)
	return rsc
}

func (rsc *RemoveStreamsComponent) initUI() {
	// Operation selector
	rsc.operationSelect = widget.NewSelect([]string{
		"Remove all streams of type",
		"Remove streams by language",
		"Remove streams by codec",
		"Keep only streams by language",
	}, func(value string) {
		rsc.updateCriteriaUI(value)
	})
	rsc.operationSelect.SetSelected("Remove all streams of type")

	// Stream type selector
	rsc.streamTypeSelect = widget.NewSelect([]string{
		"Audio",
		"Subtitle",
		"Video",
	}, nil)
	rsc.streamTypeSelect.SetSelected("Audio")

	// Criteria entry (for manual input)
	rsc.criteriaEntry = widget.NewEntry()
	rsc.criteriaEntry.SetPlaceHolder("Enter value...")
	rsc.criteriaEntry.Hide()

	// Criteria select (for predefined values)
	rsc.criteriaSelect = widget.NewSelect([]string{}, nil)
	rsc.criteriaSelect.Hide()

	// Output directory
	rsc.outputDirEntry = widget.NewEntry()
	rsc.outputDirEntry.SetPlaceHolder("Output directory")
	rsc.outputDirEntry.Text = "./processed"

	browseDirButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}
			rsc.outputDirEntry.SetText(dir.Path())
		}, rsc.window)
	})

	rsc.outputDirRow = container.NewBorder(nil, nil, nil, browseDirButton, rsc.outputDirEntry)

	// Files list
	rsc.filesList = widget.NewList(
		func() int {
			return len(rsc.selectedFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			file := rsc.selectedFiles[id]
			label.SetText(filepath.Base(file.Format.Filename))
		},
	)

	// Progress bar
	rsc.progressBar = widget.NewProgressBar()
	rsc.progressBar.Hide()

	// Status label
	rsc.statusLabel = widget.NewLabel("")
	rsc.statusLabel.Hide()

	// Process button
	rsc.processButton = widget.NewButtonWithIcon("Process Files", theme.MediaPlayIcon(), func() {
		rsc.startProcessing()
	})
	rsc.processButton.Importance = widget.HighImportance
}

func (rsc *RemoveStreamsComponent) CreateRenderer() fyne.WidgetRenderer {
	header := widget.NewLabelWithStyle(
		fmt.Sprintf("Remove/Keep Streams - %d Files", len(rsc.selectedFiles)),
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	form := container.NewVBox(
		widget.NewLabel("Operation:"),
		rsc.operationSelect,
		widget.NewLabel(""),
		widget.NewLabel("Stream Type:"),
		rsc.streamTypeSelect,
		widget.NewLabel(""),
		widget.NewLabel("Criteria:"),
		rsc.criteriaEntry,
		rsc.criteriaSelect,
		widget.NewLabel(""),
		widget.NewLabel("Output Directory:"),
		rsc.outputDirRow,
	)

	filesSection := container.NewVBox(
		widget.NewLabel("Files to process:"),
		container.NewScroll(rsc.filesList),
	)

	content := container.NewBorder(
		container.NewVBox(
			header,
			widget.NewSeparator(),
			form,
			widget.NewSeparator(),
		),
		container.NewVBox(
			widget.NewLabel(""),
			rsc.progressBar,
			rsc.statusLabel,
			widget.NewLabel(""),
			rsc.processButton,
		),
		nil,
		nil,
		filesSection,
	)

	return widget.NewSimpleRenderer(content)
}

func (rsc *RemoveStreamsComponent) updateCriteriaUI(operation string) {
	rsc.criteriaEntry.Hide()
	rsc.criteriaSelect.Hide()

	switch operation {
	case "Remove all streams of type":
		// No additional criteria needed
		return

	case "Remove streams by language", "Keep only streams by language":
		// Show language selector
		rsc.criteriaSelect.Options = []string{
			"fre", "eng", "spa", "deu", "ita", "jpn", "kor", "chi", "por", "rus", "ara", "hin",
		}
		rsc.criteriaSelect.PlaceHolder = "Select language..."
		rsc.criteriaSelect.Show()

	case "Remove streams by codec":
		// Show codec selector based on stream type
		streamType := rsc.streamTypeSelect.Selected
		if streamType == "Audio" {
			rsc.criteriaSelect.Options = []string{
				"aac", "mp3", "ac3", "eac3", "dts", "flac", "opus", "vorbis", "pcm",
			}
		} else if streamType == "Video" {
			rsc.criteriaSelect.Options = []string{
				"h264", "h265", "hevc", "vp9", "av1", "mpeg4", "mpeg2video", "xvid",
			}
		} else {
			rsc.criteriaEntry.SetPlaceHolder("Enter codec name...")
			rsc.criteriaEntry.Show()
			return
		}
		rsc.criteriaSelect.PlaceHolder = "Select codec..."
		rsc.criteriaSelect.Show()
	}

	rsc.Refresh()
}

func (rsc *RemoveStreamsComponent) startProcessing() {
	// Validate inputs
	outputDir := rsc.outputDirEntry.Text
	if outputDir == "" {
		dialog.ShowError(fmt.Errorf("please specify an output directory"), rsc.window)
		return
	}

	operation := rsc.getOperationType()
	criteria := rsc.getCriteria()

	// Disable UI during processing
	rsc.processButton.Disable()
	rsc.operationSelect.Disable()
	rsc.streamTypeSelect.Disable()
	rsc.criteriaEntry.Disable()
	rsc.criteriaSelect.Disable()
	rsc.outputDirEntry.Disable()
	rsc.progressBar.Show()
	rsc.progressBar.SetValue(0)
	rsc.statusLabel.SetText("Processing files...")
	rsc.statusLabel.Show()

	// Start processing in background
	go func() {
		ctx := context.Background()
		results, err := rsc.ffmpegService.BatchRemoveStreams(
			ctx,
			rsc.selectedFiles,
			operation,
			criteria,
			outputDir,
			func(progress float64, message string) {
				rsc.progressBar.SetValue(progress)
				rsc.statusLabel.SetText(message)
			},
		)

		// Re-enable UI
		rsc.processButton.Enable()
		rsc.operationSelect.Enable()
		rsc.streamTypeSelect.Enable()
		rsc.criteriaEntry.Enable()
		rsc.criteriaSelect.Enable()
		rsc.outputDirEntry.Enable()

		if err != nil {
			logger.Errorf("Processing failed: %v", err)
			rsc.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			dialog.ShowError(err, rsc.window)
		} else {
			rsc.statusLabel.SetText(fmt.Sprintf("Successfully processed %d files", len(results)))
			dialog.ShowInformation(
				"Success",
				fmt.Sprintf("Successfully processed %d/%d files!\n\nOutput directory: %s", len(results), len(rsc.selectedFiles), outputDir),
				rsc.window,
			)

			if rsc.onComplete != nil {
				rsc.onComplete(results)
			}
		}
	}()
}

func (rsc *RemoveStreamsComponent) getOperationType() string {
	switch rsc.operationSelect.Selected {
	case "Remove all streams of type":
		return "remove_by_type"
	case "Remove streams by language":
		return "remove_by_language"
	case "Remove streams by codec":
		return "remove_by_codec"
	case "Keep only streams by language":
		return "keep_language"
	default:
		return "remove_by_type"
	}
}

func (rsc *RemoveStreamsComponent) getCriteria() map[string]string {
	criteria := make(map[string]string)
	criteria["type"] = rsc.streamTypeSelect.Selected

	operation := rsc.operationSelect.Selected
	if operation == "Remove streams by language" || operation == "Keep only streams by language" {
		if rsc.criteriaSelect.Visible() {
			criteria["language"] = rsc.criteriaSelect.Selected
		} else {
			criteria["language"] = rsc.criteriaEntry.Text
		}
	} else if operation == "Remove streams by codec" {
		if rsc.criteriaSelect.Visible() {
			criteria["codec"] = rsc.criteriaSelect.Selected
		} else {
			criteria["codec"] = rsc.criteriaEntry.Text
		}
	}

	return criteria
}

