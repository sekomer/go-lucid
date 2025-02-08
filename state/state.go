package state

import (
	"time"

	"github.com/holiman/uint256"
)

type State struct {
	// Current hash target for the blockchain
	nBits uint32

	// Current mining difficulty for the blockchain
	difficulty int

	// Mining reward amount for successfully mining a block
	miningReward float64

	// Maximum number of transactions allowed in a block
	maxTransactions int

	// Current block height of the chain
	currentBlockHeight int

	// Version of the blockchain
	version int

	// Current timestamp of the chain
	currentTimestamp int64
}

var GlobalState = &State{
	nBits:              0x210000ff,        // Initial nbits
	difficulty:         1,                 // Initial difficulty
	miningReward:       100.0,             // Initial mining reward
	maxTransactions:    1000,              // Maximum 1000 transactions per block
	currentBlockHeight: 0,                 // Start at height 0
	version:            0,                 // Initial version
	currentTimestamp:   time.Now().Unix(), // Current timestamp
}

func (s *State) GetNBits() uint32 {
	return s.nBits
}

func (s *State) GetTarget() *uint256.Int {
	return NBitsToTarget(s.nBits)
}

func (s *State) GetDifficulty() uint32 {
	return uint32(s.difficulty)
}

func (s *State) GetMiningReward() float64 {
	return s.miningReward
}

func (s *State) GetMaxTransactions() int {
	return s.maxTransactions
}

func (s *State) GetCurrentBlockHeight() uint32 {
	return uint32(s.currentBlockHeight)
}

func (s *State) GetVersion() int32 {
	return int32(s.version)
}

func (s *State) GetCurrentTimestamp() int64 {
	return s.currentTimestamp
}

func GetState() *State {
	return GlobalState
}
