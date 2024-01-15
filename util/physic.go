package util

import (
	"math/rand"
	"time"
)

var rd = rand.New(rand.NewSource(time.Now().UnixMicro()))

type Object interface {
	Move()
	Collide(Object)
}

func RandomFactor(mean, std float64) float64 {
	return rd.NormFloat64()*std + mean
}
