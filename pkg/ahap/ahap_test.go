package ahap

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	a := New("test description", "test creator")
	if a.Version != 1.0 {
		t.Errorf("Expected version 1.0, got %f", a.Version)
	}
	if a.Metadata.Description != "test description" {
		t.Errorf("Expected description 'test description', got '%s'", a.Metadata.Description)
	}
	if a.Metadata.CreatedBy != "test creator" {
		t.Errorf("Expected creator 'test creator', got '%s'", a.Metadata.CreatedBy)
	}
	if len(a.Pattern) != 0 {
		t.Errorf("Expected empty pattern, got %d elements", len(a.Pattern))
	}
}

func TestAddHapticTransient(t *testing.T) {
	a := New("test", "test")
	a.AddHapticTransient(0.5, 1.0, 0.8)

	if len(a.Pattern) != 1 {
		t.Fatalf("Expected 1 pattern, got %d", len(a.Pattern))
	}

	event := a.Pattern[0].Event
	if event == nil {
		t.Fatal("Expected event, got nil")
	}
	if event.Time != 0.5 {
		t.Errorf("Expected time 0.5, got %f", event.Time)
	}
	if event.EventType != EventTypeHapticTransient {
		t.Errorf("Expected type HapticTransient, got %s", event.EventType)
	}
	if len(event.EventParameters) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(event.EventParameters))
	}
}

func TestAddHapticContinuous(t *testing.T) {
	a := New("test", "test")
	a.AddHapticContinuous(1.0, 2.0, 0.7, 0.6)

	if len(a.Pattern) != 1 {
		t.Fatalf("Expected 1 pattern, got %d", len(a.Pattern))
	}

	event := a.Pattern[0].Event
	if event == nil {
		t.Fatal("Expected event, got nil")
	}
	if event.EventType != EventTypeHapticContinuous {
		t.Errorf("Expected type HapticContinuous, got %s", event.EventType)
	}
	if event.EventDuration == nil || *event.EventDuration != 2.0 {
		t.Errorf("Expected duration 2.0, got %v", event.EventDuration)
	}
}

func TestCreateCurve(t *testing.T) {
	points := CreateCurve(0.0, 1.0, 0.0, 1.0, 10)
	if len(points) != 10 {
		t.Errorf("Expected 10 points, got %d", len(points))
	}
	// Check first and last points
	if points[0].Time != 0.1 {
		t.Errorf("Expected first point time 0.1, got %f", points[0].Time)
	}
	if points[9].Time != 1.0 {
		t.Errorf("Expected last point time 1.0, got %f", points[9].Time)
	}
	if points[0].ParameterValue != 0.1 {
		t.Errorf("Expected first point value 0.1, got %f", points[0].ParameterValue)
	}
	if points[9].ParameterValue != 1.0 {
		t.Errorf("Expected last point value 1.0, got %f", points[9].ParameterValue)
	}
}

func TestFreqToSharpness(t *testing.T) {
	tests := []struct {
		freq      float64
		normalize bool
		wantErr   bool
	}{
		{80, false, false},
		{155, false, false},
		{230, false, false},
		{79, false, true},
		{231, false, true},
		{50, true, false},  // normalized to 80
		{300, true, false}, // normalized to 230
	}

	for _, tt := range tests {
		result, err := FreqToSharpness(tt.freq, tt.normalize)
		if (err != nil) != tt.wantErr {
			t.Errorf("FreqToSharpness(%f, %v) error = %v, wantErr %v", tt.freq, tt.normalize, err, tt.wantErr)
		}
		if !tt.wantErr && (result < 0 || result > 1) {
			t.Errorf("FreqToSharpness(%f, %v) = %f, want value between 0 and 1", tt.freq, tt.normalize, result)
		}
	}
}

func TestMusicalContext(t *testing.T) {
	mc := NewMusicalContext(120, 4, 4)

	// At 120 BPM, one beat = 0.5 seconds
	beatDur := mc.BeatDuration()
	if beatDur != 0.5 {
		t.Errorf("Expected beat duration 0.5, got %f", beatDur)
	}

	// One bar = 4 beats = 2 seconds
	barDur := mc.BarDuration()
	if barDur != 2.0 {
		t.Errorf("Expected bar duration 2.0, got %f", barDur)
	}

	// Convert beats to seconds
	beatTime := mc.BeatToSeconds(4)
	if beatTime != 2.0 {
		t.Errorf("Expected 4 beats = 2.0 seconds, got %f", beatTime)
	}

	// Convert bars to seconds
	barTime := mc.BarToSeconds(2)
	if barTime != 4.0 {
		t.Errorf("Expected 2 bars = 4.0 seconds, got %f", barTime)
	}
}

func TestBuilder(t *testing.T) {
	b := NewBuilder("test", "test creator")
	ahap := b.
		Transient(0.0).Intensity(1.0).Sharpness(0.5).Add().
		Continuous(1.0, 2.0).Intensity(0.8).Sharpness(0.7).Add().
		Build()

	if len(ahap.Pattern) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(ahap.Pattern))
	}
}

func TestBuilderMusical(t *testing.T) {
	b := NewBuilder("musical test", "test creator").
		WithBPM(120).
		WithTimeSignature(4, 4)

	ahap := b.
		AtBeat(0).Transient().Intensity(1.0).Add().
		AtBeat(1).Transient().Intensity(0.8).Add().
		AtBar(1).Continuous(1.0).Intensity(0.5).Add().
		Build()

	if len(ahap.Pattern) != 3 {
		t.Errorf("Expected 3 patterns, got %d", len(ahap.Pattern))
	}

	// First event should be at time 0
	if ahap.Pattern[0].Event.Time != 0.0 {
		t.Errorf("Expected first event at time 0.0, got %f", ahap.Pattern[0].Event.Time)
	}

	// Second event should be at time 0.5 (beat 1 at 120 BPM)
	if ahap.Pattern[1].Event.Time != 0.5 {
		t.Errorf("Expected second event at time 0.5, got %f", ahap.Pattern[1].Event.Time)
	}

	// Third event should be at time 2.0 (bar 1 at 120 BPM in 4/4)
	if ahap.Pattern[2].Event.Time != 2.0 {
		t.Errorf("Expected third event at time 2.0, got %f", ahap.Pattern[2].Event.Time)
	}
}

func TestBuilderCurve(t *testing.T) {
	b := NewBuilder("curve test", "test creator")
	ahap := b.
		Continuous(0.0, 2.0).Intensity(0.5).Sharpness(0.5).Add().
		Curve(CurveHapticSharpness).At(0.0).From(0.0, 0.3).To(2.0, 0.8).Steps(10).Add().
		Build()

	if len(ahap.Pattern) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(ahap.Pattern))
	}

	curve := ahap.Pattern[1].ParameterCurve
	if curve == nil {
		t.Fatal("Expected parameter curve, got nil")
	}
	if len(curve.ParameterCurveControlPoints) != 10 {
		t.Errorf("Expected 10 control points, got %d", len(curve.ParameterCurveControlPoints))
	}
}

func TestExport(t *testing.T) {
	a := New("test export", "test creator")
	a.AddHapticTransient(0.0, 1.0, 0.5)

	tmpFile := "/tmp/test_export.ahap"
	defer os.Remove(tmpFile)

	err := a.Export(tmpFile, true)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Read back and verify
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var decoded AHAP
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to decode exported JSON: %v", err)
	}

	if decoded.Version != 1.0 {
		t.Errorf("Expected version 1.0, got %f", decoded.Version)
	}
	if len(decoded.Pattern) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(decoded.Pattern))
	}
}
