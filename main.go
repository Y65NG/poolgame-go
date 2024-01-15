package main

import (
	"billiards/util"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(util.ScreenWidth, util.ScreenHeight)
	ebiten.SetWindowTitle("Billiards")
	g := util.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}

}
