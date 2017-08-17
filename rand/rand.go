package rand

import (
	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"
)

func Seed() {
	var buf [8]byte
	_, err := crand.Read(buf[:])
	if err != nil {
		panic(err.Error())
	}

	seedValue := binary.BigEndian.Uint64(buf[:])
	mrand.Seed(int64(seedValue))
}

func Between(min, max int64) int64 {
	r := max - min
	n := mrand.Int63n(r)
	return n + min
}

func Coin(p float64) bool {
	c := mrand.Float64()
	return p <= c
}
