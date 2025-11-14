package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/denizsincar29/apple_haptic_creator/pkg/ahap"
)

func main() {
	output := flag.String("output", "bike.ahap", "output filename")
	indent := flag.Bool("indent", false, "indent JSON output for readability")
	flag.Parse()

	fmt.Println("Creating motorcycle sound haptic pattern...")

	// Create the AHAP with builder
	builder := ahap.NewBuilder("bike sound", "Deniz Sincar")

	// Initial rumble
	time := 0.0
	dur := 0.4
	builder.Continuous(time, dur).Intensity(0.5).Sharpness(0.4).Add()
	builder.Curve(ahap.CurveHapticSharpness).At(time).
		From(0.0, 0.4).To(0.4, 0.75).Steps(10).Add()

	// Series of quick transients (gear shift)
	time = 0.45
	for i := 0; i < 7; i++ {
		builder.Transient(time).Intensity(1.0).Sharpness(0.3).Add()
		time += 0.05
	}

	// Main engine running (15 seconds of continuous vibration with rapid transients)
	builder.Continuous(time, 15.0).Intensity(0.75).Sharpness(0.0).Add()
	
	// 300 rapid transients to simulate engine vibration
	for i := 0; i < 300; i++ {
		builder.Transient(time + float64(i)*0.05).Intensity(1.0).Sharpness(1.0).Add()
	}

	// Sharpness curves for engine acceleration
	builder.Curve(ahap.CurveHapticSharpness).At(time).
		From(0.0, 0.0).To(0.4, 0.75).Steps(10).Add()
	time += 0.4

	builder.Curve(ahap.CurveHapticSharpness).At(time).
		From(0.0, 0.75).To(0.8, 0.2).Steps(10).Add()
	time += 0.8

	builder.Curve(ahap.CurveHapticSharpness).At(time).
		From(0.0, 0.0).To(3.0, 0.5).Steps(10).Add()

	builder.Curve(ahap.CurveHapticSharpness).At(time + 3).
		From(0.0, 0.2).To(3.0, 0.65).Steps(10).Add()

	builder.Curve(ahap.CurveHapticSharpness).At(time + 6).
		From(0.0, 0.4).To(4.0, 1.0).Steps(10).Add()

	builder.Curve(ahap.CurveHapticSharpness).At(time + 10).
		From(0.0, 1.0).To(2.0, 0.0).Steps(10).Add()

	// Export
	err := builder.Export(*output, *indent)
	if err != nil {
		log.Fatalf("Failed to export AHAP: %v", err)
	}

	fmt.Printf("Successfully created %s\n", *output)
}
