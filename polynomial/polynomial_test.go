package polynomial_test

import (
	"fmt"
	"main/polynomial"
	"testing"
)

var Cons1, Cons2, Line1, Line2, Quad1, Quad2, Cube1, Cube2 *polynomial.RealPolynomial

func init() {
	// Testing polynomials
	var err error

	Cons2, err = polynomial.NewRealPolynomial([]float64{18})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Line1, err = polynomial.NewRealPolynomial([]float64{0, 5})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Line2, err = polynomial.NewRealPolynomial([]float64{4, 3})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Quad1, err = polynomial.NewRealPolynomial([]float64{-1, 1, 1, 0, 0, 0, 0, 0})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Quad2, err = polynomial.NewRealPolynomial([]float64{12, 42, -56})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Cube1, err = polynomial.NewRealPolynomial([]float64{11, 7, 3, 5})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	Cube2, err = polynomial.NewRealPolynomial([]float64{1, -12, 1, 1})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

func TestFindRootWithin(t *testing.T) {
	root, err := Quad2.FindRootWithin(-5, 5)
	Quad2.PrintExpr()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Root: %f\n", root)
}

func TestFindRootsWithin(t *testing.T) {
	var roots []float64
	for i := 0; i < 10000; i++ {
		roots = Cube2.FindRootsWithin(-5, 5)
	}
	Cube2.PrintExpr()
	fmt.Printf("Roots: %v\n", roots)
}

func TestCountRoots(t *testing.T) {
	nRoots := Quad2.CountRootsWithin(-5, 5)
	Quad2.PrintExpr()
	fmt.Printf("Number of roots: %d\n", nRoots)
}
