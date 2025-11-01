package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

// FFmpegService handles FFmpeg operations
type FFmpegService struct {
	ffmpegPath string
}

func (fs *FFmpegService) LocateFFmpeg() (string, error) {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", err
	}
	return path, nil
}

// NewFFmpegService creates a new FFmpeg service
func NewFFmpegService() *FFmpegService {
	return &FFmpegService{
		ffmpegPath: "ffmpeg", // Assume ffmpeg is in PATH
	}
}

// SetFFmpegPath sets a custom path for ffmpeg executable
func (fs *FFmpegService) SetFFmpegPath(path string) {
	if path != "" {
		fs.ffmpegPath = path
	}
}

// GetFFmpegPath returns the current ffmpeg path
func (fs *FFmpegService) GetFFmpegPath() string {
	return fs.ffmpegPath
}

// ProgressCallback is called during FFmpeg operations
type ProgressCallback func(progress float64, message string)

// MergeVideos concatenates multiple video files into one
func (fs *FFmpegService) MergeVideos(ctx context.Context, inputFiles []string, outputPath string, progress ProgressCallback) error {
	if len(inputFiles) == 0 {
		return fmt.Errorf("no input files provided")
	}

	logger.Infof("Merging %d videos into %s", len(inputFiles), outputPath)

	// Create a temporary file list for FFmpeg concat
	listFile, err := fs.createConcatList(inputFiles)
	if err != nil {
		return fmt.Errorf("failed to create concat list: %w", err)
	}
	defer os.Remove(listFile)

	// Build FFmpeg command
	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy", // Copy streams without re-encoding
		outputPath,
		"-y", // Overwrite output file
	}

	cmd := exec.CommandContext(ctx, fs.ffmpegPath, args...)

	// Capture output for progress tracking
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg merge failed: %w\nOutput: %s", err, string(output))
	}

	if progress != nil {
		progress(1.0, "Merge complete")
	}

	logger.Infof("Successfully merged videos to %s", outputPath)
	return nil
}

// RemoveStreamsByType removes all streams of a specific type from a video
func (fs *FFmpegService) RemoveStreamsByType(ctx context.Context, inputFile, outputPath, streamType string, progress ProgressCallback) error {
	logger.Infof("Removing %s streams from %s", streamType, inputFile)

	var args []string
	switch strings.ToLower(streamType) {
	case "audio":
		args = []string{
			"-i", inputFile,
			"-c:v", "copy",
			"-an", // Remove all audio
			outputPath,
			"-y",
		}
	case "subtitle":
		args = []string{
			"-i", inputFile,
			"-c:v", "copy",
			"-c:a", "copy",
			"-sn", // Remove all subtitles
			outputPath,
			"-y",
		}
	case "video":
		args = []string{
			"-i", inputFile,
			"-vn", // Remove all video
			"-c:a", "copy",
			outputPath,
			"-y",
		}
	case "metadata":
		args = []string{
			"-i", inputFile,
			"-map", "0",
			"-map_metadata", "-1", // Remove all metadata
			"-c", "copy",
			outputPath,
			"-y",
		}
	case "attachments":
		args = []string{
			"-i", inputFile,
			"-map", "0",
			"-map", "-0:t", // Remove all attachments
			"-c", "copy",
			outputPath,
			"-y",
		}
	default:
		return fmt.Errorf("unsupported stream type: %s", streamType)
	}

	cmd := exec.CommandContext(ctx, fs.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg stream removal failed: %w\nOutput: %s", err, string(output))
	}

	if progress != nil {
		progress(1.0, fmt.Sprintf("Removed %s streams", streamType))
	}

	logger.Infof("Successfully removed %s streams", streamType)
	return nil
}

// RemoveStreamsByLanguage removes streams matching a specific language
func (fs *FFmpegService) RemoveStreamsByLanguage(ctx context.Context, inputFile, outputPath, streamType, language string, progress ProgressCallback) error {
	logger.Infof("Removing %s streams with language %s from %s", streamType, language, inputFile)

	probeResult, err := fs.probeFile(ctx, inputFile)
	if err != nil {
		return err
	}

	if !fs.hasStreamWithLanguage(probeResult, streamType, language) {
		return fs.copyFile(ctx, inputFile, outputPath, progress, "No streams removed; file copied")
	}

	args := fs.buildRemoveByLanguageArgs(inputFile, outputPath, streamType, language, probeResult)

	if err := fs.runFFmpeg(ctx, args); err != nil {
		return fmt.Errorf("ffmpeg language removal failed: %w", err)
	}

	if progress != nil {
		progress(1.0, fmt.Sprintf("Removed %s streams with language %s", streamType, language))
	}

	logger.Infof("Successfully removed %s streams with language %s", streamType, language)
	return nil
}

