# MediaTools

A powerful GUI application for managing and manipulating video files with advanced filtering, merging, and stream management capabilities.

![License](https://img.shields.io/github/license/Developpeur-du-dimanche/MediaTools)
![Go Version](https://img.shields.io/github/go-mod/go-version/Developpeur-du-dimanche/MediaTools)
![Release](https://img.shields.io/github/v/release/Developpeur-du-dimanche/MediaTools)

## Features

- **Bulk Video Scanning**: Recursively scan folders to analyze video files
- **Advanced Filtering**: Filter videos by codec, bitrate, resolution, duration, language, and more
- **Video Merging**: Merge multiple videos into a single file
- **Stream Management**: Remove or keep specific audio, video, or subtitle streams
- **Video Integrity Check**: Verify video file integrity
- **FFmpeg Integration**: Leverages FFmpeg for all media operations
- **Localization**: Supports multiple languages (English, French)
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Download

Download the latest version from the [releases page](https://github.com/Developpeur-du-dimanche/MediaTools/releases).

## Installation from Source

### Prerequisites

- **Go 1.19+**: Install from [golang.org](https://golang.org/doc/install)
- **FFmpeg**: Required for media operations
  - Windows: Download from [ffmpeg.org](https://ffmpeg.org/download.html)
  - macOS: `brew install ffmpeg`
  - Linux: `sudo apt install ffmpeg` (Ubuntu/Debian) or equivalent

### Build

1. Clone the repository:
```bash
git clone https://github.com/Developpeur-du-dimanche/MediaTools.git
cd MediaTools
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
# Windows
go build -o mediatools.exe ./cmd/mediatools/main.go

# macOS/Linux
go build -o mediatools ./cmd/mediatools/main.go
```

**Note**: The first build may take several minutes as it compiles the GUI dependencies.

### Run

```bash
# Windows
./mediatools.exe

# macOS/Linux
./mediatools
```

## Development

### Hot Reload with Air (Optional)

For faster development, use [Air](https://github.com/air-verse/air) to automatically rebuild on file changes.

Install Air:
```bash
go install github.com/air-verse/air@latest
```

Run with Air:
```bash
air
```

### Project Structure

```
MediaTools/
├── cmd/mediatools/       # Application entry point
├── internal/
│   ├── components/       # UI components
│   ├── filters/          # Filter implementations
│   ├── mediatools/       # Main application logic
│   ├── services/         # Business logic services
│   └── utils/            # Utility functions
├── pkg/
│   ├── logger/           # Logging package
│   └── medias/           # Media handling (ffprobe integration)
└── readme.md
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Issues**: When creating an issue, assign it to the person responsible for resolving it to ensure proper notifications
2. **Pull Requests**: Ensure your code follows Go conventions and includes appropriate tests
3. **Localization**: New language translations are appreciated

## License

This project is licensed under the terms specified in the [LICENSE](LICENSE) file.

## Built With

- [Go](https://golang.org/) - Programming language
- [Fyne](https://fyne.io/) - Cross-platform GUI framework
- [FFmpeg](https://ffmpeg.org/) - Media processing
