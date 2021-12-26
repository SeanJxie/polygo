# PolyGo
A collection of basic tools that make working with polynomials easier in Go.

## Installation
```
go get -u github.com/SeanJxie/polygo
```

## Documentation
You can find the full list of functions through godoc: https://pkg.go.dev/github.com/SeanJxie/polygo

## Examples

- Create a simple quadratic and find the root:
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
	
	quad.PrintExpr() // 0.000000x^0 + 0.000000x^1 + 1.000000x^2
	fmt.Printf("Root: %f\n", root) // Root: 0.000000
}
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
	poly.PrintExpr() // 5.000000x^0 + 2.000000x^1 + 5.000000x^2 + 2.000000x^3 + 63.000000x^4 + 1.000000x^5 + 2.000000x^6 + 5.000000x^7 + 1.000000x^8
	deriv.PrintExpr() // 2.000000x^0 + 10.000000x^1 + 6.000000x^2 + 252.000000x^3 + 5.000000x^4 + 12.000000x^5 + 35.000000x^6 + 8.000000x^7
}
```

