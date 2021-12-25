package polynomial

import (
	"errors"
	"fmt"
	"math"
)

/*
A real RealPolynomial is represented as a slice of coefficients ordered increasingly by degree.

For example:
5 x^0 + 4 x^1 + (-2) x^2 + ...
*/
type RealPolynomial struct {
	coeffs []float64
}

/* --- BEGIN GLOBAL SETTINGS --- */
var globalNewtonIterations = 100

/* --- END GLOBAL SETTINGS --- */

/* --- BEGIN STRUCT METHODS --- */

/*
Returns the number of coefficients of the current instance.
*/
func (rp *RealPolynomial) NumCoeffs() int {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	return len(rp.coeffs)
}

/*
Returns the int degree of the current instance.
*/
func (rp *RealPolynomial) Degree() int {
	if rp == nil {
		panic("received nil RealPolynomial")
	}
	// Coefficients should be maintained in such a way that allow the
	// number of coefficients to be one less than the degree of the polynomial.
	return rp.NumCoeffs() - 1
}

/*
Returns the float64 value of the current instance evaluated at x.
*/
func (rp *RealPolynomial) At(x float64) float64 {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	var out float64 // Zero value 0.0
	for d, c := range rp.coeffs {
		out += c * math.Pow(x, float64(d))
	}
	return out
}

/*
Returns the derivative of the current instance.
*/
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

/*
Returns the float64 coefficient of the highest degree term of the current instance.
*/
func (rp *RealPolynomial) LeadCoeff() float64 {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	return rp.coeffs[len(rp.coeffs)-1]
}

/*
Returns a RealPolynomial which has been multiplied by x^offset.
*/
func (rp *RealPolynomial) ShiftRight(offset int) *RealPolynomial {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	if offset < 0 {
		panic("invalid offset")
	}
	shiftedCoeffs := make([]float64, rp.NumCoeffs()+offset)
	for i, c := range rp.coeffs {
		shiftedCoeffs[i+offset] = c
	}
	rp, _ = NewRealPolynomial(shiftedCoeffs) // safe call
	return rp
}

/*
Checks if the current instance is equal to the RealPolynomial input and returns true if so. Otherwise, false is returned.
*/
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

/*
Checks if the current instance is equal to the zero RealPolynomial (only one coefficient of 0). Otherwise, false is returned.
*/
func (rp *RealPolynomial) IsZero() bool {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}
	return rp.Degree() == 0 && rp.coeffs[0] == 0.0
}

/*
Adds the current instance and the RealPolynomial input and returns the sum. The current instance is also set to the sum.
*/
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
	return rp1
}

/*
Subtracts the current instance and the RealPolynomial input and returns the difference. The current instance is also set to the difference.
*/
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
	return rp1
}

/*
Multiplies the current instance and the RealPolynomial input and returns the product. The current instance is also set to the product.
*/
func (rp1 *RealPolynomial) Mul(rp2 *RealPolynomial) *RealPolynomial {
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

/*
Multiplies the current instance and the float64 input and returns the product. The current instance is also set to the product.
*/
func (rp *RealPolynomial) MulS(s float64) *RealPolynomial {
	for i := 0; i < len(rp.coeffs); i++ {
		rp.coeffs[i] *= s
	}
	return rp
}

/*
Divides the current instance by the RealPolynomial input and returns the result as a quotient, remainder pair. The current instance is also set to the quotient.
*/
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

/*
Returns the int number of roots of the current instance on the closed interval [a, b].

Note: if there are an infinite amount of roots, -1 is returned.
*/
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

	return rp.countRootsWithinFromSturmChain(a, b, rp.sturmChain())

}

