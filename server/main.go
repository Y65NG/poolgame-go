package main

import (
	"log"
	"net/http"
	"poolball/util"

	"github.com/gorilla/websocket"
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

func serveWs()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func upgradeHttp(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return conn, err
	}
	return conn, nil
}
