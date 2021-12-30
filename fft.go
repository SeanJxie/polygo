package polygo

/*
An implmentation of the Fast Fourier Transform.
*/

import (
	"math"
	"math/cmplx"
)

// Implements FFT
func fastFourierTransform(a []complex128) []complex128 {
	n := len(a)

	if n == 1 {
		return a
	}

	halfn := n / 2

	ae := make([]complex128, halfn)
	ao := make([]complex128, halfn)

	for i := 0; i < halfn; i++ {
		ae[i], ao[i] = a[i*2], a[i*2+1]
	}

	ye := fastFourierTransform(ae)
	yo := fastFourierTransform(ao)

	wn := cmplx.Exp(complex((2.0*math.Pi)/float64(n), 0.0) * 1.0i)
	w := complex(1.0, 0.0)

	y := make([]complex128, n)
	for k := 0; k < halfn; k++ {
		y[k] = ye[k] + w*yo[k]
		y[k+halfn] = ye[k] - w*yo[k]
		w *= wn
	}

	return y
}

// Implements IFFT
func inverseFastFourierTransform(a []complex128) []complex128 {
	n := len(a)

	if n == 1 {
		return a
	}

	halfn := n / 2

	ae := make([]complex128, halfn)
	ao := make([]complex128, halfn)

	for i := 0; i < halfn; i++ {
		ae[i], ao[i] = a[i*2], a[i*2+1]
	}

	ye := inverseFastFourierTransform(ae)
	yo := inverseFastFourierTransform(ao)

	wnInv := cmplx.Exp(complex((-2.0*math.Pi)/float64(n), 0.0) * 1.0i)
	w := 1.0 + 0.0i

	y := make([]complex128, n)
	for k := 0; k < halfn; k++ {
		y[k] = ye[k] + w*yo[k]
		y[k+halfn] = ye[k] - w*yo[k]
		w *= wnInv
	}

	return y
}
