package ahap

import (
	"fmt"
	"math"
)

// CreateCurve creates interpolated control points between start and end
func CreateCurve(startTime, endTime, startValue, endValue float64, steps int) []ControlPoint {
	if steps < 1 {
		steps = 1
	}

	timeDiff := endTime - startTime
	valueDiff := endValue - startValue
	timeStep := timeDiff / float64(steps)
	valueStep := valueDiff / float64(steps)

	points := make([]ControlPoint, steps)
	for i := 0; i < steps; i++ {
		points[i] = ControlPoint{
			Time:           startTime + timeStep*float64(i+1),
			ParameterValue: startValue + valueStep*float64(i+1),
		}
	}
	return points
}

// LinearInterpolation creates a linear curve between two points
func LinearInterpolation(start, end ControlPoint, steps int) []ControlPoint {
	return CreateCurve(start.Time, end.Time, start.ParameterValue, end.ParameterValue, steps)
}

// EaseInOut creates an ease-in-out curve
func EaseInOut(start, end ControlPoint, steps int) []ControlPoint {
	if steps < 1 {
		steps = 1
	}

	timeDiff := end.Time - start.Time
	valueDiff := end.ParameterValue - start.ParameterValue

	points := make([]ControlPoint, steps)
	for i := 0; i < steps; i++ {
		t := float64(i+1) / float64(steps)
		// Ease in-out using smoothstep
		smoothT := t * t * (3.0 - 2.0*t)
		points[i] = ControlPoint{
			Time:           start.Time + timeDiff*t,
			ParameterValue: start.ParameterValue + valueDiff*smoothT,
		}
	}
	return points
}

// Exponential creates an exponential curve
func Exponential(start, end ControlPoint, steps int, exponent float64) []ControlPoint {
	if steps < 1 {
		steps = 1
	}

	timeDiff := end.Time - start.Time
	valueDiff := end.ParameterValue - start.ParameterValue

	points := make([]ControlPoint, steps)
	for i := 0; i < steps; i++ {
		t := float64(i+1) / float64(steps)
		expT := math.Pow(t, exponent)
		points[i] = ControlPoint{
			Time:           start.Time + timeDiff*t,
			ParameterValue: start.ParameterValue + valueDiff*expT,
		}
	}
	return points
}

// FreqToSharpness converts frequency (Hz) to sharpness value (0-1)
func FreqToSharpness(freq float64, normalize bool) (float64, error) {
	if normalize {
		if freq > 230 {
			freq = 230
		}
		if freq < 80 {
			freq = 80
		}
	}

	if freq < 80 || freq > 230 {
		return 0, fmt.Errorf("incorrect frequency: frequency must be between 80 and 230, but it is %.2f", freq)
	}

	r := (math.Log(freq) - math.Log(80)) / (math.Log(230) - math.Log(80))

	if r < 0 || r > 1 {
		return 0, fmt.Errorf("the calculated normalized frequency is out of range: result must be between 0 and 1")
	}

	return r, nil
}

// AddCurve adds a parameter curve to the AHAP
func (a *AHAP) AddCurve(parameterID string, startTime float64, controlPoints []ControlPoint) {
	curve := &ParameterCurve{
		ParameterID:                 parameterID,
		Time:                        startTime,
		ParameterCurveControlPoints: controlPoints,
	}
	a.AddParameterCurve(curve)
}
