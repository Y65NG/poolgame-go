package util

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 840

	_powerBarLength = 100
	_powerBarHeight = 10
)

const DEBUG = true

type GameState int

const (
	StatePlaying GameState = iota
	StateGameOver
)

type Game struct {
	Board        *Board
	Players      [2]*Player
	CurPlayerIdx int
	Station      *Station
}

func NewGame() *Game {
	var (
		station = NewStation()
		players = [2]*Player{
			NewPlayer("Player 1", &station),
			NewPlayer("Player 2", &station),
		}
	)
	return &Game{
		Board:   NewBoard(&station),
		Players: players,
		Station: &station,
	}
}

func (g *Game) Loop() {

	for {
		select {
		case player := <-g.Station.ChanGameOver:
			g.Station.GameState = StateGameOver
			g.Station.Winner = player

		case <-g.Station.ChanFoul:
			if g.Station.GameState == StateGameOver {
				continue
			}
			g.TurnToNextPlayer()
			g.Station.FreeMode = true

		case ball := <-g.Station.ChanBallIn:
			if g.Station.GameState == StateGameOver {
				continue
			}
			if ball.kind == _kindWhite {
				g.Station.ChanFoul <- struct{}{}
				continue
			}
			if ball.kind == _kindBlack {
				if g.CurPlayer().Score == 7 {
					g.Station.ChanGameOver <- g.CurPlayer()
				} else {
					g.Station.ChanGameOver <- g.NextPlayer()
				}
				continue
			}
			g.Station.BallIn = true
			if g.CurPlayer().BallKind == _kindWhite {
				g.CurPlayer().BallKind = ball.kind
				g.Station.FirstCollidedBall = ball
				if ball.kind == _kindSolid {
					g.NextPlayer().BallKind = _kindStripe
				} else if ball.kind == _kindStripe {
					g.NextPlayer().BallKind = _kindSolid
				}
			}
			if ball.kind == g.CurPlayer().BallKind {
				g.CurPlayer().Score++
			} else {
				g.NextPlayer().Score++
				g.Station.ChanFoul <- struct{}{}
			}

		case <-g.Station.ChanBallsStatic:
			if g.Station.GameState == StateGameOver {
				continue
			}
			if g.Station.Shot {
				g.Station.Shot = false

				// check if the white ball collide with other balls
				if ball := g.Station.FirstCollidedBall; ball != nil {
					// check if the collided ball isn't the same kind as the current player's
					if ball.kind != _kindBlack && g.CurPlayer().BallKind != _kindWhite && ball.kind != g.CurPlayer().BallKind {
						g.Station.FirstCollidedBall = nil

						g.Station.ChanFoul <- struct{}{}
					} else if !g.Station.BallIn {
						g.TurnToNextPlayer()
					}

				} else {
					g.Station.ChanFoul <- struct{}{}
				}

			}
			g.Station.FirstCollidedBall = nil
			g.Station.BallIn = false
			if !g.Board.stick.selected && !g.Board.stick.targetBall.catched {
				g.Board.stick.selected = true
				g.Board.stick.angleToPos()
				g.Board.stick.arrow.angleToPos()
			}
		}
	}
}

func (g *Game) CurPlayer() *Player {
	return g.Players[g.CurPlayerIdx]
}

func (g *Game) NextPlayer() *Player {
	idx := g.CurPlayerIdx + 1
	if idx == len(g.Players) {
		idx = 0
	}
	return g.Players[idx]
}

func (g *Game) TurnToNextPlayer() {
	g.CurPlayerIdx++
	if g.CurPlayerIdx == len(g.Players) {
		g.CurPlayerIdx = 0
	}
	g.Station.Reset()
}

