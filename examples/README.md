# Examples

This directory contains example Go programs demonstrating how to use the Apple Haptic Creator library.

## Running the Examples

```bash
cd examples
go run simple_example.go
```

This will generate three AHAP files demonstrating different features:

### simple.ahap
Basic example showing:
- Transient events (sharp taps)
- Continuous events (sustained vibrations)

### musical.ahap
Musical timing example showing:
- BPM-based timing (120 BPM)
- Time signature (4/4)
- Beats and bars
- Accented first beat

### curves.ahap
Parameter curves example showing:
- Linear interpolation
- Ease-in-out interpolation
- Dynamic parameter changes over time

## Code Example

```go
package main

import "github.com/denizsincar29/apple_haptic_creator/pkg/ahap"

func main() {
    // Simple example with fluent API
    builder := ahap.NewBuilder("My Haptic", "Creator")
    builder.
        Transient(0.0).Intensity(1.0).Sharpness(0.5).Add().
        Continuous(1.0, 2.0).Intensity(0.8).Sharpness(0.7).Add().
        Export("output.ahap", true)
}
```

## Musical Timing Example

```go
// Create a drum beat at 120 BPM
builder := ahap.NewBuilder("Drum Beat", "Creator").
    WithBPM(120).
    WithTimeSignature(4, 4)

// Add beats
for i := 0; i < 4; i++ {
    builder.AtBeat(ahap.Beat(i)).
        Transient().Intensity(1.0).Sharpness(0.5).Add()
}

builder.Export("drumbeat.ahap", true)
```

## Curves Example

```go
builder := ahap.NewBuilder("Curve", "Creator")

// Add continuous event
builder.Continuous(0.0, 2.0).Intensity(0.5).Sharpness(0.5).Add()

// Add sharpness curve
builder.Curve(ahap.CurveHapticSharpness).
    From(0.0, 0.3).To(2.0, 0.8).Steps(20).Add()

builder.Export("curves.ahap", true)
```

## More Examples

See the main [README.md](../README.md) and [README_GO.md](../README_GO.md) for more examples and complete API documentation.
