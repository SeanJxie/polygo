package polygo

import (
	"log"
	"math"
)

var (
	// The maximum absolute or relative error for two values to be considered equal.
	epsilon = 1e-5
)

const (
	ln2 = 0.693147180559945309417232121458176568075500134360255254120680009
)

// removeTrailingZeroes returns a copy of s with all trailing zeroes removed.
//
// If the entire slice is filled with 0.0, []float64{0.0} is returned.
func removeTrailingZeroes(s []float64) []float64 {
	if len(s) == 0 {
		return s
	}

	for approxEqual(s[len(s)-1], 0) && len(s) > 1 {
		s = s[:len(s)-1]
	}

	return s
}

// reverse returns a copy of s with reversed order.
func reverse(s []float64) []float64 {
	ls := len(s)

	ret := make([]float64, ls)

	for i := 0; i < ls; i++ {
		ret[i] = s[ls-i-1]
	}

	return ret
}

// plusOrMinus returns true if c is '+' or '-', else false.
func plusOrMinus(c rune) bool {
	return c == '+' || c == '-'
}

// expand returns s padded with trailing zeroes to reach length n.
//
// If n <= len(s), nothing is changed.
func expand(s []float64, n int) []float64 {

	if n <= len(s) {
		return s
	}

	expanded := make([]float64, n)
	copy(expanded, s)

	return expanded
}

// isPOT returns true if n is a natural (including 0) power of two, else false.
func isPOT(n int) bool {
	return n&(n-1) == 0 && n >= 1
}

// nextPOT returns the least power of 2 greater than or equal to n.
func nextPOT(n int) int {
	if isPOT(n) {
		return n
	}

	if n <= 0 {
		return 1
	}

	return int(math.Pow(2, math.Ceil(math.Log(float64(n))/ln2)))
}

// toComplex128 returns the []float64 slice s as a []complex128 slice.
func toComplex128(s []float64) []complex128 {

	c := make([]complex128, len(s))
	for i, v := range s {
		c[i] = complex(v, 0)
	}

	return c
}

// toFloat64 returns the []complex128 slice s as a []float64 slice.
//
// Specifically, the real parts are taken from the []complex128 slice to form the
// []float64 slice.
func toFloat64(s []complex128) []float64 {

	f := make([]float64, len(s))
	for i, v := range s {
		f[i] = real(v)
	}

	return f
}

// max returns the maximum value in s.
func max(s []float64) float64 {

	if len(s) == 0 {
		log.Panic("max: empty slice.")
	}

	max := s[0]
	for _, v := range s[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// min returns the minimum value in s.
func min(s []float64) float64 {

	if len(s) == 0 {
		log.Panic("min: empty slice.")
	}

	min := s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// approxEqual returns true if a and b are approximately equal, else false.
func approxEqual(a, b float64) bool {

	if a == b {
		return true
	}

	// For a visualization of the error check: https://www.desmos.com/calculator/b75lmu3wd6.

	absErr := math.Abs(a - b)
	relErr := math.Abs(a-b) / (math.Abs(a) + math.Abs(b))

	return absErr <= epsilon || relErr <= epsilon
}

// sign returns the sign of a.
func sign(a float64) int {

	if approxEqual(a, 0) {
		return 0
	}

	if a > 0 {
		return 1
	}

	return -1
}

// SetEpsilon sets a variable named epsilon, which is the absolute or relative error for two values
// to be considered equal in Polygo.
//
// Default is set to 1e-5 = 0.00001.
//
// Epsilon is used in equality checks all over the Polygo library, so the user should be sure to
// test that the value they set works for their use case.
//
// Panics for negative v.
func SetEpsilon(v float64) {

	if v < 0 {
		log.Panic("SetEpsilon: negative epsilon.")
	}

	epsilon = v
}
