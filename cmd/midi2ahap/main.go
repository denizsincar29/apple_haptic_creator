package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizsincar29/apple_haptic_creator/pkg/ahap"
	"gitlab.com/gomidi/midi/v2/smf"
)

// midiNoteToFreq converts MIDI note number to frequency in Hz
func midiNoteToFreq(note uint8) float64 {
	// A4 (note 69) = 440 Hz
	return 440.0 * math.Pow(2.0, (float64(note)-69.0)/12.0)
}

// DrumMapping represents haptic characteristics for a drum sound
type DrumMapping struct {
	Name      string
	Intensity float64
	Sharpness float64
}

// General MIDI drum mappings (channel 10)
// Based on GM standard percussion map
var drumMappings = map[uint8]DrumMapping{
	// Bass drums
	35: {"Acoustic Bass Drum", 1.0, 0.2},
	36: {"Bass Drum 1", 1.0, 0.2},
	
	// Snares
	38: {"Acoustic Snare", 0.95, 0.85},
	40: {"Electric Snare", 0.9, 0.9},
	
	// Toms
	41: {"Low Floor Tom", 0.85, 0.4},
	43: {"High Floor Tom", 0.85, 0.45},
	45: {"Low Tom", 0.85, 0.5},
	47: {"Low-Mid Tom", 0.85, 0.55},
	48: {"Hi-Mid Tom", 0.85, 0.6},
	50: {"High Tom", 0.85, 0.65},
	
	// Hi-hats
	42: {"Closed Hi-Hat", 0.5, 1.0},
	44: {"Pedal Hi-Hat", 0.55, 0.95},
	46: {"Open Hi-Hat", 0.6, 0.9},
	
	// Cymbals
	49: {"Crash Cymbal 1", 0.9, 0.85},
	51: {"Ride Cymbal 1", 0.7, 0.75},
	52: {"Chinese Cymbal", 0.85, 0.8},
	53: {"Ride Bell", 0.65, 0.7},
	55: {"Splash Cymbal", 0.8, 0.9},
	57: {"Crash Cymbal 2", 0.9, 0.85},
	59: {"Ride Cymbal 2", 0.7, 0.75},
	
	// Percussion
	37: {"Side Stick", 0.7, 0.95},
	39: {"Hand Clap", 0.75, 0.8},
	54: {"Tambourine", 0.65, 0.85},
	56: {"Cowbell", 0.7, 0.7},
	58: {"Vibraslap", 0.7, 0.75},
	60: {"Hi Bongo", 0.75, 0.6},
	61: {"Low Bongo", 0.75, 0.5},
	62: {"Mute Hi Conga", 0.75, 0.65},
	63: {"Open Hi Conga", 0.75, 0.6},
	64: {"Low Conga", 0.75, 0.55},
	65: {"High Timbale", 0.8, 0.7},
	66: {"Low Timbale", 0.8, 0.65},
	67: {"High Agogo", 0.7, 0.8},
	68: {"Low Agogo", 0.7, 0.75},
	69: {"Cabasa", 0.65, 0.7},
	70: {"Maracas", 0.6, 0.85},
	71: {"Short Whistle", 0.6, 0.9},
	72: {"Long Whistle", 0.6, 0.85},
	73: {"Short Guiro", 0.65, 0.75},
	74: {"Long Guiro", 0.65, 0.7},
	75: {"Claves", 0.7, 0.95},
	76: {"Hi Wood Block", 0.7, 0.8},
	77: {"Low Wood Block", 0.7, 0.75},
	78: {"Mute Cuica", 0.65, 0.7},
	79: {"Open Cuica", 0.65, 0.75},
	80: {"Mute Triangle", 0.6, 0.9},
	81: {"Open Triangle", 0.6, 0.95},
}

// isDrumChannel checks if a channel is a drum channel (channel 10 in GM)
func isDrumChannel(channel uint8) bool {
	return channel == 9 // MIDI channels are 0-indexed, so channel 10 = 9
}

// getDrumMapping returns the drum mapping for a given note
func getDrumMapping(note uint8) (DrumMapping, bool) {
	mapping, ok := drumMappings[note]
	return mapping, ok
}

