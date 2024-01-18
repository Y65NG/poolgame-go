package util

type Player struct {
	Name     string
	Score    int
	BallKind ballKind
	InTurn   bool

	station *Station
}

func NewPlayer(name string, station *Station) *Player {
	return &Player{
		Name:     name,
		Score:    0,
		BallKind: _kindWhite,

		station: station,
	}
}