func (rp *RealPolynomial) countRootsWithinFromSturmChain(a, b float64, chain []*RealPolynomial) int {
	// Generate sequence A and B to count sign variations
	var seqA, seqB []float64

	for _, p := range chain {
		seqA = append(seqA, p.At(a))
		seqB = append(seqB, p.At(b))
	}

	halfOpenCount := countSignVariations(seqA) - countSignVariations(seqB)

	if rp.At(a) == 0.0 { // Manually check open lower bound
		return halfOpenCount + 1
	} else {
		return halfOpenCount
	}
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
		// The fact that EuclideanDiv changes the instance is kind of annoying here.
		tmp = *sturmChain[i-1]
		_, rem = tmp.EuclideanDiv(sturmChain[i])
		sturmChain = append(sturmChain, rem.MulS(-1))
	}

	return sturmChain
}

/*
Returns any float64 root of the current instance existing on the closed interval [a, b].

Note: if there are no roots on the provided interval, an error is set.
*/
func (rp *RealPolynomial) FindRootWithin(a, b float64) (float64, error) {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}

	nRootsWithin := rp.CountRootsWithin(a, b)

	if nRootsWithin == 0 {
		return 0.0, errors.New("the polynomial has no solutions in the provided interval")
	}

	if nRootsWithin < 0 { // Infinite amount of roots
		return 0.0, nil
	}

	// Implement Newton's Method
	var deriv *RealPolynomial
	lower, upper := math.Min(a, b), math.Max(a, b)
	guess := (lower + upper) / 2

	for i := 0; i < globalNewtonIterations; i++ {
		deriv = rp.Derivative()
		guess -= rp.At(guess) / deriv.At(guess)
	}

	return guess, nil
}

func (rp *RealPolynomial) findRootWithinFromSturmChain(a, b float64, chain []*RealPolynomial) (float64, error) {
	if rp == nil {
		panic("received nil *RealPolynomial")
	}

	nRootsWithin := rp.countRootsWithinFromSturmChain(a, b, chain)

	if nRootsWithin == 0 {
		return 0.0, errors.New("the polynomial has no solutions in the provided interval")
	}

	if nRootsWithin < 0 { // Infinite amount of roots
		return 0.0, nil
	}

	// Implement Newton's Method
	var deriv *RealPolynomial
	lower, upper := math.Min(a, b), math.Max(a, b)
	guess := (lower + upper) / 2

	for i := 0; i < globalNewtonIterations; i++ {
		deriv = rp.Derivative()
		guess -= rp.At(guess) / deriv.At(guess)
	}

	return guess, nil
}

func (rp *RealPolynomial) FindRootsWithin(a, b float64) []float64 {
	return rp.findRootsWithinAcc(a, b, nil, rp.sturmChain())
}

func (rp *RealPolynomial) findRootsWithinAcc(a, b float64, roots []float64, chain []*RealPolynomial) []float64 {
	// Implement a hybrid Bisection Method through accumulative recursion

	nRoots := rp.countRootsWithinFromSturmChain(a, b, chain)
	if nRoots > 1 {
		mp := (a + b) / 2
		return append(
			rp.findRootsWithinAcc(a, mp, roots, chain),
			rp.findRootsWithinAcc(mp, b, roots, chain)...,
		)

	} else if nRoots == 1 {
		root, _ := rp.findRootWithinFromSturmChain(a, b, chain)
		roots = append(roots, root)
	}

	return roots
}

/*
Returns a string representation of the current instance in increasing sum form.
*/
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

/*
Prints the string expression of the polynomial in increasing sum form to standard output.
*/
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

/*
Initializes and returns a new RealPolynomial struct with the given coefficients.
*/
func NewRealPolynomial(coeffs []float64) (*RealPolynomial, error) {
	if len(coeffs) == 0 {
		return nil, errors.New("cannot create polynomial with no coefficients")
	}

	var newPolynomial RealPolynomial
	newPolynomial.coeffs = stripTailingZeroes(coeffs)

	return &newPolynomial, nil
}

/*
Set the number of iterations that Newton's Method will use in root finding functions.
*/
func SetNewtonIterations(n int) error {
	if n < 0 {
		return errors.New("cannot set negative iterations for Newton's Method")
	}
	globalNewtonIterations = n
	return nil
}

// --- END UTILITY FUNCTIONS ---
