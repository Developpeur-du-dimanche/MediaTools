package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

// FFmpegService handles FFmpeg operations
type FFmpegService struct {
	ffmpegPath string
}

// NewFFmpegService creates a new FFmpeg service
func NewFFmpegService() *FFmpegService {
	return &FFmpegService{
		ffmpegPath: "ffmpeg", // Assume ffmpeg is in PATH
	}
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

	// This requires mapping streams and is more complex
	// For now, we'll implement a basic version
	var streamSelector string
	switch strings.ToLower(streamType) {
	case "audio":
		streamSelector = "a"
	case "subtitle":
		streamSelector = "s"
	default:
		return fmt.Errorf("unsupported stream type for language removal: %s", streamType)
	}

	// Build mapping to exclude streams with specific language
	args := []string{
		"-i", inputFile,
		"-map", "0:v", // Keep all video
		"-map", fmt.Sprintf("0:%s:m:language:%s?", streamSelector, language), // Map streams NOT matching language
		"-map_metadata", "0",
		"-c", "copy",
		outputPath,
		"-y",
	}

	cmd := exec.CommandContext(ctx, fs.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg language removal failed: %w\nOutput: %s", err, string(output))
	}

	if progress != nil {
		progress(1.0, fmt.Sprintf("Removed %s streams with language %s", streamType, language))
	}

	logger.Infof("Successfully removed %s streams with language %s", streamType, language)
	return nil
}

// RemoveStreamsByCodec removes streams matching a specific codec
func (fs *FFmpegService) RemoveStreamsByCodec(ctx context.Context, inputFile, outputPath, streamType, codec string, progress ProgressCallback) error {
	logger.Infof("Removing %s streams with codec %s from %s", streamType, codec, inputFile)

	// For this operation, we need to analyze streams first and build a custom map
	// This is a simplified version - a full implementation would need stream inspection

	var streamSelector string
	switch strings.ToLower(streamType) {
	case "audio":
		streamSelector = "a"
	case "subtitle":
		streamSelector = "s"
	case "video":
		streamSelector = "v"
	default:
		return fmt.Errorf("unsupported stream type: %s", streamType)
	}

	args := []string{
		"-i", inputFile,
		"-map", "0",
		"-map", fmt.Sprintf("-0:%s:codec:%s", streamSelector, codec), // Exclude streams with this codec
		"-c", "copy",
		outputPath,
		"-y",
	}

	cmd := exec.CommandContext(ctx, fs.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg codec removal failed: %w\nOutput: %s", err, string(output))
	}

	if progress != nil {
		progress(1.0, fmt.Sprintf("Removed %s streams with codec %s", streamType, codec))
	}

	logger.Infof("Successfully removed %s streams with codec %s", streamType, codec)
	return nil
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
