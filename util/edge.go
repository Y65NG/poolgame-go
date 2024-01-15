package util

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	_edgeLength float64 = _boardWidth - 2*math.Sqrt2*_basketRadius
	_edgeWidth          = math.Sqrt2 / 2 * _basketRadius
)

type Edge struct {
	*ebiten.Image
	*resolv.Object
}

type EdgePosition int

const (
	_edgeTop EdgePosition = iota + 1
	_edgeTopLeft
	_edgeBottomLeft
	_edgeBottom
	_edgeBottomRight
	_edgeTopRight
)

func NewEdge(pos EdgePosition) *Edge {
	w, l := _edgeWidth, _edgeLength
	ll := _edgeLength + _ballRadius/2
	var edge *Edge
	var x, y float64
	switch pos {
	case _edgeTop:
		x, y = math.Sqrt2*_basketRadius, 0
		edge = &Edge{
			Image:  ebiten.NewImage(int(l+1), int(w+1)),
			Object: resolv.NewObject(x, y, l, w, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, 0,
				l, 0,
				l-w, w,
				w, w,
			))
	case _edgeTopLeft:
		x, y = 0, math.Sqrt2*_basketRadius
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+1), int(ll+2)),
			Object: resolv.NewObject(x, y, w, l+8, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, 0,
				w, w,
				w, ll-w*math.Tan(20./180*math.Pi),
				0, ll,
			))
	case _edgeBottomLeft:
		x, y = 0, _boardHeight-math.Sqrt2*_basketRadius-ll
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+1), int(ll+2)),
			Object: resolv.NewObject(x, y, w, ll, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, 0,
				w, w*math.Tan(20./180*math.Pi),
				w, ll-w,
				0, ll,
			))
	case _edgeBottom:
		x, y = math.Sqrt2*_basketRadius, _boardHeight-w
		edge = &Edge{
			Image:  ebiten.NewImage(int(l+1), int(w+1)),
			Object: resolv.NewObject(x, y, l, w, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				w, 0,
				l-w, 0,
				l, w,
				0, w,
			))
	case _edgeBottomRight:
		x, y = _boardWidth-w, _boardHeight-math.Sqrt2*_basketRadius-ll
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+1), int(ll+2)),
			Object: resolv.NewObject(x, y, w, ll, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, w*math.Tan(20./180*math.Pi),
				w, 0,
				w, ll,
				0, ll-w,
			))
	case _edgeTopRight:
		x, y = _boardWidth-w, math.Sqrt2*_basketRadius
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+1), int(ll+2)),
			Object: resolv.NewObject(x, y, w, ll, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, w,
				w, 0,

				w, ll,
				0, ll-w*math.Tan(20./180*math.Pi),
			))
	default:
		panic("invalid edge position")
	}
	return edge
}

func (e *Edge) Move() {}

func (e *Edge) Collide(other Object) {
	if ball, ok := other.(*Ball); ok {
		if intersection := ball.Shape.Intersection(ball.velocity.X, ball.velocity.Y, e.Shape); intersection != nil {
			pts := intersection.Points
			rp := Vec{pts[0].X() - pts[1].X(), pts[0].Y() - pts[1].Y()}
			v := ball.velocity
			vp := v.Project(rp)
			vn := v.Sub(vp)
			ball.velocity = vp.Add(vn.Negate()) // normal velocity reversed, parallel velocity unchanged
		}
	}

}

func (e *Edge) draw() {
	// b1, b2 := e.Shape.Bounds()
	// vector.StrokeRect(e.Image, 0, 0, float32(b2.X()-b1.X()), float32(b2.Y()-b1.Y()), 1, color.White, true)
	vector.StrokeRect(e.Image, 0.1, 0, float32(e.W), float32(e.H), 1, color.White, true)
}
