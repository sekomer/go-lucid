package simple

import (
	"go-lucid/miner/core"
)

type SimpleMiner struct {
	Resources *core.Resources
}

// interface validation
var _ core.MinerInterface = &SimpleMiner{}
