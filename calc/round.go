package calc

import "math"

// Rounds floats to 'points' decimal points after the comma.
func Round[T float32 | float64](f T, points uint) T {
	multiplier := math.Pow10(int(points))

	return T(math.Round(float64(f)*multiplier) / multiplier)
}
