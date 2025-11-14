package main

import (
	"fmt"
	"log"

	"github.com/denizsincar29/apple_haptic_creator/pkg/ahap"
)

func main() {
	// Example 1: Simple transient and continuous events
	fmt.Println("Creating simple haptic pattern...")
	builder := ahap.NewBuilder("Simple Example", "Go Developer")
	
	// Add a sharp tap at the start
	builder.Transient(0.0).Intensity(1.0).Sharpness(1.0).Add()
	
	// Add a gentle continuous haptic from 0.5s to 2.5s
	builder.Continuous(0.5, 2.0).Intensity(0.5).Sharpness(0.3).Add()
	
	// Add another sharp tap at 3.0s
	builder.Transient(3.0).Intensity(1.0).Sharpness(1.0).Add()
	
	if err := builder.Export("simple.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created simple.ahap")

	// Example 2: Musical timing
	fmt.Println("\nCreating musical haptic pattern...")
	musical := ahap.NewBuilder("Musical Example", "Go Developer").
		WithBPM(120).
		WithTimeSignature(4, 4)
	
	// Create a 4-beat pattern (bar 0, beats 0-3)
	for beat := 0; beat < 4; beat++ {
		intensity := 1.0
		if beat == 0 {
			intensity = 1.0 // Accent on first beat
		} else {
			intensity = 0.7
		}
		musical.At(0, beat).
			Transient().
			Intensity(intensity).
			Sharpness(0.5).
			Add()
	}
	
	if err := musical.Export("musical.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created musical.ahap")

	// Example 3: Parameter curves
	fmt.Println("\nCreating haptic with curves...")
	curves := ahap.NewBuilder("Curve Example", "Go Developer")
	
	// Add continuous event
	curves.Continuous(0.0, 3.0).Intensity(0.5).Sharpness(0.5).Add()
	
	// Add curve that gradually increases sharpness
	curves.Curve(ahap.CurveHapticSharpness).
		From(0.0, 0.2).To(3.0, 0.9).Steps(20).Add()
	
	// Add curve that varies intensity with ease-in-out
	curves.Curve(ahap.CurveHapticIntensity).
		From(0.0, 0.3).To(3.0, 1.0).EaseInOut(15).Add()
	
	if err := curves.Export("curves.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created curves.ahap")

	// Example 4: Frequency to sharpness conversion
	fmt.Println("\nFrequency to sharpness examples:")
	frequencies := []float64{80, 100, 155, 200, 230}
	for _, freq := range frequencies {
		sharpness, err := ahap.FreqToSharpness(freq, false)
		if err != nil {
			log.Printf("Error converting %d Hz: %v", int(freq), err)
			continue
		}
		fmt.Printf("  %d Hz → sharpness %.3f\n", int(freq), sharpness)
	}

	fmt.Println("\n✓ All examples created successfully!")
	fmt.Println("\nGenerated files:")
	fmt.Println("  - simple.ahap    : Basic transient and continuous events")
	fmt.Println("  - musical.ahap   : Musical timing with BPM")
	fmt.Println("  - curves.ahap    : Dynamic parameter curves")
}
