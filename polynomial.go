package polygo

import (
	"errors"
	"fmt"
)

// A RealPolynomial is represented as a slice of coefficients ordered increasingly by degree.
// For example, one can imagine: 5x^0 + 4x^1 + (-2)x^2 + ...
type RealPolynomial struct {
	coeffs []float64
}

// A point in R^2.
type Point struct {
	X, Y float64
}

// The number of iterations used in Newton's Method implmentation in root solving functions.
var (
	globalNewtonIterations = 25
)

// NumCoeffs returns the number of coefficients the polynomial contains.
func (rp *RealPolynomial) NumCoeffs() int {
	return len(rp.coeffs)
}

// LeadCoeff Returns the coefficient of the highest degree term of the current instance.
func (rp *RealPolynomial) LeadCoeff() float64 {
	return rp.coeffs[len(rp.coeffs)-1]
}

// Degree returns the degree of the polynomial.
func (rp *RealPolynomial) Degree() int {
	// Coefficients should be maintained in such a way that allow the
	// number of coefficients to be one less than the degree of the polynomial.
	return rp.NumCoeffs() - 1
}

// At returns the value of the polynomial evaluated at x.
func (rp *RealPolynomial) At(x float64) float64 {
	// Implement Horner's Method
	length := len(rp.coeffs)
	out := rp.coeffs[length-1]
	for i := length - 2; i >= 0; i-- {
		out = out*x + rp.coeffs[i]
	}
	return out
}

// Derivative returns the derivative of the polynomial.
func (rp *RealPolynomial) Derivative() *RealPolynomial {
	// In the case that the polynomial is constant, the derivative has the same number of terms.
	// We deal with this case knowing that the derivative of any real constant is 0.
	if rp.Degree() == 0 {
		deriv, _ := NewRealPolynomial([]float64{0})
		return deriv
	}

	nDerivativeCoeffs := len(rp.coeffs) - 1
	derivativeCoeffs := make([]float64, nDerivativeCoeffs)
	for i := 0; i < nDerivativeCoeffs; i++ {
		derivativeCoeffs[i] = rp.coeffs[i+1] * float64(i+1)
	}

	deriv, _ := NewRealPolynomial(derivativeCoeffs)
	return deriv
}

// ShiftRight returns a new polynomial in which all coefficients of each term
// are shifted to the right by a specified offset based on the current polynomial.
//
// A right shift by N is equivalent to multipliying the current polynomial by x^N.
func (rp *RealPolynomial) ShiftRight(offset int) *RealPolynomial {
	if offset < 0 {
		panic("invalid offset")
	}
	shiftedCoeffs := make([]float64, rp.NumCoeffs()+offset)
	copy(shiftedCoeffs[offset:], rp.coeffs)
	rp, _ = NewRealPolynomial(shiftedCoeffs)
	return rp
}

// Equal returns true if the polynomial is equal to rp2. Otherwise, false is returned.
func (rp1 *RealPolynomial) Equal(rp2 *RealPolynomial) bool {
	if rp1.NumCoeffs() != rp2.NumCoeffs() {
		return false
	}

	for i := 0; i < rp1.NumCoeffs(); i++ {
		if rp1.coeffs[i] != rp2.coeffs[i] {
			return false
		}
	}

	return true
}

// IsZero returns true if polynomial is equal to the zero polynomial. Otherwise, false is returned.
func (rp *RealPolynomial) IsZero() bool {
	return rp.Degree() == 0 && rp.coeffs[0] == 0.0
}

// IsDegree returns true if polynomial is of degree n. Otherwise, false is returned.
func (rp *RealPolynomial) IsDegree(n int) bool {
	return rp.Degree() == n
}

// CoeffAtDegree returns the coefficient of the polynomial at degree n.
func (rp *RealPolynomial) CoeffAtDegree(n int) float64 {
	if n < 0 {
		panic("invalid degree")
	}
	return rp.coeffs[n]
}

// Add adds the polynomial to rp2 and returns the sum.
//
// The current instance is also set to the sum.
func (rp1 *RealPolynomial) Add(rp2 *RealPolynomial) *RealPolynomial {
	var maxNumCoeffs int

	// Pad "shorter" polynomial with 0s.
	if rp1.NumCoeffs() >= rp2.NumCoeffs() {
		maxNumCoeffs = rp1.NumCoeffs()
		for rp2.NumCoeffs() < maxNumCoeffs {
			rp2.coeffs = append(rp2.coeffs, 0.0)
		}

	} else if rp1.NumCoeffs() < rp2.NumCoeffs() {
		maxNumCoeffs = len(rp2.coeffs)
		for rp1.NumCoeffs() < maxNumCoeffs {
			rp1.coeffs = append(rp1.coeffs, 0.0)
		}
	} else {
		maxNumCoeffs = len(rp1.coeffs)
	}

	// Add coefficients with matching degrees.
	sumCoeffs := make([]float64, maxNumCoeffs)

	for i := 0; i < maxNumCoeffs; i++ {
		sumCoeffs[i] = rp1.coeffs[i] + rp2.coeffs[i]
	}

	rp1.coeffs = stripTailingZeroes(sumCoeffs)
	rp2.coeffs = stripTailingZeroes(rp2.coeffs)
	return rp1
}