// probeFile probes a video file and returns the result
func (fs *FFmpegService) probeFile(ctx context.Context, inputFile string) (*medias.FfprobeResult, error) {
	ffprobeData := medias.NewFfprobe(inputFile,
		medias.FFPROBE_LOGLEVEL_FATAL,
		medias.PRINT_FORMAT_JSON,
		medias.SHOW_FORMAT,
		medias.SHOW_STREAMS,
		medias.EXPERIMENTAL,
	)
	probeResult, err := ffprobeData.Probe(ctx)
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}
	return probeResult, nil
}

// hasStreamWithLanguage checks if the file has streams of the specified type and language
func (fs *FFmpegService) hasStreamWithLanguage(probeResult *medias.FfprobeResult, streamType, language string) bool {
	switch strings.ToLower(streamType) {
	case "audio":
		for _, stream := range probeResult.Audios {
			if strings.EqualFold(stream.Language, language) {
				return true
			}
		}
	case "subtitle":
		for _, stream := range probeResult.Subtitles {
			if strings.EqualFold(stream.Language, language) {
				return true
			}
		}
	}
	return false
}

// copyFile copies a file without modification
func (fs *FFmpegService) copyFile(ctx context.Context, inputFile, outputPath string, progress ProgressCallback, message string) error {
	args := []string{
		"-i", inputFile,
		"-c", "copy",
		outputPath,
		"-y",
	}
	if err := fs.runFFmpeg(ctx, args); err != nil {
		return fmt.Errorf("ffmpeg copy failed: %w", err)
	}
	if progress != nil {
		progress(1.0, message)
	}
	return nil
}

// buildRemoveByLanguageArgs builds FFmpeg arguments for removing streams by language
func (fs *FFmpegService) buildRemoveByLanguageArgs(inputFile, outputPath, streamType, language string, probeResult *medias.FfprobeResult) []string {
	args := []string{
		"-i", inputFile,
		"-map", "0:v", // Keep all video
	}

	switch strings.ToLower(streamType) {
	case "audio":
		for i, stream := range probeResult.Audios {
			if !strings.EqualFold(stream.Language, language) {
				args = append(args, "-map", fmt.Sprintf("0:a:%d", i))
			}
		}
	case "subtitle":
		for i, stream := range probeResult.Subtitles {
			if !strings.EqualFold(stream.Language, language) {
				args = append(args, "-map", fmt.Sprintf("0:s:%d", i))
			}
		}
	}

	args = append(args,
		"-map_metadata", "0",
		"-c", "copy",
		outputPath,
		"-y",
	)

	return args
}

// runFFmpeg runs FFmpeg with the specified arguments
func (fs *FFmpegService) runFFmpeg(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, fs.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w\nOutput: %s", err, string(output))
	}
	return nil
}

// RemoveStreamsByCodec removes streams matching a specific codec
func (fs *FFmpegService) RemoveStreamsByCodec(ctx context.Context, inputFile, outputPath, streamType, codec string, progress ProgressCallback) error {
	logger.Infof("Removing %s streams with codec %s from %s", streamType, codec, inputFile)

	probeResult, err := fs.probeFile(ctx, inputFile)
	if err != nil {
		return err
	}

	if !fs.hasStreamWithCodec(probeResult, streamType, codec) {
		logger.Infof("No %s streams with codec %s found in %s; skipping removal", streamType, codec, inputFile)
		return fs.copyFile(ctx, inputFile, outputPath, progress, "No streams removed; file copied")
	}

	args := fs.buildRemoveByCodecArgs(inputFile, outputPath, streamType, codec, probeResult)

	if err := fs.runFFmpeg(ctx, args); err != nil {
		return fmt.Errorf("ffmpeg codec removal failed: %w", err)
	}

	if progress != nil {
		progress(1.0, fmt.Sprintf("Removed %s streams with codec %s", streamType, codec))
	}

	logger.Infof("Successfully removed %s streams with codec %s", streamType, codec)
	return nil
}

// hasStreamWithCodec checks if the file has streams of the specified type and codec
func (fs *FFmpegService) hasStreamWithCodec(probeResult *medias.FfprobeResult, streamType, codec string) bool {
	switch strings.ToLower(streamType) {
	case "audio":
		for _, stream := range probeResult.Audios {
			if strings.EqualFold(stream.CodecName, codec) {
				return true
			}
		}
	case "subtitle":
		for _, stream := range probeResult.Subtitles {
			if strings.EqualFold(stream.CodecName, codec) {
				return true
			}
		}
	case "video":
		for _, stream := range probeResult.Videos {
			if strings.EqualFold(stream.CodecName, codec) {
				return true
			}
		}
	}
	return false
}

