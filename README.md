# polygo
A polynomial library written in Go.

## Features
- Methods of initialization:
	- From a slice of coefficients
	- Parse from string
	- Zero polynomial

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

- Solvers:
	- Work in progress.

- Grapher:
	- Work in progress.

## Installation
```
go get -u github.com/SeanJxie/polygo
```

## Documentation
Godoc: https://pkg.go.dev/github.com/seanjxie/polygo

