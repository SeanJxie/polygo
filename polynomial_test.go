package polygo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	Basic white-box tests and benchmarks for functions and methods defined in polynomial.go.
*/

func Test_NewPanic(t *testing.T) {

	assert.Panics(t, func() { NewPoly([]float64{}) })
	assert.Panics(t, func() { NewPolyFromString("") })
}

func Test_NewPoly(t *testing.T) {

	testCases := []struct {
		name      string
		arg       []float64
		wantCoefs []float64
		wantLen   int
		wantDeg   int
	}{
		{
			name:      "zero",
			arg:       []float64{0},
			wantCoefs: []float64{0},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "constant",
			arg:       []float64{3.14159265358},
			wantCoefs: []float64{3.14159265358},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "linear",
			arg:       []float64{314, 512},
			wantCoefs: []float64{512, 314},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "quadratic",
			arg:       []float64{691, 2666, 12},
			wantCoefs: []float64{12, 2666, 691},
			wantLen:   3,
			wantDeg:   2,
		},
		{
			name:      "redundant zeroes",
			arg:       []float64{0, 0, 0, 1, 2, 0, 3},
			wantCoefs: []float64{3, 0, 2, 1},
			wantLen:   4,
			wantDeg:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPoly(tc.arg)

			assert.Equal(t, tc.wantCoefs, got.coef)
			assert.Equal(t, tc.wantLen, got.len)
			assert.Equal(t, tc.wantDeg, got.deg)
		})
	}
}

func Test_newPolyNoReverse(t *testing.T) {

	testCases := []struct {
		name      string
		arg       []float64
		wantCoefs []float64
		wantLen   int
		wantDeg   int
	}{
		{
			name:      "zero",
			arg:       []float64{0},
			wantCoefs: []float64{0},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "constant",
			arg:       []float64{3.14159265358},
			wantCoefs: []float64{3.14159265358},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "linear",
			arg:       []float64{314, 512},
			wantCoefs: []float64{314, 512},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "quadratic",
			arg:       []float64{691, 2666, 12},
			wantCoefs: []float64{691, 2666, 12},
			wantLen:   3,
			wantDeg:   2,
		},
		{
			name:      "faux redundant zeroes",
			arg:       []float64{0, 0, 0, 1, 2, 0, 3},
			wantCoefs: []float64{0, 0, 0, 1, 2, 0, 3},
			wantLen:   7,
			wantDeg:   6,
		},
		{
			name:      "actual redundant zeroes",
			arg:       []float64{1, 2, 0, 3, 0, 0, 0},
			wantCoefs: []float64{1, 2, 0, 3},
			wantLen:   4,
			wantDeg:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := newPolyNoReverse(tc.arg)

			assert.Equal(t, tc.wantCoefs, got.coef)
			assert.Equal(t, tc.wantLen, got.len)
			assert.Equal(t, tc.wantDeg, got.deg)
		})
	}
}

func Test_parseTermPanic(t *testing.T) {

	assert.Panics(t, func() { parseTerm("+n") },
		"parseTerm: could not parse deg 0 term coefficient \"n\" "+
			"(strconv.ParseFloat: parsing \"n\": invalid syntax).")

	assert.Panics(t, func() { parseTerm("-nx") },
		"parseTerm: could not parse deg 1 term coefficient \"n\" "+
			"(strconv.ParseFloat: parsing \"n\": invalid syntax).")

	assert.Panics(t, func() { parseTerm("+nx^k") },
		"could not parse exponent \"k\" "+
			"(strconv.ParseInt: parsing \"k\": invalid syntax).")
}

