package miner

import (
	"go-lucid/miner/core"
	"go-lucid/miner/simple"
)

var (
	globalMiner core.MinerInterface = simple.NewSimpleMiner()
)

func SetGlobalMiner(miner core.MinerInterface) {
	globalMiner = miner
}

func GetGlobalMiner() core.MinerInterface {
	return globalMiner
}
