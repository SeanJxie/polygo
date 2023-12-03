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

func Test_NewPolyConst(t *testing.T) {

	zp := NewPolyConst(3.1415)

	assert.Equal(t, []float64{3.1415}, zp.coef)
	assert.Equal(t, 1, zp.len)
	assert.Equal(t, 0, zp.deg)
}

func Test_NewPolyZero(t *testing.T) {

	zp := NewPolyZero()

	assert.Equal(t, []float64{0}, zp.coef)
	assert.Equal(t, 1, zp.len)
	assert.Equal(t, 0, zp.deg)
}

func Test_NewPolyWilkinson(t *testing.T) {

	wp := NewPolyWilkinson()

	wantCoef := reverse([]float64{
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
		2432902008176640000,
	})

	assert.Equal(t, wantCoef, wp.coef)
	assert.Equal(t, 21, wp.len)
	assert.Equal(t, 20, wp.deg)
}

func Test_NewPolyFactoredPanic(t *testing.T) {

	assert.Panics(t, func() { NewPolyFactored(123, []float64{}) })
}

func Test_NewPolyFactored(t *testing.T) {

	testCases := []struct {
		name      string
		argA      float64
		argR      []float64
		wantCoefs []float64
		wantLen   int
		wantDeg   int
	}{
		{
			name:      "zero",
			argA:      0,
			argR:      []float64{0},
			wantCoefs: []float64{0},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "linear",
			argA:      5,
			argR:      []float64{-1},
			wantCoefs: []float64{5, 5},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "quadratic",
			argA:      1,
			argR:      []float64{2, 3},
			wantCoefs: []float64{6, -5, 1},
			wantLen:   3,
			wantDeg:   2,
		},
		{
			name:      "cubic",
			argA:      2,
			argR:      []float64{5, 4, 3},
			wantCoefs: []float64{-120, 94, -24, 2},
			wantLen:   4,
			wantDeg:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPolyFactored(tc.argA, tc.argR)
			t.Log(got)

			assert.Equal(t, tc.wantCoefs, got.coef)
			assert.Equal(t, tc.wantLen, got.len)
			assert.Equal(t, tc.wantDeg, got.deg)
		})
	}
}

func Test_NewPolyTaylorSinPanic(t *testing.T) {

	assert.Panics(t, func() { NewPolyTaylorSin(-1, 1) })
	assert.Panics(t, func() { NewPolyTaylorSin(-1248, -63) })
}

func Test_NewPolyTaylorSin(t *testing.T) {

	testCases := []struct {
		name      string
		argN      int
		argA      float64
		wantCoefs []float64
		wantLen   int
		wantDeg   int
	}{
		{
			name:      "Maclaurin deg 0",
			argN:      0,
			argA:      0,
			wantCoefs: []float64{math.Sin(0)},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "Taylor at 1 deg 0",
			argN:      0,
			argA:      1,
			wantCoefs: []float64{math.Sin(1)},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "Maclaurin deg 1",
			argN:      1,
			argA:      0,
			wantCoefs: []float64{0, 1},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "Maclaurin deg 2",
			argN:      2,
			argA:      0,
			wantCoefs: []float64{0, 1},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "Maclaurin deg 3",
			argN:      3,
			argA:      0,
			wantCoefs: []float64{0, 1, 0, -1 / fact(3)},
			wantLen:   4,
			wantDeg:   3,
		},
		{
			name: "Maclaurin deg 10",
			argN: 10,
			argA: 0,
			wantCoefs: []float64{0, 1, 0, -0.16666666666666666, 0, 0.008333333333333333, 0,
				-0.0001984126984126984, 0, 2.7557319223985893e-06},
			wantLen: 10,
			wantDeg: 9, // For default epsilon, best we can do is deg 9.
		},
		{
			name:      "Taylor at 1 deg 1",
			argN:      1,
			argA:      1,
			wantCoefs: []float64{0.30116867893975674, 0.5403023058681398},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "Taylor at 1 deg 2",
			argN:      2,
			argA:      1,
			wantCoefs: []float64{-0.11956681346419151, 1.3817732906760363, -0.42073549240394825},
			wantLen:   3,
			wantDeg:   2,
		},
		{
			name: "Taylor at 1 deg 3",
			argN: 3,
			argA: 1,
			wantCoefs: []float64{-0.02951642915283488, 1.1116221377419664, -0.15058433946987837,
				-0.09005038431135663},
			wantLen: 4,
			wantDeg: 3,
		},
		{
			name: "Taylor at 1 deg 10",
			argN: 10,
			argA: 1,
			wantCoefs: []float64{-1.5196462559423447e-08, 1.0000001687171585,
				-8.528093474674307e-07, -0.16666407491564922, -5.265293160941084e-06,
				0.00834084856576574, -7.70119705670489e-06, -0.00019273352650938828,
				-2.9654467625047508e-06, 3.807796766633698e-06, -2.3188684546072984e-07},
			wantLen: 11,
			wantDeg: 10,
		},
		{
			name: "Taylor at 1 deg 50",
			argN: 50,
			argA: 1,
			wantCoefs: []float64{
				-3.9195392259798626e-17, 1.0000000000000002, 1.60198421891246e-17,
				-0.16666666666666663, -7.769784824795162e-18, 0.008333333333333335,
				-1.469069077880964e-20, -0.00019841269841269836, -1.7536242580586098e-21,
				2.7557319223985897e-06, 3.331078944386683e-23, -2.505210838544173e-08,
				-1.5173629625498664e-25, 1.6059043836821619e-10, -1.058951257234695e-28,
				-7.647163731819818e-13, 1.2476377341865343e-31, 2.8114572543455206e-15,
				8.178785303513588e-33, -8.220635246624333e-18, 5.873966728766431e-36,
				1.957294106339126e-20, 1.1487350833336406e-38, -3.8681701706306824e-23,
				-1.0934801595580858e-40, 6.446950284384472e-26, 4.4612856285855765e-44,
				-9.183689863795543e-29, -2.729810664078527e-47, 1.1309962886447714e-31,
				1.2198621477645828e-50, -1.216125041553518e-34, -2.6342522847488516e-52,
				1.1516335620771955e-37, -5.541906320036663e-54, -9.67759295863162e-41,
				-1.214647260443147e-54, 7.265460179202573e-44, -1.8348734750871204e-55,
				-4.902469750354594e-47, -1.8625049518972674e-56, 2.989311331490434e-50,
				-1.2139581980144593e-57, -1.6551851283754923e-53, -4.753729342811867e-59,
				8.367189813483244e-57, -1.0076762895559066e-60, -3.755953903954799e-60,
				-9.63167604716772e-63, 2.2716003424987766e-63, -2.7667140336128526e-65},
			wantLen: 51,
			wantDeg: 50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPolyTaylorSin(tc.argN, tc.argA)

			assert.Equal(t, tc.wantCoefs, got.coef)
			assert.Equal(t, tc.wantLen, got.len)
			assert.Equal(t, tc.wantDeg, got.deg)
		})
	}
}

