// Package polygo is a collection of tools that make working with polynomials easier in Go.
package polygo

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
)

// A real RealPolynomial is represented as a slice of coefficients ordered increasingly by degree.
// For example, one can imagine: 5x^0 + 4x^1 + (-2)x^2 + ...
type RealPolynomial struct {
	coeffs []float64
}

/* --- BEGIN GLOBAL SETTINGS --- */
// The
var globalNewtonIterations = 100

/* --- END GLOBAL SETTINGS --- */

/* --- BEGIN STRUCT METHODS --- */

// NumCoeffs returns the number of coefficients of the current instance.
func (rp *RealPolynomial) NumCoeffs() int {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	return len(rp.coeffs)
}

// Degree returns the degree of the current instance.
func (rp *RealPolynomial) Degree() int {
	if rp == nil {
		panic("received nil RealPolynomial")
	}
	// Coefficients should be maintained in such a way that allow the
	// number of coefficients to be one less than the degree of the polynomial.
	return len(rp.coeffs) - 1
}

// At returns the value of the current instance evaluated at x.
func (rp *RealPolynomial) At(x float64) float64 {

	// Implement Horner's Method
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	length := len(rp.coeffs)
	out := rp.coeffs[length-1]
	for i := length - 2; i >= 0; i-- {
		out = out*x + rp.coeffs[i]
	}
	return out
}

// Derivative returns the derivative of the current instance.
// The current instance is not modified.
func (rp *RealPolynomial) Derivative() *RealPolynomial {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	// In the case that the polynomial is constant, the derivative has the same number of terms.
	// We deal with this case knowing that the derivative of any real constant is 0.
	if rp.Degree() == 0 {
		deriv, _ := NewRealPolynomial([]float64{0}) // safe call
		return deriv
	}

	nDerivativeCoeffs := len(rp.coeffs) - 1
	derivativeCoeffs := make([]float64, nDerivativeCoeffs)
	for i := 0; i < nDerivativeCoeffs; i++ {
		derivativeCoeffs[i] = rp.coeffs[i+1] * float64(i+1)
	}

	deriv, _ := NewRealPolynomial(derivativeCoeffs) // safe call
	return deriv
}

// LeadCoeff Returns the coefficient of the highest degree term of the current instance.
func (rp *RealPolynomial) LeadCoeff() float64 {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	return rp.coeffs[len(rp.coeffs)-1]
}

// ShiftRight shifts the coefficients of each term in the current instance rightwards by offset and returns the resulting polynomial.
// The current instance is not modified.
// A right shift by N is equivalent to multipliying the current instance by x^N.
func (rp *RealPolynomial) ShiftRight(offset int) *RealPolynomial {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	if offset < 0 {
		panic("invalid offset")
	}
	shiftedCoeffs := make([]float64, rp.NumCoeffs()+offset)
	copy(shiftedCoeffs[offset:], rp.coeffs)
	rp, _ = NewRealPolynomial(shiftedCoeffs) // safe call
	return rp
}

