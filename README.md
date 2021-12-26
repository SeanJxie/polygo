# PolyGo
A collection of basic tools that make working with polynomials easier in Go.

## Installation
```
go get -u github.com/SeanJxie/polygo
```

## Documentation
- godoc: https://pkg.go.dev/github.com/SeanJxie/polygo

## Examples

Below, we create a simple quadratic and find its root:
```go
package main

import (
	"fmt"
	"github.com/SeanJxie/polygo"
)

func main() {
	quadCoefficients := []float64{0, 0, 2} // x^2
	quad, _ := polygo.NewRealPolynomial(quadCoefficients)

	root, _ := quad.FindRootWithin(-1, 1)

	fmt.Printf("Root: %f\n", root) // Root: 0.000000
}
```