func Test_parseTerm(t *testing.T) {

	testCases := []struct {
		name     string
		arg      string
		wantCoef float64
		wantDeg  int
	}{
		{
			name:     "+zero",
			arg:      "+0",
			wantCoef: 0,
			wantDeg:  0,
		},
		{
			name:     "-zero",
			arg:      "-0",
			wantCoef: 0,
			wantDeg:  0,
		},
		{
			name:     "+const",
			arg:      "+000.121151251612128576",
			wantCoef: 0.121151251612128576,
			wantDeg:  0,
		},
		{
			name:     "-const",
			arg:      "-000.121151251612128576",
			wantCoef: -0.121151251612128576,
			wantDeg:  0,
		},
		{
			name:     "linear",
			arg:      "-12x",
			wantCoef: -12,
			wantDeg:  1,
		},
		{
			name:     "linear implicit pos",
			arg:      "+x",
			wantCoef: 1,
			wantDeg:  1,
		},
		{
			name:     "linear implicit neg",
			arg:      "-x",
			wantCoef: -1,
			wantDeg:  1,
		},
		{
			name:     "quadratic",
			arg:      "+24.151x^2",
			wantCoef: 24.151,
			wantDeg:  2,
		},
		{
			name:     "large deg",
			arg:      "-11111.151x^1296124",
			wantCoef: -11111.151,
			wantDeg:  1296124,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotCoef, gotDeg := parseTerm(tc.arg)

			assert.Equal(t, tc.wantCoef, gotCoef)
			assert.Equal(t, tc.wantDeg, gotDeg)
		})
	}
}

func Test_NewPolyFromString(t *testing.T) {

	testCases := []struct {
		name      string
		arg       string
		wantCoefs []float64
		wantLen   int
		wantDeg   int
	}{
		{
			name:      "zero",
			arg:       "   -0     ",
			wantCoefs: []float64{0},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "constants",
			arg:       "   1+2+3+4+5+6-2.15-2     ",
			wantCoefs: []float64{16.85},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "linear pos leading",
			arg:       "   x +4    ",
			wantCoefs: []float64{4, 1},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "linear neg leading",
			arg:       "   -x +4    ",
			wantCoefs: []float64{4, -1},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "multi cancel out",
			arg:       "  2 -x +4 -   1  -x^2 + 5 + x -  2 +x   ^ 2 -3 -  5  ",
			wantCoefs: []float64{0},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "unordered",
			arg:       " x^3 - x^4 + x^1 + x^1",
			wantCoefs: []float64{0, 2, 0, 1, -1},
			wantLen:   5,
			wantDeg:   4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPolyFromString(tc.arg)

			assert.Equal(t, tc.wantCoefs, got.coef)
			assert.Equal(t, tc.wantLen, got.len)
			assert.Equal(t, tc.wantDeg, got.deg)
		})
	}
}

func Test_NewZeroPoly(t *testing.T) {

	zp := NewZeroPoly()

	assert.Equal(t, []float64{0}, zp.coef)
	assert.Equal(t, 1, zp.len)
	assert.Equal(t, 0, zp.deg)
}

