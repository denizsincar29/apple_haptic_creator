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

func main() {
	input := flag.String("input", "", "input MIDI file (required)")
	output := flag.String("output", "", "output AHAP file (default: <input>.ahap)")
	indent := flag.Bool("indent", false, "indent JSON output for readability")
	flag.Parse()

	if *input == "" {
		fmt.Println("Usage: midi2ahap -input <file.mid> [-output <file.ahap>] [-indent]")
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
				// Note on
				noteState[key] = NoteInfo{
					startTime: currentTime,
					velocity:  velocity,
				}
			} else if event.Message.GetNoteEnd(&channel, &key) {
				// Note off
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
					}
					delete(noteState, key)
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
}