func (g *Game) Reset() {
	g.Board.Reset()
	g.Station.Reset()
	for i, player := range g.Players {
		g.Players[i] = NewPlayer(player.Name, g.Station)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	switch g.Station.GameState {
	case StatePlaying:
		g.Board.Update()
		g.Board.Image.Clear()
	case StateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			log.Println("Restart")
			g.Reset()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	// draw board
	op.GeoM.Reset()
	op.GeoM.Translate(float64(ScreenWidth-_boardWidth)/2, float64(ScreenHeight-_boardHeight)/2)
	g.Board.Draw(screen)
	screen.DrawImage(g.Board.Image, op)

	// draw baskets
	for _, basket := range g.Board.baskets {
		op.GeoM.Reset()
		op.GeoM.Translate(basket.X+(ScreenWidth-_boardWidth)/2-_basketRadius, basket.Y+(ScreenHeight-_boardHeight)/2-_basketRadius)
		basket.draw()
		screen.DrawImage(basket.Image, op)
	}

	// draw balls
	for _, ball := range g.Board.balls {
		if ball.catched {
			continue
		}
		op.GeoM.Reset()
		op.GeoM.Translate((ScreenWidth-_boardWidth)/2+ball.X-_ballRadius, (ScreenHeight-_boardHeight)/2+ball.Y-_ballRadius)
		ball.draw(false)
		screen.DrawImage(ball.Image, op)
		if ball == g.Board.stick.targetBall && g.Board.stick.selected {
			vector.StrokeRect(screen, (ScreenWidth-_boardWidth)/2+float32(ball.X-_ballRadius), (ScreenHeight-_boardHeight)/2+float32(ball.Y-_ballRadius), float32(2*_ballRadius), float32(2*_ballRadius), 1, color.White, true)
		}
	}

	// draw free white ball
	if g.Station.FreeMode {
		whiteBall := g.Board.balls[0]
		op.GeoM.Reset()
		cx, cy := ebiten.CursorPosition()
		op.GeoM.Translate(float64(cx)-_ballRadius, float64(cy)-_ballRadius)

		whiteBall.draw(true)
		screen.DrawImage(whiteBall.Image, op)
	}

	// draw cue stick
	op.GeoM.Reset()
	op.GeoM.Translate(-_stickWidth/2, 0)
	op.GeoM.Rotate(g.Board.stick.angle)
	op.GeoM.Translate(_stickWidth/2, 0)

	op.GeoM.Translate(g.Board.stick.X+(ScreenWidth-_boardWidth)/2-_stickWidth/2, g.Board.stick.Y+(ScreenHeight-_boardHeight)/2)
	g.Board.stick.draw()
	screen.DrawImage(g.Board.stick.Image, op)

	// draw arrow
	op.GeoM.Reset()
	op.GeoM.Translate(-_arrowWidth/2, 0)
	op.GeoM.Rotate(g.Board.stick.arrow.angle)
	op.GeoM.Translate(_arrowWidth/2, 0)

	op.GeoM.Translate(g.Board.stick.X+(ScreenWidth-_boardWidth)/2-_arrowWidth/2, g.Board.stick.Y+(ScreenHeight-_boardHeight)/2)
	g.Board.stick.arrow.draw()
	screen.DrawImage(g.Board.stick.arrow.Image, op)
	// draw power bar
	vector.StrokeRect(screen, ScreenWidth-_powerBarLength-10, ScreenHeight-_powerBarHeight-10, _powerBarLength, _powerBarHeight, 1, color.White, true)
	vector.DrawFilledRect(screen, ScreenWidth-_powerBarLength-10, ScreenHeight-_powerBarHeight-10, float32(g.Board.stick.powerLevel)*.2*_powerBarLength, _powerBarHeight, color.White, true)

	// draw player info
	infoStr := g.Players[0].Name + ": " + strconv.Itoa(g.Players[0].Score)
	if g.CurPlayer() == g.Players[0] {
		infoStr += " <-"
	}
	var kindStr string
	if g.Players[0].BallKind == _kindWhite {
		kindStr = ""
	} else {
		kindStr = "<" + g.Players[0].BallKind.String() + ">"
	}
	ebitenutil.DebugPrintAt(screen, infoStr, 10, 10)
	ebitenutil.DebugPrintAt(screen, kindStr, 10, 30)

	if g.Players[1].BallKind == _kindWhite {
		kindStr = ""
	} else {
		kindStr = "<" + g.Players[1].BallKind.String() + ">"
	}
	infoStr = g.Players[1].Name + ": " + strconv.Itoa(g.Players[1].Score)
	if g.CurPlayer() == g.Players[1] {
		infoStr += " <-"
	}
	ebitenutil.DebugPrintAt(screen, infoStr, ScreenWidth-10-100, 10)
	ebitenutil.DebugPrintAt(screen, kindStr, ScreenWidth-10-100, 30)

	if g.Station.FreeMode {
		ebitenutil.DebugPrintAt(screen, "Click to place the white ball", ScreenWidth/2.5, 10)
	}

	if g.Station.GameState == StateGameOver {
		ebitenutil.DebugPrintAt(screen, "Game Over", ScreenWidth/2.5, 10)
		ebitenutil.DebugPrintAt(screen, "Winner: "+g.Station.Winner.Name, ScreenWidth/2.5, 25)
		ebitenutil.DebugPrintAt(screen, "Press SPACE to restart", ScreenWidth/2.5, 40)

	}
}
