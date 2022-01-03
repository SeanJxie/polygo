package polygo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGraphBasic(t *testing.T) {
	p1, err := NewRealPolynomial([]float64{0, -2, 0, 1})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	p2, err := NewRealPolynomial([]float64{-5, -2, 5, 1})
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	t.Log(p1)
	t.Log(p2)

	graphOptions := GraphOptions{
		ShowIntersections:      true,
		ShowAxis:               true,
		ShowAxisLabels:         true,
		ShowIntersectionLabels: true,
		ShowRootLabels:         true,
		ShowRoots:              true,
		ShowYintercepts:        true,
		ShowGrid:               true,
	}

	graph, err := NewGraph([]*RealPolynomial{p2, p1}, Point{X: 0, Y: 0}, 1000, 1000, 10, 10, 0.01, 1.0, &graphOptions)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = graph.SaveAsPNG("TestGraphBasic.png")
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
}

func TestGraphStress(t *testing.T) {
	start := time.Now()

	rand.Seed(time.Now().UnixNano())

	nPolynomials := 500
	nCoeffs := 2
	coeffMax := 10.0
	coeffMin := -10.0

	polynomials := make([]*RealPolynomial, nPolynomials)

	var err error

	for i := 0; i < nPolynomials; i++ {
		tmpCoeffs := make([]float64, nCoeffs)

		for j := 0; j < nCoeffs; j++ {
			tmpCoeffs[j] = coeffMin + rand.Float64()*(coeffMax-coeffMin)
		}

		polynomials[i], err = NewRealPolynomial(tmpCoeffs)
		if err != nil {
			t.Fatalf("error: %v\n", err)
		}
	}

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

	graph, err := NewGraph(polynomials, Point{X: 0, Y: 0}, 5000, 5000, 50, 50, 0.001, 1, &graphOptions)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}

	elapsed := time.Since(start)
	t.Logf("init runtime: %s\n", elapsed)

	start = time.Now()
	err = graph.SaveAsPNG("TestGraphStress.png")
	elapsed = time.Since(start)
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}

	t.Logf("runtime: %s\n", elapsed)
}

func TestGraphFrameAnimation(t *testing.T) {
	graphOptions := GraphOptions{
		ShowIntersections:      true,
		ShowAxis:               true,
		ShowAxisLabels:         true,
		ShowIntersectionLabels: true,
		ShowRootLabels:         true,
		ShowRoots:              true,
		ShowYintercepts:        true,
		ShowGrid:               true,
	}

	frameCount := 0
	for a := -10.0; a <= 10.0; a += 0.1 {
		p1, err := NewRealPolynomial([]float64{0, -2, 0, a})
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		graph, err := NewGraph([]*RealPolynomial{p1}, Point{X: 0, Y: 0}, 1000, 1000, 10, 10, 0.01, 1.0, &graphOptions)
		if err != nil {
			t.Fatalf("error: %v", err)
		}

		err = graph.SaveAsPNG(fmt.Sprintf("frame_%d.png", frameCount))
		if err != nil {
			t.Fatalf("error: %v\n", err)
		}
		frameCount++
	}
}
