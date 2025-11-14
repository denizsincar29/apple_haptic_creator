package ahap

// Builder provides a fluent API for creating AHAP files
type Builder struct {
	ahap    *AHAP
	musical *MusicalContext
}

// NewBuilder creates a new AHAP builder
func NewBuilder(description, creator string) *Builder {
	return &Builder{
		ahap:    New(description, creator),
		musical: nil,
	}
}

// WithBPM sets the BPM for musical timing
func (b *Builder) WithBPM(bpm float64) *Builder {
	if b.musical == nil {
		b.musical = NewMusicalContext(bpm, 4, 4) // Default 4/4 time
	} else {
		b.musical.BPM = bpm
	}
	return b
}

// WithTimeSignature sets the time signature
func (b *Builder) WithTimeSignature(numerator, denominator int) *Builder {
	if b.musical == nil {
		b.musical = NewMusicalContext(120, numerator, denominator) // Default 120 BPM
	} else {
		b.musical.TimeSignature = TimeSignature{
			Numerator:   numerator,
			Denominator: denominator,
		}
	}
	return b
}

// WithMusicalContext sets the complete musical context
func (b *Builder) WithMusicalContext(musical *MusicalContext) *Builder {
	b.musical = musical
	return b
}

// Transient creates a transient event builder
func (b *Builder) Transient(time float64) *TransientBuilder {
	return &TransientBuilder{
		builder:   b,
		time:      time,
		intensity: 0.5,
		sharpness: 0.5,
	}
}

// Continuous creates a continuous event builder
func (b *Builder) Continuous(time, duration float64) *ContinuousBuilder {
	return &ContinuousBuilder{
		builder:   b,
		time:      time,
		duration:  duration,
		intensity: 0.5,
		sharpness: 0.5,
	}
}

// AtBeat creates an event builder at a specific beat
func (b *Builder) AtBeat(beat Beat) *EventBuilder {
	if b.musical == nil {
		b.musical = NewMusicalContext(120, 4, 4) // Default context
	}
	return &EventBuilder{
		builder: b,
		time:    b.musical.BeatToSeconds(beat),
	}
}

// AtBar creates an event builder at a specific bar
func (b *Builder) AtBar(bar Bar) *EventBuilder {
	if b.musical == nil {
		b.musical = NewMusicalContext(120, 4, 4) // Default context
	}
	return &EventBuilder{
		builder: b,
		time:    b.musical.BarToSeconds(bar),
	}
}

// Curve creates a parameter curve builder
func (b *Builder) Curve(parameterID string) *CurveBuilder {
	return &CurveBuilder{
		builder:     b,
		parameterID: parameterID,
		startTime:   0,
		points:      make([]ControlPoint, 0),
	}
}

// Build returns the final AHAP
func (b *Builder) Build() *AHAP {
	return b.ahap
}

// Export writes the AHAP to a file
func (b *Builder) Export(filename string, indent bool) error {
	return b.ahap.Export(filename, indent)
}

// TransientBuilder builds a transient event
type TransientBuilder struct {
	builder   *Builder
	time      float64
	intensity float64
	sharpness float64
}

// Intensity sets the intensity
func (tb *TransientBuilder) Intensity(intensity float64) *TransientBuilder {
	tb.intensity = intensity
	return tb
}

// Sharpness sets the sharpness
func (tb *TransientBuilder) Sharpness(sharpness float64) *TransientBuilder {
	tb.sharpness = sharpness
	return tb
}

// Add adds the transient event and returns the builder
func (tb *TransientBuilder) Add() *Builder {
	tb.builder.ahap.AddHapticTransient(tb.time, tb.intensity, tb.sharpness)
	return tb.builder
}

// ContinuousBuilder builds a continuous event
type ContinuousBuilder struct {
	builder   *Builder
	time      float64
	duration  float64
	intensity float64
	sharpness float64
}

// Intensity sets the intensity
func (cb *ContinuousBuilder) Intensity(intensity float64) *ContinuousBuilder {
	cb.intensity = intensity
	return cb
}

// Sharpness sets the sharpness
func (cb *ContinuousBuilder) Sharpness(sharpness float64) *ContinuousBuilder {
	cb.sharpness = sharpness
	return cb
}

