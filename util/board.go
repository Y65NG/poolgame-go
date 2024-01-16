package util

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	_boardHeight = ScreenHeight - 120
	_boardWidth  = _boardHeight / 2
)

const _friction = 0.012

type Board struct {
	*ebiten.Image
	*resolv.Space

	balls   []*Ball
	stick   *CueStick
	baskets []*Basket
	edges   []*Edge
}

func NewBoard() *Board {
	balls := make([]*Ball, 16)
	balls[0] = NewWhiteball(balls)

	// Positions of balls relative to the black ball from bottom right to top left.
	relativePoses := [][2]float64{
		{0 * _ballRadius, 4 * _ballRadius},
		{-1 * _ballRadius, 2 * _ballRadius},
		{-2 * _ballRadius, 0 * _ballRadius},
		{-3 * _ballRadius, -2 * _ballRadius},
		{-4 * _ballRadius, -4 * _ballRadius},

		{1 * _ballRadius, 2 * _ballRadius},
		{0 * _ballRadius, 0 * _ballRadius}, // the black ball
		{-1 * _ballRadius, -2 * _ballRadius},
		{-2 * _ballRadius, -4 * _ballRadius},

		{2 * _ballRadius, 0 * _ballRadius},
		{1 * _ballRadius, -2 * _ballRadius},
		{0 * _ballRadius, -4 * _ballRadius},

		{3 * _ballRadius, -2 * _ballRadius},
		{2 * _ballRadius, -4 * _ballRadius},

		{4 * _ballRadius, -4 * _ballRadius},
	}

	for i := 1; i <= len(relativePoses); i++ {
		balls[i] = NewBall(
			balls,
			ballID(i),
			_boardWidth/2+relativePoses[i-1][0]*1.5,
			_boardHeight/2-240+relativePoses[i-1][1]*1.5,
		)
	}

	baskets := make([]*Basket, 6)
	for id := _idBasketTopLeft; id <= _idBasketBottomRight; id++ {
		baskets[id-1] = newBasket(id)
	}

	edges := make([]*Edge, 6)
	for pos := _edgeTop; pos <= _edgeTopRight; pos++ {
		edges[pos-1] = NewEdge(pos)
	}
	b := &Board{
		Image: ebiten.NewImage(_boardWidth, _boardHeight),
		Space: resolv.NewSpace(_boardWidth*5, _boardHeight*5, _ballRadius*10, _ballRadius*10),

		balls:   balls,
		stick:   NewCueStick(balls[0]),
		baskets: baskets,
		edges:   edges,
	}
	b.Space.Add(b.stick.Object)
	for _, ball := range b.balls {
		b.Space.Add(ball.Object)
	}
	for _, edge := range b.edges {
		b.Space.Add(edge.Object)
	}
	return b
}

func (b *Board) Draw(screen *ebiten.Image) {
	for _, ball := range b.balls {
		if ball.catched {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(ball.X-_ballRadius, ball.Y-_ballRadius)
		ball.draw()
		b.DrawImage(ball.Image, op)
		// vector.StrokeRect(b.Image, float32(ball.X-_ballRadius), float32(ball.Y-_ballRadius), float32(2*_ballRadius), float32(2*_ballRadius), 1, color.White, true)
		// v1, v2 := b.Stick.Shape.Bounds()
		// vector.StrokeRect(b.Image, float32(v1.X()), float32(v1.Y()), float32(v2.X()-v1.X()), float32(v2.Y()-v1.Y()), 1, color.White, true)
		// vector.StrokeRect(b.Image, b.Stick.Shape.)
		if ball == b.stick.targetBall {
			vector.StrokeRect(b.Image, float32(ball.X-_ballRadius), float32(ball.Y-_ballRadius), float32(2*_ballRadius), float32(2*_ballRadius), 1, color.White, true)
		}
	}

	for _, edge := range b.edges {
		edge.draw()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(edge.X, edge.Y)
		b.DrawImage(edge.Image, op)

		// v1, v2 := edge.Shape.Bounds()
		// // fmt.Println("edge bounds:", v1, v2)
		// vector.StrokeRect(b.Image, float32(v1.X()), float32(v1.Y()), float32(v2.X()-v1.X()), float32(v2.Y()-v1.Y()), 1, color.White, true)
	}
	vector.StrokeRect(b.Image, .1, .1, float32(_boardWidth)-.2, float32(_boardHeight)-.2, 1, color.White, false)

}

func (b *Board) Update() {
	b.stick.Move()
	for i := 0; i < len(b.balls); i++ {
		b.balls[i].Move()
		b.stick.Collide(b.balls[i])
		// b.edges[0].Collide(b.balls[i])
		for _, edge := range b.edges {
			if edge == nil {
				continue
			}
			edge.Collide(b.balls[i])
		}

		for _, basket := range b.baskets {
			basket.catchBall(b.balls[i])
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			x -= (ScreenWidth - _boardWidth) / 2
			y -= (ScreenHeight - _boardHeight) / 2
			if b.balls[i].containsPos(float64(x), float64(y)) {
				// fmt.Println(b.balls[i].velocity, b.balls[i].velocity.Size())
				if b.balls[i].velocity.Size() <= .005 {
					b.stick.targetBall = b.balls[i]
					b.stick.angleToPos()
					b.stick.arrow.angleToPos()
				}
			}
		}
	}
}
