package main

import (
	"math/big"
	"math/rand"
	"time"
)

func GetPseodoRandomString() string {
	// seed
	rand.Seed(time.Now().UnixNano())

	// random float
	r := rand.Float64()

	// to int64
	rInt := big.NewInt(int64(r * 1000000000000000000))

	// to string with radix 36
	s := rInt.Text(36)

	return s
}
