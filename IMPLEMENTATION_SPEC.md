# Apple Haptic Creator - Go Implementation Specification

## Overview
This document provides a comprehensive specification for rewriting the Apple Haptic Creator from Python to Go, with enhanced features including musical timing (BPM, bars, beats, time signatures) and a beautiful, syntax-sugared API.

## Current Python Implementation Analysis

### Core Components
1. **ahap.py**: Main library containing:
   - `HapticCurve`: Represents control points for parameter curves
   - `CurveParamID`: Enum for parameter curve types (intensity, sharpness, etc.)
   - `ParamID`: Enum for event parameter types
   - `AHAP`: Main class for creating AHAP files
   - `create_curve()`: Helper for creating parameter curves
   - `freq()`: Converts frequency (Hz) to sharpness value

2. **makeahap.py**: Example creating motorcycle sound with haptics
3. **music.py**: MIDI to haptics converter using librosa and mido
4. **test.py**: Basic unit tests for freq() function

### AHAP File Format
AHAP files are JSON with the structure:
```json
{
  "Version": 1.0,
  "Metadata": {
    "Project": "Basis",
    "Created": "timestamp",
    "Description": "description",
    "Created By": "creator"
  },
  "Pattern": [
    {
      "Event": {
        "Time": 0.0,
        "EventType": "HapticTransient|HapticContinuous|AudioCustom|AudioContinuous",
        "EventParameters": [...],
        "EventDuration": 1.0,  // optional
        "EventWaveformPath": "path"  // optional
      }
    },
    {
      "ParameterCurve": {
        "ParameterID": "HapticIntensityControl",
        "Time": 0.0,
        "ParameterCurveControlPoints": [...]
      }
    }
  ]
}
```

## Go Implementation Design

### Project Structure
```
/
├── go.mod
├── go.sum
├── README.md
├── IMPLEMENTATION_SPEC.md
├── .gitignore
├── pkg/
│   └── ahap/
│       ├── ahap.go          # Core AHAP types and creation
│       ├── events.go        # Event creation methods
│       ├── curves.go        # Parameter curves
│       ├── musical.go       # BPM, beats, bars, time signatures
│       ├── builder.go       # Fluent API builder pattern
│       └── ahap_test.go     # Tests
├── cmd/
│   ├── makeahap/
│   │   └── main.go         # Motorcycle sound example
│   ├── midi2ahap/
│   │   └── main.go         # MIDI to haptics converter
│   └── ahapgen/
│       └── main.go         # General purpose CLI tool
├── examples/
│   └── *.ahap              # Example AHAP files
└── demo/
    └── *.ahap              # Demo AHAP files
```

### Core Package: pkg/ahap

#### Types

```go
// AHAP represents a complete Apple Haptic pattern
type AHAP struct {
    Version  float64  `json:"Version"`
    Metadata Metadata `json:"Metadata"`
    Pattern  []Pattern `json:"Pattern"`
}

// Metadata contains file metadata
type Metadata struct {
    Project     string `json:"Project"`
    Created     string `json:"Created"`
    Description string `json:"Description"`
    CreatedBy   string `json:"Created By"`
}

// Pattern can be either an Event or a ParameterCurve
type Pattern struct {
    Event          *Event          `json:"Event,omitempty"`
    ParameterCurve *ParameterCurve `json:"ParameterCurve,omitempty"`
}

// Event represents a haptic or audio event
type Event struct {
    Time              float64          `json:"Time"`
    EventType         string           `json:"EventType"`
    EventParameters   []EventParameter `json:"EventParameters"`
    EventDuration     *float64         `json:"EventDuration,omitempty"`
    EventWaveformPath *string          `json:"EventWaveformPath,omitempty"`
}

// EventParameter represents a parameter of an event
type EventParameter struct {
    ParameterID    string  `json:"ParameterID"`
    ParameterValue float64 `json:"ParameterValue"`
}

// ParameterCurve represents a dynamic parameter change over time
type ParameterCurve struct {
    ParameterID                  string         `json:"ParameterID"`
    Time                         float64        `json:"Time"`
    ParameterCurveControlPoints []ControlPoint `json:"ParameterCurveControlPoints"`
}

// ControlPoint represents a point in a parameter curve
type ControlPoint struct {
    Time           float64 `json:"Time"`
    ParameterValue float64 `json:"ParameterValue"`
}

// EventType constants
const (
    EventTypeHapticTransient  = "HapticTransient"
    EventTypeHapticContinuous = "HapticContinuous"
    EventTypeAudioCustom      = "AudioCustom"
    EventTypeAudioContinuous  = "AudioContinuous"
)

// ParameterID constants
const (
    ParamHapticIntensity  = "HapticIntensity"
    ParamHapticSharpness  = "HapticSharpness"
    ParamHapticAttackTime = "HapticAttackTime"
    ParamHapticDecayTime  = "HapticDecayTime"
    ParamHapticReleaseTime = "HapticReleaseTime"
    ParamAudioVolume      = "AudioVolume"
    ParamAudioPitch       = "AudioPitch"
    ParamAudioPan         = "AudioPan"
    ParamAudioBrightness  = "AudioBrightness"
    ParamAudioAttackTime  = "AudioAttackTime"
    ParamAudioDecayTime   = "AudioDecayTime"
    ParamAudioReleaseTime = "AudioReleaseTime"
)

// CurveParameterID constants
const (
    CurveHapticIntensity  = "HapticIntensityControl"
    CurveHapticSharpness  = "HapticSharpnessControl"
    CurveHapticAttackTime = "HapticAttackTimeControl"
    CurveHapticDecayTime  = "HapticDecayTimeControl"
    CurveHapticReleaseTime = "HapticReleaseTimeControl"
    CurveAudioVolume      = "AudioVolumeControl"
    CurveAudioPitch       = "AudioPitchControl"
    CurveAudioPan         = "AudioPanControl"
    CurveAudioBrightness  = "AudioBrightnessControl"
    CurveAudioAttackTime  = "AudioAttackTimeControl"
    CurveAudioDecayTime   = "AudioDecayTimeControl"
    CurveAudioReleaseTime = "AudioReleaseTimeControl"
)
```

