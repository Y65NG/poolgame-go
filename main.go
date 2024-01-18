package main

import (
	"poolball/util"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(util.ScreenWidth, util.ScreenHeight)
	ebiten.SetWindowTitle("Pool Ball")
	g := util.NewGame()
	go g.Loop()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}

}
