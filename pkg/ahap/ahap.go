package ahap

import (
	"encoding/json"
	"os"
	"time"
)

// AHAP represents a complete Apple Haptic pattern
type AHAP struct {
	Version  float64   `json:"Version"`
	Metadata Metadata  `json:"Metadata"`
	Pattern  []Pattern `json:"Pattern"`
}

// Metadata contains file metadata
type Metadata struct {
	Project     string `json:"Project"`
	Created     string `json:"Created"`
	Description string `json:"Description"`
	CreatedBy   string `json:"Created By"`
}

// Pattern can be either an Event or a ParameterCurve
type Pattern struct {
	Event          *Event          `json:"Event,omitempty"`
	ParameterCurve *ParameterCurve `json:"ParameterCurve,omitempty"`
}

// Event represents a haptic or audio event
type Event struct {
	Time              float64          `json:"Time"`
	EventType         string           `json:"EventType"`
	EventParameters   []EventParameter `json:"EventParameters"`
	EventDuration     *float64         `json:"EventDuration,omitempty"`
	EventWaveformPath *string          `json:"EventWaveformPath,omitempty"`
}

// EventParameter represents a parameter of an event
type EventParameter struct {
	ParameterID    string  `json:"ParameterID"`
	ParameterValue float64 `json:"ParameterValue"`
}

// ParameterCurve represents a dynamic parameter change over time
type ParameterCurve struct {
	ParameterID                 string         `json:"ParameterID"`
	Time                        float64        `json:"Time"`
	ParameterCurveControlPoints []ControlPoint `json:"ParameterCurveControlPoints"`
}

// ControlPoint represents a point in a parameter curve
type ControlPoint struct {
	Time           float64 `json:"Time"`
	ParameterValue float64 `json:"ParameterValue"`
}

// EventType constants
const (
	EventTypeHapticTransient  = "HapticTransient"
	EventTypeHapticContinuous = "HapticContinuous"
	EventTypeAudioCustom      = "AudioCustom"
	EventTypeAudioContinuous  = "AudioContinuous"
)

// ParameterID constants for event parameters
const (
	ParamHapticIntensity   = "HapticIntensity"
	ParamHapticSharpness   = "HapticSharpness"
	ParamHapticAttackTime  = "HapticAttackTime"
	ParamHapticDecayTime   = "HapticDecayTime"
	ParamHapticReleaseTime = "HapticReleaseTime"
	ParamAudioVolume       = "AudioVolume"
	ParamAudioPitch        = "AudioPitch"
	ParamAudioPan          = "AudioPan"
	ParamAudioBrightness   = "AudioBrightness"
	ParamAudioAttackTime   = "AudioAttackTime"
	ParamAudioDecayTime    = "AudioDecayTime"
	ParamAudioReleaseTime  = "AudioReleaseTime"
)

// CurveParameterID constants for parameter curves
const (
	CurveHapticIntensity   = "HapticIntensityControl"
	CurveHapticSharpness   = "HapticSharpnessControl"
	CurveHapticAttackTime  = "HapticAttackTimeControl"
	CurveHapticDecayTime   = "HapticDecayTimeControl"
	CurveHapticReleaseTime = "HapticReleaseTimeControl"
	CurveAudioVolume       = "AudioVolumeControl"
	CurveAudioPitch        = "AudioPitchControl"
	CurveAudioPan          = "AudioPanControl"
	CurveAudioBrightness   = "AudioBrightnessControl"
	CurveAudioAttackTime   = "AudioAttackTimeControl"
	CurveAudioDecayTime    = "AudioDecayTimeControl"
	CurveAudioReleaseTime  = "AudioReleaseTimeControl"
)

// New creates a new AHAP with default metadata
func New(description, createdBy string) *AHAP {
	return &AHAP{
		Version: 1.0,
		Metadata: Metadata{
			Project:     "Basis",
			Created:     time.Now().Format("2006-01-02 15:04:05.000000"),
			Description: description,
			CreatedBy:   createdBy,
		},
		Pattern: make([]Pattern, 0),
	}
}

// AddEvent adds a raw event to the pattern
func (a *AHAP) AddEvent(event *Event) {
	a.Pattern = append(a.Pattern, Pattern{Event: event})
}

// AddParameterCurve adds a parameter curve to the pattern
func (a *AHAP) AddParameterCurve(curve *ParameterCurve) {
	a.Pattern = append(a.Pattern, Pattern{ParameterCurve: curve})
}

// Export writes the AHAP to a file
func (a *AHAP) Export(filename string, indent bool) error {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(a, "", "  ")
	} else {
		data, err = json.Marshal(a)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ToJSON returns the AHAP as a JSON string
func (a *AHAP) ToJSON(indent bool) (string, error) {
	var data []byte
	var err error

	if indent {
		data, err = json.MarshalIndent(a, "", "  ")
	} else {
		data, err = json.Marshal(a)
	}

	if err != nil {
		return "", err
	}

	return string(data), nil
}
