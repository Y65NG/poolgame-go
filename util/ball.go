package util

import (
	"image/color"

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

type Ball struct {
	*ebiten.Image
	// *resolv.Circle
	*resolv.Object

	id       ballID
	mass     float64
	velocity Vec
	balls    []*Ball
	catched  bool
}

func NewBall(balls []*Ball, id ballID, x, y float64) *Ball {
	r := _ballRadius
	b := &Ball{
		Image: ebiten.NewImage(2*int(r)+1, 2*int(r)+1),
		// Circle: resolv.NewCircle(x, y, _ballRadius),
		Object: resolv.NewObject(x, y, 2*_ballRadius, 2*_ballRadius, "ball"),

		id:    id,
		mass:  _ballMass,
		balls: balls,
	}
	b.Object.SetShape(resolv.NewCircle(0, 0, _ballRadius))
	return b
}

func NewWhiteball(balls []*Ball) *Ball {
	return NewBall(balls, _idWhite, _boardWidth/2, _boardHeight/2+240)
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
	if b.mass == 0 {
		return
	}

	if other, ok := other.(*Ball); ok {
		var (
			vx1, vy1 = b.velocity.X, b.velocity.Y
			vx2, vy2 = other.velocity.X, other.velocity.Y
		)
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
	c := color.RGBA{255, 255, 255, 255}
	if b.id != _idWhite {
		c = color.RGBA{100, 150, 150, 255}
	}

	vector.DrawFilledCircle(
		b.Image,
		_ballRadius,
		_ballRadius,
		float32(_ballRadius),
		c,
		true,
	)
}
