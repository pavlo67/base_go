package strlib

import (
	"math/rand"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(strlen int) string {
	if strlen < 1 {
		return ""
	}

	result := make([]byte, strlen)

	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}
