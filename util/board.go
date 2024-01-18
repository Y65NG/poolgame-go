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

const _friction = 0.013

type Board struct {
	*ebiten.Image
	*resolv.Space

	balls   []*Ball
	stick   *CueStick
	baskets []*Basket
	edges   []*Edge

	station *Station
}

func NewBoard(station *Station) *Board {
	balls := make([]*Ball, 16)
	balls[0] = NewWhiteball(balls, station)

	// Positions of balls relative to the black ball from bottom right to top left.
	relativePoses := [][2]float64{

		{3 * _ballRadius, -2 * _ballRadius},
		{1 * _ballRadius, 2 * _ballRadius},
		{-2 * _ballRadius, -4 * _ballRadius},
		{2 * _ballRadius, -4 * _ballRadius},
		{1 * _ballRadius, -2 * _ballRadius},
		{-3 * _ballRadius, -2 * _ballRadius},
		{-1 * _ballRadius, 2 * _ballRadius},
		{0 * _ballRadius, 0 * _ballRadius}, // the black ball

		{0 * _ballRadius, 4 * _ballRadius},
		{-2 * _ballRadius, 0 * _ballRadius},
		{-4 * _ballRadius, -4 * _ballRadius},

		{-1 * _ballRadius, -2 * _ballRadius},

		{2 * _ballRadius, 0 * _ballRadius},
		{0 * _ballRadius, -4 * _ballRadius},

		{4 * _ballRadius, -4 * _ballRadius},
	}

	for i := 1; i <= len(relativePoses); i++ {
		balls[i] = NewBall(
			balls,
			ballID(i),
			_boardWidth/2+relativePoses[i-1][0]*1.3,
			_boardHeight/4+relativePoses[i-1][1]*1.3,
			station,
		)
	}

	baskets := make([]*Basket, 6)
	for id := _idBasketTopLeft; id <= _idBasketBottomRight; id++ {
		baskets[id-1] = newBasket(id, station)
	}

	edges := make([]*Edge, 6)
	for pos := _edgeTop; pos <= _edgeTopRight; pos++ {
		edges[pos-1] = NewEdge(pos, station)
	}
	b := &Board{
		Image: ebiten.NewImage(_boardWidth, _boardHeight),
		Space: resolv.NewSpace(_boardWidth*5, _boardHeight*5, _ballRadius*10, _ballRadius*10),

		balls:   balls,
		stick:   NewCueStick(balls[0], station),
		baskets: baskets,
		edges:   edges,

		station: station,
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
	vector.DrawFilledRect(b.Image, 0, 0, float32(_boardWidth), float32(_boardHeight), _boardColor, false)

	for _, edge := range b.edges {
		edge.draw()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(edge.X, edge.Y)
		b.DrawImage(edge.Image, op)

	}
	vector.StrokeRect(b.Image, .1, .1, float32(_boardWidth)-.2, float32(_boardHeight)-.2, 1, color.White, false)

}

func (b *Board) Update() {
	b.stick.Move(b.balls)
	for i := 0; i < len(b.balls); i++ {
		b.balls[i].Move()
		b.stick.Collide(b.balls[i])
		for _, edge := range b.edges {
			edge.Collide(b.balls[i])
		}

		for _, basket := range b.baskets {
			basket.catchBall(b.balls[i])
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !b.station.FreeMode {
			x, y := ebiten.CursorPosition()
			x -= (ScreenWidth - _boardWidth) / 2
			y -= (ScreenHeight - _boardHeight) / 2
			if b.balls[i].containsPos(float64(x), float64(y)) {
				// fmt.Println(b.balls[i].velocity, b.balls[i].velocity.Size())
				if b.balls[i].velocity.Size() <= .005 {
					b.stick.targetBall = b.balls[i]
					// b.stick.selected = true
					b.stick.angleToPos()
					b.stick.arrow.angleToPos()
				}
			}
		}
	}
	if b.BallsStatic() {
		b.station.ChanBallsStatic <- struct{}{}
	}

}

func (b *Board) Reset() {
	b = NewBoard(b.station)
}

func (b *Board) BallsStatic() bool {
	for _, ball := range b.balls {
		if ball.velocity.Size() > 0 {
			return false
		}
	}
	return true
}
