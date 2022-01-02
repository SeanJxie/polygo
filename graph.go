package polygo

/*
This file contains a small graphing library built on top of the polygo core.
*/

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg" // For graphics.
)

// A RealPolynomialGraph represents the graph of a set of polynomials.
type RealPolynomialGraph struct {

	// The polynomials to be plotted.
	elements []*RealPolynomial

	// The following slices are used to store computed values so that no re-computation is needed.
	intersections []Point
	yIntercepts   []float64
	roots         []float64

	// Rendering options.
	center      Point
	xResolution int
	yResolution int
	viewX       float64
	viewY       float64
	hViewX      float64
	hViewY      float64
	xRenderStep float64
	gridStep    float64

	// Visual options.
	options *GraphOptions

	// Context which handles the actual graphics.
	context *gg.Context
}

// GraphOptions specify some visual options that the RealPolynomialGraph can have.
//
// The options are:
//  ShowAxis               // show the x = 0 and y = 0 axis lines.
//  ShowGrid               // show the grid.
//  ShowIntersections      // highlight the intersection points of all polynomials.
//  ShowRoots              // highlight the roots of all polynomials.
//  ShowYintercepts        // highlight the y-intercepts of all polynomials.
//  ShowAxisLabels         // label axis values.
//  ShowIntersectionLabsls // label intersection points.
//  ShowRootLabels         // label roots.
//  ShowYinterceptLabels   // label y-intercepts.
type GraphOptions struct {
	ShowAxis          bool
	ShowGrid          bool
	ShowIntersections bool
	ShowRoots         bool
	ShowYintercepts   bool

	ShowAxisLabels         bool
	ShowIntersectionLabels bool
	ShowRootLabels         bool
	ShowYinterceptLabels   bool
}

// Some colour definitions
var colBlack = color.RGBA{0x0, 0x0, 0x0, 0xFF}
var colBlackTrans = color.RGBA{0x0, 0x0, 0x0, 0x10}
var colWhite = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
var colGray = color.RGBA{0xCC, 0xCC, 0xCC, 0xFF}
var colBlue = color.RGBA{0x0, 0x0, 0xFF, 0xFF}
var colGreen = color.RGBA{0x0, 0xFF, 0x0, 0xFF}
var colRed = color.RGBA{0xFF, 0x0, 0x0, 0xFF}
var colMagenta = color.RGBA{0xFF, 0x0, 0xFF, 0xFF}

// NewGraph returns a new *RealPolynomialGraph instance with the provided settings.
//
// If any settings are invalid, an appropriate error is set.
//
// The settings are:
//  center      // the point at which the graph is centered.
//  xResolution // the width of the graph in pixels.
//  yResolution // the height of the graph in pixels.
//  viewX       // the width of the viewing area. For example, a viewX of 1.0 will provide a graph spanning the horizontally closed interval [center.X - 1.0, center.X + 1.0].
//  viewY       // the height of the viewing area. For example, a viewY of 1.0 will provide a graph spanning the vertically closed interval [center.Y - 1.0, center.Y + 1.0].
//  xRenderStep // the detail the polynomial curves are rendered at. The closer this positive value is to 0.0, the more precise the curves will be (recommended to be 0.01).
//  gridStep    // the gap between consecutive axis lines.
//  options     // a *GraphOptions instance.
func NewGraph(elements []*RealPolynomial, center Point, xResolution, yResolution int, viewX, viewY, xRenderStep, gridStep float64, options *GraphOptions) (*RealPolynomialGraph, error) {
	if xResolution <= 0 || yResolution <= 0 {
		return nil, errors.New("xResolution and yResolution must be positive")
	}

	if viewX <= 0.0 || viewY <= 0.0 {
		return nil, errors.New("viewX and viewY must be positive")
	}

	if xRenderStep <= 0.0 {
		return nil, errors.New("xRenderStep must be positive")
	}

	if gridStep <= 0.0 {
		return nil, errors.New("gridStep must be positive")
	}

	var newGraph RealPolynomialGraph

	newGraph.elements = elements
	newGraph.center = center
	newGraph.xResolution = xResolution
	newGraph.yResolution = yResolution
	newGraph.viewX = viewX
	newGraph.viewY = viewY
	newGraph.xRenderStep = xRenderStep
	newGraph.gridStep = gridStep

	newGraph.options = options

	newGraph.context = gg.NewContext(xResolution, yResolution)
	newGraph.hViewX = viewX / 2.0
	newGraph.hViewY = viewY / 2.0

	return &newGraph, nil
}

