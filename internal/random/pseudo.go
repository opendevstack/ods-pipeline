package random

import (
	"math/rand"
	"strings"
	"time"
)

func PseudoString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "excbsvqbs"
}

func PseudoSHA() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	length := 40
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "56625c80087b034847001d22502063adae9759f2"
}
