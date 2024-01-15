package util

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 600
	ScreenHeight = 840

	_powerBarLength = 100
	_powerBarHeight = 10
)

const DEBUG = true

type Game struct {
	Board *Board
}

func NewGame() *Game {
	return &Game{
		Board: NewBoard(),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	g.Board.Update()
	g.Board.Image.Clear()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw baskets
	op := &ebiten.DrawImageOptions{}
	for _, basket := range g.Board.baskets {
		op.GeoM.Reset()
		op.GeoM.Translate(basket.X+(ScreenWidth-_boardWidth)/2-_basketRadius, basket.Y+(ScreenHeight-_boardHeight)/2-_basketRadius)
		basket.draw()
		screen.DrawImage(basket.Image, op)
	}

	// draw board
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(ScreenWidth-_boardWidth)/2, float64(ScreenHeight-_boardHeight)/2)
	g.Board.Draw(screen)
	screen.DrawImage(g.Board.Image, op)

	// draw cue stick
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-_stickWidth/2, 0)
	op.GeoM.Rotate(g.Board.stick.angle)
	op.GeoM.Translate(_stickWidth/2, 0)

	op.GeoM.Translate(g.Board.stick.X+(ScreenWidth-_boardWidth)/2-_stickWidth/2, g.Board.stick.Y+(ScreenHeight-_boardHeight)/2)
	g.Board.stick.draw()
	screen.DrawImage(g.Board.stick.Image, op)

	// draw arrow
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-_arrowWidth/2, 0)
	op.GeoM.Rotate(g.Board.stick.arrow.angle)
	op.GeoM.Translate(_arrowWidth/2, 0)

	op.GeoM.Translate(g.Board.stick.X+(ScreenWidth-_boardWidth)/2-_arrowWidth/2, g.Board.stick.Y+(ScreenHeight-_boardHeight)/2)
	g.Board.stick.arrow.draw()
	screen.DrawImage(g.Board.stick.arrow.Image, op)
	// screen.DrawTriangles(g.Board.stick.arrow.Vertices(), []uint16{0, 1, 2}, g.Board.stick.arrow.Image, nil)
	// draw power bar
	vector.StrokeRect(screen, ScreenWidth-_powerBarLength-10, ScreenHeight-_powerBarHeight-10, _powerBarLength, _powerBarHeight, 1, color.White, true)
	vector.DrawFilledRect(screen, ScreenWidth-_powerBarLength-10, ScreenHeight-_powerBarHeight-10, float32(g.Board.stick.powerLevel)*.2*_powerBarLength, _powerBarHeight, color.White, true)

	// draw edge
	// edge := NewSideEdge(false)

}
