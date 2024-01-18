package util

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	_ballRadius = 12
	_ballMass   = 1
)

type ballID int

const (
	_idWhite ballID = 0
	_idBlack ballID = 8
)

type ballKind int

const (
	_kindSolid ballKind = iota + 1
	_kindStripe
	_kindBlack
	_kindWhite
)

func (k ballKind) String() string {
	switch k {
	case _kindSolid:
		return "solid"
	case _kindStripe:
		return "stripe"
	case _kindBlack:
		return "black"
	case _kindWhite:
		return "undecided"
	}
	return ""
}

var (
	stripeColorImage    = ebiten.NewImage(3, 3)
	stripeColorSubImage = stripeColorImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	stripeColorImage.Fill(_colorWhite)
}

type Ball struct {
	*ebiten.Image
	*resolv.Object

	id       ballID
	kind     ballKind
	mass     float64
	velocity Vec
	balls    []*Ball
	catched  bool

	station *Station
}

func NewBall(balls []*Ball, id ballID, x, y float64, station *Station) *Ball {
	r := _ballRadius
	var kind ballKind
	switch id {
	case _idWhite:
		kind = _kindWhite
	case _idBlack:
		kind = _kindBlack
	case 1, 2, 3, 4, 5, 6, 7:
		kind = _kindSolid
	case 9, 10, 11, 12, 13, 14, 15:
		kind = _kindStripe
	}
	b := &Ball{
		Image: ebiten.NewImage(2*int(r)+5, 2*int(r)+5),
		// Circle: resolv.NewCircle(x, y, _ballRadius),
		Object: resolv.NewObject(x, y, 2*_ballRadius, 2*_ballRadius, "ball"),

		id:    id,
		kind:  kind,
		mass:  _ballMass,
		balls: balls,

		station: station,
	}
	b.Object.SetShape(resolv.NewCircle(0, 0, _ballRadius))
	return b
}

func NewWhiteball(balls []*Ball, station *Station) *Ball {
	return NewBall(balls, _idWhite, _boardWidth/2, _boardHeight/4*3, station)
}

func (b *Ball) Move() {

	var (
		x, y   = b.X, b.Y
		vx, vy = b.velocity.X, b.velocity.Y
	)
	if vx*vx+vy*vy < 0.005 {
		b.velocity.X, b.velocity.Y = 0, 0
	}

	nx, ny := x+vx+RandomFactor(0, 0.0001), y+vy+RandomFactor(0, 0.0001)

	b.X, b.Y = nx, ny
	b.Update()

	b.velocity.X *= 1 - _friction
	b.velocity.Y *= 1 - _friction

	for _, ball := range b.balls {
		if b.id != ball.id {
			b.Collide(ball)
		}
	}
}

func (b *Ball) Collide(other Object) {

	if other, ok := other.(*Ball); ok {
		if b.mass == 0 || b.catched || other.mass == 0 || other.catched {
			return
		}
		var (
			vx1, vy1 = b.velocity.X, b.velocity.Y
			vx2, vy2 = other.velocity.X, other.velocity.Y

			// dx, dy = vx1, vy1
		)

		// if vx1 > 1 {
		// 	dx = 1
		// } else if vx1 < -1 {
		// 	dx = -1
		// }

		// if vy1 > 1 {
		// 	dy = 1
		// } else if vy1 < -1 {
		// 	dy = -1
		// }

		// newB1, newB2 := b.Clone(), other.Clone()

		if intersection := b.Clone().Shape.Intersection(vx1, vy1, other.Clone().Shape); intersection != nil {
			var (
				p1, p2   = Vec{b.X, b.Y}, Vec{other.X, other.Y}
				v1, v2   = Vec{vx1, vy1}, Vec{vx2, vy2}
				rp       = p2.Add(p1.Negate()).Normalize()
				v1p, v2p = v1.Project(rp), v2.Project(rp)
				v1n, v2n = v1.Add(v1p.Negate()), v2.Add(v2p.Negate())
				m1, m2   = b.mass, other.mass
			)
			v1 = v1p.Scale((m1 - m2) / (m1 + m2)).Add(v2p.Scale(2 * m2 / (m1 + m2))).Add(v1n)
			v2 = v2p.Scale((m2 - m1) / (m1 + m2)).Add(v1p.Scale(2 * m1 / (m1 + m2))).Add(v2n)
			b.velocity.X, b.velocity.Y = v1.X, v1.Y
			other.velocity.X, other.velocity.Y = v2.X, v2.Y
			if b.station.FirstCollidedBall == nil {
				b.station.FirstCollidedBall = other
			}
		}

	}

}

func (b *Ball) containsPos(x, y float64) bool {
	b1, b2 := b.Shape.Bounds()
	if x >= b1.X() && x <= b2.X() && y >= b1.Y() && y <= b2.Y() {
		return true
	}
	return false
}

func (b *Ball) draw() {

	var c color.Color
	switch b.id {
	case _idWhite:
		c = _colorWhite

	case _idBlack:
		c = _colorBlack
	case 1, 9:
		c = _colorYellow
	case 2, 10:
		c = _colorBlue
	case 3, 11:
		c = _colorRed
	case 4, 12:
		c = _colorPurple
	case 5, 13:
		c = _colorOrange
	case 6, 14:
		c = _colorGreen
	case 7, 15:
		c = _colorMaroon
	}

	vector.DrawFilledCircle(
		b.Image,
		_ballRadius,
		_ballRadius,
		float32(_ballRadius),
		c,
		true,
	)
	if b.kind == _kindStripe {
		var path vector.Path
		path.MoveTo(float32(_ballRadius)*1.5, float32(_ballRadius*(1-math.Sqrt(3)/2)))
		path.Arc(_ballRadius, _ballRadius, _ballRadius, 60*math.Pi/180, -60*math.Pi/180, vector.CounterClockwise)
		path.Close()

		path.MoveTo(float32(_ballRadius)*0.5, float32(_ballRadius*(1-math.Sqrt(3)/2)))
		path.Arc(_ballRadius, _ballRadius, _ballRadius, 120*math.Pi/180, -120*math.Pi/180, vector.Clockwise)
		path.Close()
		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		op := &ebiten.DrawTrianglesOptions{AntiAlias: true}
		b.Image.DrawTriangles(vs, is, stripeColorSubImage, op)

	}
}
