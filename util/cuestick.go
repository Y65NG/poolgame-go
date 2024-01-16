package util

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	_stickWidth  = 10
	_stickLength = 150

	_distanceToTarget = 12 + _ballRadius
)

type CueStick struct {
	*ebiten.Image
	// Massive
	*resolv.Object // the tip of the cue stick

	targetBall *Ball
	// Pivot      Vec
	speed      float64
	powerLevel int
	mass       float64
	angle      float64
	arrow      *Arrow
}

func NewCueStick(target *Ball) *CueStick {
	obj := resolv.NewObject(_boardWidth/2, _boardHeight/2+265, _stickWidth, _stickWidth*2, "cue")
	obj.SetShape(resolv.NewRectangle(0, 0, .1, _stickWidth*2))
	// obj.SetShape(resolv.NewCircle(0, 0, 1))
	stick := &CueStick{
		Image:  ebiten.NewImage(_stickWidth, _stickLength),
		Object: obj,

		targetBall: target,
		speed:      30.,
	}
	stick.arrow = NewArrow(stick)

	return stick
}

func (c *CueStick) Move() {
	if DEBUG {
		b := c.targetBall
		if b == nil {
			return
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			b.velocity.X -= 0.1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			b.velocity.X += 0.1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			b.velocity.Y -= 0.1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			b.velocity.Y += 0.1
		}
	}

	a := -0.001
	aa := a * 20
	if c.targetBall == nil {
		return
	}
	if inpututil.KeyPressDuration(ebiten.KeyA) > 0 && c.X > -2*_stickWidth && c.powerLevel == 0 {
		if inpututil.KeyPressDuration(ebiten.KeyShift) > 0 {
			c.rotateByTarget(a)
			c.arrow.rotateByTarget(a)
			// c.addPos(-v, 0)
		} else {
			c.rotateByTarget(aa)
			c.arrow.rotateByTarget(aa)
			// c.addPos(-20*v, 0)
		}
	} else if inpututil.KeyPressDuration(ebiten.KeyD) > 0 && c.X < _boardWidth+2*_stickWidth && c.powerLevel == 0 {
		if inpututil.KeyPressDuration(ebiten.KeyShift) > 0 {
			// c.addPos(v, 0)
			c.rotateByTarget(-a)
			c.arrow.rotateByTarget(-a)
		} else {
			c.rotateByTarget(-aa)
			c.arrow.rotateByTarget(-aa)
			// c.addPos(20*v, 0)
		}
	} else if inpututil.KeyPressDuration(ebiten.KeyQ) > 0 {
		if inpututil.KeyPressDuration(ebiten.KeyShift) > 0 {
			c.moveByTarget(a)
			c.arrow.moveByTarget(a)
		} else {
			c.moveByTarget(aa)
			c.arrow.moveByTarget(aa)

		}
	} else if inpututil.KeyPressDuration(ebiten.KeyE) > 0 {
		if inpututil.KeyPressDuration(ebiten.KeyShift) > 0 {
			c.moveByTarget(-a)
			c.arrow.moveByTarget(-a)
		} else {
			c.moveByTarget(-aa)
			c.arrow.moveByTarget(-aa)

		}
	} else if dt := inpututil.KeyPressDuration(ebiten.KeySpace); dt > 0 {
		// fmt.Println("Pressed")
		if dt > 20 {
			c.powerLevel = 1
			switch {
			case dt > 100:
				c.powerLevel = 5
			case dt > 80:
				c.powerLevel = 4
			case dt > 60:
				c.powerLevel = 3
			case dt > 40:
				c.powerLevel = 2
			}
		}

	} else if inpututil.IsKeyJustReleased(ebiten.KeySpace) && c.powerLevel > 0 {
		// fmt.Println("c.PowerLevel:", c.PowerLevel)
		switch c.powerLevel {
		case 5:
			c.mass = .6
		case 4:
			c.mass = .4
		case 3:
			c.mass = .25
		case 2:
			c.mass = .12
		case 1:
			c.mass = .05
		}
		u := Vec{0, -1}.Rotate(c.angle).Scale(c.speed)

		// c.addPos(u.X, u.Y)
		c.X += u.X
		c.Y += u.Y

		c.powerLevel = 0
	}
	c.Update()
}

func (c *CueStick) Collide(ball Object) {
	if c.mass == 0 {
		return
	}
	if other, ok := ball.(*Ball); ok {
		var (
			v1       = Vec{0, -1}.Rotate(c.angle).Scale(c.speed)
			vx1, vy1 = v1.X, v1.Y
			vx2, vy2 = other.velocity.X, other.velocity.Y
		)
		// newB1, newB2 := c.Clone(), other.Clone()
		// newB1.X, newB1.Y = c.X+vx1, c.Y+vy1
		// newB2.X, newB2.Y = other.X+vx2, other.Y+vy2

		// fmt.Println(v1.Normalize().X, v1.Normalize().Y)
		if intersection := c.Clone().Shape.Intersection(1, 1, other.Clone().Shape); intersection != nil {
			var (
				p1, p2   = Vec{c.X, c.Y}, Vec{other.X, other.Y}
				v1, v2   = Vec{vx1, vy1}, Vec{vx2, vy2}
				rp       = p2.Add(p1.Negate()).Normalize()
				v1p, v2p = v1.Project(rp), v2.Project(rp)
				v2n      = v2.Add(v2p.Negate())
				m1, m2   = c.mass, other.mass
			)
			// v1 = v1p.Scale((m1 - m2) / (m1 + m2)).Add(v2p.Scale(2 * m2 / (m1 + m2))).Add(v1n)
			v2 = v2p.Scale((m2 - m1) / (m1 + m2)).Add(v1p.Scale(2 * m1 / (m1 + m2))).Add(v2n)
			// c.Velocity.X, c.Velocity.Y = v1.X, v1.Y
			other.velocity.X, other.velocity.Y = v2.X, v2.Y
			fmt.Println(v2.X, v2.Y)
			c.mass = 0
			c.targetBall = nil
		}

	}
}

func (c *CueStick) moveByTarget(angle float64) {
	// r := Vec{c.TargetBall.X, c.TargetBall.Y}.Sub(Vec{c.X, c.Y})
	// dr := r.Sub(r.Rotate(angle))
	// c.X, c.Y = c.X+dr.X, c.Y+dr.Y
	// c.Update()
}
func (c *CueStick) rotateByTarget(angle float64) {
	r := Vec{c.targetBall.X, c.targetBall.Y}.Sub(Vec{c.X, c.Y})
	dr := r.Sub(r.Rotate(angle))
	c.X, c.Y = c.X+dr.X, c.Y+dr.Y
	c.rotate(angle)
	c.Update()
}

func (c *CueStick) angleToPos() {
	if c.targetBall == nil {
		return
	}
	c.X = c.targetBall.X + _distanceToTarget*math.Cos(c.angle+math.Pi/2)
	c.Y = c.targetBall.Y + _distanceToTarget*math.Sin(c.angle+math.Pi/2)
	c.Update()
}

func (c *CueStick) rotate(angle float64) {
	c.angle += angle
}

func (c *CueStick) draw() {
	vector.DrawFilledRect(c.Image, 0, _stickWidth+1, _stickWidth, _stickLength, color.RGBA{185, 135, 64, 255}, true)
	vector.DrawFilledRect(c.Image, 2, 0, _stickWidth-4, _stickWidth+1, color.RGBA{240, 240, 240, 255}, true)
	vector.DrawFilledRect(c.Image, -3, 2*_stickLength/3, _stickWidth+6, _stickLength, color.RGBA{82, 52, 31, 255}, true)
}
