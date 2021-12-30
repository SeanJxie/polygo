package polygo

import (
	"testing"
)

// test instance of QuadCoefficents, print Expression and root value
func TestQuadCoefficients(t *testing.T) {
	quadCoefficients := []float64{0, 0, 2}
	quad, _ := NewRealPolynomial(quadCoefficients)
	root, _ := quad.FindRootWithin(-1, 1)

	quadraticExpression := quad.String()
	if quadraticExpression != "0.000000x^0 + 0.000000x^1 + 2.000000x^2" {
		t.Fatalf("error in calculation -'%v'", quadraticExpression)
	}

	if root != 0.0 {
		t.Fatalf("error root calculation expected 0.00000 found -'%v'", root)
	}
}
