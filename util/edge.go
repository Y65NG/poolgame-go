package util

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	_edgeLength = _boardWidth - 2*math.Sqrt2*_basketRadius
	_edgeWidth  = math.Sqrt2 / 2 * _basketRadius
)

var (
	edgeColorImage    = ebiten.NewImage(3, 3)
	edgeColorSubImage = edgeColorImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	edgeColorImage.Fill(_edgeColor)
}

type Edge struct {
	*ebiten.Image
	*resolv.Object
	position EdgePosition

	station *Station
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

func NewEdge(pos EdgePosition, station *Station) *Edge {
	w, l := _edgeWidth, _edgeLength
	ll := _edgeLength + _ballRadius/2 + 2
	var edge *Edge
	var x, y float64
	wScale := 1000.
	switch pos {
	case _edgeTop:
		x, y = math.Sqrt2*_basketRadius, 0
		edge = &Edge{
			Image:  ebiten.NewImage(int(l+2), int(w+2)),
			Object: resolv.NewObject(x, y, l, (wScale+1)*w, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, 0,
				0, -wScale*w,
				l, -wScale*w,
				l, 0,
				l-w, w,
				w, w,
			))
	case _edgeTopLeft:
		x, y = 0, math.Sqrt2*_basketRadius
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+2), int(ll+2)),
			Object: resolv.NewObject(x, y, 2*w, l+8, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, 0,
				w, w,
				w, ll-w*math.Tan(20./180*math.Pi),
				0, ll,
				-wScale*w, ll,
				-wScale*w, 0,
			))
	case _edgeBottomLeft:
		x, y = 0, _boardHeight-math.Sqrt2*_basketRadius-ll
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+2), int(ll+2)),
			Object: resolv.NewObject(x, y, 2*w, ll, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, 0,
				w, w*math.Tan(20./180*math.Pi),
				w, ll-w,
				0, ll,
				-wScale*w, ll,
				-wScale*w, 0,
			))
	case _edgeBottom:
		x, y = math.Sqrt2*_basketRadius, _boardHeight-w
		edge = &Edge{
			Image:  ebiten.NewImage(int(l+2), int(w+2)),
			Object: resolv.NewObject(x, y, l, 2*w, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, w,
				w, 0,
				l-w, 0,
				l, w,
				l, wScale*w,
				0, wScale*w,
			))
	case _edgeBottomRight:
		x, y = _boardWidth-w, _boardHeight-math.Sqrt2*_basketRadius-ll
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+2), int(ll+2)),
			Object: resolv.NewObject(x, y, 2*w, ll, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, w*math.Tan(20./180*math.Pi),
				w, 0,
				wScale*w, 0,
				wScale*w, ll,
				w, ll,
				0, ll-w,
			))
	case _edgeTopRight:
		x, y = _boardWidth-w, math.Sqrt2*_basketRadius
		edge = &Edge{
			Image:  ebiten.NewImage(int(w+2), int(ll+2)),
			Object: resolv.NewObject(x, y, 2*w, ll, "edge"),
		}
		edge.Object.SetShape(
			resolv.NewConvexPolygon(
				x, y,
				0, w,
				w, 0,
				wScale*w, 0,
				wScale*w, ll,
				w, ll,
				0, ll-w*math.Tan(20./180*math.Pi),
			))
	default:
		panic("invalid edge position")
	}
	edge.position = pos
	edge.station = station
	return edge
}

func (e *Edge) Move() {}

func (e *Edge) Collide(other Object) {
	if ball, ok := other.(*Ball); ok {
		var vx, vy = ball.velocity.X, ball.velocity.Y
		maxV := 1.
		if vx > maxV {
			vx = maxV
		} else if vx < -maxV {
			vx = -maxV
		}
		if vy > maxV {
			vy = maxV
		} else if vy < -maxV {
			vy = -maxV
		}
		if intersection := ball.Shape.Intersection(vx, vy, e.Shape); intersection != nil {
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

	var path vector.Path

	w, l := float32(_edgeWidth), float32(_edgeLength)
	c := float32(_ballRadius/2 + 2)

	switch e.position {
	case _edgeTop:
		path.MoveTo(0, 0)
		path.LineTo(l, 0)
		path.LineTo(l-w, w)
		path.LineTo(w, w)
		// path.LineTo(0, 0)
		path.Close()
	case _edgeTopLeft:
		path.MoveTo(0, 0)
		path.LineTo(w, w)
		path.LineTo(w, l+c-w*float32(math.Tan(20./180*math.Pi)))
		path.LineTo(0, l+c)
		path.Close()
	case _edgeBottomLeft:
		path.MoveTo(0, 0)
		path.LineTo(w, w*float32(math.Tan(20./180*math.Pi)))
		path.LineTo(w, l+c-w)
		path.LineTo(0, l+c)
		path.Close()
	case _edgeBottom:
		path.MoveTo(w, 0)
		path.LineTo(l-w, 0)
		path.LineTo(l, w)
		path.LineTo(0, w)
		path.Close()
	case _edgeBottomRight:
		path.MoveTo(0, w*float32(math.Tan(20./180*math.Pi)))
		path.LineTo(w, 0)
		path.LineTo(w, l+c)
		path.LineTo(0, l+c-w)
		path.Close()
	case _edgeTopRight:
		path.MoveTo(0, w)
		path.LineTo(w, 0)
		path.LineTo(w, l+c)
		path.LineTo(0, l+c-w*float32(math.Tan(20./180*math.Pi)))
		path.Close()

	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	// log.Println(vs, is)

	// for i := range vs {
	// 	vs[i].SrcX, vs[i].SrcY = 1, 1
	// 	vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = float32(_edgeColor.R)/0xff, float32(_edgeColor.G)/0xff, float32(_edgeColor.B)/0xff, 1
	// }
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true}

	e.DrawTriangles(vs, is, edgeColorSubImage, op)
	edgeColorImage.Fill(_edgeColor)
	edgeColorSubImage.Fill(_edgeColor)

}
