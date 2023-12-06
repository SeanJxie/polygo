package polygo

import (
	"log"
	"math"
	"math/rand"
)

// CauchyBound returns Cauchy's root bound of p.
//
// If p(x) = 0, then |x| <= p.CauchyBound().
//
// Panics for constant p.
func (p Poly) CauchyBound() float64 {

	if p.deg == 0 {
		log.Panic("CauchyBound: constant polynomial.")
	}

	// Compute Cauchy's bound.

	leadrecip := 1 / p.coef[p.deg]
	maxi := math.Abs(p.coef[0] * leadrecip)
	var tmp float64

	for i := 1; i < p.deg; i++ {
		tmp = math.Abs(p.coef[i] * leadrecip)

		if tmp > maxi {
			maxi = tmp
		}
	}

	return 1 + maxi
}

// CountSturm returns the number of distinct real roots of p on the interval (a, b].
func (p Poly) CountSturm(a, b float64) int {

	return cacheSturmChain(p).count(a, b)
}

// SolveNewtonRaphson implements the Newton-Raphson method for a single root with the intial guess
// and given number of iterations. An approximated root is returned.
//
// Panics for negative iterations.
func (p Poly) SolveNewtonRaphson(guess float64, iterations int) float64 {

	if iterations < 0 {
		log.Panicf("SolveNewtonRaphson: negative iterations %d.", iterations)
	}

	if iterations == 0 {
		return guess
	}

	d := p.Derivative()
	var dvalue float64

	for ; iterations > 0; iterations-- {

		dvalue = d.At(guess)

		if equalAbs(dvalue, 0, 0.001) {
			// Exit iteration if dvalue is close enough to zero.
			break
		}

		guess -= p.At(guess) / dvalue
	}

	return guess
}

// SolveBisect returns all distinct real roots of p on the interval [left, right] with given
// precision.
//
// Panics for invalid intervals and negative precision.
func (p Poly) SolveBisect(left, right, precision float64) float64 {

	if left > right {
		log.Panicf("SolveBisect: invalid interval [%f, %f].", left, right)
	}

	if precision < 0 {
		log.Panicf("SolveBisect: negative precision %f.", precision)
	}

	var mid float64

	for right-left > bisectPrecision {
		mid = 0.5 * (left + right)

		if p.CountSturm(left, mid) == 1 {
			right = mid
		} else {
			// If the root doesn't lie on the left-mid side, then it must lie on the
			// mid-right side.
			left = mid
		}
	}

	return mid
}

type CountAlgorithm int
type IsolateAlgorithm int
type SearchAlgorithm int

const (
	// ALG_COUNT represents the algorithm used for counting the number of roots on a half-open
	// interval (a, b].
	ALG_COUNT_STURM CountAlgorithm = iota

	// ALG_ISOLATE represents the algorithm used for isolating single roots to a sequence of
	// non-overlapping half-open intervals with form (a, b].
	ALG_ISOLATE_BISECT IsolateAlgorithm = iota

	// ALG_SEARCH represents the algorithm used for finding a single root on a half-open interval
	// (a, b].
	ALG_SEARCH_NEWTON SearchAlgorithm = iota
	ALG_SEARCH_BISECT
)

var (
	newtonIterations = 500
	bisectPrecision  = 1e-6
)

func (a CountAlgorithm) String() string {
	switch a {
	case ALG_COUNT_STURM:
		return "ALG_COUNT_STURM"
	}
	return "ALG_COUNT_UNKNOWN"
}

func (a IsolateAlgorithm) String() string {
	switch a {
	case ALG_ISOLATE_BISECT:
		return "ALG_ISOLATE_BISECT"
	}
	return "ALG_ISOLATE_UNKNOWN"
}

func (a SearchAlgorithm) String() string {
	switch a {
	case ALG_SEARCH_NEWTON:
		return "ALG_SEARCH_NEWTON"
	case ALG_SEARCH_BISECT:
		return "ALG_SEARCH_BISECT"
	}
	return "ALG_SEARCH_UNKNOWN"
}

// HalfOpenInterval represents a half-open interval (L, R].
type HalfOpenInterval struct {
	L, R float64
}

// Solver represents a collection of methods used to obtain information about polynomial equations.
type Solver struct {
	counter  CountAlgorithm
	isolator IsolateAlgorithm
	searcher SearchAlgorithm

	// Optional attributes (depends on algorithms used).
	chainCache map[uint32]sturmChain
}

// NewSolver returns a Solver equipped with the given root counting, isolation, and searching
// algorithms.
func NewSolver(counter CountAlgorithm, isolator IsolateAlgorithm, searcher SearchAlgorithm) Solver {

	return Solver{
		counter:    counter,
		isolator:   isolator,
		searcher:   searcher,
		chainCache: make(map[uint32]sturmChain),
	}
}

// NewSolverDefault returns a default solver.
//
// Specifically, the defaults are:
//   - Root counting algorithm:  ALG_COUNT_STURM
//   - Root isolation algorithm: ALG_ISOLATE_BISECT
//   - Root search algorithm:    ALG_SEARCH_BISECT
func NewSolverDefault() Solver {

	return NewSolver(ALG_COUNT_STURM, ALG_ISOLATE_BISECT, ALG_SEARCH_BISECT)
}

func (s Solver) cacheSturmChain(p Poly) sturmChain {

	id := p.id()
	cache := s.chainCache

	for key := range cache {
		if id == key {
			return cache[id]
		}
	}

	// Sturm chain has not been cached.
	cache[id] = new_sturmChain(p)

	return cache[id]
}

