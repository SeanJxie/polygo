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

	// TODO: check when a, b are multiple roots of p.

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := new_sturmChain(tc.argP).count(tc.argA, tc.argB)

			assert.Equal(t, tc.want, got)
		})
	}
}