// buildRemoveByCodecArgs builds FFmpeg arguments for removing streams by codec
func (fs *FFmpegService) buildRemoveByCodecArgs(inputFile, outputPath, streamType, codec string, probeResult *medias.FfprobeResult) []string {
	args := []string{
		"-i", inputFile,
	}

	switch strings.ToLower(streamType) {
	case "audio":
		// Map video streams
		args = append(args, "-map", "0:v")
		// Map audio streams excluding the codec
		for i, stream := range probeResult.Audios {
			if !strings.EqualFold(stream.CodecName, codec) {
				args = append(args, "-map", fmt.Sprintf("0:a:%d", i))
			}
		}
		// Map subtitle streams
		args = append(args, "-map", "0:s?")
	case "subtitle":
		// Map video and audio streams
		args = append(args, "-map", "0:v", "-map", "0:a")
		// Map subtitle streams excluding the codec
		for i, stream := range probeResult.Subtitles {
			if !strings.EqualFold(stream.CodecName, codec) {
				args = append(args, "-map", fmt.Sprintf("0:s:%d", i))
			}
		}
	case "video":
		// Map video streams excluding the codec
		for i, stream := range probeResult.Videos {
			if !strings.EqualFold(stream.CodecName, codec) {
				args = append(args, "-map", fmt.Sprintf("0:v:%d", i))
			}
		}
		// Map audio and subtitle streams
		args = append(args, "-map", "0:a", "-map", "0:s?")
	}

	args = append(args,
		"-map_metadata", "0",
		"-c", "copy",
		outputPath,
		"-y",
	)

	return args
}

// KeepOnlyStreamsByLanguage keeps only streams matching a specific language
func (fs *FFmpegService) KeepOnlyStreamsByLanguage(ctx context.Context, inputFile, outputPath, streamType, language string, progress ProgressCallback) error {
	logger.Infof("Keeping only %s streams with language %s from %s", streamType, language, inputFile)

	var streamSelector string
	switch strings.ToLower(streamType) {
	case "audio":
		streamSelector = "a"
	case "subtitle":
		streamSelector = "s"
	default:
		return fmt.Errorf("unsupported stream type for language keeping: %s", streamType)
	}

	args := []string{
		"-i", inputFile,
		"-map", "0:v", // Keep all video
		"-map", fmt.Sprintf("0:%s:m:language:%s", streamSelector, language), // Keep only streams matching language
		"-map_metadata", "0",
		"-c", "copy",
		outputPath,
		"-y",
	}

	cmd := exec.CommandContext(ctx, fs.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg language keeping failed: %w\nOutput: %s", err, string(output))
	}

	if progress != nil {
		progress(1.0, fmt.Sprintf("Kept only %s streams with language %s", streamType, language))
	}

	logger.Infof("Successfully kept only %s streams with language %s", streamType, language)
	return nil
}

// BatchRemoveStreams applies stream removal to multiple files
func (fs *FFmpegService) BatchRemoveStreams(ctx context.Context, files []*medias.FfprobeResult, operation string, criteria map[string]string, outputDir string, progress ProgressCallback) ([]string, error) {
	results := make([]string, 0, len(files))

	for i, file := range files {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		// Generate output path
		inputPath := file.Format.Filename
		outputPath := filepath.Join(outputDir, fmt.Sprintf("processed_%s", filepath.Base(inputPath)))

		var err error
		switch operation {
		case "remove_by_type":
			err = fs.RemoveStreamsByType(ctx, inputPath, outputPath, criteria["type"], nil)
		case "remove_by_language":
			err = fs.RemoveStreamsByLanguage(ctx, inputPath, outputPath, criteria["type"], criteria["language"], nil)
		case "remove_by_codec":
			err = fs.RemoveStreamsByCodec(ctx, inputPath, outputPath, criteria["type"], criteria["codec"], nil)
		case "keep_language":
			err = fs.KeepOnlyStreamsByLanguage(ctx, inputPath, outputPath, criteria["type"], criteria["language"], nil)
		default:
			err = fmt.Errorf("unknown operation: %s", operation)
		}

		if err != nil {
			logger.Warnf("Failed to process %s: %v", inputPath, err)
			continue
		}

		results = append(results, outputPath)

		if progress != nil {
			progress(float64(i+1)/float64(len(files)), fmt.Sprintf("Processed %d/%d files", i+1, len(files)))
		}
	}

	return results, nil
}

// VideoCheckResult contains the result of a video integrity check
type VideoCheckResult struct {
	FilePath  string
	IsValid   bool
	Error     string
	Duration  float64
	HasErrors bool
}

