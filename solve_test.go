package polygo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PolyCauchyBoundPanic(t *testing.T) {

	assert.Panics(t, func() { NewPolyZero().CauchyBound() })
	assert.Panics(t, func() { NewPoly([]float64{3.14}).CauchyBound() })
}

func Test_PolyCauchyBound(t *testing.T) {
	testCases := []struct {
		name string
		arg  Poly
		want float64
	}{
		{
			name: "linear",
			arg:  NewPoly([]float64{125124, 1.1715}),
			want: 1.0000093627121895,
		},
		{
			name: "quadratic",
			arg:  NewPoly([]float64{1, -1, 0}),
			want: 2,
		},
		{
			name: "wilkinson",
			arg:  NewPolyWilkinson(),
			want: 1.3803759753640704e19,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.arg.CauchyBound()

			assert.Equal(t, tc.want, got)
		})
	}
}

// func Test_NewSolver(t *testing.T) {

// 	testCases := []struct {
// 		name string
// 		argC CountAlgorithm
// 		argI IsolateAlgorithm
// 		argS SearchAlgorithm

// 		wantP Poly
// 		wantC CountAlgorithm
// 		wantI IsolateAlgorithm
// 		wantS SearchAlgorithm

// 		wantChainCache map[uint32]sturmChain
// 	}{
// 		{
// 			name:           "linear",
// 			argC:           ALG_COUNT_STURM,
// 			argI:           ALG_ISOLATE_BISECT,
// 			argS:           ALG_SEARCH_NEWTON,
// 			wantC:          ALG_COUNT_STURM,
// 			wantI:          ALG_ISOLATE_BISECT,
// 			wantS:          ALG_SEARCH_NEWTON,
// 			wantChainCache: map[uint32]sturmChain{},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got := NewSolver(tc.argC, tc.argI, tc.argS)

// 			assert.Equal(t, tc.wantC, got.counter)
// 			assert.Equal(t, tc.wantI, got.isolator)
// 			assert.Equal(t, tc.wantS, got.searcher)
// 			assert.Equal(t, tc.wantChainCache, got.chainCache)
// 		})
// 	}
// }

// func Test_NewSolverDefault(t *testing.T) {

// 	s := NewSolverDefault()

// 	assert.Equal(t, ALG_COUNT_STURM, s.counter)
// 	assert.Equal(t, ALG_ISOLATE_BISECT, s.isolator)
// 	assert.Equal(t, ALG_SEARCH_BISECT, s.searcher)
// }

// func Test_SolverCountRootsWithin(t *testing.T) {

// 	testCases := []struct {
// 		name string
// 		argS Solver
// 		argP Poly
// 		argA float64
// 		argB float64
// 		want int
// 	}{
// 		{
// 			name: "linear, 0 in (-1, 0]",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, 0}),
// 			argA: -1,
// 			argB: 0,
// 			want: 1,
// 		},
// 		{
// 			name: "linear, 0 not in (0, 1]",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, 0}),
// 			argA: 0,
// 			argB: 1,
// 			want: 0,
// 		},
// 		{
// 			name: "cubic 3 roots",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, -1, -1, 0}),
// 			argA: -2,
// 			argB: 2,
// 			want: 3,
// 		},
// 		{
// 			name: "cubic 2 roots",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, -1, -1, 0}),
// 			argA: -2,
// 			argB: 0,
// 			want: 2,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got := tc.argS.CountRootsWithin(tc.argP, tc.argA, tc.argB)

// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }

// func Test_SolverIsolateRootsWithin(t *testing.T) {

// 	testCases := []struct {
// 		name string
// 		argS Solver
// 		argP Poly
// 		argA float64
// 		argB float64
// 		want []HalfOpenInterval
// 	}{
// 		{
// 			name: "linear on interval (-1, 0]",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, 0}),
// 			argA: -1,
// 			argB: 0,
// 			want: []HalfOpenInterval{{-1, 0}},
// 		},
// 		{
// 			name: "linear off interval (0, 1]",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, 0}),
// 			argA: 0,
// 			argB: 1,
// 			want: []HalfOpenInterval{},
// 		},
// 		{
// 			name: "quartic",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, 1, 0, -1, -1}),
// 			argA: -2,
// 			argB: 2,
// 			want: []HalfOpenInterval{{-2, 0}, {0, 2}},
// 		},
// 		{
// 			name: "cubic",
// 			argS: NewSolverDefault(),
// 			argP: NewPoly([]float64{1, -1, -1, 0}),
// 			argA: -2,
// 			argB: 2,
// 			want: []HalfOpenInterval{{-1, -0.5}, {-0.5, 0}, {0, 2}},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got := tc.argS.IsolateRootsWithin(tc.argP, tc.argA, tc.argB)

// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }

// func Test_SolverFindRoots(t *testing.T) {

// 	p := NewPolyFromString("x^21 - 86400x + 86399")
// 	//q := NewPolyFromString("x^2-4x+4")
// 	s := NewSolver(ALG_COUNT_STURM, ALG_ISOLATE_BISECT, ALG_SEARCH_BISECT)

// 	//SetNewtonSearchIterations(500)

// 	fmt.Println(s.FindRootsWithin(p, 1, 100))
// 	t.Log(s.FindRootsWithin(p, 0, 100))
// 	t.Log(s.FindRoots(p))
// 	// t.Log(s.FindRoots(q))
// 	// tmp := NewPolyWilkinson()
// 	// t.Log(s.FindRoots(tmp))
// 	// t.Log(s.FindRoots(tmp))

// 	// a := NewPolyFromString("x^3 - x^2 - x")
// 	// b := NewPolyFromString("x^2")
// 	// t.Log(s.FindIntersections(a, b))
// }
