package ahap

// Beat represents a single beat in musical time
type Beat float64

// Bar represents a bar/measure in musical time
type Bar float64

// TimeSignature represents musical time signature
type TimeSignature struct {
	Numerator   int // Beats per bar
	Denominator int // Note value (4 = quarter note)
}

// MusicalContext provides musical timing functionality
type MusicalContext struct {
	BPM           float64
	TimeSignature TimeSignature
}

// NewMusicalContext creates a new musical context
func NewMusicalContext(bpm float64, numerator, denominator int) *MusicalContext {
	return &MusicalContext{
		BPM: bpm,
		TimeSignature: TimeSignature{
			Numerator:   numerator,
			Denominator: denominator,
		},
	}
}

// BeatToSeconds converts beat to seconds
func (m *MusicalContext) BeatToSeconds(beat Beat) float64 {
	// 60 seconds per minute / BPM = seconds per beat
	return float64(beat) * (60.0 / m.BPM)
}

// BarToSeconds converts bar to seconds
func (m *MusicalContext) BarToSeconds(bar Bar) float64 {
	beatsPerBar := float64(m.TimeSignature.Numerator)
	return float64(bar) * beatsPerBar * (60.0 / m.BPM)
}

// BeatDuration returns the duration of one beat in seconds
func (m *MusicalContext) BeatDuration() float64 {
	return m.BeatToSeconds(1)
}

// BarDuration returns the duration of one bar in seconds
func (m *MusicalContext) BarDuration() float64 {
	return m.BarToSeconds(1)
}

// BeatsPerBar returns the number of beats per bar
func (m *MusicalContext) BeatsPerBar() int {
	return m.TimeSignature.Numerator
}

// SecondsToBeats converts seconds to beats
func (m *MusicalContext) SecondsToBeats(seconds float64) Beat {
	return Beat(seconds / (60.0 / m.BPM))
}

// SecondsToBars converts seconds to bars
func (m *MusicalContext) SecondsToBars(seconds float64) Bar {
	beatsPerBar := float64(m.TimeSignature.Numerator)
	return Bar(seconds / (beatsPerBar * (60.0 / m.BPM)))
}
