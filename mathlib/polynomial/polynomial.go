package polynomial

import (
	"errors"
	"fmt"
	"math"
)

/*
A real RealPolynomial is represented as a sum of increasing
degrees of a x.

For example:
5x^0 + 4x^1 + (-2)x^2 + ...
*/
type RealPolynomial struct {
	// Store coefficients in a slice of floats.
	coeffs []float64
}

// --- BEGIN GLOBAL CONSTS ---
const globalNewtonIterations = 100

// ---

// --- BEGIN STRUCT METHODS ---

func (rp *RealPolynomial) NumCoeffs() int {
	return len(rp.coeffs)

}

func (rp *RealPolynomial) Degree() int {
	// Coefficients should be maintained in such a way that allow the
	// number of coefficients to be one less than the degree of the polynomial.
	return rp.NumCoeffs() - 1
}

// Evaluates the RealPolynomial at x and returns the computed value.
func (rp *RealPolynomial) At(x float64) float64 {
	var out float64 // Zero value 0.0
	for d, c := range rp.coeffs {
		out += c * math.Pow(x, float64(d))
	}
	return out
}

// Returns the a pointer to the RealPolynomial derivative of the current instance of RealPolynomial.
func (rp *RealPolynomial) Derivative() (*RealPolynomial, error) {
	if rp == nil {
		return nil, errors.New("RealPolynomial instance cannot be nil")
	}

	// In the case that the polynomial is constant, the derivative has the same number of terms.
	// We deal with this case knowing that the derivative of any real constant is 0.
	if rp.Degree() == 0 {
		return NewRealPolynomial([]float64{0})
	}

	nDerivativeCoeffs := len(rp.coeffs) - 1
	derivativeCoeffs := make([]float64, nDerivativeCoeffs)
	for i := 0; i < nDerivativeCoeffs; i++ {
		derivativeCoeffs[i] = rp.coeffs[i+1] * float64(i+1)
	}

	return NewRealPolynomial(derivativeCoeffs)
}

func (rp *RealPolynomial) LeadCoeff() float64 {
	return rp.coeffs[len(rp.coeffs)-1]
}

func (rp *RealPolynomial) ShiftRight(offset int) *RealPolynomial {
	shiftedCoeffs := make([]float64, rp.NumCoeffs()+offset)
	for i, c := range rp.coeffs {
		shiftedCoeffs[i+offset] = c
	}
	rp, _ = NewRealPolynomial(shiftedCoeffs)
	return rp
}

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

// Adds rp2 to the current instance and returns the sum.
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
	return rp1
}

// Subtracts rp2 to the current instance and returns the difference.
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
	return rp1
}

