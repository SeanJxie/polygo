# polygo
`polygo` is a light-weight library that makes working with polynomials possible and easy in Go.

## What's New
- A robust graphing tool

## What's to Come
- More polynomial tools (integration, critial point finding, binomial expansion, etc.)

## Installation
```
go get -u github.com/SeanJxie/polygo
```

## Documentation
You can find the full list of functions through godoc: https://pkg.go.dev/github.com/seanjxie/polygo

## Examples
- Graph functions:
```go
package main

import (
	"fmt"
	"github.com/SeanJxie/polygo"
)

func main() {
	p1, _ := polygo.NewRealPolynomial([]float64{0, -2, 0, 1})
	p2, _ := polygo.NewRealPolynomial([]float64{-5, -2, 5, 1})

	graphOptions := GraphOptions{
		ShowIntersections:      true,
		ShowIntersectionLabels: true,

		ShowRoots:      true,
		ShowRootLabels: true,

		ShowYintercepts:      true,
		ShowYinterceptLabels: true,

		ShowAxis:       true,
		ShowAxisLabels: true,
		ShowGrid:       true,

		DarkMode: false,
	}

	graph, _ := polygo.NewGraph([]*polygo.RealPolynomial{p2, p1}, polygo.Point{X: 0, Y: 0}, 1000, 1000, 5, 5, 0.01, 1.0, &graphOptions)
	graph.SaveAsPNG("graph1.png")
}
```
Output:
![graph1.png](https://github.com/SeanJxie/polygo/blob/main/graph_samples/graph1.png)

>You can find some more image examples in the [`graph_samples`](https://github.com/SeanJxie/polygo/tree/main/graph_samples) folder.

- Create a simple quadratic and find its root:
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
	
	fmt.Println(quad)
	fmt.Printf("Root: %f\n", root) 
}
```
Output:
```
0.000000x^0 + 0.000000x^1 + 1.000000x^2
Root: 0.000000
```

- Create a polynomial and find its derivative:
```go
package main

import (
	"github.com/SeanJxie/polygo"
)

func main() {
	coeffs := []float64{5, 2, 5, 2, 63, 1, 2, 5, 1}
	poly, _ := polygo.NewRealPolynomial(coeffs)
	deriv := poly.Derivative()
	fmt.Println(poly)
	fmt.Println(deriv)
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
	fmt.Println(cubic)
	fmt.Println(affine)
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
