package main

import (
	"fmt"
	"log"

	"github.com/denizsincar29/apple_haptic_creator/pkg/ahap"
)

func main() {
	fmt.Println("Creating sequence pattern examples...")

	// Example 1: Transients on beats 1 and 3 from bars 5-8
	builder1 := ahap.NewBuilder("Sequence Example 1", "Go Developer").
		WithBPM(120).
		WithTimeSignature(4, 4)

	// Add transients on beats 0 and 2 (1st and 3rd in musical terms) from bars 5 to 8
	builder1.Sequence().TransientsOnBeats([]ahap.Beat{0, 2}, 5, 8, 1.0, 0.5)

	if err := builder1.Export("sequence1.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created sequence1.ahap - transients on beats 1 and 3, bars 5-8")

	// Example 2: Different pattern on beats 2 and 4 with different sharpness
	builder2 := ahap.NewBuilder("Sequence Example 2", "Go Developer").
		WithBPM(120).
		WithTimeSignature(4, 4)

	// Kick drum on beats 0 and 2
	builder2.Sequence().TransientsOnBeats([]ahap.Beat{0, 2}, 0, 3, 1.0, 0.2)
	// Snare on beats 1 and 3 with different sharpness
	builder2.Sequence().TransientsOnBeats([]ahap.Beat{1, 3}, 0, 3, 0.9, 0.8)

	if err := builder2.Export("sequence2.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created sequence2.ahap - kick and snare pattern")

	// Example 3: Every beat pattern
	builder3 := ahap.NewBuilder("Every Beat", "Go Developer").
		WithBPM(140).
		WithTimeSignature(4, 4)

	builder3.Sequence().EveryBeat(0, 2, 0.7, 0.6)

	if err := builder3.Export("sequence3.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created sequence3.ahap - every beat for 3 bars")

	// Example 4: Every 2nd beat (on-beat pattern)
	builder4 := ahap.NewBuilder("Every 2nd Beat", "Go Developer").
		WithBPM(120).
		WithTimeSignature(4, 4)

	builder4.Sequence().EveryNthBeat(2, 0, 3, 0.8, 0.5)

	if err := builder4.Export("sequence4.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created sequence4.ahap - every 2nd beat")

	// Example 5: Custom pattern using Pattern function
	builder5 := ahap.NewBuilder("Custom Pattern", "Go Developer").
		WithBPM(120).
		WithTimeSignature(4, 4)

	// Create a pattern where each bar has a different feel
	builder5.Sequence().Pattern(0, 3, func(b *ahap.Builder, bar int) {
		// First beat of each bar is always strong (beat 0)
		barStart := float64(bar) * 4 // 4 beats per bar at 120 BPM
		b.AtBeat(ahap.Beat(barStart)).Transient().Intensity(1.0).Sharpness(0.8).Add()
		
		// Other beats vary by bar
		if bar%2 == 0 {
			// Even bars: add on beat 2
			b.AtBeat(ahap.Beat(barStart + 2)).Transient().Intensity(0.7).Sharpness(0.5).Add()
		} else {
			// Odd bars: add on beats 1 and 3
			b.AtBeat(ahap.Beat(barStart + 1)).Transient().Intensity(0.6).Sharpness(0.6).Add()
			b.AtBeat(ahap.Beat(barStart + 3)).Transient().Intensity(0.6).Sharpness(0.6).Add()
		}
	})

	if err := builder5.Export("sequence5.ahap", true); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Created sequence5.ahap - custom pattern with function")

	fmt.Println("\n✓ All sequence examples created successfully!")
}
