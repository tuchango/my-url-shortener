package random

import (
	"math/rand/v2"
)

var (
	chars    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	charsLen = len(chars)
)

func NewRandomString(len int) string {
	b := make([]byte, len)
	for i := range b {
		b[i] = chars[rand.IntN(charsLen)]
	}

	return string(b)
}
