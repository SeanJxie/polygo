package polygo

/*
This file contains polynomial solvers and related algorithms.
*/

import (
	"errors"
)

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
	return removeDuplicateFloat(rp.findRootsWithinAcc(a, b, nil, rp.sturmChain())), nil
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

// FindIntersectionWithin returns ANY intersection point of the current instance and rp2 existing on the closed interval [a, b].
// If there are no intersections on the provided interval, an error is set.
func (rp *RealPolynomial) FindIntersectionWithin(a, b float64, rp2 *RealPolynomial) (Point, error) {
	tmp := *rp

	root, err := (&tmp).Sub(rp2).FindRootWithin(a, b)
	if err != nil {
		return Point{}, err
	}

	point := Point{root, rp.At(root)}
	return point, nil
}

// FindIntersectionsWithin returns ALL intersection point of the current instance and rp2 existing on the closed interval [a, b].
// Unlike FindIntersectionWithin, no error is set if there are no intersections on the provided interval. Instead, an empty slice is returned.
// If there are an infinite number or solutions, an error is set.
func (rp *RealPolynomial) FindIntersectionsWithin(a, b float64, rp2 *RealPolynomial) ([]Point, error) {
	if rp == nil || rp2 == nil {
		panic("received nil *RealPolynomial")
	}

	tmp := *rp
	roots, err := tmp.Sub(rp2).FindRootsWithin(a, b)
	if err != nil {
		return nil, err
	}

	roots = removeDuplicateFloat(roots)
	points := make([]Point, len(roots))

	for i, x := range roots {
		points[i] = Point{x, rp.At(x)}
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

End "WSC"-suffixed (and related) functions.

*/