// Add adds the continuous event and returns the builder
func (cb *ContinuousBuilder) Add() *Builder {
	cb.builder.ahap.AddHapticContinuous(cb.time, cb.duration, cb.intensity, cb.sharpness)
	return cb.builder
}

// EventBuilder builds events at musical positions
type EventBuilder struct {
	builder *Builder
	time    float64
}

// Transient creates a transient event at this time
func (eb *EventBuilder) Transient() *TransientBuilder {
	return &TransientBuilder{
		builder:   eb.builder,
		time:      eb.time,
		intensity: 0.5,
		sharpness: 0.5,
	}
}

// Continuous creates a continuous event at this time
func (eb *EventBuilder) Continuous(duration float64) *ContinuousBuilder {
	return &ContinuousBuilder{
		builder:   eb.builder,
		time:      eb.time,
		duration:  duration,
		intensity: 0.5,
		sharpness: 0.5,
	}
}

// ContinuousBars creates a continuous event at this time with duration in bars
func (eb *EventBuilder) ContinuousBars(bars Bar) *ContinuousBuilder {
	if eb.builder.musical == nil {
		eb.builder.musical = NewMusicalContext(120, 4, 4)
	}
	duration := eb.builder.musical.BarToSeconds(bars)
	return &ContinuousBuilder{
		builder:   eb.builder,
		time:      eb.time,
		duration:  duration,
		intensity: 0.5,
		sharpness: 0.5,
	}
}

// ContinuousBeats creates a continuous event at this time with duration in beats
func (eb *EventBuilder) ContinuousBeats(beats Beat) *ContinuousBuilder {
	if eb.builder.musical == nil {
		eb.builder.musical = NewMusicalContext(120, 4, 4)
	}
	duration := eb.builder.musical.BeatToSeconds(beats)
	return &ContinuousBuilder{
		builder:   eb.builder,
		time:      eb.time,
		duration:  duration,
		intensity: 0.5,
		sharpness: 0.5,
	}
}

// CurveBuilder builds parameter curves
type CurveBuilder struct {
	builder     *Builder
	parameterID string
	startTime   float64
	points      []ControlPoint
}

// At sets the start time
func (cb *CurveBuilder) At(time float64) *CurveBuilder {
	cb.startTime = time
	return cb
}

// From sets the start point and begins defining the curve
func (cb *CurveBuilder) From(time, value float64) *CurveFromBuilder {
	return &CurveFromBuilder{
		curveBuilder: cb,
		startTime:    time,
		startValue:   value,
	}
}

// AddPoint adds a single control point
func (cb *CurveBuilder) AddPoint(time, value float64) *CurveBuilder {
	cb.points = append(cb.points, ControlPoint{
		Time:           time,
		ParameterValue: value,
	})
	return cb
}

// Points sets all control points at once
func (cb *CurveBuilder) Points(points []ControlPoint) *CurveBuilder {
	cb.points = points
	return cb
}

// Add adds the curve to the AHAP and returns the builder
func (cb *CurveBuilder) Add() *Builder {
	cb.builder.ahap.AddCurve(cb.parameterID, cb.startTime, cb.points)
	return cb.builder
}

// CurveFromBuilder helps build interpolated curves
type CurveFromBuilder struct {
	curveBuilder *CurveBuilder
	startTime    float64
	startValue   float64
}

// To defines the end point and creates a linear interpolation
func (cfb *CurveFromBuilder) To(endTime, endValue float64) *CurveToBuilder {
	return &CurveToBuilder{
		curveBuilder: cfb.curveBuilder,
		startTime:    cfb.startTime,
		startValue:   cfb.startValue,
		endTime:      endTime,
		endValue:     endValue,
	}
}

// CurveToBuilder completes the curve definition
type CurveToBuilder struct {
	curveBuilder *CurveBuilder
	startTime    float64
	startValue   float64
	endTime      float64
	endValue     float64
}

// Steps creates linear interpolation with specified steps
func (ctb *CurveToBuilder) Steps(steps int) *CurveBuilder {
	points := CreateCurve(ctb.startTime, ctb.endTime, ctb.startValue, ctb.endValue, steps)
	ctb.curveBuilder.points = append(ctb.curveBuilder.points, points...)
	return ctb.curveBuilder
}

