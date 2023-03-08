package utils

import "math"

// Unique integer values
const (
	MAX_INT_64  = math.MaxInt64
	MAX_UINT_64 = math.MaxUint64
)

// IsFloat64AlmostEqual2 compares two floats and returns true if absolute diff is less OR EQ than threshold
func IsFloat64AlmostEqual2(a, b, threshold float64) bool {
	return math.Abs(a-b) <= threshold
}

// IsFloat64AlmostEqualByFraction compares two floats and returns true if relative diff is less than the threshold
func IsFloat64AlmostEqualByFraction(a, b, thresholdFraction float64) bool {
	// use threshold 0.01 for $ (otherwise 1e-9 should be sufficiently small)

	d := math.Max(math.Abs(a), math.Abs(b))
	if d == 0 {
		return true // both ==0
	}
	return math.Abs(a-b)/d < thresholdFraction
}
