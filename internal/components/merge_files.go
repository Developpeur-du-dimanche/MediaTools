package components

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
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

type ffmpegProgress struct {
	bitrate    float64
	totalSize  int64
	outTimeUs  int64
	outTimeMs  int64
	outTime    string
	dupFrames  int
	dropFrames int
	speed      float64
	progress   string
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
}

func (f *MergeFiles) MergeFiles(files []inputFiles, output string) error {
	output = strings.ReplaceAll(output, "\\", "/")

	output += strings.TrimSuffix(files[0].filename, filepath.Ext(files[0].filename)) + "_merged" + filepath.Ext(files[0].filename)

	totalSize, err := f.calculateTotalSize(files)
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
		"-c",
		"copy",
		output,
	}
	err = f.runCommandWithProgress(cmd, totalSize)
	if err != nil {
		return err
	}

	return os.Remove(txtFile.Name())
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

func (f *MergeFiles) calculateTotalSize(files []inputFiles) (int64, error) {
	var totalSize int64
	for _, file := range files {
		// Normaliser le chemin avant d'accéder au fichier
		fileInfo, err := os.Stat(file.dir + file.filename)
		if err != nil {
			return 0, err
		}
		totalSize += fileInfo.Size()
	}
	return totalSize, nil
}

func (f *MergeFiles) runCommandWithProgress(cmd []string, totalSize int64) error {
	path := fyne.CurrentApp().Preferences().String("ffmpeg")
	if path == "" {
		return errors.New(lang.L("ffmpeg_not_found"))
	}

	// Normaliser le chemin de FFmpeg
	path = strings.ReplaceAll(path, "\\", "/")

	command := exec.Command(path, cmd...)

	// print command
	fmt.Printf("Running command: %s %s\n", path, strings.Join(command.Args, " "))

	helper.RunCmdBackground(command)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	command.Stderr = os.Stderr

	if err := command.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	progress := &ffmpegProgress{}

	var lines []string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "progress=continue" {
			f.parseProgressBlock(lines, progress)

			if totalSize > 0 {
				percentComplete := float64(progress.totalSize) / float64(totalSize) * 100

				f.processText.SetText(fmt.Sprintf("%s %.1f%% (%s: %.1fx, Bitrate: %.1f kbits/s)",
					lang.L("merging_files"),
					percentComplete,
					lang.L("speed"),
					progress.speed,
					progress.bitrate))
			}

			lines = []string{}
		} else if line == "progress=end" {
			f.processText.SetText(lang.L("merging_files") + " 100%")
			break
		} else {
			lines = append(lines, line)
		}
	}

	return command.Wait()
}

func (f *MergeFiles) parseProgressBlock(lines []string, progress *ffmpegProgress) {
	for _, line := range lines {
		// Diviser sur le premier "=" seulement
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "bitrate":
			// Extraire uniquement les chiffres et le point décimal
			value = extractNumber(value)
			progress.bitrate, _ = strconv.ParseFloat(value, 64)
		case "total_size":
			value = extractNumber(value)
			progress.totalSize, _ = strconv.ParseInt(value, 10, 64)
		case "out_time_us":
			value = extractNumber(value)
			progress.outTimeUs, _ = strconv.ParseInt(value, 10, 64)
		case "out_time_ms":
			value = extractNumber(value)
			progress.outTimeMs, _ = strconv.ParseInt(value, 10, 64)
		case "out_time":
			progress.outTime = value
		case "dup_frames":
			value = extractNumber(value)
			progress.dupFrames, _ = strconv.Atoi(value)
		case "drop_frames":
			value = extractNumber(value)
			progress.dropFrames, _ = strconv.Atoi(value)
		case "speed":
			value = extractNumber(value)
			progress.speed, _ = strconv.ParseFloat(value, 64)
		case "progress":
			progress.progress = value
		}
	}
}

// extractNumber extrait uniquement les chiffres et le point décimal d'une chaîne
func extractNumber(s string) string {
	var result strings.Builder
	hasDecimal := false

	// Gérer le signe négatif au début
	if strings.HasPrefix(strings.TrimSpace(s), "-") {
		result.WriteRune('-')
		s = strings.TrimPrefix(strings.TrimSpace(s), "-")
	}

	for _, c := range s {
		if c >= '0' && c <= '9' {
			result.WriteRune(c)
		} else if c == '.' && !hasDecimal {
			result.WriteRune(c)
			hasDecimal = true
		}
	}

	return result.String()
}
