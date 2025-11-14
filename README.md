# Apple Haptic Creator

A powerful Go library and command-line tools for creating Apple Haptic pattern files (AHAP). Features include a beautiful fluent API, musical timing support (BPM, bars, beats, time signatures), and MIDI to haptics conversion.

## What are AHAP files?

AHAP files are JSON-formatted Apple Haptic pattern files used in iOS games and applications to create immersive haptic experiences. They can be played directly from the Files app or any apps that support Apple's Quick Look API, making them shareable via WhatsApp, Telegram, and other platforms.

For more information, see this [article on AppleVis](https://applevis.com/forum/ios-ipados/now-possible-ios-17-can-play-haptic-signals-vibrations-special-ahap-apple-haptic).

## ✨ Features

- **Beautiful Fluent API** - Chain methods for intuitive haptic creation
- **Musical Timing** - Support for BPM, bars, beats, and time signatures
- **MIDI Conversion** - Convert MIDI files to haptic patterns
- **High Performance** - Fast Go implementation
- **Zero Dependencies** - Core library uses only Go standard library
- **Clean Architecture** - Reusable package with multiple CLI utilities

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/denizsincar29/apple_haptic_creator.git
cd apple_haptic_creator

# Build all commands
go build -o bin/makeahap cmd/makeahap/main.go
go build -o bin/midi2ahap cmd/midi2ahap/main.go
go build -o bin/ahapgen cmd/ahapgen/main.go
```

### Using as a Library

```go
package main

import "github.com/denizsincar29/apple_haptic_creator/pkg/ahap"

func main() {
    // Simple example
    builder := ahap.NewBuilder("My Haptic", "Me")
    builder.
        Transient(0.0).Intensity(1.0).Sharpness(0.5).Add().
        Continuous(1.0, 2.0).Intensity(0.8).Sharpness(0.7).Add().
        Export("example.ahap", true)
}
```

### Musical Timing Example

```go
// Create a drum beat pattern at 120 BPM
builder := ahap.NewBuilder("Drum Beat", "Creator").
    WithBPM(120).
    WithTimeSignature(4, 4)

// Add kick drum on beats 0 and 2
builder.AtBeat(0).Transient().Intensity(1.0).Sharpness(0.2).Add()
builder.AtBeat(2).Transient().Intensity(1.0).Sharpness(0.2).Add()

// Add snare on beats 1 and 3
builder.AtBeat(1).Transient().Intensity(0.9).Sharpness(0.8).Add()
builder.AtBeat(3).Transient().Intensity(0.9).Sharpness(0.8).Add()

builder.Export("drumbeat.ahap", true)
```

## 🎯 Command Line Tools

### makeahap - Motorcycle Sound Example

```bash
go run cmd/makeahap/main.go -output bike.ahap -indent
```

Creates a realistic motorcycle engine sound using haptics.

### midi2ahap - MIDI to Haptics Converter

```bash
go run cmd/midi2ahap/main.go -input song.mid -output song.ahap -indent
```

Converts MIDI files to haptic patterns, mapping notes to sharpness and velocity to intensity.

### ahapgen - Interactive Haptic Generator

```bash
go run cmd/ahapgen/main.go -bpm 120 -time 4/4 -o output.ahap
```

Interactive command-line tool with support for musical timing.

## 📚 Documentation

For complete API documentation and advanced examples, see [README_GO.md](README_GO.md).

## 📁 Project Structure

```
.
├── pkg/ahap/              # Core library package
│   ├── ahap.go           # Core AHAP types
│   ├── events.go         # Event creation
│   ├── curves.go         # Parameter curves
│   ├── musical.go        # Musical timing (BPM, bars, beats)
│   ├── builder.go        # Fluent API builder
│   └── ahap_test.go      # Tests
├── cmd/                   # Command-line utilities
│   ├── makeahap/         # Motorcycle sound example
│   ├── midi2ahap/        # MIDI converter
│   └── ahapgen/          # Interactive generator
├── ahaps/                 # Example AHAP files
└── demo/                  # Demo files (including MIDI)
```

## 🧪 Testing

```bash
go test ./pkg/ahap/        # Run tests
go test -cover ./pkg/ahap/ # With coverage
```

## 📖 Examples

The `ahaps/` folder contains example AHAP files:
- `bike.ahap` - Motorcycle engine sound
- `interval.ahap` - Simple interval pattern
- `music.ahap` - Musical pattern
- `notes.ahap` - Musical notes

Demo MIDI files for testing:
- `demo/themeters.mid` - Example MIDI file
- `demo/donnalee.mid` - Example MIDI file

## 🤝 Contributing

Contributions are welcome! Areas for improvement:
- Additional curve interpolation methods
- More MIDI conversion options
- Haptic pattern templates library
- Visualization tools
- Pattern analysis utilities

## 📄 License

See [IMPLEMENTATION_SPEC.md](IMPLEMENTATION_SPEC.md) for complete implementation details.

## 🙏 Credits

- Original Python implementation by Deniz Sincar
- Go rewrite with enhanced features and musical timing support