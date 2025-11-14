# Examples

This directory contains example Go programs and haptrack pattern files demonstrating how to use the Apple Haptic Creator library.

## Running the Examples

### Go Examples

```bash
cd examples
go run simple_example.go    # Basic examples
go run sequence_example.go  # Sequence builder examples
go run bike.go              # Motorcycle sound example
```

### Haptrack Pattern Files

```bash
cd examples
go run ../cmd/haptrack/main.go -input drum_pattern.hap -output drum_pattern.ahap
go run ../cmd/haptrack/main.go -input complex_pattern.hap -output complex_pattern.ahap
```

## Example Files

### Go Examples

#### simple_example.go
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

#### sequence_example.go
Generates multiple examples demonstrating the sequence builder API:
- Transients on specific beats across multiple bars
- Kick and snare drum patterns
- Every beat and every nth beat patterns
- Custom pattern functions

#### bike.go
Complete motorcycle sound example showing:
- Complex haptic patterns
- Parameter curves
- Engine acceleration simulation

### Haptrack Pattern Files

#### drum_pattern.hap
Basic drum pattern demonstrating haptrack DSL syntax:
- BPM and time signature definition
- Haptic instrument definitions
- Multiple track composition
- Note durations and rests

#### complex_pattern.hap
Advanced rock drum pattern with:
- Multiple drum instruments
- Complex rhythmic patterns
- Simultaneous tracks

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

// Add beats (bar 0, beats 0-3)
for i := 0; i < 4; i++ {
    builder.At(0, i).
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

## Sequence Builder Example

```go
// Create patterns across multiple bars easily
builder := ahap.NewBuilder("Pattern", "Creator").
    WithBPM(120).
    WithTimeSignature(4, 4)

// Add transients on beats 0 and 2 (1st and 3rd) from bars 5 to 8
builder.Sequence().TransientsOnBeats([]int{0, 2}, 5, 8, 1.0, 0.5)

// Add different pattern on beats 1 and 3 (2nd and 4th)
builder.Sequence().TransientsOnBeats([]int{1, 3}, 5, 8, 0.9, 0.8)

builder.Export("sequence.ahap", true)
```

## Haptrack DSL Example

Create a file `mypattern.hap`:

```
# Define settings
bpm = 140
time = 4/4

# Define haptic sounds
s = snare, 1.0, 0.9, down, 60
k = kick, 1.0, 0.2
h = hihat, 0.6, 1.0

# Begin tracks
begin

track1
k8-8k8-8k8-8k8-8

track2
-4s8-8-4s8-8

track3
h8h8h8h8h8h8h8h8
```

Then compile it:
```bash
haptrack -input mypattern.hap -output mypattern.ahap
```

## More Examples

See the main [README.md](../README.md) and [README_GO.md](../README_GO.md) for more examples and complete API documentation.
