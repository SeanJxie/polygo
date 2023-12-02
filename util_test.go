package polygo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	Basic white-box tests and benchmarks for functions and methods defined in util.go.
*/

func Test_removeTrailingZeroes(t *testing.T) {

	testCases := []struct {
		name string
		arg  []float64
		want []float64
	}{
		{
			name: "empty",
			arg:  []float64{},
			want: []float64{},
		},
		{
			name: "no trailing zeroes",
			arg:  []float64{1, 2, 3},
			want: []float64{1, 2, 3},
		},
		{
			name: "single trailing zero",
			arg:  []float64{1, 2, 3, 0},
			want: []float64{1, 2, 3},
		},
		{
			name: "multiple trailing zeros",
			arg:  []float64{0, 1, 2, 3, 0, 0, 0, 0, 0},
			want: []float64{0, 1, 2, 3},
		},
		{
			name: "single zero",
			arg:  []float64{0},
			want: []float64{0},
		},
		{
			name: "multiple zeros",
			arg:  []float64{0, 0, 0, 0, 0},
			want: []float64{0},
		},
		{
			name: "leading zeros",
			arg:  []float64{0, 0, 0, 1, 2, 3},
			want: []float64{0, 0, 0, 1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := removeTrailingZeroes(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_reverse(t *testing.T) {

	testCases := []struct {
		name string
		arg  []float64
		want []float64
	}{
		{
			name: "empty",
			arg:  []float64{},
			want: []float64{},
		},
		{
			name: "single element",
			arg:  []float64{5},
			want: []float64{5},
		},
		{
			name: "swap",
			arg:  []float64{1000, 1001},
			want: []float64{1001, 1000},
		},
		{
			name: "multiple elements",
			arg:  []float64{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want: []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, -1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := reverse(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_plusOrMinus(t *testing.T) {

	assert.True(t, plusOrMinus('+'))
	assert.True(t, plusOrMinus('-'))
	assert.False(t, plusOrMinus('x'))
}

func Test_expand(t *testing.T) {

	testCases := []struct {
		name string
		argS []float64
		argN int
		want []float64
	}{
		{
			name: "empty expand zero",
			argS: []float64{},
			argN: 0,
			want: []float64{},
		},
		{
			name: "empty expand nonzero",
			argS: []float64{},
			argN: 5,
			want: []float64{0, 0, 0, 0, 0},
		},
		{
			name: "empty expand negative",
			argS: []float64{},
			argN: -24,
			want: []float64{},
		},
		{
			name: "filled expand bigger",
			argS: []float64{1, 2, 3, 4},
			argN: 10,
			want: []float64{1, 2, 3, 4, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "filled expand equal",
			argS: []float64{1, 2, 3, 4},
			argN: 4,
			want: []float64{1, 2, 3, 4},
		},
		{
			name: "filled expand smaller",
			argS: []float64{1, 2, 3, 4},
			argN: 1,
			want: []float64{1, 2, 3, 4},
		},
		{
			name: "filled expand negative",
			argS: []float64{1, 2, 3, 4},
			argN: -125,
			want: []float64{1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := expand(tc.argS, tc.argN)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_isPOT(t *testing.T) {

	testCases := []struct {
		name string
		arg  int
		want bool
	}{
		{
			name: "negative",
			arg:  -124,
			want: false,
		},
		{
			name: "zero",
			arg:  0,
			want: false,
		},
		{
			name: "one",
			arg:  1,
			want: true,
		},
		{
			name: "two",
			arg:  2,
			want: true,
		},
		{
			name: "non POT",
			arg:  2361,
			want: false,
		},
		{
			name: "big POT 2^62",
			arg:  4611686018427387904,
			want: true,
		},
		{
			name: "big POT signed int64 max 2^63-1",
			arg:  9223372036854775807,
			want: false,
		},
		{
			name: "big non POT",
			arg:  4611686018227387401,
			want: false,
		},
		{
			name: "max non POT",
			arg:  math.MaxInt64,
			want: false,
		},
		{
			name: "min non POT",
			arg:  math.MinInt64,
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isPOT(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_nextPOT(t *testing.T) {

	testCases := []struct {
		name string
		arg  int
		want int
	}{
		{
			name: "negative",
			arg:  -124,
			want: 1,
		},
		{
			name: "zero",
			arg:  0,
			want: 1,
		},
		{
			name: "one",
			arg:  1,
			want: 1,
		},
		{
			name: "two",
			arg:  2,
			want: 2,
		},
		{
			name: "big POT 2^62",
			arg:  4611686018427387904,
			want: 4611686018427387904,
		},
		{
			name: "big POT signed int64 overflow 2^63-1",
			arg:  9223372036854775807,
			want: -9223372036854775808,
		},
		{
			name: "big non POT",
			arg:  4611686018227387401,
			want: 4611686018427387904,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := nextPOT(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_toComplex128(t *testing.T) {

	testCases := []struct {
		name string
		arg  []float64
		want []complex128
	}{
		{
			name: "empty",
			arg:  []float64{},
			want: []complex128{},
		},
		{
			name: "nonempty",
			arg:  []float64{0, 1, 2, 3, 4},
			want: []complex128{0, 1, 2, 3, 4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := toComplex128(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_toFloat64(t *testing.T) {

	testCases := []struct {
		name string
		arg  []complex128
		want []float64
	}{
		{
			name: "empty",
			arg:  []complex128{},
			want: []float64{},
		},
		{
			name: "nonempty",
			arg:  []complex128{0, 1, 2, 3, 4},
			want: []float64{0, 1, 2, 3, 4},
		},
		{
			name: "complex",
			arg:  []complex128{0, 1i, 2i + 7, 3, 4i},
			want: []float64{0, 0, 7, 3, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := toFloat64(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_max_minPanic(t *testing.T) {

	assert.Panics(t, func() { max([]float64{}) })
	assert.Panics(t, func() { min([]float64{}) })
}

func Test_max(t *testing.T) {

	testCases := []struct {
		name string
		arg  []float64
		want float64
	}{
		{
			name: "single",
			arg:  []float64{1},
			want: 1,
		},
		{
			name: "multiple",
			arg:  []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4},
			want: 4,
		},
		{
			name: "extrema",
			arg:  []float64{math.MaxFloat64, -math.MaxFloat64},
			want: math.MaxFloat64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := max(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}
func Test_min(t *testing.T) {

	testCases := []struct {
		name string
		arg  []float64
		want float64
	}{
		{
			name: "single",
			arg:  []float64{1},
			want: 1,
		},
		{
			name: "multiple",
			arg:  []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4},
			want: -4,
		},
		{
			name: "extrema",
			arg:  []float64{math.MaxFloat64, -math.MaxFloat64},
			want: -math.MaxFloat64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := min(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_factPanic(t *testing.T) {

	assert.Panics(t, func() { fact(-1) })
	assert.Panics(t, func() { fact(-124) })
}

func Test_fact(t *testing.T) {

	testCases := []struct {
		name string
		arg  int
		want float64
	}{
		{
			name: "zero",
			arg:  0,
			want: 1,
		},
		{
			name: "nonzero",
			arg:  25,
			want: 15511210043330985984000000,
		},
		{
			name: "big",
			arg:  50,
			want: 3.0414093201713376e+64,
		},
		{
			name: "inf",
			arg:  1000,
			want: math.Inf(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fact(tc.arg)
			assert.Equal(t, tc.want, got)
		})
	}
}
