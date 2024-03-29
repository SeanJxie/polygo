package polygo

import (
	"log"
	"math"
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

	for s[len(s)-1] == 0 && len(s) > 1 {
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

// equalAbs returns true if the absolute error between a and b is at most delta, else false.
func equalAbs(a, b, delta float64) bool {

	if a == b {
		return true
	}

	return math.Abs(a-b) <= delta
}

// equalRel returns true if the relative error between a and b is at most epsilon, else false.
func equalRel(a, b, epsilon float64) bool {

	if a == b {
		return true
	}

	// The "true" value is taken to be the smaller one.
	return math.Abs(a-b)/math.Min(math.Abs(a), math.Abs(b)) <= epsilon
}

// sign returns the sign of a.
func sign(a float64) int {

	if equalRel(a, 0, 0.0001) {
		return 0
	}

	if a > 0 {
		return 1
	}

	return -1
}

// fact returns n factorial.
//
// Panics for negative n.
func fact(n int) float64 {

	if n < 0 {
		log.Panic("fact: negative n.")
	}

	if n == 0 {
		return 1
	}

	return float64(n) * fact(n-1)
}

// choose returns the binomial coefficient n choose k.
//
// Panics for negative n, k or if k > n.
func choose(n, k int) float64 {

	if n < 0 || k < 0 || k > n {
		log.Panicf("choose: invalid binomial coefficient: %dC%d.", n, k)
	}

	return fact(n) / (fact(k) * fact(n-k))
}
