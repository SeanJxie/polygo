package polygo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_new_sturmChain(t *testing.T) {

	testCases := []struct {
		name      string
		arg       Poly
		wantChain []Poly
		wantLen   int
	}{
		{
			name:      "linear",
			arg:       NewPoly([]float64{1, 0}),
			wantChain: []Poly{NewPoly([]float64{1, 0}), NewPoly([]float64{1})},
			wantLen:   2,
		},
		{
			name: "quartic",
			arg:  NewPoly([]float64{1, 1, 0, -1, -1}),
			wantChain: []Poly{
				NewPoly([]float64{1, 1, 0, -1, -1}),
				NewPoly([]float64{4, 3, 0, -1}),
				NewPoly([]float64{3. / 16, 3. / 4, 15. / 16}),
				NewPoly([]float64{-32, -64}),
				NewPoly([]float64{-3. / 16}),
			},
			wantLen: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := new_sturmChain(tc.arg)

			assert.Equal(t, tc.wantChain, got.c)
			assert.Equal(t, tc.wantLen, got.len)
		})
	}
}

func Test_sturmChain_countPanic(t *testing.T) {

	assert.Panics(t, func() { new_sturmChain(NewPoly([]float64{1, 0})).count(2, 0) })
	assert.Panics(t, func() { new_sturmChain(NewPoly([]float64{1, 0})).count(1, 0) })
}

func Test_sturmChain_count(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argA float64
		argB float64
		want int
	}{
		{
			name: "linear, 0 in (-1, 0]",
			argP: NewPoly([]float64{1, 0}),
			argA: -1,
			argB: 0,
			want: 1,
		},
		{
			name: "linear, 0 not in (0, 1]",
			argP: NewPoly([]float64{1, 0}),
			argA: 0,
			argB: 1,
			want: 0,
		},
		{
			name: "cubic 3 roots",
			argP: NewPoly([]float64{1, -1, -1, 0}),
			argA: -2,
			argB: 2,
			want: 3,
		},
		{
			name: "cubic 2 roots",
			argP: NewPoly([]float64{1, -1, -1, 0}),
			argA: -2,
			argB: 0,
			want: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := new_sturmChain(tc.argP).count(tc.argA, tc.argB)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_NewSolver(t *testing.T) {

	testCases := []struct {
		name string
		argC CountAlgorithm
		argI IsolateAlgorithm
		argS SearchAlgorithm

		wantP Poly
		wantC CountAlgorithm
		wantI IsolateAlgorithm
		wantS SearchAlgorithm

		wantChainCache map[uint32]sturmChain
	}{
		{
			name:           "linear",
			argC:           ALG_COUNT_STURM,
			argI:           ALG_ISOLATE_BISECT,
			argS:           ALG_SEARCH_NEWTON,
			wantC:          ALG_COUNT_STURM,
			wantI:          ALG_ISOLATE_BISECT,
			wantS:          ALG_SEARCH_NEWTON,
			wantChainCache: map[uint32]sturmChain{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewSolver(tc.argC, tc.argI, tc.argS)

			assert.Equal(t, tc.wantC, got.counter)
			assert.Equal(t, tc.wantI, got.isolator)
			assert.Equal(t, tc.wantS, got.searcher)
			assert.Equal(t, tc.wantChainCache, got.chainCache)
		})
	}
}

func Test_NewSolverDefault(t *testing.T) {

	s := NewSolverDefault()

	assert.Equal(t, ALG_COUNT_STURM, s.counter)
	assert.Equal(t, ALG_ISOLATE_BISECT, s.isolator)
	assert.Equal(t, ALG_ISOLATE_BISECT, s.searcher)
}

func Test_SolverCountRootsWithin(t *testing.T) {

	testCases := []struct {
		name string
		argS Solver
		argP Poly
		argA float64
		argB float64
		want int
	}{
		{
			name: "linear, 0 in (-1, 0]",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, 0}),
			argA: -1,
			argB: 0,
			want: 1,
		},
		{
			name: "linear, 0 not in (0, 1]",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, 0}),
			argA: 0,
			argB: 1,
			want: 0,
		},
		{
			name: "cubic 3 roots",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, -1, -1, 0}),
			argA: -2,
			argB: 2,
			want: 3,
		},
		{
			name: "cubic 2 roots",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, -1, -1, 0}),
			argA: -2,
			argB: 0,
			want: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argS.CountRootsWithin(tc.argP, tc.argA, tc.argB)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_SolverIsolateRootsWithin(t *testing.T) {

	testCases := []struct {
		name string
		argS Solver
		argP Poly
		argA float64
		argB float64
		want []HalfOpenInterval
	}{
		{
			name: "linear on interval (-1, 0]",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, 0}),
			argA: -1,
			argB: 0,
			want: []HalfOpenInterval{{-1, 0}},
		},
		{
			name: "linear off interval (0, 1]",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, 0}),
			argA: 0,
			argB: 1,
			want: []HalfOpenInterval{},
		},
		{
			name: "quartic",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, 1, 0, -1, -1}),
			argA: -2,
			argB: 2,
			want: []HalfOpenInterval{{-2, 0}, {0, 2}},
		},
		{
			name: "cubic",
			argS: NewSolverDefault(),
			argP: NewPoly([]float64{1, -1, -1, 0}),
			argA: -2,
			argB: 2,
			want: []HalfOpenInterval{{-1, -0.5}, {-0.5, 0}, {0, 2}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argS.IsolateRootsWithin(tc.argP, tc.argA, tc.argB)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_SolverFindRoots(t *testing.T) {

	p := NewPolyFromString("x^21 - 86400x + 86399")
	q := NewPolyFromString("x^2-4x+4")
	s := NewSolver(ALG_COUNT_STURM, ALG_ISOLATE_BISECT, ALG_SEARCH_BISECT)

	t.Log(s.FindRoots(p))
	t.Log(s.FindRoots(q))
	tmp := NewPolyWilkinson()
	t.Log(s.FindRoots(tmp))
	t.Log(s.FindRoots(tmp))

	a := NewPolyFromString("x^3 - x^2 - x")
	b := NewPolyFromString("x^2")
	t.Log(s.FindIntersections(a, b))
}
