package util_test

import (
	"go-lucid/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomBytes32(t *testing.T) {
	t.Parallel()

	bytes := util.GenerateRandomBytes(32)
	assert.Len(t, bytes, 32)
}

func TestGenerateRandomBytes64(t *testing.T) {
	t.Parallel()

	bytes := util.GenerateRandomBytes(64)
	assert.Len(t, bytes, 64)
}

func TestRandRange(t *testing.T) {
	t.Parallel()

	min := 10
	max := 20
	randNum := util.RandRange(min, max)
	assert.GreaterOrEqual(t, randNum, min)
	assert.Less(t, randNum, max)
}

func TestTrimLeadingZeroesWithoutLeadingZeroes(t *testing.T) {
	t.Parallel()

	input := "42"
	trimmed := util.TrimLeadingZeroes(input)
	assert.Equal(t, trimmed, "42")
}

func TestTrimLeadingZeroesWithNonZeroes(t *testing.T) {
	t.Parallel()

	input := "00000000000000000000042"
	trimmed := util.TrimLeadingZeroes(input)
	assert.Equal(t, trimmed, "42")
}

func TestTrimLeadingZeroesAllZeroes(t *testing.T) {
	t.Parallel()

	input := "0000000000000000000000000000000000000000000000000000000000000000"
	trimmed := util.TrimLeadingZeroes(input)
	assert.Equal(t, trimmed, "")
}
