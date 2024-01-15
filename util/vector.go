package util

import (
	"math"
)

var (
	VecI = Vec{
		X: 1,
		Y: 0,
	}
	VecJ = Vec{
		X: 0,
		Y: 1,
	}
)

type Vec struct {
	X, Y float64
}

func (v Vec) Value() (float64, float64) {
	return v.X, v.Y
}

func (v Vec) Add(other Vec) Vec {
	return Vec{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vec) Sub(other Vec) Vec {
	return Vec{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vec) Dot(other Vec) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec) Cross(other Vec) Vec {
	return Vec{
		X: v.X*other.Y - v.Y*other.X,
		Y: v.Y*other.X - v.X*other.Y,
	}
}

func (v Vec) Scale(s float64) Vec {
	return Vec{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vec) Negate() Vec {
	return Vec{
		X: -v.X,
		Y: -v.Y,
	}
}

func (v Vec) Rotate(angle float64) Vec {
	return Vec{
		X: v.X*math.Cos(angle) - v.Y*math.Sin(angle),
		Y: v.X*math.Sin(angle) + v.Y*math.Cos(angle),
	}
}

func (v Vec) Normalize() Vec {
	if v.Size() == 0 {
		return Vec{}
	}
	return v.Scale(1 / v.Size())
}

func (v Vec) Project(other Vec) Vec {
	if other.Size() == 0 {
		return Vec{}
	}
	return other.Scale(v.Dot(other) / other.Dot(other))
}

func (v Vec) Size() float64 {
	vSquare := v.Dot(v)
	return math.Sqrt(vSquare)
}