// EaseInOut creates ease-in-out interpolation
func (ctb *CurveToBuilder) EaseInOut(steps int) *CurveBuilder {
	start := ControlPoint{Time: ctb.startTime, ParameterValue: ctb.startValue}
	end := ControlPoint{Time: ctb.endTime, ParameterValue: ctb.endValue}
	points := EaseInOut(start, end, steps)
	ctb.curveBuilder.points = append(ctb.curveBuilder.points, points...)
	return ctb.curveBuilder
}

// Exponential creates exponential interpolation
func (ctb *CurveToBuilder) Exponential(steps int, exponent float64) *CurveBuilder {
	start := ControlPoint{Time: ctb.startTime, ParameterValue: ctb.startValue}
	end := ControlPoint{Time: ctb.endTime, ParameterValue: ctb.endValue}
	points := Exponential(start, end, steps, exponent)
	ctb.curveBuilder.points = append(ctb.curveBuilder.points, points...)
	return ctb.curveBuilder
}

// SequenceBuilder helps build sequences of events across multiple bars/beats
type SequenceBuilder struct {
	builder *Builder
}

// Sequence creates a new sequence builder for creating patterns
func (b *Builder) Sequence() *SequenceBuilder {
	if b.musical == nil {
		b.musical = NewMusicalContext(120, 4, 4) // Default context
	}
	return &SequenceBuilder{
		builder: b,
	}
}

// TransientsOnBeats adds transient events on specific beats across a range of bars
// Example: TransientsOnBeats([]Beat{0, 2}, 5, 8) adds transients on beats 0 and 2 in bars 5-8
func (sb *SequenceBuilder) TransientsOnBeats(beats []Beat, startBar, endBar int, intensity, sharpness float64) *Builder {
	for bar := startBar; bar <= endBar; bar++ {
		barTime := sb.builder.musical.BarToSeconds(Bar(bar))
		for _, beat := range beats {
			beatTime := sb.builder.musical.BeatToSeconds(beat)
			sb.builder.Transient(barTime + beatTime).
				Intensity(intensity).
				Sharpness(sharpness).
				Add()
		}
	}
	return sb.builder
}

// TransientsOnBeatsInBar adds transient events on specific beats within a single bar
func (sb *SequenceBuilder) TransientsOnBeatsInBar(beats []Beat, bar int, intensity, sharpness float64) *Builder {
	return sb.TransientsOnBeats(beats, bar, bar, intensity, sharpness)
}

// EveryBeat adds a transient on every beat for a range of bars
func (sb *SequenceBuilder) EveryBeat(startBar, endBar int, intensity, sharpness float64) *Builder {
	beatsPerBar := sb.builder.musical.BeatsPerBar()
	beats := make([]Beat, beatsPerBar)
	for i := 0; i < beatsPerBar; i++ {
		beats[i] = Beat(i)
	}
	return sb.TransientsOnBeats(beats, startBar, endBar, intensity, sharpness)
}

// EveryNthBeat adds a transient on every nth beat for a range of bars (e.g., every 2nd beat = on-beat)
func (sb *SequenceBuilder) EveryNthBeat(n int, startBar, endBar int, intensity, sharpness float64) *Builder {
	beatsPerBar := sb.builder.musical.BeatsPerBar()
	totalBars := endBar - startBar + 1
	totalBeats := totalBars * beatsPerBar
	
	for i := 0; i < totalBeats; i += n {
		bar := i / beatsPerBar
		beatInBar := i % beatsPerBar
		actualBar := startBar + bar
		if actualBar <= endBar {
			barTime := sb.builder.musical.BarToSeconds(Bar(actualBar))
			beatTime := sb.builder.musical.BeatToSeconds(Beat(beatInBar))
			sb.builder.Transient(barTime + beatTime).
				Intensity(intensity).
				Sharpness(sharpness).
				Add()
		}
	}
	return sb.builder
}

// Pattern applies a custom pattern function to each bar
// The function receives the bar number and should add events to the builder
func (sb *SequenceBuilder) Pattern(startBar, endBar int, fn func(b *Builder, bar int)) *Builder {
	for bar := startBar; bar <= endBar; bar++ {
		fn(sb.builder, bar)
	}
	return sb.builder
}