// SaveAsPNG renders and saves the current instance as a PNG image file to the provided path.
//
// If any rendering errors occur, an error addressing the problem is set.
func (g *RealPolynomialGraph) SaveAsPNG(path string) error {
	var err error

	if err = g.renderAxisAndGrid(); err != nil {
		return err
	}
	if err = g.renderPolynomials(); err != nil {
		return err
	}
	if err = g.renderPointIndicators(); err != nil {
		return err
	}
	if err = g.renderLabels(); err != nil {
		return err
	}
	if err = g.context.SavePNG(path); err != nil {
		return err
	}

	fmt.Printf("Your graph has been successfully saved as %s!\n", path)
	return nil
}

// Render the axis lines that form the grid.
func (g *RealPolynomialGraph) renderAxisAndGrid() error {
	ctx := g.context

	// Set background colour.
	ctx.SetColor(colWhite)
	ctx.Clear()

	ctx.SetColor(colGray)

	if g.options.ShowGrid {
		// Left center to left x.
		for x := math.Floor(g.center.X); x > g.center.X-g.hViewX; x -= g.gridStep {
			tmpPt := g.mapPointToViewport(Point{x, 0.0})

			// Skip x == 0.0 (if the axis lines are actually being drawn) since we want the axis lines to go on top of everything.
			if !g.options.ShowAxis || roundToNearestUnit(x, g.gridStep) != 0.0 {
				drawLineBresenham(ctx, int(tmpPt.X), 0, int(tmpPt.X), g.yResolution)
			}
		}

		// Right center to right x.
		for x := math.Ceil(g.center.X); x < g.center.X+g.hViewX; x += g.gridStep {
			tmpPt := g.mapPointToViewport(Point{x, 0.0})

			if !g.options.ShowAxis || roundToNearestUnit(x, g.gridStep) != 0.0 {
				drawLineBresenham(ctx, int(tmpPt.X), 0, int(tmpPt.X), g.yResolution)
			}
		}

		// Bottom center to bottom y.
		for y := math.Floor(g.center.Y); y > g.center.Y-g.hViewY; y -= g.gridStep {
			tmpPt := g.mapPointToViewport(Point{0.0, y})

			if !g.options.ShowAxis || roundToNearestUnit(y, g.gridStep) != 0.0 {
				drawLineBresenham(ctx, 0, int(tmpPt.Y), g.xResolution, int(tmpPt.Y))
			}
		}

		// Top center to top y.
		for y := math.Ceil(g.center.Y); y < g.center.Y+g.viewY; y += g.gridStep {
			tmpPt := g.mapPointToViewport(Point{0.0, y})

			// Skip y == 0.0 (if the axis lines are actually being drawn) since we want the axis lines to go on top of everything.
			if !g.options.ShowAxis || roundToNearestUnit(y, g.gridStep) != 0.0 {
				drawLineBresenham(ctx, 0, int(tmpPt.Y), g.xResolution, int(tmpPt.Y))
			}
		}
	}

	if g.options.ShowAxis {

		// If the closed interval [a, b] contains or touches 0, we must have that ab <= 0.
		xAxisInViewport := (g.center.X-g.hViewX)*(g.center.X+g.hViewX) <= 0.0
		yAxisInViewport := (g.center.Y-g.hViewY)*(g.center.Y+g.hViewY) <= 0.0

		origin := g.mapPointToViewport(Point{0.0, 0.0})

		ctx.SetColor(colBlackTrans) // Axis lines will be darker

		if xAxisInViewport { // Draw x-axis if visible
			drawLineBresenham(ctx, int(origin.X), 0, int(origin.X), g.yResolution)
		}

		if yAxisInViewport { // Draw y-axis if visible
			drawLineBresenham(ctx, 0, int(origin.Y), g.xResolution, int(origin.Y))
		}
	}

	return nil
}

// Render the polynomials.
func (g *RealPolynomialGraph) renderPolynomials() error {
	ctx := g.context
	rand.Seed(time.Now().UnixNano())

	var currGraphY, prevGraphY float64
	var prevPt Point
	var currIsInView, prevIsInView bool

	for _, p := range g.elements {

		// We'll use a random dark colour for each polynomial.
		ctx.SetRGB(
			rand.Float64()*0.8,
			rand.Float64()*0.8,
			rand.Float64()*0.8,
		)

		prevGraphY = p.At(g.center.X - g.hViewX)
		prevPt = g.mapPointToViewport(Point{g.center.X - g.hViewX, prevGraphY})

		for x := g.center.X - g.hViewX + g.xRenderStep; x < g.center.X+g.hViewX; x += g.xRenderStep {
			currGraphY = p.At(x)

			currIsInView = (g.center.Y-g.hViewY <= currGraphY && currGraphY <= g.center.Y+g.hViewY)
			prevIsInView = (g.center.Y-g.hViewY <= prevGraphY && prevGraphY <= g.center.Y+g.hViewY)

			if (currIsInView && prevIsInView) || (currIsInView && !prevIsInView) || (!currIsInView && prevIsInView) {
				currPt := g.mapPointToViewport(Point{x, currGraphY})
				drawLineBresenham(ctx, int(prevPt.X), int(prevPt.Y), int(currPt.X), int(currPt.Y))
				prevPt = currPt
				prevGraphY = currGraphY
			}
		}
	}

	return nil
}