#### Musical Timing Types

```go
// TimeSignature represents musical time signature
type TimeSignature struct {
    Numerator   int // Beats per bar
    Denominator int // Note value (4 = quarter note)
}

// MusicalContext provides musical timing functionality
type MusicalContext struct {
    BPM           float64
    TimeSignature TimeSignature
}

// Beat represents a single beat in musical time
type Beat float64

// Bar represents a bar/measure in musical time
type Bar float64
```

### Beautiful API Design

#### Fluent Builder Pattern

```go
// New creates a new AHAP builder
func New(description, creator string) *Builder

// Builder provides a fluent API for creating AHAP files
type Builder struct {
    ahap *AHAP
    musical *MusicalContext
}

// WithBPM sets the BPM for musical timing
func (b *Builder) WithBPM(bpm float64) *Builder

// WithTimeSignature sets the time signature
func (b *Builder) WithTimeSignature(numerator, denominator int) *Builder

// Transient adds a haptic transient event
func (b *Builder) Transient(time float64) *TransientBuilder

// Continuous adds a haptic continuous event
func (b *Builder) Continuous(time, duration float64) *ContinuousBuilder

// At adds event at specific bar and beat
func (b *Builder) At(bar, beat int) *EventBuilder

// Build returns the final AHAP
func (b *Builder) Build() *AHAP

// Export writes AHAP to file
func (b *Builder) Export(filename string, indent bool) error
```

#### Example Usage

```go
// Simple example
ahap := ahap.New("My Haptic", "Creator").
    Transient(0.0).Intensity(1.0).Sharpness(0.5).Add().
    Continuous(1.0, 2.0).Intensity(0.8).Sharpness(0.7).Add().
    Build()

// Musical example with BPM and beats
ahap := ahap.New("Musical Haptic", "Creator").
    WithBPM(120).
    WithTimeSignature(4, 4).
    At(0, 0).Transient().Intensity(1.0).Sharpness(1.0).Add().
    At(0, 1).Transient().Intensity(0.8).Sharpness(0.8).Add().
    At(0, 2).Transient().Intensity(0.8).Sharpness(0.8).Add().
    At(0, 3).Transient().Intensity(0.8).Sharpness(0.8).Add().
    At(1, 0).Continuous(1.0).Intensity(0.5).Sharpness(0.3).Add().
    Build()

// With curves
ahap := ahap.New("Curve Example", "Creator").
    Continuous(0.0, 2.0).Intensity(0.5).Sharpness(0.5).Add().
    Curve(ahap.CurveHapticSharpness).
        From(0.0, 0.3).To(2.0, 0.8).Steps(20).Add().
    Build()

ahap.Export("output.ahap", true)
```

### Command Line Utilities

#### cmd/makeahap
Recreates the motorcycle sound example:
```bash
go run cmd/makeahap/main.go -output bike.ahap
```

#### cmd/midi2ahap
Converts MIDI files to haptics:
```bash
go run cmd/midi2ahap/main.go -input song.mid -output song.ahap
```

