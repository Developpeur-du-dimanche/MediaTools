package components

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components/customs"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type MergeFilesContract interface {
	MergeFiles(files []string, output string) error
}

type inputFiles struct {
	dir      string
	filename string
	position int
	enabled  bool
	metadata *helper.FileMetadata
}

type MergeFiles struct {
	MergeFilesContract
	listFiles   *list.List[*helper.FileMetadata]
	inputFiles  []inputFiles
	outputFile  string
	window      *fyne.Window
	procesPopup *widget.PopUp
	processText *widget.Label
}

func NewMergeFilesComponent(window *fyne.Window, fileList *list.List[*helper.FileMetadata]) Component {
	return &MergeFiles{
		listFiles:   fileList,
		window:      window,
		processText: widget.NewLabel(""),
	}
}

func (f *MergeFiles) Content() fyne.CanvasObject {

	mergeButton := widget.NewButton(lang.L("merge"), func() {

		err := f.checkFileIntegrity()

		if err != nil {
			dialog.NewCustomConfirm("Error", "Confirmer", "Annuler", widget.NewLabel(err.Error()), func(b bool) {
				if !b {
					return
				}
				f.Merge()
			}, *f.window).Show()
			return
		} else {
			f.Merge()
		}
	})

	outputFile := widget.NewEntry()
	outputFile.SetPlaceHolder(lang.L("output_folder"))
	outputFile.OnChanged = func(s string) {
		f.outputFile = s
	}

	listFiles := widget.NewList(
		func() int {
			return len(f.inputFiles)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("", nil), customs.NewNumericalEntry(), widget.NewLabel(""))
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Check).Checked = f.inputFiles[i].enabled
			item.(*fyne.Container).Objects[0].(*widget.Check).OnChanged = func(b bool) {
				f.inputFiles[i].enabled = b
			}
			item.(*fyne.Container).Objects[1].(*customs.NumericalEntry).SetText(fmt.Sprintf("%d", f.inputFiles[i].position))
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(f.inputFiles[i].filename)
			// refresh ui for check box
			item.Refresh()
		},
	)

	refreshButton := widget.NewButton(lang.L("refresh"), func() {
		f.inputFiles = make([]inputFiles, len(f.listFiles.GetItems()))
		for i, file := range f.listFiles.GetItems() {
			if !strings.HasSuffix(file.Directory, "/") {
				file.Directory = file.Directory + "/"
			}
			file.Directory = strings.ReplaceAll(file.Directory, "\\", "/")
			f.inputFiles[i] = inputFiles{file.Directory, file.FileName, i, true, file}
		}
		listFiles.Refresh()
	})

	content := container.NewBorder(
		container.NewHBox(
			refreshButton,
		),
		container.NewBorder(
			nil,
			mergeButton,
			widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
				dialogFolder := dialog.NewFolderOpen(func(uc fyne.ListableURI, err error) {
					if err != nil {
						dialog.ShowError(err, *f.window)
						return
					}

					if uc == nil {
						return
					}

					f.outputFile = uc.Path()
					outputFile.SetText(f.outputFile)
				}, *f.window)
				size := (*f.window).Canvas().Size()
				dialogFolder.Resize(fyne.NewSize(size.Width-150, size.Height-150))
				dialogFolder.Show()
			}),
			nil, outputFile,
		),
		nil,
		nil,
		listFiles,
	)

	f.procesPopup = widget.NewModalPopUp(f.processText, (*f.window).Canvas())

	return content
}

func (f *MergeFiles) checkFileIntegrity() error {
	// check if width and height are the same
	// check if codec is the same
	var width, height int
	var codec string

	for i, file := range f.inputFiles {
		if i == 0 {
			width = file.metadata.GetVideoStreams()[0].Width
			height = file.metadata.GetVideoStreams()[0].Height
			codec = file.metadata.GetVideoStreams()[0].CodecName
			continue
		}

		if width != file.metadata.GetVideoStreams()[0].Width || height != file.metadata.GetVideoStreams()[0].Height {
			return errors.New(lang.L("width_height_not_same"))
		}

		if codec != file.metadata.GetVideoStreams()[0].CodecName {
			return errors.New(lang.L("codec_not_same"))
		}
	}

	return nil
}