func main() {
	input := flag.String("input", "", "input MIDI file (required)")
	output := flag.String("output", "", "output AHAP file (default: <input>.ahap)")
	indent := flag.Bool("indent", false, "indent JSON output for readability")
	drums := flag.Bool("drums", true, "enable drum detection (channel 10 becomes transients)")
	flag.Parse()

	if *input == "" {
		fmt.Println("Usage: midi2ahap -input <file.mid> [-output <file.ahap>] [-indent] [-drums=true]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Set default output filename
	if *output == "" {
		ext := filepath.Ext(*input)
		*output = strings.TrimSuffix(*input, ext) + ".ahap"
	}

	fmt.Printf("Converting MIDI file %s to AHAP...\n", *input)

	// Read MIDI file
	smfData, err := smf.ReadFile(*input)
	if err != nil {
		log.Fatalf("Failed to read MIDI file: %v", err)
	}

	// Create AHAP
	builder := ahap.NewBuilder(fmt.Sprintf("midi file %s", filepath.Base(*input)), "midi to haptic generator")

	// Track note on/off events
	type NoteInfo struct {
		startTime float64
		velocity  uint8
	}
	noteState := make(map[uint8]NoteInfo)
	
	// Statistics
	drumCount := 0
	melodicCount := 0
	unknownDrumCount := 0

	// Get tick resolution
	var ticksPerQuarterNote uint32 = 480 // default
	if mt, isMT := smfData.TimeFormat.(smf.MetricTicks); isMT {
		ticksPerQuarterNote = uint32(mt)
	}

	// Default tempo (120 BPM = 500000 microseconds per quarter note)
	microsecondsPerQuarterNote := 500000.0

	// Process each track
	for _, track := range smfData.Tracks {
		currentTime := 0.0
		deltaAccum := uint32(0)

		for _, event := range track {
			// Update accumulated delta
			deltaAccum += event.Delta
			
			// Convert delta to seconds
			currentTime = float64(deltaAccum) / float64(ticksPerQuarterNote) * (microsecondsPerQuarterNote / 1000000.0)

			// Check for tempo changes
			var bpm float64
			if event.Message.GetMetaTempo(&bpm) {
				microsecondsPerQuarterNote = 60000000.0 / bpm
			}

			// Handle note events
			var channel, key, velocity uint8
			if event.Message.GetNoteStart(&channel, &key, &velocity) {
				// Check if this is a drum channel and drum mode is enabled
				if *drums && isDrumChannel(channel) {
					// Handle as drum - create transient event immediately
					if drumMap, ok := getDrumMapping(key); ok {
						// Scale intensity by velocity
						intensity := drumMap.Intensity * (float64(velocity) / 127.0)
						builder.Transient(currentTime).
							Intensity(intensity).
							Sharpness(drumMap.Sharpness).
							Add()
						drumCount++
					} else {
						// Unknown drum, use defaults
						intensity := float64(velocity) / 127.0
						builder.Transient(currentTime).
							Intensity(intensity).
							Sharpness(0.7).
							Add()
						drumCount++
						unknownDrumCount++
					}
				} else {
					// Note on for melodic instruments
					noteState[key] = NoteInfo{
						startTime: currentTime,
						velocity:  velocity,
					}
				}
			} else if event.Message.GetNoteEnd(&channel, &key) {
				// Note off - only for melodic instruments (or when drums disabled)
				if !(*drums && isDrumChannel(channel)) {
					if info, ok := noteState[key]; ok {
						duration := currentTime - info.startTime
						if duration > 0 {
							freq := midiNoteToFreq(key)
							sharpness, err := ahap.FreqToSharpness(freq, true)
							if err != nil {
								sharpness = 0.5 // Default if out of range
							}
							intensity := float64(info.velocity) / 127.0
							builder.Continuous(info.startTime, duration).
								Intensity(intensity).
								Sharpness(sharpness).
								Add()
							melodicCount++
						}
						delete(noteState, key)
					}
				}
			}
		}
	}

	// Export
	err = builder.Export(*output, *indent)
	if err != nil {
		log.Fatalf("Failed to export AHAP: %v", err)
	}

	fmt.Printf("Successfully created %s\n", *output)
	fmt.Printf("Conversion statistics:\n")
	fmt.Printf("  Drum events (transients): %d\n", drumCount)
	if unknownDrumCount > 0 {
		fmt.Printf("    (including %d unmapped drum notes)\n", unknownDrumCount)
	}
	fmt.Printf("  Melodic events (continuous): %d\n", melodicCount)
	fmt.Printf("  Total haptic events: %d\n", drumCount+melodicCount)
}
