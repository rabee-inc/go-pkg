package randutil

import (
	"math/rand"
	"time"
)

func seed() {
	rand.Seed(time.Now().UnixNano())
}

// Bool ... 指定確率でbool値を生成する
func Bool(rate float32) bool {
	seed()
	return rand.Float32()*100 < rate
}

// Int ... 指定範囲の乱数を生成する
func Int(min int, max int) int {
	seed()
	return rand.Intn((max+1)-min) + min
}
