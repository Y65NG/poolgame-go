package util

import (
	"fmt"
	"math"
	"testing"
)

func TestRotate(t *testing.T) {
	v := Vec{X: 1, Y: 0}
	v = v.Rotate(math.Pi)
	fmt.Println("v", v)
}

func TestScale(t *testing.T) {
	v := Vec{X: 1, Y: 1}
	v = v.Normalize()
	fmt.Println("v", v)

}
