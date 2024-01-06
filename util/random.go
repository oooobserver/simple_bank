package util

import (
	"math/rand"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

const letters = "abcdefghijklmnopqrstuvwxyz"

var currs = [...]string{"EUR", "USD", "RMB"}

func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	l := len(letters)

	for i := 0; i < n; i++ {
		char := letters[rng.Intn(l)]
		sb.WriteByte(char)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	l := len(currs)
	return currs[rng.Intn(l)]
}

func RandomMoney() int64 {
	return RandomInt(0, 100000)
}
