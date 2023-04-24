package utils

import (
	"math/rand"
	"time"
	"unsafe"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func RandomInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Int()
}

func RandomIntRange(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	return int(rand.Int31n(int32(end-start))) + start
}

func RandomInt64() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63()
}

func RandomInt64Range(start int64, end int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(end-start) + start
}
