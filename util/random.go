package util

import (
	"math/rand"
	"strings"
)

// generates a random number between min and max including both
func GetRandomInt(min int, max int) int64 {
	return int64(min) + rand.Int63n(int64(max-min+1))
}

var alphabets = "abcdefghijklmnopqrst"

func GetRandomString(l int) string {
	var strbul strings.Builder
	for i := 0; i < l; i++ {
		rnum := GetRandomInt(0, len(alphabets)-1)
		rchar := alphabets[rnum]
		strbul.WriteByte(rchar)
	}
	return strbul.String()
}
