package components

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components/customs"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type MergeFilesContract interface {
	MergeFiles(files []string, output string) error
}

type inputFiles struct {
	file     string
	position int
	enabled  bool
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
	listFiles   *list.List[fileinfo.FileInfo]
	inputFiles  []inputFiles
	outputFile  string
	window      *fyne.Window
	procesPopup *widget.PopUp
	processText *widget.Label
}

func NewMergeFilesComponent(window *fyne.Window, fileList *list.List[fileinfo.FileInfo]) *MergeFiles {
	return &MergeFiles{
		listFiles:   fileList,
		window:      window,
		processText: widget.NewLabel(""),
	}
}

func (f *MergeFiles) Content() fyne.CanvasObject {

	mergeButton := widget.NewButton(lang.L("merge"), func() {
		f.Merge()
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
			item.(*fyne.Container).Objects[1].(*customs.NumericalEntry).SetText(fmt.Sprintf("%d", f.inputFiles[i].position))
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(f.inputFiles[i].file)
			// refresh ui for check box
			item.Refresh()
		},
	)

	refreshButton := widget.NewButton(lang.L("refresh"), func() {
		f.inputFiles = []inputFiles{}
		for i, file := range f.listFiles.GetItems() {
			f.inputFiles = append(f.inputFiles, inputFiles{file.GetPath(), i, true})
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

func (f *MergeFiles) Merge() {
	f.processText.SetText(lang.L("merging_files") + " 0%")
	f.procesPopup.Show()

	defer f.procesPopup.Hide()

	var finalInputFiles []string
	sort.SliceStable(f.inputFiles, func(i, j int) bool {
		return f.inputFiles[i].position < f.inputFiles[j].position
	})
	for _, file := range f.inputFiles {
		if file.enabled {
			finalInputFiles = append(finalInputFiles, file.file)
		}
	}

	if len(finalInputFiles) == 0 {
		return
	}

	err := f.MergeFiles(finalInputFiles, f.outputFile+"output_"+finalInputFiles[0])
	if err != nil {
		dialog.ShowError(err, *f.window)
	}
}

func (f *MergeFiles) MergeFiles(files []string, output string) error {
	output = strings.ReplaceAll(output, "\\", "/")
	totalSize, err := f.calculateTotalSize(files)
	if err != nil {
		return err
	}

	txtFile, err := f.createTempFile(files, output)
	if err != nil {
		return err
	}
	defer txtFile.Close()

	normalizedTxtPath := strings.ReplaceAll(txtFile.Name(), "\\", "/")
	cmd := fmt.Sprintf("-progress pipe:1 -y -f concat -safe 0 -i %s -c copy %s", normalizedTxtPath, output)
	err = f.runCommandWithProgress(cmd, totalSize)
	if err != nil {
		return err
	}

	return os.Remove(txtFile.Name())
}

func (f *MergeFiles) createTempFile(files []string, output string) (*os.File, error) {
	txtFile, err := os.Create(output + ".txt")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		normalizedPath := strings.ReplaceAll(file, "\\", "/")
		_, err := txtFile.WriteString(fmt.Sprintf("%s '%s'\n", lang.L("file"), normalizedPath))
		if err != nil {
			return nil, err
		}
	}

	return txtFile, nil
}

func (f *MergeFiles) calculateTotalSize(files []string) (int64, error) {
	var totalSize int64
	for _, file := range files {
		// Normaliser le chemin avant d'accéder au fichier
		normalizedPath := strings.ReplaceAll(file, "\\", "/")
		fileInfo, err := os.Stat(normalizedPath)
		if err != nil {
			return 0, err
		}
		totalSize += fileInfo.Size()
	}
	return totalSize, nil
}

func (f *MergeFiles) runCommandWithProgress(cmd string, totalSize int64) error {
	path, err := exec.LookPath("FFMPEG")
	if err != nil {
		return err
	}

	if path == "" {
		return errors.New(lang.L("ffmpeg_not_found"))
	}

	// Normaliser le chemin de FFmpeg
	path = strings.ReplaceAll(path, "\\", "/")

	command := exec.Command(path, strings.Split(cmd, " ")...)

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
