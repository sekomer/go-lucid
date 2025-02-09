package hexutil

import (
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"strings"
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
	return mrand.Intn(max-min) + min
}

func HashToBytes(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return decoded
}

func PadOrTrimTo32Bytes(input []byte) []byte {
	if len(input) >= 32 {
		return input[:32]
	}
	result := make([]byte, 32)
	copy(result, input)
	return result
}

func TrimLeadingZeroes(input string) string {
	// Handle case where input is all zeroes
	for i := 0; i < len(input); i++ {
		if input[i] != '0' {
			return input[i:]
		}
	}
	return ""
}

func AddPrefix(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input
	}
	return "0x" + input
}

func Encode(input []byte, options ...func(string) string) string {
	encoded := hex.EncodeToString(input)
	for _, option := range options {
		encoded = option(encoded)
	}
	return encoded
}

func EncodeToPrefixedHex(input []byte) string {
	encoded := Encode(input, TrimLeadingZeroes, AddPrefix)
	return encoded
}

func DecodeFromPrefixedHex(input string) ([]byte, error) {
	input = strings.TrimPrefix(input, "0x")
	decoded, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
