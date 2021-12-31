package polygo

import (
	"testing"
)

// test instance of QuadCoefficents, print Expression and root value
func TestQuadCoefficients(t *testing.T) {
	quadCoefficients := []float64{0, 0, 2}
	quad, _ := NewRealPolynomial(quadCoefficients)
	root, _ := quad.FindRootWithin(-1, 1)

	quadraticExpression := quad.Expr()
	if quadraticExpression != "0.000000x^0 + 0.000000x^1 + 2.000000x^2" {
		t.Fatalf("error in calculation -'%v'", quadraticExpression)
	}

	if root != 0.0 {
		t.Fatalf("error root calculation expected 0.00000 found -'%v'", root)
	}
}

// test find Derivate
func TestFindDerivate(t *testing.T) {
	coeffs := []float64{5, 2, 5, 2, 63, 1, 2, 5, 1}
	poly, _ := NewRealPolynomial(coeffs)
	derivate := poly.Derivative()

	polyExpress := poly.Expr()
	if polyExpress != "5.000000x^0 + 2.000000x^1 + 5.000000x^2 + 2.000000x^3 + 63.000000x^4 + 1.000000x^5 + 2.000000x^6 + 5.000000x^7 + 1.000000x^8" {
		t.Fatalf("error in calculation -'%v'", polyExpress)
	}

	valDerivate := derivate.Expr()
	if valDerivate != "2.000000x^0 + 10.000000x^1 + 6.000000x^2 + 252.000000x^3 + 5.000000x^4 + 12.000000x^5 + 35.000000x^6 + 8.000000x^7" {
		t.Fatalf("error in calculation -'%v'", valDerivate)
	}

}

// test find intersection
func TestIntersection(t *testing.T) {
	cubic, _ := NewRealPolynomial([]float64{0, -2, 0, 1})
	affine, _ := NewRealPolynomial([]float64{3, 5})

	_, err := cubic.FindIntersectionsWithin(-10, 10, affine)

	if err != nil {
		t.Fatalf("error in find intersection")
	}

}