// Multiplies rp2 with the current instance and returns the product.
func (rp1 *RealPolynomial) Mul(rp2 *RealPolynomial) *RealPolynomial {
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

// Multiplies scalar s with the current instance and returns the product.
func (rp *RealPolynomial) MulS(s float64) *RealPolynomial {
	for i := 0; i < len(rp.coeffs); i++ {
		rp.coeffs[i] *= s
	}
	return rp
}

// Divides the current instance by rp2 and returns the quotient and remainder.
func (rp1 *RealPolynomial) EuclideanDiv(rp2 *RealPolynomial) (*RealPolynomial, *RealPolynomial) {

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

	//rp1.coeffs = quotCoeffs
	return rp1, &rem
}

// Returns the number of roots on the closed interval [a, b].
func (rp *RealPolynomial) CountRootsWithin(a, b float64) (int, error) {
	if rp == nil {
		return 0, errors.New("instance cannot be nil")
	}

	if rp.Degree() == 0 && rp.coeffs[0] == 0.0 {
		return 0, errors.New("infinite number of roots on the provided interval")
	}

	// Implement Sturm's Theorem

	// Generate Sturm chain
	var sturmChain []*RealPolynomial
	sturmChain = append(sturmChain, rp)

	deriv, _ := rp.Derivative()
	sturmChain = append(sturmChain, deriv)

	for i := 1; i < rp.Degree(); i++ {
		if sturmChain[i].Degree() == 0 {
			break
		}
		_, rem := sturmChain[i-1].EuclideanDiv(sturmChain[i])
		sturmChain = append(sturmChain, rem.MulS(-1))
	}

	// Generate sequence A and B to count sign variations
	var seqA, seqB []float64

	for _, p := range sturmChain {
		seqA = append(seqA, p.At(a))
		seqB = append(seqB, p.At(b))
	}

	halfOpenCount := countSignVariations(seqA) - countSignVariations(seqB)
	if rp.At(a) == 0.0 { // Manually check open interval
		return halfOpenCount + 1, nil
	} else {
		return halfOpenCount, nil
	}
}

/*
FindRootWithin(a, b float64) returns a float64 root (out of potentially many) of the
current instance on the closed interval [a, b].

Examples:
>>> NewRealPolynomial([]float64{0, 1}).FindRootWithin(-1, 1)
0.0

*/
func (rp *RealPolynomial) FindRootWithin(a, b float64) (float64, error) {
	if rp == nil {
		return 0, errors.New("instance cannot be nil")
	}

	nRootsWithin, err := rp.CountRootsWithin(a, b)

	if err != nil { //
		return 0, nil
	}

	if nRootsWithin == 0 {
		return 0, errors.New("the polynomial has no solutions in the provided interval")
	}

	// Implement Newton's Method
	var deriv *RealPolynomial
	lower, upper := math.Min(a, b), math.Max(a, b)
	guess := (lower + upper) / 2

	for i := 0; i < globalNewtonIterations; i++ {
		deriv, err = rp.Derivative() // rp is not nil. This is a safe call.
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		guess -= rp.At(guess) / deriv.At(guess)
	}

	return guess, nil
}

func (rp *RealPolynomial) FindRootsBestAttempt() {
	// switch rp.Degree() {

	// 	case 0: // Case: constant
	// 		{
	// 			if rp.coeffs[0] == 0.0 {
	// 				solutions = nil
	// 				nSolutions = int(math.Inf(1))
	// 			} else {
	// 				solutions = []float64{}
	// 				nSolutions = 0
	// 			}
	// 			break
	// 		}

	// 	case 1: // Case: linear
	// 		{
	// 			return []float64{-rp.coeffs[0] / rp.coeffs[1]}, 1
	// 		}

	// 	case 2: // Case: quadratic (quadratic formula)
	// 		{
	// 			a, b, c := rp.coeffs[2], rp.coeffs[1], rp.coeffs[0]
	// 			disc := b*b - 4*a*c
	// 			sqrtDisc := math.Sqrt(disc)
	// 			if disc > 0 {
	// 				solutions = []float64{(-b + sqrtDisc) / (2 * a), (-b - sqrtDisc) / (2 * a)}
	// 				nSolutions = 2
	// 			} else if disc < 0 {
	// 				solutions = []float64{}
	// 				nSolutions = 0
	// 			} else {
	// 				solutions = []float64{-b / (2 * a)}
	// 				nSolutions = 1
	// 			}

	// 		}
	// 	}

	// 	return solutions, nSolutions
}

// Returns a string expression of the polynomial in increasing sum form.
func (rp *RealPolynomial) Expr() string {
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

// Prints the string expression of the polynomial in increasing sum form to standard output.
func (rp *RealPolynomial) PrintExpr() {
	fmt.Print(rp.Expr())
}

// --- END STRUCT METHODS ---

// --- BEGIN UTILITY FUNCTIONS ---

func stripTailingZeroes(s []float64) []float64 {
	for s[len(s)-1] == 0.0 && len(s) > 1 {
		s = s[:len(s)-1]
	}
	return s
}

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

// Initializes and returns a new RealPolynomial struct with the given coefficients.
func NewRealPolynomial(coeffs []float64) (*RealPolynomial, error) {
	if len(coeffs) == 0 {
		return nil, errors.New("cannot create polynomial with no coefficients")
	}

	var newPolynomial RealPolynomial
	newPolynomial.coeffs = stripTailingZeroes(coeffs)

	return &newPolynomial, nil
}

// --- END UTILITY FUNCTIONS ---