// CheckVideoIntegrity checks if a video file is corrupted
func (fs *FFmpegService) CheckVideoIntegrity(ctx context.Context, inputFile string, progress ProgressCallback) (*VideoCheckResult, error) {
	logger.Infof("Checking video integrity for %s", inputFile)

	result := &VideoCheckResult{
		FilePath: inputFile,
		IsValid:  true,
	}

	// First, get video duration using ffprobe for progress calculation
	duration, err := fs.getVideoDuration(inputFile)
	if err != nil {
		logger.Warnf("Could not get duration for %s: %v", inputFile, err)
		duration = 0
	}
	result.Duration = duration

	// Use ffmpeg to decode the entire video and check for errors with progress
	args := []string{
		"-progress", "pipe:2", // Send progress to stderr
		"-i", inputFile,
		"-f", "null", // No output file
		"-",
	}

	cmd := exec.CommandContext(ctx, fs.ffmpegPath, args...)

	// Capture stderr for progress and errors
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	errorOutput := fs.captureProgress(stderr, duration, progress)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		// Check if there were actual errors in the output
		if strings.Contains(errorOutput, "error") || strings.Contains(errorOutput, "Error") {
			result.IsValid = false
			result.HasErrors = true
			result.Error = errorOutput
		}
	}

	// Check for errors in output
	if strings.Contains(errorOutput, "error") || strings.Contains(errorOutput, "Error") {
		result.IsValid = false
		result.HasErrors = true
		result.Error = errorOutput
	}

	if progress != nil {
		if result.IsValid {
			progress(1.0, "Video is valid")
		} else {
			progress(1.0, "Video has errors")
		}
	}

	logger.Infof("Video check complete for %s: valid=%v", inputFile, result.IsValid)
	return result, nil
}

func (fs *FFmpegService) captureProgress(stderrPipe io.ReadCloser, duration float64, progress ProgressCallback) string {
	// Parse progress from stderr
	errorOutput := ""
	buffer := make([]byte, 4096)
	for {
		n, err := stderrPipe.Read(buffer)
		if n > 0 {
			output := string(buffer[:n])
			errorOutput += output

			// Parse progress information
			if progress != nil && duration > 0 {
				// Look for "out_time_ms=" to get current position
				if strings.Contains(output, "out_time_ms=") {
					lines := strings.Split(output, "\n")
					for _, line := range lines {
						if timeStr, found := strings.CutPrefix(line, "out_time_ms="); found {
							timeMicros, parseErr := strconv.ParseInt(strings.TrimSpace(timeStr), 10, 64)
							if parseErr == nil {
								currentTime := float64(timeMicros) / 1000000.0 // Convert to seconds
								progressPercent := currentTime / duration
								if progressPercent > 1.0 {
									progressPercent = 1.0
								}
								progress(progressPercent, fmt.Sprintf("Checking... %.1f%%", progressPercent*100))
							}
						}
					}
				}
			}
		}
		if err != nil {
			break
		}
	}

	return errorOutput
}

// getVideoDuration gets the duration of a video file using ffprobe
func (fs *FFmpegService) getVideoDuration(inputFile string) (float64, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputFile,
	}

	cmd := exec.Command("ffprobe", args...)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

// BatchCheckVideos checks multiple video files for corruption
func (fs *FFmpegService) BatchCheckVideos(ctx context.Context, files []*medias.FfprobeResult, progress ProgressCallback) ([]*VideoCheckResult, error) {
	results := make([]*VideoCheckResult, 0, len(files))

	for i, file := range files {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}

		inputPath := file.Format.Filename

		// Create a progress callback for individual file
		fileProgress := func(fileProgressPercent float64, message string) {
			if progress != nil {
				// Calculate overall progress: (completed files + current file progress) / total files
				overallProgress := (float64(i) + fileProgressPercent) / float64(len(files))
				progress(overallProgress, fmt.Sprintf("[%d/%d] %s: %s", i+1, len(files), filepath.Base(inputPath), message))
			}
		}

		result, err := fs.CheckVideoIntegrity(ctx, inputPath, fileProgress)
		if err != nil {
			logger.Warnf("Failed to check %s: %v", inputPath, err)
			result = &VideoCheckResult{
				FilePath:  inputPath,
				IsValid:   false,
				HasErrors: true,
				Error:     err.Error(),
			}
		}

		results = append(results, result)

		if progress != nil {
			status := "✓ OK"
			if !result.IsValid {
				status = "✗ CORRUPTED"
			}
			progress(float64(i+1)/float64(len(files)), fmt.Sprintf("Checked %d/%d files - %s: %s", i+1, len(files), filepath.Base(inputPath), status))
		}
	}

	return results, nil
}

// createConcatList creates a temporary file list for FFmpeg concat demuxer
func (fs *FFmpegService) createConcatList(files []string) (string, error) {
	tmpFile, err := os.CreateTemp("", "ffmpeg_concat_*.txt")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	for _, file := range files {
		// Escape single quotes and wrap in quotes
		escapedPath := strings.ReplaceAll(file, "'", "'\\''")
		_, err := fmt.Fprintf(tmpFile, "file '%s'\n", escapedPath)
		if err != nil {
			return "", err
		}
	}

	return tmpFile.Name(), nil
}
