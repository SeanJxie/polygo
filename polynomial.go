package polygo

import (
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/mjibson/go-dsp/fft"
)

// A Poly represents a univariate real polynomial.
//
// Note: in the documentation for each method of Poly, we refer to the receiver instance as "p".
type Poly struct {
	coef []float64
	len  int
	deg  int
}

// NewPoly returns a polynomial p with the given coefficients.
//
// Let c = coefficients and let n = len(c). Then, p is defined by
//
//   - p(x) = c[0]x^(n-1) + c[1]x^(n-2) + ... + c[n-2]x^1 + c[n-1]x^0.
//
// # Examples:
//   - NewPoly([]float64{3, -1, 4}) represents p(x) = 3x^2 - x + 4.
//   - NewPoly([]float64{0}) represents p(x) = 0.
//   - NewPoly([]float64{0, 0, 0, 2, 0, 0, 7, 0, 0, 1}) represents p(x) = 2x^6 + 7x^3 + 1.
//
// Panics if coefficients slice is empty.
func NewPoly(coefficients []float64) Poly {

	if len(coefficients) == 0 {
		// When dealing with invalid inputs, polygo will not use the
		// "return error" convention in order to keep user code less
		// cluttered. Instead, functions will panic (as opposed to Fatal,
		// since Fatal calls os.exit(1), whereas panic works it's way up
		// the call stack and returns a useful stacktrace so we know where
		// things are going wrong).
		log.Panic("NewPoly: empty coefficients slice.")
	}

	// Makes things easier internally to have the degree of a term be the index of its coefficient,
	// so we reverse the coefficients slice.
	//
	// Also, we don't want to deal with arbitrary lengths of leading zeroes in the code. So, we
	// strip the leading zeroes. This also guarantees that the degree of the polynomial is the
	// length of the coefficient slice minus 1 (which also happens to be the largest index of the
	// coefficient slice).
	coefficients = removeTrailingZeroes(reverse(coefficients))
	coefLen := len(coefficients)

	ret := Poly{
		coef: coefficients,
		len:  coefLen,
		deg:  coefLen - 1,
	}

	return ret
}

// newPolyNoReverse is just NewPoly but with no coefficient slice reversal.
//
// This is needed because we internally represent coefs as a slice with increasing degree, whereas
// the user interacts with coefs slices with decreasing degree. So, having this function allows us
// to return new Polys without having to compensate for the difference in representation.
//
// Doesn't do the empty panic check like in NewPoly().
func newPolyNoReverse(coefficients []float64) Poly {

	coefficients = removeTrailingZeroes(coefficients)
	coefLen := len(coefficients)

	ret := Poly{
		coef: coefficients,
		len:  coefLen,
		deg:  coefLen - 1,
	}

	return ret
}

// parseTerm returns the coefficient and exponent of a term in string form "(+|-)cx^n".
//
// Panics for invalid terms.
func parseTerm(t string) (float64, int) {

	var xpos, caratpos int
	var sign, coef float64
	var deg int64
	var err error

	if t[0] == '+' {
		sign = float64(1.0)
	} else {
		sign = float64(-1.0)
	}

	t = t[1:]

	xpos = strings.IndexByte(t, 'x')
	caratpos = strings.IndexByte(t, '^')

	// deg(t) = 0 or 1 (at least '^' is missing).
	if caratpos == -1 {

		// deg(t) = 0 ('^' and 'x' are missing).
		if xpos == -1 {
			if coef, err = strconv.ParseFloat(t, 64); err != nil {
				log.Panicf("parseTerm: could not parse deg 0 term coefficient \"%s\" (%v).",
					t, err)
			}

			return sign * coef, 0
		}

		// deg(t) = 1 (Only '^' is missing).
		if t[:xpos] == "" {
			coef = 1
		} else if coef, err = strconv.ParseFloat(t[:xpos], 64); err != nil {
			log.Panicf("parseTerm: could not parse deg 1 term coefficient \"%s\" (%v).",
				t[:xpos], err)
		}

		return sign * coef, 1
	}

	// deg(t) > 1.
	if deg, err = strconv.ParseInt(t[caratpos+1:], 10, 64); err != nil {
		log.Panicf("parseTerm: could not parse exponent \"%s\" (%v).", t[caratpos+1:], err)
	}

	if t[:xpos] == "" {
		coef = 1
	} else if coef, err = strconv.ParseFloat(t[:xpos], 64); err != nil {
		log.Panicf("parseTerm: could not parse deg %d term coefficient %s (%v).", deg, t, err)
	}

	return sign * coef, int(deg)
}