func Test_NewPolyChebyshevPanic(t *testing.T) {

	assert.Panics(t, func() { NewPolyChebyshev1(-1) })
	assert.Panics(t, func() { NewPolyChebyshev2(-1) })
}

func Test_NewPolyChebyshev1(t *testing.T) {

	testCases := []struct {
		name      string
		arg       int
		wantCoefs []float64
		wantLen   int
		wantDeg   int
	}{
		{
			name:      "n = 0",
			arg:       0,
			wantCoefs: []float64{1},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "n = 1",
			arg:       1,
			wantCoefs: []float64{0, 1},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "n = 2",
			arg:       2,
			wantCoefs: []float64{-1, 0, 2},
			wantLen:   3,
			wantDeg:   2,
		},
		{
			name:      "n = 3",
			arg:       3,
			wantCoefs: []float64{0, -3, 0, 4},
			wantLen:   4,
			wantDeg:   3,
		},
		{
			name: "n = 20",
			arg:  20,
			wantCoefs: []float64{1, 0, -200, 0, 6600, 0, -84480, 0, 549120, 0, -2.050048e+06, 0,
				4.6592e+06, 0, -6.5536e+06, 0, 5.57056e+06, 0, -2.62144e+06, 0, 524288},
			wantLen: 21,
			wantDeg: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPolyChebyshev1(tc.arg)

			t.Log(got.Stringn(0))

			assert.Equal(t, tc.wantCoefs, got.coef)
			assert.Equal(t, tc.wantLen, got.len)
			assert.Equal(t, tc.wantDeg, got.deg)
		})
	}
}

