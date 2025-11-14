# Apple Haptic Creator (Go)

A powerful Go library and command-line tools for creating Apple Haptic pattern files (AHAP). Features include a beautiful fluent API, musical timing support (BPM, bars, beats, time signatures), and MIDI to haptics conversion.

## What are AHAP files?

AHAP files are JSON-formatted Apple Haptic pattern files used in iOS games and applications to create immersive haptic experiences. They can be played directly from the Files app or any apps that support Apple's Quick Look API, making them shareable via WhatsApp, Telegram, and other platforms.

For more information, see this [article on AppleVis](https://applevis.com/forum/ios-ipados/now-possible-ios-17-can-play-haptic-signals-vibrations-special-ahap-apple-haptic).

## Features

✨ **Beautiful Fluent API** - Chain methods for intuitive haptic creation  
🎵 **Musical Timing** - Support for BPM, bars, beats, and time signatures  
🎹 **MIDI Conversion** - Convert MIDI files to haptic patterns  
🚀 **High Performance** - Fast Go implementation  
🔧 **Zero Dependencies** - Core library uses only Go standard library  
📦 **Clean Architecture** - Reusable package with multiple CLI utilities

## Installation

```bash
# Clone the repository
git clone https://github.com/denizsincar29/apple_haptic_creator.git
cd apple_haptic_creator

# Build all commands
go build -o bin/makeahap cmd/makeahap/main.go
go build -o bin/midi2ahap cmd/midi2ahap/main.go
go build -o bin/ahapgen cmd/ahapgen/main.go

# Or install globally
go install ./cmd/makeahap
go install ./cmd/midi2ahap
go install ./cmd/ahapgen
```

## Quick Start

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

// Add kick drum on beats 0 and 2 of bar 0
builder.At(0, 0).Transient().Intensity(1.0).Sharpness(0.2).Add()
builder.At(0, 2).Transient().Intensity(1.0).Sharpness(0.2).Add()

// Add snare on beats 1 and 3 of bar 0
builder.At(0, 1).Transient().Intensity(0.9).Sharpness(0.8).Add()
builder.At(0, 3).Transient().Intensity(0.9).Sharpness(0.8).Add()

// Hi-hat on every half beat (using absolute beat positions across bars)
// For simplicity, use Transient with time calculation
for i := 0; i < 8; i++ {
    // Calculate bar and beat for half-beat positions
    bar := (i / 2) / 4  // Which bar (8 half-beats = 4 beats = 1 bar at 4/4)
    beat := (i / 2) % 4 // Which beat in the bar
    builder.At(bar, beat).Transient().Intensity(0.5).Sharpness(1.0).Add()
}

builder.Export("drumbeat.ahap", true)
```

### Parameter Curves

```go
// Create a haptic with dynamic parameter changes
builder := ahap.NewBuilder("Curve Example", "Creator")

// Add a continuous event
builder.Continuous(0.0, 2.0).Intensity(0.5).Sharpness(0.5).Add()

// Add a sharpness curve that goes from 0.3 to 0.8 over 2 seconds
builder.Curve(ahap.CurveHapticSharpness).
    From(0.0, 0.3).To(2.0, 0.8).Steps(20).Add()

// Ease-in-out curve
builder.Curve(ahap.CurveHapticIntensity).
    From(2.0, 0.5).To(4.0, 1.0).EaseInOut(15).Add()

builder.Export("curves.ahap", true)
```

## Command Line Tools

### makeahap - Motorcycle Sound Example

Creates a realistic motorcycle engine sound using haptics:

```bash
go run cmd/makeahap/main.go -output bike.ahap -indent
```

### midi2ahap - MIDI to Haptics Converter

Converts MIDI files to haptic patterns:

```bash
go run cmd/midi2ahap/main.go -input song.mid -output song.ahap -indent
```

The converter maps:
- MIDI notes to haptic sharpness (based on frequency)
- Note velocity to haptic intensity
- Note duration to continuous haptic events

### ahapgen - Interactive Haptic Generator

Interactive command-line tool for creating haptics:

```bash
# With musical timing
go run cmd/ahapgen/main.go -bpm 120 -time 4/4 -o output.ahap

# Interactive mode commands:
# > t 0.0 1.0 0.5              # Transient at 0s
# > c 1.0 2.0 0.8 0.7          # Continuous at 1s for 2s
# > beat 0 1.0 0.8             # Transient at beat 0
# > bar 1 0.8 0.5              # Transient at bar 1
# > export                      # Save and exit
```

## API Reference

### Core Types

```go
type AHAP struct {
    Version  float64
    Metadata Metadata
    Pattern  []Pattern
}

// Create new AHAP
ahap := ahap.New("description", "creator")
```

### Builder API

```go
builder := ahap.NewBuilder("description", "creator")

// Musical timing
builder.WithBPM(120)
builder.WithTimeSignature(4, 4)

// Events
builder.Transient(time).Intensity(i).Sharpness(s).Add()
builder.Continuous(time, duration).Intensity(i).Sharpness(s).Add()

// Musical events (bar, beat)
builder.At(bar, beat).Transient().Intensity(i).Sharpness(s).Add()
builder.At(bar, beat).Continuous(duration).Intensity(i).Sharpness(s).Add()
builder.At(bar, beat).ContinuousBeats(2).Intensity(i).Sharpness(s).Add()

// Curves
builder.Curve(paramID).From(start, val1).To(end, val2).Steps(n).Add()
builder.Curve(paramID).From(start, val1).To(end, val2).EaseInOut(n).Add()
builder.Curve(paramID).From(start, val1).To(end, val2).Exponential(n, exp).Add()

// Export
builder.Export("output.ahap", true)
```

### Event Parameters

#### Haptic Parameters
- `ParamHapticIntensity` - Intensity (0.0 to 1.0)
- `ParamHapticSharpness` - Sharpness (0.0 to 1.0)
- `ParamHapticAttackTime` - Attack time
- `ParamHapticDecayTime` - Decay time
- `ParamHapticReleaseTime` - Release time

#### Audio Parameters
- `ParamAudioVolume` - Volume (0.0 to 1.0)
- `ParamAudioPitch` - Pitch
- `ParamAudioPan` - Pan
- `ParamAudioBrightness` - Brightness

### Curve Parameters

Append "Control" to parameter names:
- `CurveHapticIntensity` → `HapticIntensityControl`
- `CurveHapticSharpness` → `HapticSharpnessControl`
- etc.

### Helper Functions

```go
// Convert frequency (80-230 Hz) to sharpness (0-1)
sharpness, err := ahap.FreqToSharpness(frequency, normalize)

// Create interpolated curve
points := ahap.CreateCurve(startTime, endTime, startVal, endVal, steps)

// Ease-in-out interpolation
points := ahap.EaseInOut(startPoint, endPoint, steps)

// Exponential interpolation
points := ahap.Exponential(startPoint, endPoint, steps, exponent)
```

### Musical Timing

```go
mc := ahap.NewMusicalContext(120, 4, 4) // 120 BPM, 4/4 time

// Convert musical time to seconds
seconds := mc.BeatToSeconds(ahap.Beat(4))    // 4 beats
seconds := mc.BarToSeconds(ahap.Bar(2))      // 2 bars

// Get durations
beatDur := mc.BeatDuration()  // Seconds per beat
barDur := mc.BarDuration()    // Seconds per bar

// Convert back
beats := mc.SecondsToBeats(seconds)
bars := mc.SecondsToBars(seconds)
```

## Project Structure

```
.
├── pkg/ahap/              # Core library package
│   ├── ahap.go           # Core AHAP types
│   ├── events.go         # Event creation
│   ├── curves.go         # Parameter curves
│   ├── musical.go        # Musical timing
│   ├── builder.go        # Fluent API builder
│   └── ahap_test.go      # Tests
├── cmd/                   # Command-line utilities
│   ├── makeahap/         # Motorcycle sound example
│   ├── midi2ahap/        # MIDI converter
│   └── ahapgen/          # Interactive generator
├── examples/              # Example AHAP files
└── demo/                  # Demo files (including MIDI)
```

## Examples Directory

The repository includes several example AHAP files:
- `examples/bike.ahap` - Motorcycle engine sound
- `examples/interval.ahap` - Simple interval pattern
- `examples/music.ahap` - Musical pattern
- `examples/notes.ahap` - Musical notes

Demo MIDI files for testing:
- `demo/themeters.mid` - Example MIDI file
- `demo/donnalee.mid` - Example MIDI file

## Testing

```bash
# Run tests
go test ./pkg/ahap/

# Run tests with coverage
go test -cover ./pkg/ahap/

# Run tests with verbose output
go test -v ./pkg/ahap/
```

## Advanced Features

### Combining Multiple Patterns

```go
// Create base pattern
base := ahap.NewBuilder("Base", "Creator")
base.Continuous(0, 5.0).Intensity(0.3).Sharpness(0.2).Add()

// Add accents on top
for i := 0.0; i < 5.0; i += 0.5 {
    base.Transient(i).Intensity(1.0).Sharpness(0.8).Add()
}

base.Export("combined.ahap", true)
```

### Complex Curves

```go
builder := ahap.NewBuilder("Complex Curves", "Creator")

// Multi-segment curve
builder.Curve(ahap.CurveHapticSharpness).
    From(0.0, 0.0).To(1.0, 0.5).Steps(10).
    From(1.0, 0.5).To(2.0, 1.0).EaseInOut(10).
    From(2.0, 1.0).To(3.0, 0.0).Exponential(10, 2.0).
    Add()
```

### Working with Bars and Beats

```go
builder := ahap.NewBuilder("Musical", "Creator").
    WithBPM(140).
    WithTimeSignature(3, 4)  // 3/4 time (waltz)

// Waltz pattern - accent on first beat
for bar := 0; bar < 4; bar++ {
    builder.At(bar, 0).
        Transient().Intensity(1.0).Sharpness(0.5).Add()
    
    builder.At(bar, 1).
        Transient().Intensity(0.6).Sharpness(0.5).Add()
    
    builder.At(bar, 2).
        Transient().Intensity(0.6).Sharpness(0.5).Add()
}
```

## Performance Considerations

- The library uses efficient JSON marshaling from Go's standard library
- Large AHAP files (1000+ events) are handled efficiently
- MIDI conversion processes all tracks in a single pass
- Builder pattern allows for memory-efficient construction

## Compatibility

- **Go Version**: 1.19 or higher
- **AHAP Version**: 1.0 (standard Apple format)
- **Output Format**: JSON (compatible with iOS 13+)

## Contributing

Contributions are welcome! Areas for improvement:
- Additional curve interpolation methods
- More MIDI conversion options
- Haptic pattern templates library
- Visualization tools
- Pattern analysis utilities

## License

This project maintains compatibility with the original Python implementation while adding significant enhancements.

## Credits

- Original Python implementation by Deniz Sincar
- Go rewrite with enhanced features and musical timing support

## See Also

- [Apple Haptic Documentation](https://developer.apple.com/documentation/corehaptics)
- [AHAP File Format](https://developer.apple.com/documentation/corehaptics/representing_haptic_patterns_in_ahap_files)
