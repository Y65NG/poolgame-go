package util

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	_basketRadius = _ballRadius * 1.5
)

type basketID int

const (
	_idBasketTopLeft basketID = iota + 1
	_idBasketMidLeft
	_idBasketBottomLeft
	_idBasketTopRight
	_idBasketMidRight
	_idBasketBottomRight
)

type Basket struct {
	*ebiten.Image
	*resolv.Object

	id    basketID
	balls []*Ball

	station *Station
}

func newBasket(id basketID, station *Station) *Basket {
	w := _basketRadius
	var x, y float64
	switch id {
	case _idBasketTopLeft:
		x, y = (1-math.Sqrt2/2)*_basketRadius, (1-math.Sqrt2/2)*_basketRadius
		// x, y = 0, 0
	case _idBasketMidLeft:
		x, y = -(_basketRadius - _edgeWidth), _boardHeight/2
	case _idBasketBottomLeft:
		x, y = (1-math.Sqrt2/2)*_basketRadius, _boardHeight-(1-math.Sqrt2/2)*_basketRadius
	case _idBasketTopRight:
		x, y = _boardWidth-(1-math.Sqrt2/2)*_basketRadius, (1-math.Sqrt2/2)*_basketRadius
	case _idBasketMidRight:
		x, y = _boardWidth+(_basketRadius-_edgeWidth), _boardHeight/2
	case _idBasketBottomRight:
		x, y = _boardWidth-(1-math.Sqrt2/2)*_basketRadius, _boardHeight-(1-math.Sqrt2/2)*_basketRadius
	}
	b := &Basket{
		Image:  ebiten.NewImage(int(w*2)+2, int(w*2)+2),
		Object: resolv.NewObject(x, y, _basketRadius, _basketRadius, "basket"),
		id:     id,

		station: station,
	}
	b.Object.SetShape(resolv.NewCircle(x, y, _basketRadius))
	return b
}

func (b *Basket) catchBall(ball *Ball) {
	d := Vec{ball.X - b.X, ball.Y - b.Y}.Size()
	nextD := Vec{ball.X + ball.velocity.X - b.X, ball.Y + ball.velocity.Y - b.Y}.Size()
	if (d <= (_basketRadius-_ballRadius/2.5) || nextD <= (_basketRadius-_ballRadius/2.5)) && !ball.catched {

		ball.catched = true
		b.station.ChanBallIn <- ball
		b.balls = append(b.balls, ball)
		ball.velocity = Vec{0, 0}
		ball.X, ball.Y = b.X, b.Y
	}
}

func (b *Basket) draw() {
	vector.StrokeCircle(b.Image, _basketRadius+1, _basketRadius+1, float32(_basketRadius), 2, color.White, false)
	vector.DrawFilledCircle(b.Image, _basketRadius+1, _basketRadius+1, float32(_basketRadius)-1, color.Black, false)
}
