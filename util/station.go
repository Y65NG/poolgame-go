package util

type Station struct {
	FreeMode          bool
	FirstCollidedBall *Ball
	Shot              bool
	BallIn            bool

	ChanGameOver    chan struct{}
	ChanFoul        chan struct{}
	ChanNextTurn    chan struct{}
	ChanBallIn      chan *Ball
	ChanBallsStatic chan struct{}
}

func NewStation() Station {
	return Station{
		ChanGameOver:    make(chan struct{}, 1),
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
}
