package polygo

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {
	p1, _ := NewRealPolynomial([]float64{0, -2, 0, 1})
	p2, _ := NewRealPolynomial([]float64{-5, -2, 5, 1})
	fmt.Println(p1)
	fmt.Println(p2)

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

	graph, err := NewGraph([]*RealPolynomial{p2, p1}, Point{X: 0, Y: 0}, 1000, 1000, 5, 5, 0.01, 1.0, &graphOptions)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = graph.SaveAsPNG("test.png")
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
}
