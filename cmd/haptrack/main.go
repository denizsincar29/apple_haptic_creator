package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/denizsincar29/apple_haptic_creator/pkg/ahap"
)

// HapticDefinition defines what a letter represents
type HapticDefinition struct {
	Name      string
	Intensity float64
	Sharpness float64
	HasCurve  bool
	CurveDown bool // if true, curve goes down
	CurveDur  float64
}

// HaptrackParser parses and executes haptrack DSL
type HaptrackParser struct {
	definitions map[rune]HapticDefinition
	bpm         float64
	numerator   int
	denominator int
	builder     *ahap.Builder
}

// NewHaptrackParser creates a new parser
func NewHaptrackParser() *HaptrackParser {
	return &HaptrackParser{
		definitions: make(map[rune]HapticDefinition),
		bpm:         120,
		numerator:   4,
		denominator: 4,
	}
}

// ParseFile parses a haptrack file
func (p *HaptrackParser) ParseFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDefinitions := true
	trackNumber := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for begin marker
		if strings.ToLower(line) == "begin" {
			inDefinitions = false
			// Initialize builder now that we have all settings
			p.builder = ahap.NewBuilder("Haptrack Pattern", "Haptrack DSL").
				WithBPM(p.bpm).
				WithTimeSignature(p.numerator, p.denominator)
			continue
		}

		if inDefinitions {
			if err := p.parseDefinitionLine(line); err != nil {
				return fmt.Errorf("error parsing definition: %v", err)
			}
		} else {
			// Parse track line
			if strings.HasPrefix(strings.ToLower(line), "track") {
				trackNumber++
				continue
			}
			// Parse pattern
			if err := p.parseTrack(line, trackNumber); err != nil {
				return fmt.Errorf("error parsing track %d: %v", trackNumber, err)
			}
		}
	}

	return scanner.Err()
}

// parseDefinitionLine parses a definition line
func (p *HaptrackParser) parseDefinitionLine(line string) error {
	// Format: letter = name, intensity, sharpness [, curve_down, duration_ms]
	// Example: s = snare, 1.0, 0.9, down, 60
	// Or: bpm = 120
	// Or: time = 4/4

	if strings.Contains(line, "=") {
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Handle special keys
		if key == "bpm" {
			bpm, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid BPM: %v", err)
			}
			p.bpm = bpm
			return nil
		}

		if key == "time" {
			// Parse time signature like "4/4"
			timeParts := strings.Split(value, "/")
			if len(timeParts) != 2 {
				return fmt.Errorf("invalid time signature: %s", value)
			}
			num, err := strconv.Atoi(timeParts[0])
			if err != nil {
				return fmt.Errorf("invalid time signature numerator: %v", err)
			}
			denom, err := strconv.Atoi(timeParts[1])
			if err != nil {
				return fmt.Errorf("invalid time signature denominator: %v", err)
			}
			p.numerator = num
			p.denominator = denom
			return nil
		}

		// Parse letter definition
		if len(key) == 1 {
			letter := rune(key[0])
			def, err := p.parseHapticDefinition(value)
			if err != nil {
				return err
			}
			p.definitions[letter] = def
		}
	}

	return nil
}

// parseHapticDefinition parses a haptic definition
func (p *HaptrackParser) parseHapticDefinition(value string) (HapticDefinition, error) {
	def := HapticDefinition{}
	parts := strings.Split(value, ",")
	
	if len(parts) < 3 {
		return def, fmt.Errorf("definition needs at least name, intensity, sharpness")
	}

	def.Name = strings.TrimSpace(parts[0])
	
	intensity, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return def, fmt.Errorf("invalid intensity: %v", err)
	}
	def.Intensity = intensity

	sharpness, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return def, fmt.Errorf("invalid sharpness: %v", err)
	}
	def.Sharpness = sharpness

	// Optional curve
	if len(parts) >= 5 {
		curveDir := strings.TrimSpace(strings.ToLower(parts[3]))
		if curveDir == "down" || curveDir == "up" {
			def.HasCurve = true
			def.CurveDown = (curveDir == "down")
			
			durMs, err := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
			if err != nil {
				return def, fmt.Errorf("invalid curve duration: %v", err)
			}
			def.CurveDur = durMs / 1000.0 // Convert ms to seconds
		}
	}

	return def, nil
}

