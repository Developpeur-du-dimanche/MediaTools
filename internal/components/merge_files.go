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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components/customs"
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
	listFiles   *list.List[*fileinfo.FileInfo]
	inputFiles  []inputFiles
	outputFile  string
	window      fyne.Window
	procesPopup *widget.PopUp
	processText *widget.Label
}

func NewMergeFilesComponent(window *fyne.Window, fileList *list.List[*fileinfo.FileInfo]) *MergeFiles {
	return &MergeFiles{
		listFiles:   fileList,
		window:      *window,
		processText: widget.NewLabel(""),
	}
}

func (f *MergeFiles) Content() fyne.CanvasObject {

	mergeButton := widget.NewButton("Merge", func() {
		f.Merge()
	})
	outputFile := widget.NewEntry()
	outputFile.SetPlaceHolder("Output folder")
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

		},
	)

	refreshButton := widget.NewButton("Refresh", func() {
		f.inputFiles = []inputFiles{}
		for i, file := range f.listFiles.GetItems() {
			f.inputFiles = append(f.inputFiles, inputFiles{file.Path, i, true})
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
						dialog.ShowError(err, f.window)
						return
					}

					if uc == nil {
						return
					}

					f.outputFile = uc.Path()
					outputFile.SetText(f.outputFile)
				}, f.window)
				size := (f.window).Canvas().Size()
				dialogFolder.Resize(fyne.NewSize(size.Width-150, size.Height-150))
				dialogFolder.Show()
			}),
			nil, outputFile,
		),
		nil,
		nil,
		listFiles,
	)

	f.procesPopup = widget.NewModalPopUp(f.processText, f.window.Canvas())

	return content
}

func (f *MergeFiles) Merge() {
	f.processText.SetText("Merging files... 0%")
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

	err := f.MergeFiles(finalInputFiles, f.outputFile+"/output.mkv")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		dialog.ShowError(err, f.window)
	}
}

func (f *MergeFiles) MergeFiles(files []string, output string) error {
	// Normaliser le chemin de sortie
	output = strings.ReplaceAll(output, "\\", "/")

	// Calculer la taille totale des fichiers d'entrée
	totalSize, err := f.calculateTotalSize(files)
	if err != nil {
		return err
	}

	// Créer le fichier temporaire avec le chemin normalisé
	txtFile, err := os.Create(output + ".txt")
	if err != nil {
		return err
	}
	defer txtFile.Close()

	// Écrire les chemins normalisés dans le fichier temporaire
	for _, file := range files {
		normalizedPath := strings.ReplaceAll(file, "\\", "/")
		_, err := txtFile.WriteString(fmt.Sprintf("file '%s'\n", normalizedPath))
		if err != nil {
			return err
		}
	}

	// Normaliser le chemin du fichier temporaire pour la commande FFmpeg
	normalizedTxtPath := strings.ReplaceAll(txtFile.Name(), "\\", "/")

	cmd := fmt.Sprintf("-progress pipe:1 -y -f concat -safe 0 -i %s -c copy %s",
		normalizedTxtPath,
		output)
	err = f.runCommandWithProgress(cmd, totalSize)
	if err != nil {
		return err
	}

	return os.Remove(txtFile.Name())
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
		return errors.New("ffmpeg not found")
	}

	// Normaliser le chemin de FFmpeg
	path = strings.ReplaceAll(path, "\\", "/")

	command := exec.Command(path, strings.Split(cmd, " ")...)

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

				f.processText.SetText(fmt.Sprintf("Merging files... %.1f%% (Speed: %.1fx, Bitrate: %.1f kbits/s)",
					percentComplete,
					progress.speed,
					progress.bitrate))
			}

			lines = []string{}
		} else if line == "progress=end" {
			f.processText.SetText("Merging files... 100%")
			break
		} else {
			lines = append(lines, line)
		}
	}

	return command.Wait()
}

func (f *MergeFiles) parseProgressBlock(lines []string, progress *ffmpegProgress) {
	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]

		switch key {
		case "bitrate":
			progress.bitrate, _ = strconv.ParseFloat(strings.TrimSuffix(value, "kbits/s"), 64)
		case "total_size":
			progress.totalSize, _ = strconv.ParseInt(value, 10, 64)
		case "out_time_us":
			progress.outTimeUs, _ = strconv.ParseInt(value, 10, 64)
		case "out_time_ms":
			progress.outTimeMs, _ = strconv.ParseInt(value, 10, 64)
		case "out_time":
			progress.outTime = value
		case "dup_frames":
			progress.dupFrames, _ = strconv.Atoi(value)
		case "drop_frames":
			progress.dropFrames, _ = strconv.Atoi(value)
		case "speed":
			progress.speed, _ = strconv.ParseFloat(strings.TrimSuffix(value, "x"), 64)
		case "progress":
			progress.progress = value
		}
	}
}