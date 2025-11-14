# MIDI to AHAP Converter

Converts MIDI files to Apple Haptic (AHAP) format with intelligent drum detection.

## Features

### Drum Detection

The converter automatically detects General MIDI drum channel (channel 10) and converts drum notes to transient haptic events with optimized characteristics.

**Drum Categories:**

1. **Bass Drums** (Notes 35-36)
   - High intensity (1.0), low sharpness (0.2)
   - Creates deep, powerful haptic "thump"

2. **Snares** (Notes 38, 40)
   - Very high intensity (0.9-0.95), very high sharpness (0.85-0.9)
   - Creates sharp, crisp haptic "snap"

3. **Toms** (Notes 41, 43, 45, 47, 48, 50)
   - High intensity (0.85), varying sharpness (0.4-0.65)
   - Sharpness increases with pitch

4. **Hi-Hats** (Notes 42, 44, 46)
   - Medium intensity (0.5-0.6), very high sharpness (0.9-1.0)
   - Creates light, crisp haptic

5. **Cymbals** (Notes 49, 51-53, 55, 57, 59)
   - High intensity (0.7-0.9), high sharpness (0.7-0.9)
   - Creates sustained, bright haptic

6. **Percussion** (Notes 37, 39, 54, 56-81)
   - Varied profiles optimized for each instrument
   - Includes hand claps, tambourine, cowbell, bongos, congas, etc.

### Melodic Conversion

Non-drum MIDI notes are converted to continuous haptic events:
- **Frequency → Sharpness**: Note frequency maps to haptic sharpness (80-230 Hz range)
- **Velocity → Intensity**: MIDI velocity (0-127) maps to haptic intensity (0.0-1.0)
- **Duration**: Note duration preserved from MIDI timing

## Usage

```bash
# Basic conversion
midi2ahap -input song.mid -output song.ahap

# With pretty-printed JSON
midi2ahap -input song.mid -output song.ahap -indent

# Disable drum detection (treat all as melodic)
midi2ahap -input song.mid -output song.ahap -drums=false
```

## Options

- `-input`: Input MIDI file (required)
- `-output`: Output AHAP file (default: input filename with .ahap extension)
- `-indent`: Pretty-print JSON output for readability
- `-drums`: Enable drum detection (default: true)

## Examples

### Converting a drum pattern:
```bash
midi2ahap -input drums.mid -output drums.ahap -indent
```

Output statistics:
```
Conversion statistics:
  Drum events (transients): 128
  Melodic events (continuous): 0
  Total haptic events: 128
```

### Converting a full song:
```bash
midi2ahap -input song.mid -output song.ahap
```

Output statistics:
```
Conversion statistics:
  Drum events (transients): 256
  Melodic events (continuous): 1024
  Total haptic events: 1280
```

## General MIDI Drum Map Reference

The converter uses the standard General MIDI Level 1 Percussion Key Map (channel 10):

| Note | Drum Sound           | Intensity | Sharpness |
|------|---------------------|-----------|-----------|
| 35   | Acoustic Bass Drum  | 1.0       | 0.2       |
| 36   | Bass Drum 1         | 1.0       | 0.2       |
| 38   | Acoustic Snare      | 0.95      | 0.85      |
| 40   | Electric Snare      | 0.9       | 0.9       |
| 42   | Closed Hi-Hat       | 0.5       | 1.0       |
| 44   | Pedal Hi-Hat        | 0.55      | 0.95      |
| 46   | Open Hi-Hat         | 0.6       | 0.9       |
| 49   | Crash Cymbal 1      | 0.9       | 0.85      |
| 51   | Ride Cymbal 1       | 0.7       | 0.75      |

*And 40+ more percussion sounds...*

## Technical Details

### Channel Detection
- MIDI channel 10 (index 9) is automatically detected as the drum channel
- All notes on this channel are converted to transient events
- Other channels are converted to continuous events

### Velocity Scaling
- Drum intensity is scaled by MIDI velocity: `drum_intensity * (velocity / 127.0)`
- Melodic intensity directly maps: `velocity / 127.0`

### Timing
- Respects MIDI tempo changes (Meta Tempo events)
- Accurate timing conversion from ticks to seconds
- Preserves all note timing and duration

## Notes

- Unknown drum notes (not in the mapping table) use default values: intensity=velocity/127, sharpness=0.7
- Drum events are instantaneous (transients), even if the MIDI note has duration
- The converter tracks the count of drum vs. melodic events for statistics
