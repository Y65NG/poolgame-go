package util

type Station struct {
	FreeMode          bool
	FirstCollidedBall *Ball
	Shot              bool
	BallIn            bool
	GameState         GameState
	Winner            *Player

	ChanGameOver    chan *Player
	ChanFoul        chan struct{}
	ChanNextTurn    chan struct{}
	ChanBallIn      chan *Ball
	ChanBallsStatic chan struct{}
}

func NewStation() Station {
	return Station{
		ChanGameOver:    make(chan *Player, 1),
		ChanFoul:        make(chan struct{}, 1),
		ChanNextTurn:    make(chan struct{}, 1),
		ChanBallIn:      make(chan *Ball, 1),
		ChanBallsStatic: make(chan struct{}, 1),
	}
}

func (s *Station) Reset() {
	s.FreeMode = false
	s.FirstCollidedBall = nil
	s.Shot = false
	s.BallIn = false
	s.GameState = StatePlaying
	s.Winner = nil

}
