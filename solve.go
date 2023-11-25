package polygo

import (
	"log"
	"math/rand"
)

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
	newtonIterations = 100
	bisectPrecision  = epsilon
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

// sturmChain represents the Sturm chain (or sequence) of a Poly.
type sturmChain struct {
	c   []Poly
	len int
}

// computeSturmChain computes and caches the Sturm chain of p.
func new_sturmChain(p Poly) sturmChain {

	// Constant case.
	if p.deg == 0 {
		return sturmChain{[]Poly{p}, 1}
	}

	// Construct Sturm chain.
	chain := make([]Poly, p.deg+1)
	chain[0], chain[1] = p, p.Derivative()

	i := 1
	var rem Poly

	for !chain[i].IsConstant() {
		_, rem = chain[i-1].Div(chain[i])
		chain[i+1] = rem.MulScalar(-1)
		i++
	}

	// Resize.
	i++
	chain = chain[:i]

	return sturmChain{
		c:   chain,
		len: i,
	}
}

// count returns the number of roots on the half-open interval (a, b] for p associated with s.
//
// Panics for invalid intervals.
func (s sturmChain) count(a, b float64) int {

	if a > b {
		log.Panicf("count: invalid interval (%f, %f].", a, b)
	}

	if s.len == 1 {
		return 0 // No sign changes in one value.
	}

	chain := s.c

	prevsa := sign(chain[0].At(a))
	prevsb := sign(chain[0].At(b))
	var currsa, currsb int

	// Sign change counters.
	va := 0
	vb := 0

	for i := 1; i < s.len; i++ {
		currsa = sign(chain[i].At(a))
		currsb = sign(chain[i].At(b))

		if currsa != prevsa && prevsa != 0 {
			va++
		}

		if currsb != prevsb && prevsb != 0 {
			vb++
		}

		prevsa = currsa
		prevsb = currsb
	}

	// Hacky fix for when a, b are multiple roots of p.
	if va-vb < 0 {
		return 0
	}

	return va - vb
}

// HalfOpenInterval represents a half-open interval (L, R].
type HalfOpenInterval struct {
	L, R float64
}

// Solver represents a collection of methods used to solve polynomial equations.
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

// FindRootsWithin returns the distinct roots of p on the half-open interval (a, b].
//
// Panics for invalid intervals.
func (s Solver) FindRootsWithin(p Poly, a, b float64) []float64 {

	if b < a {
		log.Panicf("FindRootsWithin: invalid interval (%f, %f].", a, b)
	}

	roots := []float64{}

	// For deg(p) = 0, 1, 2, just solve exactly.
	if p.deg == 0 {

	}

	intervals := s.IsolateRootsWithin(p, a, b)

	switch s.searcher {

	case ALG_SEARCH_NEWTON:

		// Implement Newton-Raphson with random guess sampling.

		fprime := p.Derivative()

		var root, fprimeroot, left, right float64

		for _, h := range intervals {

			left, right = h.L, h.R

		random_sample_newton:

			// Random sample on interval.
			root = left + rand.Float64()*(right-left)

			for i := 0; i < newtonIterations; i++ {
				fprimeroot = fprime.At(root)

				if fprimeroot == 0 {
					break
				}

				root -= p.At(root) / fprimeroot
			}

			// Excluding root found outside of interval.
			if !(left < root && root <= right) {
				goto random_sample_newton
			}

			roots = append(roots, root)
		}

	case ALG_SEARCH_BISECT:

		var left, right, mid float64

		for _, h := range intervals {
			left, right = h.L, h.R

			for right-left > bisectPrecision {
				mid = 0.5 * (left + right)

				if s.CountRootsWithin(p, left, mid) == 1 {
					right = mid
				} else {
					// If the root doesn't lie on the left-mid side, then it must lie on the
					// mid-right side.
					left = mid
				}
			}

			roots = append(roots, mid)
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