// Render all the graphics specified by the additional options.
func (g *RealPolynomialGraph) renderPointIndicators() error {
	ctx := g.context

	// A "hacky" way of scaling distance based on resolution.
	// We want our intersection markers to be of constant unit radius 0.08.
	markerRad := g.mapPointToViewport(Point{g.center.X - (g.hViewX) + 0.08, 0.0}).X

	for _, p := range g.elements {

		if g.options.ShowRoots {
			ctx.SetColor(colBlue)
			roots, err := p.FindRootsWithin(g.center.X-g.hViewX, g.center.X+g.hViewX)
			if err != nil {
				return err
			}
			g.roots = append(g.roots, roots...)
			for _, rt := range roots {
				tmpPt := g.mapPointToViewport(Point{rt, 0.0})
				ctx.DrawCircle(tmpPt.X, tmpPt.Y, markerRad)
			}
			ctx.Stroke()
		}

		if g.options.ShowYintercepts {
			ctx.SetColor(colGreen)
			yInt := p.At(0.0)
			g.yIntercepts = append(g.yIntercepts, yInt)
			tmpPt := g.mapPointToViewport(Point{0.0, yInt})
			ctx.DrawCircle(tmpPt.X, tmpPt.Y, markerRad)
			ctx.Stroke()
		}

		if g.options.ShowIntersections {
			ctx.SetColor(colMagenta)
			var pois []Point
			for _, p2 := range g.elements {
				if !p.Equal(p2) {
					tmp, err := p.FindIntersectionsWithin(g.center.X-g.hViewX, g.center.X+g.hViewX, p2)
					if err != nil {
						return err
					}
					pois = append(pois, tmp...)
				}
			}

			g.intersections = append(g.intersections, pois...)

			for _, poi := range pois {
				tmpPt := g.mapPointToViewport(Point{poi.X, poi.Y})
				ctx.DrawCircle(tmpPt.X, tmpPt.Y, markerRad)
			}
			ctx.Stroke()
		}
	}

	return nil
}

