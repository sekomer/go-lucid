package simple_test

import (
	"go-lucid/miner/simple"
	"testing"
)

func TestSimpleMiner(t *testing.T) {
	t.Parallel()

	miner := simple.NewSimpleMiner()
	miner.GetResources()
}
