package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Developpeur-du-dimanche/MediaTools/internal/utils"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

// ScanProgress represents the progress of a folder scan
type ScanProgress struct {
	CurrentFile  string
	FilesScanned int
	TotalFiles   int
	IsComplete   bool
}

// MediaService handles all media-related business logic
type MediaService struct {
	validExtensions []string
	probeTimeout    time.Duration
}

// NewMediaService creates a new media service instance
func NewMediaService(validExtensions []string, probeTimeout time.Duration) *MediaService {
	if probeTimeout == 0 {
		probeTimeout = 10 * time.Second
	}
	return &MediaService{
		validExtensions: validExtensions,
		probeTimeout:    probeTimeout,
	}
}

// GetMediaInfo analyzes a media file and returns its information
func (ms *MediaService) GetMediaInfo(ctx context.Context, filePath string) (*medias.FfprobeResult, error) {
	// Validate file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Errorf("File does not exist: %s", filePath)
		return nil, fmt.Errorf("%w: %s", ErrInvalidPath, filePath)
	}

	// Validate extension
	if !utils.IsValidExtensions(filePath, ms.validExtensions) {
		logger.Debugf("Invalid extension for file: %s", filePath)
		return nil, fmt.Errorf("%w: %s", ErrInvalidExtension, filepath.Ext(filePath))
	}

	// Create timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, ms.probeTimeout)
	defer cancel()

	// Probe the file
	ffprobe := medias.NewFfprobe(filePath,
		medias.FFPROBE_LOGLEVEL_FATAL,
		medias.PRINT_FORMAT_JSON,
		medias.SHOW_FORMAT,
		medias.SHOW_STREAMS,
		medias.EXPERIMENTAL,
	)

	result, err := ffprobe.Probe(timeoutCtx)
	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			logger.Warnf("Ffprobe timeout for file: %s", filePath)
			return nil, fmt.Errorf("%w: %s", ErrProbeTimeout, filePath)
		}
		logger.Errorf("Ffprobe failed for file %s: %v", filePath, err)
		return nil, fmt.Errorf("%w: %v", ErrProbeFailed, err)
	}

	logger.Infof("Successfully analyzed file: %s", filePath)
	return result, nil
}

// ScanFolder scans a folder recursively and sends progress updates
func (ms *MediaService) ScanFolder(ctx context.Context, folderPath string, progressChan chan<- ScanProgress) ([]string, error) {
	logger.Infof("Starting folder scan: %s", folderPath)

	// Validate folder exists
	if info, err := os.Stat(folderPath); os.IsNotExist(err) || !info.IsDir() {
		logger.Errorf("Invalid folder path: %s", folderPath)
		return nil, fmt.Errorf("%w: %s", ErrInvalidPath, folderPath)
	}

	var mediaFiles []string
	filesScanned := 0

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		// Check for cancellation
		select {
		case <-ctx.Done():
			logger.Info("Folder scan cancelled")
			return ctx.Err()
		default:
		}

		if err != nil {
			logger.Warnf("Error accessing path %s: %v", path, err)
			return nil // Skip this file
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip files without extension
		if filepath.Ext(path) == "" {
			return nil
		}

		filesScanned++

		// Check if valid media file
		if utils.IsValidExtensions(path, ms.validExtensions) {
			mediaFiles = append(mediaFiles, path)
			logger.Debugf("Found media file: %s", path)
		}

		// Send progress update
		if progressChan != nil {
			progressChan <- ScanProgress{
				CurrentFile:  path,
				FilesScanned: filesScanned,
				IsComplete:   false,
			}
		}

		return nil
	})

	// Send completion update
	if progressChan != nil {
		progressChan <- ScanProgress{
			FilesScanned: filesScanned,
			TotalFiles:   len(mediaFiles),
			IsComplete:   true,
		}
	}

	if err != nil {
		logger.Errorf("Folder scan failed: %v", err)
		return mediaFiles, err
	}

	logger.Infof("Folder scan complete. Found %d media files", len(mediaFiles))
	return mediaFiles, nil
}

// SetValidExtensions updates the list of valid extensions
func (ms *MediaService) SetValidExtensions(extensions []string) {
	ms.validExtensions = extensions
	logger.Infof("Updated valid extensions: %v", extensions)
}

// GetValidExtensions returns the current list of valid extensions
func (ms *MediaService) GetValidExtensions() []string {
	return ms.validExtensions
}
