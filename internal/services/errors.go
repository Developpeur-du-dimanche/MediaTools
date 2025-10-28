package services

import "errors"

var (
	// ErrInvalidPath is returned when a file or folder path is invalid
	ErrInvalidPath = errors.New("invalid file or folder path")

	// ErrInvalidExtension is returned when a file has an unsupported extension
	ErrInvalidExtension = errors.New("unsupported file extension")

	// ErrProbeTimeout is returned when ffprobe times out
	ErrProbeTimeout = errors.New("ffprobe timeout")

	// ErrProbeFailed is returned when ffprobe fails to analyze a file
	ErrProbeFailed = errors.New("ffprobe analysis failed")

	// ErrNoMediaStreams is returned when no media streams are found
	ErrNoMediaStreams = errors.New("no media streams found in file")
)