// Equal returns true if the current instance is equal to rp2. Otherwise, false is returned.
func (rp1 *RealPolynomial) Equal(rp2 *RealPolynomial) bool {
	if rp1 == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

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

// IsZero returns true if current instance is equal to the zero polynomial. Otherwise, false is returned.
func (rp *RealPolynomial) IsZero() bool {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	return rp.Degree() == 0 && rp.coeffs[0] == 0.0
}

// Add adds the current instance and rp2 and returns the sum.
// The current instance is also set to the sum.
func (rp1 *RealPolynomial) Add(rp2 *RealPolynomial) *RealPolynomial {
	if rp1 == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

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

// Sub subtracts rp2 from the current instance and returns the difference.
// The current instance is also set to the difference.
func (rp1 *RealPolynomial) Sub(rp2 *RealPolynomial) *RealPolynomial {
	if rp1 == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

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

// MulNaive multiplies the current instance with rp2 and returns the product.
// The current instance is also set to the product.
//
// It is not recommended to use this function. Use Mul instead.
func (rp1 *RealPolynomial) MulNaive(rp2 *RealPolynomial) *RealPolynomial {
	if rp1 == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

	prodCoeffs := make([]float64, rp1.Degree()+rp2.Degree()+1)

	for i := 0; i < rp1.NumCoeffs(); i++ {
		for j := 0; j < rp2.NumCoeffs(); j++ {
			// We use += since we may visit the same index multiple times
			prodCoeffs[i+j] += rp1.coeffs[i] * rp2.coeffs[j]
		}
	}

	rp1.coeffs = prodCoeffs
	return rp1
}

// Mul multiplies the current instance with rp2 and returns the product.
// The current instance is also set to the product.
func (rp1 *RealPolynomial) Mul(rp2 *RealPolynomial) *RealPolynomial {
	if rp1 == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

	lenRp1 := len(rp1.coeffs)
	lenRp2 := len(rp2.coeffs)

	var padLen int

	if lenRp1 > lenRp2 {
		padLen = nextClosestPowerOfTwo(lenRp1)
	} else {
		padLen = nextClosestPowerOfTwo(lenRp2)
	}

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

	rp1.coeffs = stripTailingZeroes(tmpCoeffs)
	return rp1
}

// MulS multiplies the current instance with the scalar s and returns the product.
// The current instance is also set to the product.
func (rp *RealPolynomial) MulS(s float64) *RealPolynomial {
	for i := 0; i < len(rp.coeffs); i++ {
		rp.coeffs[i] *= s
	}
	return rp
}

// EuclideanDiv divides the current instance by rp2 and returns the result as a quotient-remainder pair.
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

// CountRootsWithin returns the number of roots of the current instance on the closed interval [a, b].
// If there are an infinite amount of roots, -1 is returned.
func (rp *RealPolynomial) CountRootsWithin(a, b float64) int {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}

	if a > b {
		panic("invalid interval")
	}

	if rp.Degree() == 0 && rp.coeffs[0] == 0.0 {
		return -1
	}

	return rp.countRootsWithinWSC(a, b, rp.sturmChain())

}

// FindRootWithin returns ANY root of the current instance existing on the closed interval [a, b].
// If there are no roots on the provided interval, an error is set.
func (rp *RealPolynomial) FindRootWithin(a, b float64) (float64, error) {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}

	if a > b {
		panic("invalid interval")
	}

	// Since findRootsWithinAcc operates on the half-open interval (a, b], manually check if a is a root.
	if rp.At(a) == 0.0 {
		return a, nil
	}

	sturmChain := rp.sturmChain()
	nRootsWithin := rp.countRootsWithinWSC(a, b, sturmChain)

	if nRootsWithin == 0 {
		return 0.0, errors.New("the polynomial has no solutions in the provided interval")
	}

	if nRootsWithin < 0 { // Infinite amount of roots
		return 0.0, nil
	}

	return rp.findRootWithinWSC(a, b, sturmChain)
}

// FindRootsWithin returns ALL roots of the current instance existing on the closed interval [a, b].
// Unlike FindRootWithin, no error is set if there are no solutions on the provided interval. Instead, an empty slice is returned.
// If there are an infinite number of solutions on [a, b], an error is set.
func (rp *RealPolynomial) FindRootsWithin(a, b float64) ([]float64, error) {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}

	// The only polynomial with infinitely many roots is P(x) = 0
	// https://math.stackexchange.com/questions/1137190/is-there-a-polynomial-that-has-infinitely-many-roots
	if rp.IsZero() {
		return nil, errors.New("infinitely many solutions")
	}

	// Since findRootsWithinAcc operates on the half-open interval (a, b], manually check if a is a root.
	if rp.At(a) == 0.0 {
		return append(rp.findRootsWithinAcc(a, b, nil, rp.sturmChain()), a), nil
	}
	return rp.findRootsWithinAcc(a, b, nil, rp.sturmChain()), nil
}

// findRootsWithinAcc is an accumulative implmentation of a hybrid Bisection Method through recursion.
// Wrapped by FindRootsWithin.
func (rp *RealPolynomial) findRootsWithinAcc(a, b float64, roots []float64, chain []*RealPolynomial) []float64 {
	nRoots := rp.countRootsWithinWSC(a, b, chain)
	if nRoots > 1 {
		mp := (a + b) / 2
		return append(
			rp.findRootsWithinAcc(a, mp, roots, chain),
			rp.findRootsWithinAcc(mp, b, roots, chain)...,
		)

	} else if nRoots == 1 {
		root, _ := rp.findRootWithinWSC(a, b, chain)
		roots = append(roots, root)
	}

	return roots
}

