package core

import (
	"context"
	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/state"
)

// Miner defines the interface that all miners must implement
type MinerInterface interface {
	// Mine attempts to mine a new block with the given transactions and previous block
	Mine(transactions []transaction.RawTransaction, previousBlock *block.Block, state *state.State, ctx context.Context) (*block.Block, error)

	// Resources returns the resources used by the miner
	GetResources() *Resources
}

type Resources struct {
	CPUUsage    *int
	MemoryUsage *int
	DiskUsage   *int
	HashRate    *int
}
