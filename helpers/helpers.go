package helpers

import (
	"math/rand"
	"time"
)

var RAND_SOURCE *rand.Rand

func InitRandSource() {
	source := rand.NewSource(time.Now().UnixNano())
	RAND_SOURCE = rand.New(source)
}

func GetRandomNumberWithin(ceil int) int {
	return RAND_SOURCE.Intn(ceil) + 1
}
