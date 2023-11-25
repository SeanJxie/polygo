# polygo
A polynomial library written in Go.
```
go get github.com/SeanJxie/polygo@main
```
## Features
- Initialization methods:
	- From a slice of coefficients
	- From a string
	- Special polynomials:
		- Zero polynomial
		- Wilkinson's polynomial


- Binary operations:
	- Addition
	- Subtraction
	- Scalar multiplication
	- Multiplication (with fast variant using an FFT)
	- Euclidean division
	- Equality (based on absolute and relative error)

- Unary operations/properties:
	- Evaluation (using Horner's scheme)
	- Derivative
	- Coefficients (leading, largest, n-th degree, etc.)
	- Degree
	- Reciprocal 
	- Boolean checks (constant, zero, monic, etc.)
	- Cauchy's root bound

- Solver:
	- Solve polynomial equations (i.e. roots, intersections).
	- Various algorithms to count, isolate, and find roots.

- Grapher:
	- Work in progress.

## Documentation
Godoc: https://pkg.go.dev/github.com/seanjxie/polygo

