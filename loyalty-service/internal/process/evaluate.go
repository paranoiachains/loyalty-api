package process

import (
	"math/rand/v2"
	"time"
)

func Evaluate() float64 {
	time.Sleep(1 * time.Second)
	return rand.Float64() * 500
}