// FindIntersectionWithin returns ANY intersection point (as a two-element slice) of the current instance and rp2 existing on the closed interval [a, b].
// If there are no intersections on the provided interval, an error is set.
func (rp *RealPolynomial) FindIntersectionWithin(a, b float64, rp2 *RealPolynomial) ([2]float64, error) {
	tmp := *rp

	root, err := (&tmp).Sub(rp2).FindRootWithin(a, b)
	if err != nil {
		return [2]float64{}, err
	}

	point := [2]float64{root, rp.At(root)}
	return point, nil
}

// FindIntersectionsWithin returns ALL intersection point (as a two-element slice) of the current instance and rp2 existing on the closed interval [a, b].
// Unlike FindIntersectionWithin, no error is set if there are no intersections on the provided interval. Instead, an empty slice is returned.
// If there are an infinite number or solutions, an error is set.
func (rp *RealPolynomial) FindIntersectionsWithin(a, b float64, rp2 *RealPolynomial) ([][2]float64, error) {
	if rp == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

	tmp := *rp
	roots, err := tmp.Sub(rp2).FindRootsWithin(a, b)
	if err != nil {
		return nil, err
	}

	points := make([][2]float64, len(roots))

	for i, x := range roots {
		points[i] = [2]float64{x, rp.At(x)}
	}

	return points, nil
}

/*

Since a non-changing Sturm Chain is used repetitevly through multiple functions, the following private functions
with suffix "WSC" are such that the overhead caused by recomputing the Sturm chain every use is avoided by making the chain an input.

*/

func (rp *RealPolynomial) findRootWithinWSC(a, b float64, chain []*RealPolynomial) (float64, error) {
	nRootsWithin := rp.countRootsWithinWSC(a, b, chain)

	if nRootsWithin == 0 {
		return 0.0, errors.New("the polynomial has no solutions in the provided interval")
	}

	if nRootsWithin == -1 { // Infinite amount of roots
		return 0.0, nil
	}

	// Implement Newton's Method
	deriv := rp.Derivative()
	guess := (a + b) / 2
	var derivAtGuess float64
	for i := 0; i < globalNewtonIterations; i++ {
		derivAtGuess = deriv.At(guess)
		// In the case that the derivative evaluates to zero, return the current guess.
		if derivAtGuess == 0.0 {
			return guess, nil
		}
		guess -= rp.At(guess) / derivAtGuess
	}

	// Operate on a half-open interval.
	if guess == a {
		return 0.0, errors.New("the polynomial has no solutions in the provided interval")
	}

	return guess, nil
}

func (rp *RealPolynomial) countRootsWithinWSC(a, b float64, chain []*RealPolynomial) int {
	// Generate sequence A and B to count sign variations
	var seqA, seqB []float64

	for _, p := range chain {
		seqA = append(seqA, p.At(a))
		seqB = append(seqB, p.At(b))
	}

	return countSignVariations(seqA) - countSignVariations(seqB)
}

func (rp *RealPolynomial) sturmChain() []*RealPolynomial {
	// Implement Sturm's Theorem

	var sturmChain []*RealPolynomial
	var rem *RealPolynomial
	var tmp RealPolynomial
	sturmChain = append(sturmChain, rp)

	deriv := rp.Derivative()
	sturmChain = append(sturmChain, deriv)

	for i := 1; i < rp.Degree(); i++ {
		if sturmChain[i].Degree() == 0 {
			break
		}

		tmp = *sturmChain[i-1]
		_, rem = tmp.EuclideanDiv(sturmChain[i])
		sturmChain = append(sturmChain, rem.MulS(-1))
	}

	return sturmChain
}

/*

End "WSC"-suffixed (and related) functions.

*/

// Expr returns a string representation of the current instance in increasing sum form.
func (rp *RealPolynomial) Expr() string {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}

	var expr string
	for d, c := range rp.coeffs {
		if d == len(rp.coeffs)-1 {
			expr += fmt.Sprintf("%fx^%d", c, d)
		} else {
			expr += fmt.Sprintf("%fx^%d + ", c, d)
		}
	}

	return expr + "\n"
}

// PrintExpr prints the string expression of the current instance in increasing sum form to standard output.
func (rp *RealPolynomial) PrintExpr() {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}

	fmt.Print(rp.Expr())
}

// --- END STRUCT METHODS ---

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

// --- END UTILITY FUNCTIONS ---