// NewPolyFromString returns a polynomial represented by s.
//
// # Format:
//   - Terms (without sign) have the form "cx^n", with c real and n natural (including 0).
//   - Aside from the leading term, all terms must be prefixed (spaces ignored) by a "+" or a "-"
//     denoting the sign of the term. The user may choose to omit the sign on the leading term,
//     in which case it is assumed to be positive.
//   - Terms do not need to be ordered, nor does the user have to include terms with coefficient
//     zero.
//   - The user may include multiple terms of the same degree.
//
// # Examples:
//   - NewPolyFromString("5") represents p(x) = 5
//   - NewPolyFromString("- 4 + 3x^2 - 2x") represents p(x) = 3x^2 - 2x - 4.
//   - NewPolyFromString("5x^10 + 0x^9 - 6x + 3x^2 - 2x") represents p(x) = 5x^10 + 3x^2 - 8x.
//
// Panics on empty or invalid strings.
func NewPolyFromString(s string) Poly {

	if s == "" {
		log.Panic("NewPolyFromString: empty string.")
	}

	// Manually insert implicit leading plus if the first non-whitespace
	// char is not "+" or "-".
	i := 0
	for s[i] == ' ' && i < len(s) {
		i++
	}

	if !plusOrMinus(rune(s[i])) {
		s = "+" + s
	}

	coefs := []float64{}

	var coef float64
	var deg int
	var termsb strings.Builder

	for _, c := range s {
		if c != ' ' {

			if plusOrMinus(c) && termsb.Len() != 0 {

				coef, deg = parseTerm(termsb.String())
				coefs = expand(coefs, deg+1)
				coefs[deg] += coef

				termsb.Reset()
			}

			termsb.WriteRune(c)
		}
	}

	coef, deg = parseTerm(termsb.String())
	coefs = expand(coefs, deg+1)
	coefs[deg] += coef

	return newPolyNoReverse(coefs)
}

// NewPolyConst returns the polynomial p(x) = a.
func NewPolyConst(a float64) Poly {

	return newPolyNoReverse([]float64{a})
}

// NewPolyZero returns the polynomial p(x) = 0.
func NewPolyZero() Poly {

	return NewPolyConst(0)
}

// NewPolyLinear returns the polynomial p(x) = ax + b.
func NewPolyLinear(a, b float64) Poly {

	return newPolyNoReverse([]float64{b, a})
}

// NewPolyQuadratic returns the polynomial p(x) = ax^2 + bx + c.
func NewPolyQuadratic(a, b, c float64) Poly {

	return newPolyNoReverse([]float64{c, b, a})
}

// NewPolyCubic returns the polynomial p(x) = ax^3 + bx^2 + cx + d.
func NewPolyCubic(a, b, c, d float64) Poly {

	return newPolyNoReverse([]float64{d, c, b, a})
}

// NewPolyWilkinson returns Wilkinson's polynomial.
func NewPolyWilkinson() Poly {

	return NewPoly([]float64{
		1,
		-210,
		20615,
		-1256850,
		53327946,
		-1672280820,
		40171771630,
		-756111184500,
		11310276995381,
		-135585182899530,
		1307535010540395,
		-10142299865511450,
		63030812099294896,
		-311333643161390640,
		1206647803780373360,
		-3599979517947607200,
		8037811822645051776,
		-12870931245150988800,
		13803759753640704000,
		-8752948036761600000,
		2432902008176640000},
	)
}

// NewPolyFactored returns the polynomial
//
// p(x) = a(x - r[0])(x - r[1])...(x - r[n - 1]),
//
// where n = len(r).
//
// Panics for empty r.
func NewPolyFactored(a float64, r []float64) Poly {

	if len(r) == 0 {
		log.Panic("NewPolyFactored: empty r.")
	}

	if a == 0 {
		return NewPolyZero()
	}

	prod := newPolyNoReverse([]float64{-r[0], 1})

	r = r[1:]

	for _, b := range r {
		prod = prod.Mul(newPolyNoReverse([]float64{-b, 1}))
	}

	return prod.MulScalar(a)
}