func (s Solver) CountRootsWithin(p Poly, a, b float64) int {

	var ret int

	switch s.counter {

	case ALG_COUNT_STURM:
		ret = s.cacheSturmChain(p).count(a, b)
	}

	return ret
}

// IsolateRoots returns a partition of the half-open interval (a, b] such that each half-open
// subinterval of the partition contains exactly one root of p.
func (s Solver) IsolateRootsWithin(p Poly, a, b float64) []HalfOpenInterval {

	partition := []HalfOpenInterval{}

	switch s.isolator {

	case ALG_ISOLATE_BISECT:

		c := s.CountRootsWithin(p, a, b)

		if c == 0 {
			return []HalfOpenInterval{}
		}

		if c == 1 {
			return []HalfOpenInterval{{a, b}}
		}

		m := 0.5 * (a + b)

		partition = append(s.IsolateRootsWithin(p, a, m), s.IsolateRootsWithin(p, m, b)...)
	}

	return partition
}

// solve_linear returns the root of linear p.
func solve_linear(p Poly) []float64 {

	return []float64{-p.coef[0] / p.coef[1]}
}

// solve_quadratic returns the real roots of quadratic p.
func solve_quadratic(p Poly) []float64 {

	c, negb, a := p.coef[0], -p.coef[1], p.coef[2]
	d := negb*negb - 4*a*c
	recip := 1 / (2 * a)

	if d > 0 {
		sqrtd := math.Sqrt(d)
		return []float64{(negb + sqrtd) * recip, (negb - sqrtd) * recip}
	}

	if d == 0 {
		return []float64{negb * recip}
	}

	return []float64{}
}

// solve_random_sample_newton returns the approximated roots of p on the interval (left, right].
func solve_random_sample_newton(p Poly, left, right float64) float64 {

	// Implement Newton-Raphson with random guess sampling.

	pprime := p.Derivative()
	root := left

	for !(left < root && root <= right) {

		root = left + rand.Float64()*(right-left)

		for i := 0; i < newtonIterations; i++ {
			fprimeroot := pprime.At(root)

			if fprimeroot == 0 {
				break
			}

			root -= p.At(root) / fprimeroot
		}
	}

	return root
}

// solve_bisect returns the approximated roots of p on the interval (left, right].
func solve_bisect(p Poly, left, right float64, counter func(Poly, float64, float64) int) float64 {

	var mid float64

	for right-left > bisectPrecision {
		mid = 0.5 * (left + right)

		if counter(p, left, mid) == 1 {
			right = mid
		} else {
			// If the root doesn't lie on the left-mid side, then it must lie on the
			// mid-right side.
			left = mid
		}
	}

	return mid
}

// FindRootsWithin returns the distinct roots of p on the half-open interval (a, b].
//
// Panics for invalid intervals and infinite solutions.
func (s Solver) FindRootsWithin(p Poly, a, b float64) []float64 {

	if b < a {
		log.Panicf("FindRootsWithin: invalid interval (%f, %f].", a, b)
	}

	roots := []float64{}

	// For deg(p) = 0, 1, 2, just solve exactly.
	if p.deg == 0 {
		if p.coef[0] == 0 {
			log.Panicf("FindRootsWithin: infinite solutions for %v.", p)
		}
		return []float64{}
	}

	// TODO: Check bounds a, b.
	if p.deg == 1 {
		return solve_linear(p)
	}

	if p.deg == 2 {
		return solve_quadratic(p)
	}

	intervals := s.IsolateRootsWithin(p, a, b)

	switch s.searcher {

	case ALG_SEARCH_NEWTON:

		for _, h := range intervals {
			roots = append(roots, solve_random_sample_newton(p, h.L, h.R))
		}

	case ALG_SEARCH_BISECT:

		for _, h := range intervals {
			roots = append(roots, solve_bisect(p, h.L, h.R, s.CountRootsWithin))
		}
	}

	return roots
}

// FindRoots returns all distinct roots of p.
func (s Solver) FindRoots(p Poly) []float64 {

	bound := p.CauchyBound()
	return s.FindRootsWithin(p, -bound, bound)
}

// Point represents a 2D Cartesian coordinate.
type Point struct {
	X, Y float64
}

// FindIntersectionsWithin returns the intersections of p and q on the half-open interval (a, b].
//
// Panics for invalid intervals.
func (s Solver) FindIntersectionsWithin(p, q Poly, a, b float64) []Point {

	if b < a {
		log.Panicf("FindIntersectionsWithin: invalid interval (%f, %f].", a, b)
	}

	xinter := s.FindRootsWithin(p.Sub(q), a, b)
	points := make([]Point, len(xinter))

	for i, x := range xinter {
		points[i] = Point{X: x, Y: p.At(x)}
	}

	return points
}

// FindIntersections returns all intersections of p and q.
func (s Solver) FindIntersections(p, q Poly) []Point {

	xinter := s.FindRoots(p.Sub(q))

	points := make([]Point, len(xinter))

	for i, x := range xinter {
		points[i] = Point{X: x, Y: p.At(x)}
	}

	return points
}

// SetNewtonSearchIterations sets the Newton's method root search algorithm iterations to v.
//
// Panics for negative v.
func SetNewtonSearchIterations(v int) {
	if v < 0 {
		log.Panic("SetNewtonSearchIterations: negative v.")
	}

	newtonIterations = v
}

// SetBisectSearchPrecision sets the bisect root search algorithm precision to v.
//
// The closer ot zero v is, the more accurate the roots.
//
// Panics for negative v.
func SetBisectSearchPrecision(v float64) {
	if v < 0 {
		log.Panic("SetBisectSearchPrecision: negative v.")
	}

	bisectPrecision = v
}