func Test_NewPolyChebyshev2(t *testing.T) {

	testCases := []struct {
		name      string
		arg       int
		wantCoefs []float64
		wantLen   int
		wantDeg   int
	}{
		{
			name:      "n = 0",
			arg:       0,
			wantCoefs: []float64{1},
			wantLen:   1,
			wantDeg:   0,
		},
		{
			name:      "n = 1",
			arg:       1,
			wantCoefs: []float64{0, 2},
			wantLen:   2,
			wantDeg:   1,
		},
		{
			name:      "n = 2",
			arg:       2,
			wantCoefs: []float64{-1, 0, 4},
			wantLen:   3,
			wantDeg:   2,
		},
		{
			name:      "n = 3",
			arg:       3,
			wantCoefs: []float64{0, -4, 0, 8},
			wantLen:   4,
			wantDeg:   3,
		},
		{
			name: "n = 20",
			arg:  20,
			wantCoefs: []float64{1, 0, -220, 0, 7920, 0, -109824, 0, 768768, 0, -3.075072e+06, 0,
				7.45472e+06, 0, -1.114112e+07, 0, 1.0027008e+07, 0, -4.980736e+06, 0, 1.048576e+06},
			wantLen: 21,
			wantDeg: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewPolyChebyshev2(tc.arg)

			t.Log(got.Stringn(0))

			assert.Equal(t, tc.wantCoefs, got.coef)
			assert.Equal(t, tc.wantLen, got.len)
			assert.Equal(t, tc.wantDeg, got.deg)
		})
	}
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
			arg:       NewPolyZero(),
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
			argP: NewPolyZero(),
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
			argP: NewPolyZero(),
			argQ: NewPolyZero(),
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
			argP: NewPoly([]float64{125125124, 1.0000001, 4, 1.000002}),
			argQ: NewPoly([]float64{125125125, 1, 4.0000002, 1}),
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
			arg:       NewPolyZero(),
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
			argP: NewPolyZero(),
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
			argP: NewPolyZero(),
			argQ: NewPolyZero(),
			want: NewPolyZero(),
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
			argP: NewPolyZero(),
			argQ: NewPolyZero(),
			want: NewPolyZero(),
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
			argP: NewPolyZero(),
			argS: 0,
			want: NewPolyZero(),
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
			argP: NewPolyZero(),
			argQ: NewPolyZero(),
			want: NewPolyZero(),
		},
		{
			name: "vanish",
			argP: NewPolyZero(),
			argQ: NewPoly([]float64{12, 24124, 1, 2, 5, 124}),
			want: NewPolyZero(),
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
			argP: NewPolyZero(),
			argQ: NewPolyZero(),
			want: NewPolyZero(),
		},
		{
			name: "vanish",
			argP: NewPolyZero(),
			argQ: NewPoly([]float64{12, 24124, 1, 2, 5, 124}),
			want: NewPolyZero(),
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

func Test_PolyPowPanic(t *testing.T) {

	assert.Panics(t, func() { NewPolyLinear(1, 2).Pow(-1) })
}

func Test_PolyPow(t *testing.T) {

	testCases := []struct {
		name string
		argP Poly
		argN int
		want Poly
	}{
		{
			name: "to the zero",
			argP: NewPolyZero(),
			argN: 0,
			want: NewPolyConst(1),
		},
		{
			name: "to the one",
			argP: NewPolyQuadratic(1, 2, 3),
			argN: 1,
			want: NewPolyQuadratic(1, 2, 3),
		},
		{
			name: "quadratic from binom",
			argP: NewPolyLinear(1, -2),
			argN: 2,
			want: NewPolyQuadratic(1, -4, 4),
		},
		{
			name: "large n",
			argP: NewPolyLinear(1, 0),
			argN: 100,
			want: NewPoly(append([]float64{1}, make([]float64, 100)...)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.argP.Pow(tc.argN)

			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_PolyDivZeroPanic(t *testing.T) {

	assert.Panics(t, func() { NewPolyZero().Div(NewPolyZero()) })
	assert.Panics(t, func() { NewPoly([]float64{123}).Div(NewPolyZero()) })
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
			argP:    NewPolyZero(),
			argQ:    NewPoly([]float64{1}),
			wantQuo: NewPolyZero(),
			wantRem: NewPolyZero(),
		},
		{
			name:    "equal linear",
			argP:    NewPoly([]float64{125, 254}),
			argQ:    NewPoly([]float64{125, 254}),
			wantQuo: NewPoly([]float64{1}),
			wantRem: NewPolyZero(),
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
			wantQuo: NewPolyZero(),
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

func Test_PolyReciprocal(t *testing.T) {
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

func Test_PolyString(t *testing.T) {
	testCases := []struct {
		name string
		arg  Poly
		want string
	}{
		{
			name: "zero",
			arg:  NewPolyZero(),
			want: "[ 0.000000x^0 ]",
		},
		{
			name: "negative const",
			arg:  NewPoly([]float64{-21}),
			want: "[ -21.000000x^0 ]",
		},
		{
			name: "nonzero const",
			arg:  NewPoly([]float64{3.1415}),
			want: "[ 3.141500x^0 ]",
		},
		{
			name: "linear",
			arg:  NewPoly([]float64{125124, 1.1715}),
			want: "[ 125124.000000x^1 + 1.171500x^0 ]",
		},
		{
			name: "quadratic",
			arg:  NewPoly([]float64{47346346, 734334, 2342366}),
			want: "[ 47346346.000000x^2 + 734334.000000x^1 + 2342366.000000x^0 ]",
		},
		{
			name: "cubic",
			arg:  NewPoly([]float64{2152, 47346346, 734334, 2342366}),
			want: "[ 2152.000000x^3 + 47346346.000000x^2 + 734334.000000x^1 + 2342366.000000x^0 ]",
		},
		{
			name: "cubic with negative",
			arg:  NewPoly([]float64{-2152, 47346346, -734334, 2342366}),
			want: "[ -2152.000000x^3 + 47346346.000000x^2 - 734334.000000x^1 + 2342366.000000x^0 ]",
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
			arg:  NewPolyZero(),
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
