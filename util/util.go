package util

import (
	"crypto/rand"
	m "math/rand"
)

// GenerateRandomBytes returns a random byte array of the specified size
func GenerateRandomBytes(size int) []byte {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

// RandRange returns a random integer between min and max
func RandRange(min, max int) int {
	return m.Intn(max-min) + min
}