// NewPolyTaylorSin returns the Taylor polynomial of the sine function centered at a with degree n.
//
// Panics for negative n.
func NewPolyTaylorSin(n int, a float64) Poly {

	if n < 0 {
		log.Panic("NewPolyTaylorSin: negative n.")
	}

	if n == 0 {
		return NewPolyConst(math.Sin(a))
	}

	sina := math.Sin(a)
	cosa := math.Cos(a)

	derivCycle := [4]float64{
		sina,
		cosa,
		-sina,
		-cosa,
	}

	sum := NewPolyZero()
	for i := 0; i <= n; i++ {
		sum = sum.Add(NewPolyLinear(1, -a).Pow(i).MulScalar(derivCycle[i%4] / fact(i)))
	}

	return sum
}

// NewPolyChebyshev1 returns the nth Chebyshev polynomial of the first kind.
//
// Panics for negative n.
func NewPolyChebyshev1(n int) Poly {

	if n < 0 {
		log.Panic("NewPolyChebyshev1: negative n.")
	}

	if n == 0 {
		return NewPolyConst(1)
	}

	if n == 1 {
		return NewPolyLinear(1, 0)
	}

	return NewPolyLinear(2, 0).Mul(NewPolyChebyshev1(n - 1)).Sub(NewPolyChebyshev1(n - 2))
}

// NewPolyChebyshev2 returns the nth Chebyshev polynomial of the second kind.
//
// Panics for negative n.
func NewPolyChebyshev2(n int) Poly {

	if n < 0 {
		log.Panic("NewPolyChebyshev2: negative n.")
	}

	if n == 0 {
		return NewPolyConst(1)
	}

	if n == 1 {
		return NewPolyLinear(2, 0)
	}

	return NewPolyLinear(2, 0).Mul(NewPolyChebyshev2(n - 1)).Sub(NewPolyChebyshev2(n - 2))
}

// Coefficients returns the coefficients c of p ordered in decreasing degree.
func (p Poly) Coefficients() []float64 {

	return reverse(p.coef)
}

// Degree returns the degree of p.
func (p Poly) Degree() int {

	return p.deg
}

// LeadingCoefficient returns the coefficient of the highest-degreed term in p.
func (p Poly) LeadingCoefficient() float64 {

	return p.coef[p.deg]
}

// LargestCoefficient returns the largest coefficient in p.
func (p Poly) LargestCoefficient() float64 {

	return max(p.coef)
}

// SmallestCoefficient returns the smallest coefficient in p.
func (p Poly) SmallestCoefficient() float64 {

	return min(p.coef)
}

// CoefficientWithDegree returns the coefficient of the term with degree n in p.
func (p Poly) CoefficientWithDegree(n uint) float64 {

	// Coefficients of terms with degrees larger than that of p are
	// zero by definition.
	if n > uint(p.deg) {
		return 0.0
	}

	return p.coef[n]
}

// Equal returns true if the p is equal to q (all corresponding coefficients are equal), else false.
func (p Poly) Equal(q Poly) bool {

	if p.deg != q.deg {
		return false
	}

	for i := 0; i < p.len; i++ {
		if p.coef[i] == q.coef[i] {
			return false
		}
	}

	return true
}

// EqualRel returns true if the largest relative difference between corresponding coefficients is
// maxPercentErr, else false.
func (p Poly) EqualRel(q Poly, maxPercentErr float64) bool {

	if p.deg != q.deg {
		return false
	}

	for i := 0; i < p.len; i++ {
		if equalRel(p.coef[i], q.coef[i], maxPercentErr) {
			return false
		}
	}

	return true
}

// IsConstant returns true p is constant (i.e. deg(p) = 0), else false.
func (p Poly) IsConstant() bool {

	return p.deg == 0
}

// IsZero returns true if p(x) = 0, else false.
func (p Poly) IsZero() bool {

	// Check if p is a constant and if that constant is 0.
	return p.deg == 0 && p.coef[0] == 0
}

// IsZeroRel returns true if the largest relative difference between p and 0 is maxPercentErr,
// else false.
func (p Poly) IsZeroRel(maxPercentErr float64) bool {

	return p.deg == 0 && equalRel(p.coef[0], 0, maxPercentErr)
}

// IsMonic returns true p is monic (i.e. leading coefficient 1), else false.
func (p Poly) IsMonic() bool {

	return p.coef[p.deg] == 1
}

// At returns the value of p evaluated at x.
func (p Poly) At(x float64) float64 {

	// Implement Horner's scheme.
	out := p.coef[p.deg]
	for i := p.deg - 1; i >= 0; i-- {
		out = out*x + p.coef[i]
	}

	return out
}