// Sub subtracts rp2 from the polynomial and returns the difference.
//
// The current instance is also set to the difference.
func (rp1 *RealPolynomial) Sub(rp2 *RealPolynomial) *RealPolynomial {
	var maxNumCoeffs int

	// Pad "shorter" polynomial with 0s.
	if rp1.NumCoeffs() > rp2.NumCoeffs() {
		maxNumCoeffs = rp1.NumCoeffs()
		for rp2.NumCoeffs() < maxNumCoeffs {
			rp2.coeffs = append(rp2.coeffs, 0.0)
		}

	} else if rp1.NumCoeffs() < rp2.NumCoeffs() {
		maxNumCoeffs = len(rp2.coeffs)
		for rp1.NumCoeffs() < maxNumCoeffs {
			rp1.coeffs = append(rp1.coeffs, 0.0)
		}
	} else {
		maxNumCoeffs = len(rp1.coeffs)
	}

	// Subtract coefficients with matching degrees.
	diffCoeffs := make([]float64, maxNumCoeffs)

	for i := 0; i < maxNumCoeffs; i++ {
		diffCoeffs[i] = rp1.coeffs[i] - rp2.coeffs[i]
	}
	rp1.coeffs = stripTailingZeroes(diffCoeffs)
	rp2.coeffs = stripTailingZeroes(rp2.coeffs)
	return rp1
}

// MulNaive multiplies the polynomial with rp2 and returns the product.
//
// The current instance is also set to the product.
//
// It is not recommended to use this function as it is generally slow. Use Mul() instead.
func (rp1 *RealPolynomial) MulNaive(rp2 *RealPolynomial) *RealPolynomial {
	prodCoeffs := make([]float64, rp1.Degree()+rp2.Degree()+1)

	for i := 0; i < rp1.NumCoeffs(); i++ {
		for j := 0; j < rp2.NumCoeffs(); j++ {
			// We use += since we may visit the same index multiple times
			prodCoeffs[i+j] += rp1.coeffs[i] * rp2.coeffs[j]
		}
	}

	rp1.coeffs = stripTailingZeroes(prodCoeffs)
	return rp1
}

// Mul multiplies the polynomial with rp2 and returns the product.
//
// The current instance is also set to the product.
func (rp1 *RealPolynomial) Mul(rp2 *RealPolynomial) *RealPolynomial {
	lenRp1 := len(rp1.coeffs)
	lenRp2 := len(rp2.coeffs)

	padLen := nextClosestPowerOfTwo(lenRp1 + lenRp2 - 1)

	coeffs1 := make([]float64, padLen)
	coeffs2 := make([]float64, padLen)
	copy(coeffs1, rp1.coeffs)
	copy(coeffs2, rp2.coeffs)

	// With the FFT, we can run in O(n log n) time.
	fa := fastFourierTransform(complex128Slice(coeffs1))
	fb := fastFourierTransform(complex128Slice(coeffs2))

	fc := make([]complex128, padLen)
	for i := 0; i < padLen; i++ {
		fc[i] = fa[i] * fb[i]
	}

	tmpCoeffs := float64Slice(inverseFastFourierTransform(fc))

	for i, c := range tmpCoeffs {
		tmpCoeffs[i] = c / float64(padLen)
	}

	rp1.coeffs = stripTailingZeroes(tmpCoeffs[:rp1.Degree()+rp2.Degree()+1])
	return rp1
}

// MulS multiplies the polynomial with a scalar and returns the product.
//
// The current instance is also set to the product.
func (rp *RealPolynomial) MulS(s float64) *RealPolynomial {
	for i := 0; i < len(rp.coeffs); i++ {
		rp.coeffs[i] *= s
	}
	return rp
}

// EuclideanDiv divides the polynomial by rp2 and returns the result as a quotient-remainder pair.
//
// The current instance is also set to the quotient.
func (rp1 *RealPolynomial) EuclideanDiv(rp2 *RealPolynomial) (*RealPolynomial, *RealPolynomial) {
	if rp1 == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

	if rp2.IsZero() {
		panic("RealPolynomial division by zero")
	}

	// Using special properties of the ordered coefficient system, we can divide polynomials
	// via shifts:
	// https://rosettacode.org/wiki/Polynomial_long_division
	quotCoeffs := make([]float64, rp1.Degree()-rp2.Degree()+1)
	var d *RealPolynomial
	var shift int
	var factor float64

	rem := *rp1

	for rem.Degree() >= rp2.Degree() {
		shift = rem.Degree() - rp2.Degree()
		d = rp2.ShiftRight(shift)
		factor = rem.LeadCoeff() / d.LeadCoeff()
		quotCoeffs[shift] = factor
		d.MulS(factor)
		rem.Sub(d)
	}

	rp1.coeffs = quotCoeffs
	return rp1, &rem
}

// Expr returns a string representation of the polynomial in increasing sum form.
func (rp *RealPolynomial) String() string {
	var expr string
	for d, c := range rp.coeffs {
		if d == len(rp.coeffs)-1 {
			expr += fmt.Sprintf("%fx^%d", c, d)
		} else {
			expr += fmt.Sprintf("%fx^%d + ", c, d)
		}
	}

	return expr
}

func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

// PrintExpr prints the string expression of the current instance in increasing sum form to standard output.
func (rp *RealPolynomial) Print() {
	fmt.Println(rp)
}

// NewRealPolynomial returns a new *RealPolynomial instance with the given coeffs.
func NewRealPolynomial(coeffs []float64) (*RealPolynomial, error) {
	if len(coeffs) == 0 {
		return nil, errors.New("cannot create polynomial with no coefficients")
	}

	var newPolynomial RealPolynomial
	newPolynomial.coeffs = stripTailingZeroes(coeffs)

	return &newPolynomial, nil
}

// SetNewtonIterations sets the number of iterations used in Newton's Method implmentation in root solving functions.
func SetNewtonIterations(n int) error {
	if n < 0 {
		return errors.New("cannot set negative iterations for Newton's Method")
	}
	globalNewtonIterations = n
	return nil
}