func Test_PolyProperties(t *testing.T) {

	// Test_Poly the property "getters".
	//
	// Specifically:
	// 	Degree()
	// 	Coefficients()
	// 	LeadingCoefficient()
	// 	LargestCoefficient()
	// 	SmallestCoefficient()

	testCases := []struct {
		name      string
		arg       Poly
		wantCoefs []float64
		wantDeg   int
		wantLead  float64
		wantLarge float64
		wantSmall float64
	}{
		{
			name:      "zero",
			arg:       NewZeroPoly(),
			wantCoefs: []float64{0},
			wantDeg:   0,
			wantLead:  0,
			wantLarge: 0,
			wantSmall: 0,
		},
		{
			name:      "constant",
			arg:       NewPoly([]float64{3.14159265358}),
			wantCoefs: []float64{3.14159265358},
			wantDeg:   0,
			wantLead:  3.14159265358,
			wantLarge: 3.14159265358,
			wantSmall: 3.14159265358,
		},
		{
			name:      "linear",
			arg:       NewPoly([]float64{314, 512}),
			wantCoefs: []float64{314, 512},
			wantDeg:   1,
			wantLead:  314,
			wantLarge: 512,
			wantSmall: 314,
		},
		{
			name:      "quadratic",
			arg:       NewPoly([]float64{691, 2666, 12}),
			wantCoefs: []float64{691, 2666, 12},
			wantDeg:   2,
			wantLead:  691,
			wantLarge: 2666,
			wantSmall: 12,
		},
		{
			name:      "redundant zeroes",
			arg:       NewPoly([]float64{0, 0, 0, 1, 2, -1, 3}),
			wantCoefs: []float64{1, 2, -1, 3},
			wantDeg:   3,
			wantLead:  1,
			wantLarge: 3,
			wantSmall: -1,
		},
		{
			name:      "extrema",
			arg:       NewPoly([]float64{math.MaxFloat64, -math.MaxFloat64}),
			wantCoefs: []float64{math.MaxFloat64, -math.MaxFloat64},
			wantDeg:   1,
			wantLead:  math.MaxFloat64,
			wantLarge: math.MaxFloat64,
			wantSmall: -math.MaxFloat64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotDeg := tc.arg.Degree()
			gotCoefs := tc.arg.Coefficients()
			gotLead := tc.arg.LeadingCoefficient()
			gotLarge := tc.arg.LargestCoefficient()
			gotSmall := tc.arg.SmallestCoefficient()

			assert.Equal(t, tc.wantDeg, gotDeg)
			assert.Equal(t, tc.wantCoefs, gotCoefs)
			assert.Equal(t, tc.wantLead, gotLead)
			assert.Equal(t, tc.wantLarge, gotLarge)
			assert.Equal(t, tc.wantSmall, gotSmall)
		})
	}
}

func Test_PolyCoefficientWithDegree(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argN uint
		want float64
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argN: 0,
			want: 0,
		},
		{
			name: "n < deg(p)",
			argP: NewPoly([]float64{5, 4, 3, 2, 1}),
			argN: 3,
			want: 4,
		},
		{
			name: "n > deg(p)",
			argP: NewPoly([]float64{5, 4, 3, 2, 1}),
			argN: 10,
			want: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argP.CoefficientWithDegree(tc.argN)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_PolyEqual(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argQ Poly
		want bool
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argQ: NewZeroPoly(),
			want: true,
		},
		{
			name: "deg(p) != deg(q)",
			argP: NewPoly([]float64{1, 2, 3}),
			argQ: NewPoly([]float64{4, 5, 6, 7}),
			want: false,
		},
		{
			name: "deg(p) == deg(q), p != q",
			argP: NewPoly([]float64{1, 2, 3, 2}),
			argQ: NewPoly([]float64{4, 5, 6, 7}),
			want: false,
		},
		{
			name: "deg(p) == deg(q), p == q",
			argP: NewPoly([]float64{3, 1, 4, 1}),
			argQ: NewPoly([]float64{3, 1, 4, 1}),
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1 := tc.argP.Equal(tc.argQ)
			got2 := tc.argQ.Equal(tc.argP)

			assert.Equal(t, tc.want, got1)
			assert.Equal(t, tc.want, got2)
		})
	}
}

func Test_PolyApproxEqual(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argQ Poly
		want bool
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argQ: NewZeroPoly(),
			want: true,
		},
		{
			name: "deg(p) != deg(q)",
			argP: NewPoly([]float64{1, 2, 3}),
			argQ: NewPoly([]float64{4, 5, 6, 7}),
			want: false,
		},
		{
			name: "deg(p) == deg(q), p != q",
			argP: NewPoly([]float64{1, 2, 3, 2}),
			argQ: NewPoly([]float64{4, 5, 6, 7}),
			want: false,
		},
		{
			name: "deg(p) == deg(q), p == q",
			argP: NewPoly([]float64{3, 1, 4, 1}),
			argQ: NewPoly([]float64{3, 1, 4, 1}),
			want: true,
		},
		{
			name: "deg(p) == deg(q), p ~= q",
			argP: NewPoly([]float64{125125124, 1.000001, 4, 1.00002}),
			argQ: NewPoly([]float64{125125125, 1, 4.000002, 1}),
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1 := tc.argP.Equal(tc.argQ)
			got2 := tc.argQ.Equal(tc.argP)

			assert.Equal(t, tc.want, got1)
			assert.Equal(t, tc.want, got2)
		})
	}
}

