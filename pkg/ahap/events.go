package ahap

// AddHapticTransient adds a haptic transient event
func (a *AHAP) AddHapticTransient(time, intensity, sharpness float64) {
	event := &Event{
		Time:      time,
		EventType: EventTypeHapticTransient,
		EventParameters: []EventParameter{
			{ParameterID: ParamHapticIntensity, ParameterValue: intensity},
			{ParameterID: ParamHapticSharpness, ParameterValue: sharpness},
		},
	}
	a.AddEvent(event)
}

// AddHapticContinuous adds a haptic continuous event
func (a *AHAP) AddHapticContinuous(time, duration, intensity, sharpness float64) {
	event := &Event{
		Time:      time,
		EventType: EventTypeHapticContinuous,
		EventParameters: []EventParameter{
			{ParameterID: ParamHapticIntensity, ParameterValue: intensity},
			{ParameterID: ParamHapticSharpness, ParameterValue: sharpness},
		},
		EventDuration: &duration,
	}
	a.AddEvent(event)
}

// AddAudioCustom adds a custom audio event with a WAV file
func (a *AHAP) AddAudioCustom(time float64, wavFilepath string, volume float64) {
	event := &Event{
		Time:      time,
		EventType: EventTypeAudioCustom,
		EventParameters: []EventParameter{
			{ParameterID: ParamAudioVolume, ParameterValue: volume},
		},
		EventWaveformPath: &wavFilepath,
	}
	a.AddEvent(event)
}

// AddAudioContinuous adds an audio continuous event
func (a *AHAP) AddAudioContinuous(time, duration, volume float64) {
	event := &Event{
		Time:      time,
		EventType: EventTypeAudioContinuous,
		EventParameters: []EventParameter{
			{ParameterID: ParamAudioVolume, ParameterValue: volume},
		},
		EventDuration: &duration,
	}
	a.AddEvent(event)
}
