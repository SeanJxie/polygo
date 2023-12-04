package polygo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PolyDerivative(t *testing.T) {
	testCases := []struct {
		name string
		arg  Poly
		want Poly
	}{
		{
			name: "zero",
			arg:  NewPolyZero(),
			want: NewPolyZero(),
		},
		{
			name: "nonzero const",
			arg:  NewPoly([]float64{3.1415}),
			want: NewPolyZero(),
		},
		{
			name: "linear",
			arg:  NewPoly([]float64{125124, 125123}),
			want: NewPoly([]float64{125124}),
		},
		{
			name: "quadratic",
			arg:  NewPoly([]float64{47346346, 734334, 2342366}),
			want: NewPoly([]float64{47346346 * 2, 734334}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.arg.Derivative()

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_PolyDerivativeN(t *testing.T) {
	testCases := []struct {
		name string
		argP Poly
		argN int
		want Poly
	}{
		{
			name: "zero",
			argP: NewPolyZero(),
			argN: 1,
			want: NewPolyZero(),
		},
		{
			name: "nonzero const",
			argP: NewPoly([]float64{3.1415}),
			argN: 1,
			want: NewPolyZero(),
		},
		{
			name: "linear",
			argP: NewPoly([]float64{125124, 125123}),
			argN: 1,
			want: NewPoly([]float64{125124}),
		},
		{
			name: "quadratic",
			argP: NewPoly([]float64{47346346, 734334, 2342366}),
			argN: 1,
			want: NewPoly([]float64{47346346 * 2, 734334}),
		},
		{
			name: "negative n zero",
			argP: NewPolyZero(),
			argN: -2,
			want: NewPolyZero(),
		},
		{
			name: "negative n quadratic",
			argP: NewPoly([]float64{47346346, 734334, 2342366}),
			argN: -2,
			want: NewPoly([]float64{47346346, 734334, 2342366}),
		},
		{
			name: "n = deg + 1 quadratic",
			argP: NewPoly([]float64{47346346, 734334, 2342366}),
			argN: 3,
			want: NewPolyZero(),
		},
		{
			name: "n = deg quadratic",
			argP: NewPoly([]float64{47346346, 734334, 2342366}),
			argN: 2,
			want: NewPolyConst(47346346 * 2),
		},
		{
			name: "n = big quadratic",
			argP: NewPoly([]float64{47346346, 734334, 2342366}),
			argN: 100000000000000000,
			want: NewPolyZero(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argP.DerivativeN(tc.argN)
			t.Log(got)
			assert.Equal(t, tc.want, got)
		})
	}
}