func Test_PolyBooleanChecks(t *testing.T) {

	// Test_Poly the boolean properties.
	//
	// Specifically:
	// 	IsConstant()
	// 	IsZero()
	// 	IsMonic()

	testCases := []struct {
		name      string
		arg       Poly
		wantConst bool
		wantZero  bool
		wantMonic bool
	}{
		{
			name:      "zero",
			arg:       NewZeroPoly(),
			wantConst: true,
			wantZero:  true,
			wantMonic: false,
		},
		{
			name:      "monic linear",
			arg:       NewPoly([]float64{1, 2}),
			wantConst: false,
			wantZero:  false,
			wantMonic: true,
		},
		{
			name:      "nonmonic quadratic",
			arg:       NewPoly([]float64{129481, 2, 2}),
			wantConst: false,
			wantZero:  false,
			wantMonic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotConst := tc.arg.IsConstant()
			gotZero := tc.arg.IsZero()
			gotMonic := tc.arg.IsMonic()

			assert.Equal(t, tc.wantConst, gotConst)
			assert.Equal(t, tc.wantZero, gotZero)
			assert.Equal(t, tc.wantMonic, gotMonic)
		})
	}
}

func Test_PolyAt(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argX float64
		want float64
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argX: 0,
			want: 0,
		},
		{
			name: "linear",
			argP: NewPoly([]float64{12512, 2512}),
			argX: 123,
			want: 1541488,
		},
		{
			name: "large deg",
			argP: NewPoly([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}),
			argX: 3.14,
			want: 1.950054059657435e+07,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argP.At(tc.argX)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_PolyAdd(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argQ Poly
		want Poly
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argQ: NewZeroPoly(),
			want: NewZeroPoly(),
		},
		{
			name: "nonzero const",
			argP: NewPoly([]float64{12947124}),
			argQ: NewPoly([]float64{668894789}),
			want: NewPoly([]float64{12947124 + 668894789}),
		},
		{
			name: "deg(p) == deg(q)",
			argP: NewPoly([]float64{5, 6, 7, -8}),
			argQ: NewPoly([]float64{1, -2, 3, 5}),
			want: NewPoly([]float64{6, 4, 10, -3}),
		},
		{
			name: "deg(p) != deg(q)",
			argP: NewPoly([]float64{5, 1, 22, 5, 6}),
			argQ: NewPoly([]float64{566, -25, -12243}),
			want: NewPoly([]float64{5, 1, 566 + 22, 5 - 25, 6 - 12243}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1 := tc.argP.Add(tc.argQ)
			got2 := tc.argQ.Add(tc.argP)

			assert.Equal(t, tc.want, got1)
			assert.Equal(t, tc.want, got2)
		})
	}
}