func (f *MergeFiles) Merge() {

	if f.outputFile == "" && len(f.inputFiles) > 0 {
		f.outputFile = f.inputFiles[0].dir
		// normaliser le chemin
		f.outputFile = strings.ReplaceAll(f.outputFile, "\\", "/")
		if !strings.HasSuffix(f.outputFile, "/") {
			f.outputFile = f.outputFile + "/"
		}
	}

	f.processText.SetText(lang.L("merging_files") + " 0%")
	f.procesPopup.Show()

	defer f.procesPopup.Hide()

	sort.SliceStable(f.inputFiles, func(i, j int) bool {
		return f.inputFiles[i].position < f.inputFiles[j].position
	})

	var finalInputFiles []inputFiles

	for _, file := range f.inputFiles {
		if file.enabled {
			finalInputFiles = append(finalInputFiles, file)
		}
	}

	if len(finalInputFiles) == 0 {
		return
	}

	err := f.MergeFiles(finalInputFiles, f.outputFile)
	if err != nil {
		dialog.ShowError(err, *f.window)
	}

	dialog.NewConfirm("Merge succeed", "ok", func(b bool) {}, *f.window).Show()
}

func (f *MergeFiles) MergeFiles(files []inputFiles, output string) error {
	defer f.procesPopup.Hide()
	output = strings.ReplaceAll(output, "\\", "/")

	output += strings.TrimSuffix(files[0].filename, filepath.Ext(files[0].filename)) + "_merged" + filepath.Ext(files[0].filename)

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.dir+file.filename)
	}

	totalSize, err := helper.CalculateTotalSize(filenames)
	if err != nil {
		return err
	}

	txtFile, err := f.createTempFile(files, output)
	if err != nil {
		return err
	}
	defer func() {
		txtFile.Close()
		err := os.Remove(txtFile.Name())
		if err != nil {
			dialog.ShowError(err, *f.window)
		}
	}()

	normalizedTxtPath := strings.ReplaceAll(txtFile.Name(), "\\", "/")
	cmd := []string{
		"-progress",
		"pipe:1",
		"-y",
		"-f",
		"concat",
		"-safe",
		"0",
		"-i",
		normalizedTxtPath,
		"-map",
		"0",
		"-c",
		"copy",
		output,
	}

	progress := helper.NewCmd()

	command, err := progress.RunCommandWithProgress(cmd, totalSize)
	if err != nil {
		return err
	}

	for p := range progress.CProgress {
		if p.PercentComplete == 100 {
			f.processText.SetText(lang.L("merging_files") + " 100%")
			break
		}
		f.processText.SetText(fmt.Sprintf("%s %.1f%% (%s: %.1fx, Bitrate: %.1f kbits/s)",
			lang.L("merging_files"),
			p.PercentComplete,
			lang.L("speed"),
			p.Speed,
			p.Bitrate))
	}

	err = command.Wait()

	if err != nil {
		dialog.ShowError(err, *f.window)
	}

	return nil
}

func (f *MergeFiles) createTempFile(files []inputFiles, output string) (*os.File, error) {

	// get current date time
	now := time.Now()

	txtFile, err := os.Create(strings.TrimSuffix(output, filepath.Ext(output)) + now.Format("2006-01-02_15-04-05") + "_temp.txt")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// escape "'" in file path where "'" = "'\''"
		filename := strings.ReplaceAll(file.dir+file.filename, "'", "'\\''")
		_, err := txtFile.WriteString(fmt.Sprintf(`%s '%s'`+"\n", "file", filename))
		if err != nil {
			return nil, err
		}
	}

	return txtFile, nil
}
