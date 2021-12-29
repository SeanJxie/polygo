package polygo

// --- BEGIN UTILITY FUNCTIONS ---

// Strip all leading zeroes in the slice s. If the entire slice is filled with 0, the first element will be kept.
func stripTailingZeroes(s []float64) []float64 {
	for s[len(s)-1] == 0.0 && len(s) > 1 {
		s = s[:len(s)-1]
	}
	return s
}

// Cast []float64 to []complex128
func complex128Slice(s []float64) []complex128 {
	ret := make([]complex128, len(s))
	for i, e := range s {
		ret[i] = complex(e, 0)
	}
	return ret
}

// Cast []complex128 to []float64
func float64Slice(s []complex128) []float64 {
	ret := make([]float64, len(s))
	for i, e := range s {
		ret[i] = real(e)
	}
	return ret
}

// Counts sign variations in s: https://en.wikipedia.org/wiki/Budan%27s_theorem#Sign_variation
func countSignVariations(s []float64) int {
	// Filter zeroes in s.
	var filtered []float64
	for i := 0; i < len(s); i++ {
		if s[i] != 0.0 {
			filtered = append(filtered, s[i])
		}
	}

	// Count sign changes.
	var count int
	for i := 0; i < len(filtered)-1; i++ {
		if filtered[i]*filtered[i+1] < 0 {
			count++
		}
	}

	return count
}

// Finds the smallest power of 2 greater than n.
func nextClosestPowerOfTwo(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
