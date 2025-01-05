package helper

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

type FfmpegProgress struct {
	Bitrate         float64
	TotalSize       int64
	OutTimeUs       int64
	OutTimeMs       int64
	OutTime         string
	DupFrames       int
	DropFrames      int
	Speed           float64
	Progress        string
	PercentComplete float64
	CProgress       chan *FfmpegProgress
}

func CalculateTotalSize(files []string) (int64, error) {
	var totalSize int64
	for _, file := range files {
		// Normaliser le chemin avant d'accéder au fichier
		fileInfo, err := os.Stat(file)
		if err != nil {
			return 0, err
		}
		totalSize += fileInfo.Size()
	}
	return totalSize, nil
}

func NewCmd() *FfmpegProgress {
	return &FfmpegProgress{
		CProgress: make(chan *FfmpegProgress),
	}
}

func (f *FfmpegProgress) RunCommandWithProgress(cmd []string, totalSize int64) (*exec.Cmd, error) {
	path := fyne.CurrentApp().Preferences().String("ffmpeg")
	if path == "" {
		return nil, errors.New(lang.L("ffmpeg_not_found"))
	}

	// Normaliser le chemin de FFmpeg
	path = strings.ReplaceAll(path, "\\", "/")

	command := exec.Command(path, cmd...)

	// print command
	fmt.Printf("Running command: %s %s\n", path, strings.Join(command.Args, " "))

	RunCmdBackground(command)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	command.Stderr = os.Stderr

	if err := command.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)

	var lines []string

	go func() {
		for scanner.Scan() {
			line := scanner.Text()

			if line == "progress=continue" {
				parseProgressBlock(lines, f)

				if totalSize > 0 {
					percentComplete := float64(f.TotalSize) / float64(totalSize) * 100

					f.PercentComplete = percentComplete

					f.CProgress <- f

				}

				lines = []string{}
			} else if line == "progress=end" {
				f.CProgress <- &FfmpegProgress{Progress: "end", PercentComplete: 100}
				// close channel
				close(f.CProgress)
				break
			} else {
				lines = append(lines, line)
			}
		}
	}()

	return command, nil
}

func parseProgressBlock(lines []string, progress *FfmpegProgress) {
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
			progress.Bitrate, _ = strconv.ParseFloat(value, 64)
		case "total_size":
			value = extractNumber(value)
			progress.TotalSize, _ = strconv.ParseInt(value, 10, 64)
		case "out_time_us":
			value = extractNumber(value)
			progress.OutTimeUs, _ = strconv.ParseInt(value, 10, 64)
		case "out_time_ms":
			value = extractNumber(value)
			progress.OutTimeMs, _ = strconv.ParseInt(value, 10, 64)
		case "out_time":
			progress.OutTime = value
		case "dup_frames":
			value = extractNumber(value)
			progress.DupFrames, _ = strconv.Atoi(value)
		case "drop_frames":
			value = extractNumber(value)
			progress.DropFrames, _ = strconv.Atoi(value)
		case "speed":
			value = extractNumber(value)
			progress.Speed, _ = strconv.ParseFloat(value, 64)
		case "progress":
			progress.Progress = value
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
