package polygo

import "log"

var (
	// Stores already-computed Sturm chains during runtime.
	chainCache = make(map[uint32]sturmChain)
)

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
	// if va-vb < 0 {
	// 	return 0
	// }

	return va - vb
}

func cacheSturmChain(p Poly) sturmChain {

	id := p.id()

	for key := range chainCache {
		if id == key {
			return chainCache[id]
		}
	}

	// Sturm chain has not been cached.
	chainCache[id] = new_sturmChain(p)

	return chainCache[id]
}
