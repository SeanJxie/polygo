# polygo
A collection of tools that make working with polynomials easier in Go.

## What's New
- Fast polynomial multiplication through the Fast Fourier Transform.
- Better documentation.

## Installation
```
go get -u github.com/SeanJxie/polygo
```

## Documentation
You can find the full list of functions through godoc: https://pkg.go.dev/github.com/seanjxie/polygo

## Examples

- Create a simple quadratic and solve it:
```go
package main

import (
	"fmt"
	"github.com/SeanJxie/polygo"
)

func main() {
	quadCoefficients := []float64{0, 0, 2}
	quad, _ := polygo.NewRealPolynomial(quadCoefficients)

	root, _ := quad.FindRootWithin(-1, 1)
	
	quad.PrintExpr()
	fmt.Printf("Root: %f\n", root) 
}
```
Output:
```
0.000000x^0 + 0.000000x^1 + 1.000000x^2
Root: 0.000000
```

- Create a polynomial and find the derivative:
```go
package main

import (
	"github.com/SeanJxie/polygo"
)

func main() {
	coeffs := []float64{5, 2, 5, 2, 63, 1, 2, 5, 1}
	poly, _ := polygo.NewRealPolynomial(coeffs)
	deriv := poly.Derivative()
	poly.PrintExpr()
	deriv.PrintExpr()
}
```
Output:
```
5.000000x^0 + 2.000000x^1 + 5.000000x^2 + 2.000000x^3 + 63.000000x^4 + 1.000000x^5 + 2.000000x^6 + 5.000000x^7 + 1.000000x^8
2.000000x^0 + 10.000000x^1 + 6.000000x^2 + 252.000000x^3 + 5.000000x^4 + 12.000000x^5 + 35.000000x^6 + 8.000000x^7
```

- Find the intersection of two polynomials:

```go
func main() {
	cubic, _ := polygo.NewRealPolynomial([]float64{0, -2, 0, 1})
	affine, _ := polygo.NewRealPolynomial([]float64{3, 5})
	cubic.PrintExpr()
	affine.PrintExpr()
	intersections, err := cubic.FindIntersectionsWithin(-10, 10, affine)

	if err != nil {
		fmt.Printf("error %v\n", err)
	} else {
		fmt.Printf("Intersections: %v\n", intersections) 
	}
}
```
Output:
```
Intersections: [[-2.397661540892259 -8.988307704461295] [-0.44080771150488296 0.7959614424755855] [2.8384692523971413 17.1923462619857]]
```