// Render the axis and/or point indicator labels.
func (g *RealPolynomialGraph) renderLabels() error {
	ctx := g.context
	ctx.SetColor(colBlack)

	if g.options.ShowAxisLabels {
		alernate := true
		zeroDrawn := false
		var tmpPt Point

		// Center to left x.
		for x := math.Floor(g.center.X); x > g.center.X-g.hViewX-g.gridStep; x -= g.gridStep {
			if x == 0.0 {
				zeroDrawn = true
			}

			tmpPt = g.mapPointToViewport(Point{x, 0.0})
			if alernate {
				// The "%g" format removes all trailing zeroes. I hope who ever came up with that lives a long, healthy life.
				ctx.DrawStringAnchored(fmt.Sprintf("%g", x), tmpPt.X+1, tmpPt.Y-1, 0.0, 0.0)
			} else {
				ctx.DrawStringAnchored(fmt.Sprintf("%g", x), tmpPt.X+1, tmpPt.Y+ctx.FontHeight()+1, 0.0, 0.0)
			}

			alernate = !alernate
		}

		alernate = !alernate
		if !zeroDrawn {
			tmpPt = g.mapPointToViewport(Point{0.0, 0.0})
			if alernate {
				ctx.DrawStringAnchored(fmt.Sprintf("%g", 0.0), tmpPt.X+1, tmpPt.Y-1, 0.0, 0.0)
			} else {
				ctx.DrawStringAnchored(fmt.Sprintf("%g", 0.0), tmpPt.X+1, tmpPt.Y+ctx.FontHeight()+1, 0.0, 0.0)
			}
		}

		// Center to right x.
		for x := math.Ceil(g.center.X); x < g.center.X+g.hViewX; x += g.gridStep {
			if x != 0.0 {
				tmpPt = g.mapPointToViewport(Point{x, 0.0})
				if alernate {
					ctx.DrawStringAnchored(fmt.Sprintf("%g", x), tmpPt.X+1, tmpPt.Y-1, 0.0, 0.0)
				} else {
					ctx.DrawStringAnchored(fmt.Sprintf("%g", x), tmpPt.X+1, tmpPt.Y+ctx.FontHeight()+1, 0.0, 0.0)
				}
			}

			alernate = !alernate
		}

		// Center to top y.
		for y := math.Ceil(g.center.Y); y < g.center.Y+g.viewY; y += g.gridStep {
			if roundToNearestUnit(y, g.gridStep) != 0.0 { // Ignore 0.
				tmpPt := g.mapPointToViewport(Point{0.0, y})
				ctx.DrawStringAnchored(fmt.Sprintf("%g", y), tmpPt.X+1, tmpPt.Y-ctx.FontHeight()-1, 0.0, 1.0)
			}
		}

		// Center to bottom y.
		for y := math.Floor(g.center.Y); y > g.center.Y-g.hViewY-g.gridStep; y -= g.gridStep {
			if roundToNearestUnit(y, g.gridStep) != 0.0 { // Ignore 0.
				tmpPt := g.mapPointToViewport(Point{0.0, y})
				ctx.DrawStringAnchored(fmt.Sprintf("%g", y), tmpPt.X+1, tmpPt.Y-ctx.FontHeight()-1, 0.0, 1.0)
			}
		}
	}

	// Process significant points if they have not already been processed via renderPointIndicators.
	if g.intersections == nil {
		for _, p := range g.elements {
			for _, p2 := range g.elements {
				if !p.Equal(p2) {
					tmp, err := p.FindIntersectionsWithin(g.center.X-g.hViewX, g.center.X+g.hViewX, p2)
					if err != nil {
						return err
					}
					g.intersections = append(g.intersections, tmp...)
				}
			}
		}
	}

	if g.roots == nil {
		for _, p := range g.elements {
			roots, err := p.FindRootsWithin(g.center.X-g.hViewX, g.center.X+g.hViewX)
			if err != nil {
				return err
			}
			g.roots = append(g.roots, roots...)
		}
	}

	if g.yIntercepts == nil {
		for _, p := range g.elements {
			yInt := p.At(0.0)
			g.yIntercepts = append(g.yIntercepts, yInt)
		}
	}

	if g.options.ShowIntersectionLabels {
		ctx.SetColor(colMagenta)
		for _, pt := range g.intersections {
			tmpPt := g.mapPointToViewport(Point{pt.X + 0.05, pt.Y + 0.05})
			ctx.DrawStringAnchored(fmt.Sprintf("(%.2f, %.2f)", pt.X, pt.Y), tmpPt.X, tmpPt.Y, 0.0, 0.0)
		}
	}

	if g.options.ShowRootLabels {
		ctx.SetColor(colBlue)

		alternate := true

		for _, x := range g.roots {
			tmpPt := g.mapPointToViewport(Point{x + 0.05, 0.05})
			if alternate { // Alternate the root labels to minimize cluster.
				ctx.DrawStringAnchored(fmt.Sprintf("(%.2f, %f)", x, 0.0), tmpPt.X, tmpPt.Y, 0.0, 0.0)
			} else {
				ctx.DrawStringAnchored(fmt.Sprintf("(%.2f, %f)", x, 0.0), tmpPt.X, tmpPt.Y, 0.0, 1.0)
			}
			alternate = !alternate
		}
	}

	if g.options.ShowYinterceptLabels {
		ctx.SetColor(colGreen)
		for _, y := range g.yIntercepts {
			tmpPt := g.mapPointToViewport(Point{0.05, y + 0.05})
			ctx.DrawStringAnchored(fmt.Sprintf("(%f, %.2f)", 0.0, y), tmpPt.X, tmpPt.Y, 0.0, 0.0)
		}
	}

	return nil
}

// Map a point on the graph coordinate system to the corresponding pixel coordinate.
func (g *RealPolynomialGraph) mapPointToViewport(p Point) Point {
	return Point{
		(p.X + g.hViewX - g.center.X) * float64(g.xResolution) / g.viewX,
		(-p.Y + g.hViewY + g.center.Y) * float64(g.yResolution) / g.viewY,
	}
}

// Credit to https://github.com/StephaneBunel/bresenham/blob/master/drawline.go for the implementation.
func drawLineBresenham(ctx *gg.Context, x1, y1, x2, y2 int) {
	var dx, dy, e, slope int

	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		ctx.SetPixel(x1, y1)

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			ctx.SetPixel(x1, y1)
			x1++
		}
		ctx.SetPixel(x1, y1)

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1 = y2
		}
		for ; dy != 0; dy-- {
			ctx.SetPixel(x1, y1)
			y1++
		}
		ctx.SetPixel(x1, y1)

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				ctx.SetPixel(x1, y1)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				ctx.SetPixel(x1, y1)
				x1++
				y1--
			}
		}
		ctx.SetPixel(x1, y1)

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				ctx.SetPixel(x1, y1)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				ctx.SetPixel(x1, y1)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		ctx.SetPixel(x2, y2)

	// higher than wide.
	default:
		if y1 < y2 {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				ctx.SetPixel(x1, y1)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				ctx.SetPixel(x1, y1)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		ctx.SetPixel(x2, y2)
	}
}
