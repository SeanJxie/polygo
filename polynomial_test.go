package polygo

import (
	"testing"
)

var Cons1, Cons2, Line1, Line2, Line3, Line4, Quad1, Quad2, Quad3, Cube1, Cube2, Cube3 *RealPolynomial

func init() {
	// Testing polynomials
	var err error

	Cons2, err = NewRealPolynomial([]float64{18})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Line1, err = NewRealPolynomial([]float64{4, 5})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Line2, err = NewRealPolynomial([]float64{3, 5})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Line3, err = NewRealPolynomial([]float64{4, 5})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Line4, err = NewRealPolynomial([]float64{3, 5})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Quad1, err = NewRealPolynomial([]float64{-1, 1, 1, 0, 0, 0, 0, 0})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Quad2, err = NewRealPolynomial([]float64{12, 42, -56})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Quad3, err = NewRealPolynomial([]float64{0, 0, 1})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Cube1, err = NewRealPolynomial([]float64{11, 7, 3, 5})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Cube2, err = NewRealPolynomial([]float64{1, -12, 1, 1})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Cube3, err = NewRealPolynomial([]float64{0, -2, 0, 1})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

func TestFindRootWithin(t *testing.T) {
	root, err := Quad3.FindRootWithin(0, 1)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	t.Logf("Root: %f\n", root)
}

func TestFindRootsWithin(t *testing.T) {
	roots, err := Cube3.FindRootsWithin(-5, 5)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}

	t.Logf("Roots: %v\n", roots)
}

func TestCountRoots(t *testing.T) {
	nRoots := Quad2.CountRootsWithin(-5, 5)
	t.Logf("Number of roots: %d\n", nRoots)
}

func TestFindIntersectionsWithin(t *testing.T) {
	pois, err := Line1.FindIntersectionsWithin(-5, 5, Line2)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	t.Logf("POIs: %v\n", pois)
}

func TestMul(t *testing.T) {
	Line1.Mul(Line2)
	t.Log(Line1)
}

func TestMulComp(t *testing.T) {
	p1Size := 8
	p2Size := 10

	largePoly1Coeffs := make([]float64, p1Size)
	largePoly2Coeffs := make([]float64, p2Size)

	for i := 0; i < p1Size; i++ {
		largePoly1Coeffs[i] = rand.Float64() * 10
	}

	for i := 0; i < p2Size; i++ {
		largePoly2Coeffs[i] = rand.Float64() * 10
	}

	largePoly1, _ := NewRealPolynomial(largePoly1Coeffs)
	largePoly2, _ := NewRealPolynomial(largePoly2Coeffs)

	start := time.Now()
	prodNaive := largePoly1.MulNaive(largePoly2)
	elapsed := time.Since(start)
	t.Logf("Naive took: %s\n", elapsed)

	start = time.Now()
	prodFast := largePoly1.Mul(largePoly2)
	elapsed = time.Since(start)

	t.Logf("Fast took: %s\n", elapsed)

	t.Log(prodNaive)
	t.Log(prodFast)

	t.Logf("%t\n", prodFast.Equal(prodNaive))

}

func TestShift(t *testing.T) {
	s := Line1.ShiftRight(2)
	t.Log(s)
}

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
