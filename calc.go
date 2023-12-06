package polygo

// Derivative returns the derivative of p.
func (p Poly) Derivative() Poly {

	// p is a constant, whose derivative is always 0.
	if p.deg == 0 {
		return NewPoly([]float64{0})
	}

	// Formal derivative.
	// For nonconstant p, if deg(p) = n, then deg(p') = n - 1.
	derivCoef := make([]float64, p.deg)
	for i := 0; i < p.deg; i++ {
		derivCoef[i] = p.coef[i+1] * float64(i+1)
	}

	return newPolyNoReverse(derivCoef)
}

// DerivativeN returns the nth derivative of p.
//
// Negative n will be treated as 0.
func (p Poly) DerivativeN(n int) Poly {

	// Reduce the overhead.
	if n > p.len {
		n = p.len
	}

	for ; n > 0; n-- {
		p = p.Derivative()
	}

	return p
}