// parseTrack parses a track pattern
func (p *HaptrackParser) parseTrack(pattern string, trackNum int) error {
	currentBeat := 0.0
	i := 0

	for i < len(pattern) {
		char := rune(pattern[i])

		if char == '-' {
			// Rest - check for duration
			i++
			if i < len(pattern) && isDigit(rune(pattern[i])) {
				duration := parseNoteDuration(pattern, &i)
				currentBeat += beatDuration(duration, p.denominator)
			} else {
				currentBeat += beatDuration(8, p.denominator) // default eighth note
			}
		} else if def, ok := p.definitions[char]; ok {
			// Found a defined haptic
			i++
			
			// Check for duration
			duration := 8 // default eighth note
			if i < len(pattern) && isDigit(rune(pattern[i])) {
				duration = parseNoteDuration(pattern, &i)
			}

			// Add the event using At(bar, beat)
			// Calculate bar and beat from total beats
			beatsPerBar := p.builder.GetBeatsPerBar()
			bar := int(currentBeat) / beatsPerBar
			beat := int(currentBeat) % beatsPerBar
			p.builder.At(bar, beat).Transient().
				Intensity(def.Intensity).
				Sharpness(def.Sharpness).
				Add()

			// Add curve if defined
			if def.HasCurve {
				startSharp := def.Sharpness
				endSharp := def.Sharpness
				if def.CurveDown {
					endSharp = def.Sharpness * 0.3
				} else {
					endSharp = def.Sharpness * 1.5
					if endSharp > 1.0 {
						endSharp = 1.0
					}
				}
				p.builder.Curve(ahap.CurveHapticSharpness).
					From(0, startSharp).
					To(def.CurveDur, endSharp).
					Steps(5).
					Add()
			}

			currentBeat += beatDuration(duration, p.denominator)
		} else {
			// Unknown character, skip
			i++
		}
	}

	return nil
}

// Helper functions

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func parseNoteDuration(pattern string, index *int) int {
	numStr := ""
	for *index < len(pattern) && isDigit(rune(pattern[*index])) {
		numStr += string(pattern[*index])
		*index++
	}
	duration, _ := strconv.Atoi(numStr)
	return duration
}

// beatDuration calculates beat duration based on note value
// 1 = whole note = 4 beats
// 2 = half note = 2 beats
// 4 = quarter note = 1 beat
// 8 = eighth note = 0.5 beats
// 16 = sixteenth note = 0.25 beats
func beatDuration(noteValue, denominator int) float64 {
	// denominator tells us what note value gets one beat
	// noteValue tells us what note we're playing
	// E.g., in 4/4 time, quarter note (4) = 1 beat
	// In 4/4, eighth note (8) = 0.5 beat
	
	quarterNotesInWhole := 4.0
	return (quarterNotesInWhole / float64(noteValue)) * (float64(denominator) / 4.0)
}

func main() {
	input := flag.String("input", "", "input haptrack file (required)")
	output := flag.String("output", "output.ahap", "output AHAP file")
	indent := flag.Bool("indent", true, "indent JSON output")
	flag.Parse()

	if *input == "" {
		fmt.Println("Haptrack - Haptic Pattern DSL Compiler")
		fmt.Println("\nUsage: haptrack -input <file.hap> [-output <file.ahap>]")
		fmt.Println("\nHaptrack file format:")
		fmt.Println("  # Comments start with #")
		fmt.Println("  bpm = 120")
		fmt.Println("  time = 4/4")
		fmt.Println("  s = snare, 1.0, 0.9, down, 60")
		fmt.Println("  k = kick, 1.0, 0.2")
		fmt.Println("  h = hihat, 0.6, 1.0")
		fmt.Println("")
		fmt.Println("  begin")
		fmt.Println("  track1")
		fmt.Println("  k8k8s8k8k8k8s8k8")
		fmt.Println("  track2")
		fmt.Println("  h8h8h8h8h8h8h8h8")
		fmt.Println("\nNote durations: 1=whole, 2=half, 4=quarter, 8=eighth, 16=sixteenth")
		fmt.Println("Rest: - (dash)")
		fmt.Println("Example: s8-8 means snare eighth note, then rest for eighth note")
		flag.PrintDefaults()
		os.Exit(1)
	}

	parser := NewHaptrackParser()
	
	fmt.Printf("Parsing haptrack file: %s\n", *input)
	if err := parser.ParseFile(*input); err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}

	fmt.Printf("Found %d haptic definitions\n", len(parser.definitions))
	fmt.Printf("BPM: %.0f, Time Signature: %d/%d\n", parser.bpm, parser.numerator, parser.denominator)

	if parser.builder == nil {
		log.Fatal("No tracks found in file (missing 'begin' marker?)")
	}

	if err := parser.builder.Export(*output, *indent); err != nil {
		log.Fatalf("Error exporting AHAP: %v", err)
	}

	fmt.Printf("Successfully created %s\n", *output)
}
