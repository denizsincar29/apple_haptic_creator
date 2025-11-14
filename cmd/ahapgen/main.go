package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/denizsincar29/apple_haptic_creator/pkg/ahap"
)

func main() {
	output := flag.String("o", "output.ahap", "output filename")
	indent := flag.Bool("indent", true, "indent JSON output for readability")
	description := flag.String("desc", "Custom haptic pattern", "pattern description")
	creator := flag.String("creator", "AHAP Generator", "creator name")
	bpm := flag.Float64("bpm", 0, "beats per minute for musical timing (0 = disabled)")
	timeSignature := flag.String("time", "4/4", "time signature (e.g., 4/4, 3/4, 6/8)")
	flag.Parse()

	builder := ahap.NewBuilder(*description, *creator)

	// Parse time signature and set up musical context if BPM is specified
	if *bpm > 0 {
		parts := strings.Split(*timeSignature, "/")
		if len(parts) == 2 {
			numerator, _ := strconv.Atoi(parts[0])
			denominator, _ := strconv.Atoi(parts[1])
			builder.WithBPM(*bpm).WithTimeSignature(numerator, denominator)
			fmt.Printf("Musical timing enabled: %.1f BPM, %s time\n", *bpm, *timeSignature)
		}
	}

	fmt.Println("AHAP Generator - Interactive Mode")
	fmt.Println("Commands:")
	fmt.Println("  t <time> <intensity> <sharpness>       - Add transient event")
	fmt.Println("  c <time> <duration> <intensity> <sharpness> - Add continuous event")
	fmt.Println("  beat <beat> <intensity> <sharpness>    - Add transient at beat (requires -bpm)")
	fmt.Println("  bar <bar> <intensity> <sharpness>      - Add transient at bar (requires -bpm)")
	fmt.Println("  export                                 - Export to file and exit")
	fmt.Println("  quit                                   - Exit without saving")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	eventCount := 0

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]

		switch cmd {
		case "t", "transient":
			if len(parts) != 4 {
				fmt.Println("Usage: t <time> <intensity> <sharpness>")
				continue
			}
			time, _ := strconv.ParseFloat(parts[1], 64)
			intensity, _ := strconv.ParseFloat(parts[2], 64)
			sharpness, _ := strconv.ParseFloat(parts[3], 64)
			builder.Transient(time).Intensity(intensity).Sharpness(sharpness).Add()
			eventCount++
			fmt.Printf("Added transient event at %.2fs (total: %d events)\n", time, eventCount)

		case "c", "continuous":
			if len(parts) != 5 {
				fmt.Println("Usage: c <time> <duration> <intensity> <sharpness>")
				continue
			}
			time, _ := strconv.ParseFloat(parts[1], 64)
			duration, _ := strconv.ParseFloat(parts[2], 64)
			intensity, _ := strconv.ParseFloat(parts[3], 64)
			sharpness, _ := strconv.ParseFloat(parts[4], 64)
			builder.Continuous(time, duration).Intensity(intensity).Sharpness(sharpness).Add()
			eventCount++
			fmt.Printf("Added continuous event at %.2fs for %.2fs (total: %d events)\n", time, duration, eventCount)

		case "beat":
			if *bpm == 0 {
				fmt.Println("Musical timing not enabled. Use -bpm flag.")
				continue
			}
			if len(parts) != 4 {
				fmt.Println("Usage: beat <beat> <intensity> <sharpness>")
				continue
			}
			beat, _ := strconv.Atoi(parts[1])
			intensity, _ := strconv.ParseFloat(parts[2], 64)
			sharpness, _ := strconv.ParseFloat(parts[3], 64)
			// Calculate bar and beat from absolute beat number
			beatsPerBar := builder.GetBeatsPerBar()
			bar := beat / beatsPerBar
			beatInBar := beat % beatsPerBar
			builder.At(bar, beatInBar).Transient().Intensity(intensity).Sharpness(sharpness).Add()
			eventCount++
			fmt.Printf("Added transient at beat %d (bar %d, beat %d) (total: %d events)\n", beat, bar, beatInBar, eventCount)

		case "bar":
			if *bpm == 0 {
				fmt.Println("Musical timing not enabled. Use -bpm flag.")
				continue
			}
			if len(parts) != 4 {
				fmt.Println("Usage: bar <bar> <intensity> <sharpness>")
				continue
			}
			bar, _ := strconv.Atoi(parts[1])
			intensity, _ := strconv.ParseFloat(parts[2], 64)
			sharpness, _ := strconv.ParseFloat(parts[3], 64)
			builder.At(bar, 0).Transient().Intensity(intensity).Sharpness(sharpness).Add()
			eventCount++
			fmt.Printf("Added transient at bar %d (total: %d events)\n", bar, eventCount)

		case "export", "save":
			err := builder.Export(*output, *indent)
			if err != nil {
				fmt.Printf("Error exporting: %v\n", err)
				continue
			}
			fmt.Printf("Successfully exported %d events to %s\n", eventCount, *output)
			return

		case "quit", "exit", "q":
			fmt.Println("Exiting without saving.")
			return

		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}
