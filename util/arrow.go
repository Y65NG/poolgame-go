package util

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	_arrowWidth      = 1
	_arrowBaseLength = 100
)

type Arrow struct {
	*ebiten.Image
	*resolv.Object

	stick *CueStick
	angle float64
}

func NewArrow(stick *CueStick) *Arrow {
	obj := resolv.NewObject(_boardWidth/2, _boardHeight/2+265, _arrowWidth, _arrowWidth*2, "cue")
	obj.SetShape(resolv.NewRectangle(0, 0, .1, _arrowWidth*2))
	// obj.SetShape(resolv.NewCircle(0, 0, 1))
	return &Arrow{
		Image:  ebiten.NewImage(_arrowWidth, _arrowBaseLength*2),
		Object: obj,

		stick: stick,
		angle: math.Pi,
	}
}

func (c *Arrow) Length() float64 {
	return _arrowBaseLength + _arrowBaseLength/10*float64(c.stick.powerLevel)
}

func (c *Arrow) Vertices() []ebiten.Vertex {
	return []ebiten.Vertex{
		{
			DstX: float32(c.X - _arrowWidth/2),
			DstY: float32(c.Y - _arrowWidth/2),
		},
		{
			DstX: float32(c.X + _arrowWidth/2),
			DstY: float32(c.Y - _arrowWidth/2),
		},
		{
			DstX: float32(c.X - _arrowWidth/2),
			DstY: float32(c.Y + _arrowWidth/2),
		},
	}
}

func (c *Arrow) moveByTarget(angle float64) {
	// r := Vec{c.TargetBall.X, c.TargetBall.Y}.Sub(Vec{c.X, c.Y})
	// dr := r.Sub(r.Rotate(angle))
	// c.X, c.Y = c.X+dr.X, c.Y+dr.Y
	// c.Update()
}
func (c *Arrow) rotateByTarget(angle float64) {
	r := Vec{c.stick.targetBall.X, c.stick.targetBall.Y}.Sub(Vec{c.X, c.Y})
	dr := r.Sub(r.Rotate(angle))
	c.X, c.Y = c.X+dr.X, c.Y+dr.Y
	c.rotate(angle)
	c.Update()
}

func (c *Arrow) angleToPos() {
	if c.stick.targetBall == nil {
		return
	}
	c.X = c.stick.targetBall.X + _distanceToTarget*math.Cos(c.angle+math.Pi/2)
	c.Y = c.stick.targetBall.Y + _distanceToTarget*math.Sin(c.angle+math.Pi/2)
	c.Update()
}

func (c *Arrow) rotate(angle float64) {
	c.angle += angle
}

func (c *Arrow) draw() {
	if c.stick.targetBall == nil {
		c.Clear()
		return
	}
	vector.DrawFilledRect(c.Image, 0, 2*float32(c.Length())/5, _arrowWidth, float32(c.Length()), color.RGBA{100, 100, 100, 255}, true)

	// vector.DrawFilledRect(c.Image, 2, 0, _arrowWidth-4, _arrowWidth+1, color.RGBA{240, 240, 240, 255}, true)
}