func Test_PolySub(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argQ Poly
		want Poly
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argQ: NewZeroPoly(),
			want: NewZeroPoly(),
		},
		{
			name: "nonzero const",
			argP: NewPoly([]float64{12947124}),
			argQ: NewPoly([]float64{668894789}),
			want: NewPoly([]float64{12947124 - 668894789}),
		},
		{
			name: "deg(p) == deg(q)",
			argP: NewPoly([]float64{5, 6, 7, -8}),
			argQ: NewPoly([]float64{1, -2, 3, 5}),
			want: NewPoly([]float64{4, 8, 4, -13}),
		},
		{
			name: "deg(p) != deg(q)",
			argP: NewPoly([]float64{5, 1, 22, 5, 6}),
			argQ: NewPoly([]float64{566, -25, -12243}),
			want: NewPoly([]float64{5, 1, 22 - 566, 5 + 25, 6 + 12243}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argP.Sub(tc.argQ)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_PolyMulScalar(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argS float64
		want Poly
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argS: 0,
			want: NewZeroPoly(),
		},
		{
			name: "identity",
			argP: NewPoly([]float64{25, 1, 2, 53, 1, 2, 653}),
			argS: 1,
			want: NewPoly([]float64{25, 1, 2, 53, 1, 2, 653}),
		},
		{
			name: "nonzero const",
			argP: NewPoly([]float64{125}),
			argS: -125982,
			want: NewPoly([]float64{125 * -125982}),
		},
		{
			name: "linear",
			argP: NewPoly([]float64{12225, 121521}),
			argS: 512,
			want: NewPoly([]float64{12225 * 512, 121521 * 512}),
		},
		{
			name: "large deg",
			argP: NewPoly([]float64{12225, 121521, 124, 5, 1, 2, 124}),
			argS: 2,
			want: NewPoly([]float64{12225 * 2, 121521 * 2, 124 * 2, 10, 2, 4, 248}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argP.MulScalar(tc.argS)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_PolyMul(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argQ Poly
		want Poly
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argQ: NewZeroPoly(),
			want: NewZeroPoly(),
		},
		{
			name: "vanish",
			argP: NewZeroPoly(),
			argQ: NewPoly([]float64{12, 24124, 1, 2, 5, 124}),
			want: NewZeroPoly(),
		},
		{
			name: "deg(p) == deq(q)",
			argP: NewPoly([]float64{15, 214, 225, 25, 212}),
			argQ: NewPoly([]float64{12, 4, 1, 2, 5, 124, 21, 12, 636}),
			want: NewPoly([]float64{180, 2628, 3571, 1444, 3372, 4253,
				28238, 33123, 20993, 165617, 147852, 18444, 134832}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1 := tc.argP.Mul(tc.argQ)
			got2 := tc.argQ.Mul(tc.argP)

			assert.Equal(t, tc.want, got1)
			assert.Equal(t, tc.want, got2)
		})
	}
}

func Test_PolyMulFast(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argQ Poly
		want Poly
	}{
		{
			name: "zero",
			argP: NewZeroPoly(),
			argQ: NewZeroPoly(),
			want: NewZeroPoly(),
		},
		{
			name: "vanish",
			argP: NewZeroPoly(),
			argQ: NewPoly([]float64{12, 24124, 1, 2, 5, 124}),
			want: NewZeroPoly(),
		},
		{
			name: "deg(p) == deq(q)",
			argP: NewPoly([]float64{15, 214, 225, 25, 212}),
			argQ: NewPoly([]float64{12, 4, 1, 2, 5, 124, 21, 12, 636}),
			want: NewPoly([]float64{180, 2628, 3571, 1444, 3372, 4253,
				28238, 33123, 20993, 165617, 147852, 18444, 134832}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got1 := tc.argP.MulFast(tc.argQ)
			got2 := tc.argQ.MulFast(tc.argP)

			assert.True(t, tc.want.Equal(got1))
			assert.True(t, tc.want.Equal(got2))
		})
	}
}

func Test_PolyDivZeroPanic(t *testing.T) {

	assert.Panics(t, func() { NewZeroPoly().Div(NewZeroPoly()) })
	assert.Panics(t, func() { NewPoly([]float64{123}).Div(NewZeroPoly()) })
}

func Test_PolyDiv(t *testing.T) {

	testCases := []struct {
		name    string
		argP    Poly
		argQ    Poly
		wantQuo Poly
		wantRem Poly
	}{
		{
			name:    "zero",
			argP:    NewZeroPoly(),
			argQ:    NewPoly([]float64{1}),
			wantQuo: NewZeroPoly(),
			wantRem: NewZeroPoly(),
		},
		{
			name:    "equal linear",
			argP:    NewPoly([]float64{125, 254}),
			argQ:    NewPoly([]float64{125, 254}),
			wantQuo: NewPoly([]float64{1}),
			wantRem: NewZeroPoly(),
		},
		{
			name:    "cubic div linear",
			argP:    NewPoly([]float64{1, -12, 0, -42}),
			argQ:    NewPoly([]float64{1, -3}),
			wantQuo: NewPoly([]float64{1, -9, -27}),
			wantRem: NewPoly([]float64{-123}),
		},
		{
			name:    "deg(q) > deg(p)",
			argP:    NewPoly([]float64{1, -3}),
			argQ:    NewPoly([]float64{1, -12, 0, -42}),
			wantQuo: NewZeroPoly(),
			wantRem: NewPoly([]float64{1, -3}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotQuo, gotRem := tc.argP.Div(tc.argQ)

			assert.Equal(t, tc.wantQuo, gotQuo)
			assert.Equal(t, tc.wantRem, gotRem)
		})
	}
}

func Test_PolyDerivative(t *testing.T) {
	testCases := []struct {
		name string
		arg  Poly
		want Poly
	}{
		{
			name: "zero",
			arg:  NewZeroPoly(),
			want: NewZeroPoly(),
		},
		{
			name: "nonzero const",
			arg:  NewPoly([]float64{3.1415}),
			want: NewZeroPoly(),
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

func Test_PolyReciprocal(t *testing.T) {
	testCases := []struct {
		name string
		arg  Poly
		want Poly
	}{
		{
			name: "zero",
			arg:  NewZeroPoly(),
			want: NewZeroPoly(),
		},
		{
			name: "nonzero const",
			arg:  NewPoly([]float64{3.1415}),
			want: NewPoly([]float64{3.1415}),
		},
		{
			name: "linear",
			arg:  NewPoly([]float64{55124, 12488459123}),
			want: NewPoly([]float64{12488459123, 55124}),
		},
		{
			name: "quadratic",
			arg:  NewPoly([]float64{47346346, 734334, 2342366}),
			want: NewPoly([]float64{2342366, 734334, 47346346}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.arg.Reciprocal()

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_PolyCauchyBoundPanic(t *testing.T) {

	assert.Panics(t, func() { NewZeroPoly().CauchyBound() })
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

func Test_PolyString(t *testing.T) {
	testCases := []struct {
		name string
		arg  Poly
		want string
	}{
		{
			name: "zero",
			arg:  NewZeroPoly(),
			want: "[ 0.00000x^0 ]",
		},
		{
			name: "negative const",
			arg:  NewPoly([]float64{-21}),
			want: "[ -21.00000x^0 ]",
		},
		{
			name: "nonzero const",
			arg:  NewPoly([]float64{3.1415}),
			want: "[ 3.14150x^0 ]",
		},
		{
			name: "linear",
			arg:  NewPoly([]float64{125124, 1.1715}),
			want: "[ 125124.00000x^1 + 1.17150x^0 ]",
		},
		{
			name: "quadratic",
			arg:  NewPoly([]float64{47346346, 734334, 2342366}),
			want: "[ 47346346.00000x^2 + 734334.00000x^1 + 2342366.00000x^0 ]",
		},
		{
			name: "cubic",
			arg:  NewPoly([]float64{2152, 47346346, 734334, 2342366}),
			want: "[ 2152.00000x^3 + 47346346.00000x^2 + 734334.00000x^1 + 2342366.00000x^0 ]",
		},
		{
			name: "cubic with negative",
			arg:  NewPoly([]float64{-2152, 47346346, -734334, 2342366}),
			want: "[ -2152.00000x^3 + 47346346.00000x^2 - 734334.00000x^1 + 2342366.00000x^0 ]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.arg.String()

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_Poly_id(t *testing.T) {
	testCases := []struct {
		name string
		arg  Poly
	}{
		{
			name: "zero",
			arg:  NewZeroPoly(),
		},
		{
			name: "negative const",
			arg:  NewPoly([]float64{-21}),
		},
		{
			name: "nonzero const",
			arg:  NewPoly([]float64{3.1415}),
		},
		{
			name: "linear",
			arg:  NewPoly([]float64{125124, 1.1715}),
		},
		{
			name: "quadratic",
			arg:  NewPoly([]float64{47346346, 734334, 2342366}),
		},
		{
			name: "cubic",
			arg:  NewPoly([]float64{2152, 47346346, 734334, 2342366}),
		},
		{
			name: "cubic with negative",
			arg:  NewPoly([]float64{-2152, 47346346, -734334, 2342366}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.arg.id()
		})
	}
}
