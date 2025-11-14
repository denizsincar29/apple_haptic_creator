# Haptrack - Haptic Pattern DSL

Haptrack is a domain-specific language (DSL) for creating haptic patterns using a musical notation-inspired syntax. It allows you to define haptic "instruments" and compose them into tracks, similar to how you might write drum notation.

## Features

- **Letter-based notation**: Define what each letter (a-z) represents as a haptic event
- **Musical timing**: Support for BPM and time signatures
- **Multiple tracks**: Create up to multiple simultaneous haptic tracks
- **Note durations**: Use standard musical note values (1, 2, 4, 8, 16 for whole, half, quarter, eighth, sixteenth notes)
- **Parameter curves**: Optional curves that modify sharpness over time
- **Rests**: Use dashes (-) for silent beats

## File Format

A Haptrack file consists of two sections:

### 1. Definitions Section

Define your settings and haptic sounds:

```
# Settings
bpm = 120
time = 4/4

# Haptic definitions
# Format: letter = name, intensity, sharpness [, curve_direction, duration_ms]
s = snare, 1.0, 0.9, down, 60
k = kick, 1.0, 0.2
h = hihat, 0.6, 1.0
```

**Haptic Definition Format:**
- `letter`: Single letter (a-z) that will trigger this haptic
- `name`: Descriptive name (for documentation)
- `intensity`: Haptic intensity (0.0 to 1.0)
- `sharpness`: Haptic sharpness (0.0 to 1.0)
- `curve_direction` (optional): "up" or "down" - direction of sharpness curve
- `duration_ms` (optional): Duration of the curve in milliseconds

### 2. Tracks Section

After the `begin` marker, define your tracks:

```
begin

track1
k8k8s8k8k8k8s8k8

track2
h8h8h8h8h8h8h8h8
```

## Pattern Notation

### Note Durations

Numbers after letters indicate note duration:
- `1` = whole note (4 beats)
- `2` = half note (2 beats)
- `4` = quarter note (1 beat)
- `8` = eighth note (0.5 beats)
- `16` = sixteenth note (0.25 beats)

**Default**: If no number is specified, eighth note (8) is assumed.

### Rests

Use dash `-` for rests (silence):
- `-` = eighth note rest (default)
- `-4` = quarter note rest
- `-8` = eighth note rest
- `-16` = sixteenth note rest

### Examples

```
# Kick on every beat
k4k4k4k4

# Kick on beats 1 and 3, snare on beats 2 and 4
k4s4k4s4

# Hi-hat eighth notes with rests
h8-8h8-8h8-8h8-8

# Rapid sixteenth notes
s16s16s16s16s16s16s16s16
```

## Complete Example

```
# Rock Beat
bpm = 120
time = 4/4

# Define instruments
k = kick, 1.0, 0.2
s = snare, 0.95, 0.85, down, 50
h = hihat, 0.5, 0.95
c = crash, 0.95, 0.9, down, 120

begin

# Kick pattern
track1
k8-8k8-8k8-8k8-8

# Snare on 2 and 4
track2
-4s8-8-4s8-8

# Hi-hat eighth notes
track3
h8h8h8h8h8h8h8h8
```

## Usage

```bash
# Basic usage
haptrack -input pattern.hap -output output.ahap

# With custom output
haptrack -input mybeat.hap -output mybeat.ahap

# Without indentation (smaller file)
haptrack -input pattern.hap -output output.ahap -indent=false
```

## Tips

1. **Start Simple**: Begin with a basic kick and snare pattern, then add complexity
2. **Use Comments**: Document your patterns with `#` comments
3. **Test Incrementally**: Start with one track, verify it works, then add more
4. **Curve Effects**: Use curves for cymbal crashes and snare hits for more realistic feel
5. **Multiple Tracks**: iOS can play 2-3 haptic tracks simultaneously
6. **Timing**: Ensure your pattern fits within reasonable duration (iOS has limits)

## Pattern Ideas

### Basic Rock Beat
```
bpm = 120
time = 4/4
k = kick, 1.0, 0.2
s = snare, 0.95, 0.85
h = hihat, 0.5, 0.95

begin
track1
k8-8k8-8k8-8k8-8
track2
-4s8-8-4s8-8
track3
h8h8h8h8h8h8h8h8
```

### Four on the Floor (Electronic)
```
bpm = 128
time = 4/4
k = kick, 1.0, 0.15
h = hihat, 0.4, 1.0
c = clap, 0.8, 0.9

begin
track1
k4k4k4k4
track2
-4c4-4c4
track3
h8h8h8h8h8h8h8h8
```

### Waltz (3/4 time)
```
bpm = 180
time = 3/4
k = kick, 1.0, 0.2
s = snare, 0.7, 0.6

begin
track1
k4s8s8s8s8
```

## Command-Line Options

- `-input <file>`: Input haptrack file (required)
- `-output <file>`: Output AHAP file (default: output.ahap)
- `-indent`: Indent JSON output for readability (default: true)

## See Also

- [Main README](../../README.md) - General library documentation
- [Examples](../../examples/) - More example patterns
- [AHAP Specification](../../IMPLEMENTATION_SPEC.md) - AHAP file format details