// Add returns the polynomial sum p + q.
func (p Poly) Add(q Poly) Poly {

	var max int
	if p.len > q.len {
		max = p.len
	} else {
		max = q.len
	}

	// Pad the shorter polynomial with zeroes to align.
	pe := expand(p.coef, max)
	qe := expand(q.coef, max)

	sumCoef := make([]float64, max)

	// Add like terms.
	for i := 0; i < max; i++ {
		sumCoef[i] = pe[i] + qe[i]
	}

	return newPolyNoReverse(sumCoef)
}

// Sub returns the polynomial difference p - q.
func (p Poly) Sub(q Poly) Poly {

	var max int
	if p.len > q.len {
		max = p.len
	} else {
		max = q.len
	}

	pe := expand(p.coef, max)
	qe := expand(q.coef, max)

	difCoef := make([]float64, max)

	for i := 0; i < max; i++ {
		difCoef[i] = pe[i] - qe[i]
	}

	return newPolyNoReverse(difCoef)
}

// MulScalar returns the scalar-polynomial product sp.
func (p Poly) MulScalar(s float64) Poly {

	// 0 * p = 0.
	if s == 0 {
		return NewPoly([]float64{0})
	}

	prodCoef := make([]float64, p.len)
	for i, c := range p.coef {
		prodCoef[i] = s * c
	}

	return newPolyNoReverse(prodCoef)
}

// Mul returns the polynomial product pq.
func (p Poly) Mul(q Poly) Poly {

	// The product m will have deg(m) = deg(p) + deg(q).
	// We add 1 since degree is one less than length of the coefficient slice.
	prodCoef := make([]float64, p.deg+q.deg+1)

	for i := 0; i < p.len; i++ {
		for j := 0; j < q.len; j++ {
			prodCoef[i+j] += p.coef[i] * q.coef[j]
		}
	}

	return newPolyNoReverse(prodCoef)
}

// MulFast returns the polynomial product pq.
//
// This method uses an FFT algorithm to perform fast polynomial multiplication in O(n log n) time at
// the price of small floating point errors.
//
// MulFast() should be used when precision is flexible are not rigorous and speed is a requirement.
// If equality must be checked, use EqualWithin() instead of Equal().
func (p Poly) MulFast(q Poly) Poly {

	// Algorithm reference:
	// https://faculty.sites.iastate.edu/jia/files/inline-files/polymultiply.pdf

	if p.deg == 0 {
		return q.MulScalar(p.coef[0])
	}

	if q.deg == 0 {
		return p.MulScalar(q.coef[0])
	}

	// Pad the length of the product coefficient slice to a power of 2 for an efficient FFT.
	prodlen := p.deg + q.deg + 1
	potlen := nextPOT(prodlen)

	// Evaluation to point-value representation.
	// Since len(a) and len(b) are powers of 2, the call to
	// fft.FFT() implicitly calls the radix2FFT() function,
	// which implements the radix-2 DIT Cooley-Tukey algorithm
	// (with small floating point error).
	a := fft.FFT(toComplex128(expand(p.coef, potlen)))
	b := fft.FFT(toComplex128(expand(q.coef, potlen)))

	// Pointwise multiplication.
	c := make([]complex128, potlen)
	for i := 0; i < potlen; i++ {
		c[i] = a[i] * b[i]
	}

	// Interpolation to coefficient slice.
	//
	// We manually cut the slice off at the expected product length
	// since floating point error may cause coefficients that are
	// supposed to be zero to be nonzero. This trips up the call to
	// removeTrailingZeroes within newPolyNoReverse and we end up
	// with a polynomial product with nonexistent nonzero leading
	// coefficeints of degree larger than the expected product (p.deg + q.deg).
	return newPolyNoReverse(toFloat64(fft.IFFT(c))[:prodlen])
}

// Pow returns the polynomial power p^n.
//
// Panics for negative n.
func (p Poly) Pow(n int) Poly {

	if n < 0 {
		log.Panic("Pow: negative n.")
	}

	prod := NewPolyConst(1)

	for i := 0; i < n; i++ {
		prod = prod.Mul(p)
	}

	return prod
}

// PowFast returns the polynomial power p^n.
//
// Be sure to read the documentation for MulFast(), as the behaviour is the same.
//
// Panics for negative n.
func (p Poly) PowFast(n int) Poly {

	if n < 0 {
		log.Panic("PowFast: negative n.")
	}

	prod := NewPolyConst(1)

	for i := 0; i < n; i++ {
		prod = prod.MulFast(p)
	}

	return prod
}

