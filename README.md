# GoFetch

A cross-platform system information fetcher written in Go, inspired by neofetch and similar tools.

## Overview

GoFetch is a command-line system information tool that displays hardware and system details in a clean, formatted output. It provides information about your CPU, memory, storage, GPU, display, and host system.

## Features

- **CPU Information**: Processor details, cores, and architecture
- **Memory Usage**: RAM usage statistics  
- **Disk Information**: Storage device details and usage
- **GPU Information**: Graphics card details (Windows-specific implementation included)
- **Display Information**: Monitor and display configuration (Windows-specific implementation included)
- **Host Information**: System hostname and basic OS details
- **Cross-platform Support**: Works on multiple operating systems with platform-specific optimizations

## Installation

### Prerequisites

- Go 1.16 or higher

### Build from Source

1. Clone the repository:

```bash
git clone https://github.com/AtifChy/gofetch.git
cd gofetch
```

2. Build the application:

```bash
make build
```

Or manually:

```bash
go build ./cmd/gofetch
```

3. Run the application:

```bash
./gofetch
```

## Usage

Simply run the executable to display system information:

```bash
gofetch
```

The tool will automatically detect your system and display relevant information in a formatted output.

## Project Structure

```
gofetch/
├── cmd/gofetch/            # Main application code
│   ├── main.go             # Entry point
│   ├── cpu.go              # CPU information gathering
│   ├── memory.go           # Memory usage statistics
│   ├── disk.go             # Disk and storage information
│   ├── gpu.go              # GPU information (cross-platform)
│   ├── gpu_windows.go      # Windows-specific GPU implementation
│   ├── display.go          # Display information (cross-platform)
│   ├── display_windows.go  # Windows-specific display implementation
│   └── host.go             # Host system information
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── Makefile                # Build automation
└── LICENSE                 # Project license
```

## Development

### Building

Use the provided Makefile for common development tasks:

```bash
# Build the application
make build

# Run tests
make test

# Clean build artifacts
make clean
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes
5. Submit a pull request

## License

This project is licensed under the terms specified in the LICENSE file.

## Acknowledgments

Inspired by [neofetch](https://github.com/dylanaraps/neofetch) and other system information tools.
