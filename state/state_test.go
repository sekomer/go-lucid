package state_test

import (
	"encoding/hex"
	"go-lucid/state"
	"go-lucid/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialState(t *testing.T) {
	assert.Equal(t, state.GetState().GetDifficulty(), uint32(1))
	assert.Equal(t, state.GetState().GetMiningReward(), 100.0)
	assert.Equal(t, state.GetState().GetMaxTransactions(), 1000)
	assert.Equal(t, state.GetState().GetCurrentBlockHeight(), uint32(0))
	assert.Equal(t, state.GetState().GetVersion(), int32(0))
	assert.True(t, state.GetState().GetCurrentTimestamp() > 0)
}

func TestGetTarget(t *testing.T) {
	assert.Equal(t, hex.EncodeToString(state.GetState().GetTarget().Bytes()), util.TrimLeadingZeroes("00ff000000000000000000000000000000000000000000000000000000000000"))
}

func TestSatoshiNBitsToTarget(t *testing.T) {
	target := state.NBitsToTarget(0x1d00ffff)
	assert.Equal(t, hex.EncodeToString(target.Bytes()), util.TrimLeadingZeroes("00000000ffff0000000000000000000000000000000000000000000000000000"))
}