// Div returns m (polynomial quotient) and n (polynomial remainder) such that p/q = m + n/q.
//
// Panics if q = 0.
func (p Poly) Div(q Poly) (Poly, Poly) {

	// Dividing by zero.
	if q.IsZero() {
		log.Panic("Div: division by zero polynomial.")
	}

	// Dividing zero.
	if p.IsZero() {
		return NewPolyZero(), NewPolyZero()
	}

	// Dividing by larger degree.
	if p.deg < q.deg {
		return NewPolyZero(), p
	}

	// Implement expanded synthetic division for non-monic divisors.

	pRev := reverse(p.coef)
	qRev := reverse(q.coef)

	quoRemCoef := make([]float64, p.len)
	copy(quoRemCoef, pRev)

	lead := qRev[0]
	sep := p.len - q.len + 1

	for i := 0; i < sep; i++ {
		quoRemCoef[i] /= lead

		if c := quoRemCoef[i]; c != 0 {

			for j := 1; j < q.len; j++ {
				quoRemCoef[i+j] += -qRev[j] * c
			}
		}
	}

	quoCoef := reverse(quoRemCoef[:sep])
	remCoef := reverse(quoRemCoef[sep:])

	return newPolyNoReverse(quoCoef), newPolyNoReverse(remCoef)
}

// Derivative returns the derivative p' of p.
func (p Poly) Derivative() Poly {

	// p is a constant, whose derivative is always 0.
	if p.deg == 0 {
		return NewPoly([]float64{0})
	}

	// For nonconstant p, if deg(p) = n, then deg(p') = n - 1.
	derivCoef := make([]float64, p.deg)
	for i := 0; i < p.deg; i++ {
		derivCoef[i] = p.coef[i+1] * float64(i+1)
	}

	return newPolyNoReverse(derivCoef)
}

// Reciprocal returns the reciprocal p* of p.
func (p Poly) Reciprocal() Poly {

	// Since we reverse the user's coefficient slice in NewPoly(), we just pass
	// it back into NewPoly() to reverse the coefficient slice again.

	return NewPoly(p.coef)
}

// String returns a string representation of p in decreasing-degree sum form.
func (p Poly) String() string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[ %fx^{%d}", p.coef[p.deg], p.deg))

	var sgn, strCoef string
	for i := 1; i < p.len; i++ {

		strCoef = fmt.Sprintf("%f", p.coef[p.deg-i])

		if sign(p.coef[p.deg-i]) == -1 {
			sgn = " - "
			strCoef = strCoef[1:]
		} else {
			sgn = " + "
		}

		sb.WriteString(sgn)
		sb.WriteString(fmt.Sprintf("%sx^{%d}", strCoef, p.deg-i))
	}

	sb.WriteString(" ]")

	return sb.String()
}

// Stringn returns a string representation of p in decreasing-degree sum form with it's coefficients
// to precision n.
//
// All n < 0 will be treated as n = 0.
func (p Poly) Stringn(n int) string {

	var sb strings.Builder

	precisionFormat := fmt.Sprintf(".%d", n)
	if n < 0 {
		precisionFormat = ".0"
	}

	sb.WriteString(fmt.Sprintf("[ %"+precisionFormat+"fx^{%d}", p.coef[p.deg], p.deg))

	var sgn, strCoef string
	for i := 1; i < p.len; i++ {

		strCoef = fmt.Sprintf("%"+precisionFormat+"f", p.coef[p.deg-i])

		if sign(p.coef[p.deg-i]) == -1 {
			sgn = " - "
			strCoef = strCoef[1:]
		} else {
			sgn = " + "
		}

		sb.WriteString(sgn)
		sb.WriteString(fmt.Sprintf("%sx^{%d}", strCoef, p.deg-i))
	}

	sb.WriteString(" ]")

	return sb.String()
}

// Printn prints p to standard output with it's coefficients printed to precision n followed by a
// newline.
//
// All n < 0 will be treated as n = 0.
func (p Poly) Printn(n int) {

	fmt.Println(p.Stringn(n))
}

// id returns a unqiue identifier for p.
func (p Poly) id() uint32 {

	// Generate unqiue string and hash.

	var sb strings.Builder

	for _, c := range p.coef {
		sb.WriteString(fmt.Sprintf("%f,", c))
	}

	h := fnv.New32()
	h.Write([]byte(sb.String()))

	return h.Sum32()
}