#### cmd/ahapgen
General purpose haptic generator with CLI:
```bash
# Interactive mode
go run cmd/ahapgen/main.go

# Command mode
go run cmd/ahapgen/main.go --transient 0.0 --intensity 1.0 --sharpness 0.5 -o output.ahap
```

### Key Features

1. **Type Safety**: Go's strong typing prevents runtime errors
2. **Fluent API**: Chain methods for readable haptic creation
3. **Musical Timing**: BPM, beats, bars, time signatures built-in
4. **Zero Dependencies**: Core library uses only standard library
5. **Fast**: Go's performance for large AHAP files
6. **Builder Pattern**: Flexible construction of complex haptics
7. **Validation**: Built-in validation of parameters (0-1 ranges, etc.)

### Musical Timing Implementation

```go
// Convert beat to seconds
func (m *MusicalContext) BeatToSeconds(beat Beat) float64 {
    // 60 seconds per minute / BPM = seconds per beat
    return float64(beat) * (60.0 / m.BPM)
}

// Convert bar to seconds
func (m *MusicalContext) BarToSeconds(bar Bar) float64 {
    beatsPerBar := float64(m.TimeSignature.Numerator)
    return float64(bar) * beatsPerBar * (60.0 / m.BPM)
}

// Get duration of one bar
func (m *MusicalContext) BarDuration() float64 {
    return m.BarToSeconds(1)
}

// Get duration of one beat
func (m *MusicalContext) BeatDuration() float64 {
    return m.BeatToSeconds(1)
}
```

### Helper Functions

```go
// FreqToSharpness converts frequency (Hz) to sharpness value (0-1)
func FreqToSharpness(freq float64, normalize bool) (float64, error)

// CreateCurve creates interpolated control points between start and end
func CreateCurve(startTime, endTime, startValue, endValue float64, steps int) []ControlPoint

// LinearInterpolation creates a linear curve
func LinearInterpolation(start, end ControlPoint, steps int) []ControlPoint

// EaseInOut creates an ease-in-out curve
func EaseInOut(start, end ControlPoint, steps int) []ControlPoint

// Exponential creates an exponential curve
func Exponential(start, end ControlPoint, steps int, exponent float64) []ControlPoint
```

### Testing Strategy

1. **Unit Tests**: Test each component individually
2. **Integration Tests**: Test full AHAP creation workflow
3. **Golden Files**: Compare output with known-good AHAP files
4. **Musical Timing Tests**: Verify BPM/beat/bar calculations
5. **Validation Tests**: Ensure parameter ranges enforced

### Migration from Python

1. All Python functionality will be preserved
2. Enhanced with new features (BPM, musical timing)
3. Better type safety and performance
4. More intuitive API
5. Standard Go project structure
6. Better documentation

### Dependencies

**Core Library (pkg/ahap):**
- No external dependencies (pure Go stdlib)

**cmd/midi2ahap:**
- github.com/go-audio/midi (for MIDI parsing)
- Alternative: gitlab.com/gomidi/midi/v2

**Optional:**
- github.com/spf13/cobra (for CLI)
- github.com/stretchr/testify (for testing)

## Implementation Steps

1. ✅ Create this specification document
2. Initialize Go module
3. Implement core AHAP types and JSON marshaling
4. Implement basic event creation
5. Implement parameter curves
6. Implement musical timing features
7. Implement fluent builder API
8. Create cmd/makeahap utility
9. Create cmd/midi2ahap utility (with MIDI library)
10. Create cmd/ahapgen utility
11. Write comprehensive tests
12. Update README.md
13. Update .gitignore for Go
14. Remove Python implementation
15. Verify all examples work

## Backwards Compatibility

The Go implementation will:
- Generate identical AHAP files (JSON format is identical)
- Preserve all existing functionality
- Example AHAP files remain unchanged
- Can read/validate existing AHAP files

## Additional Enhancements

1. **Validation**: Validate AHAP files for correctness
2. **Merge**: Combine multiple AHAP files
3. **Transform**: Apply transformations (time shift, scale, etc.)
4. **Analyze**: Extract info from AHAP files (duration, event count, etc.)
5. **Preview**: ASCII art visualization of haptic patterns
6. **Templates**: Pre-built haptic patterns (drumbeat, heartbeat, etc.)

## Success Criteria

- [ ] All Python functionality reimplemented in Go
- [ ] Musical timing features working (BPM, beats, bars)
- [ ] Fluent API is intuitive and well-documented
- [ ] All cmd utilities working
- [ ] Tests pass with >80% coverage
- [ ] README updated with Go examples
- [ ] Python code removed
- [ ] Example AHAP files still valid
