package mempool

import (
	"go-lucid/core"
	"sync"
)

var (
	once     sync.Once
	instance *Mempool
)

type Mempool struct {
	sync.RWMutex
	txs map[string]*core.RawTransaction
}

func GetMempool() *Mempool {
	once.Do(func() {
		instance = &Mempool{
			txs: make(map[string]*core.RawTransaction),
		}
	})
	return instance
}

func (m *Mempool) AddTx(tx *core.RawTransaction) error {
	m.Lock()
	defer m.Unlock()

	m.txs[string(tx.Hash)] = tx
	return nil
}

func (m *Mempool) RemoveTx(hash string) {
	m.Lock()
	defer m.Unlock()

	delete(m.txs, hash)
}

func (m *Mempool) GetTx(hash string) *core.RawTransaction {
	m.RLock()
	defer m.RUnlock()

	return m.txs[hash]
}

func (m *Mempool) GetTxs() []*core.RawTransaction {
	m.RLock()
	defer m.RUnlock()

	txs := make([]*core.RawTransaction, 0, len(m.txs))
	for _, tx := range m.txs {
		txs = append(txs, tx)
	}

	return txs
}

func (m *Mempool) Size() int32 {
	m.RLock()
	defer m.RUnlock()

	return int32(len(m.txs))
}

func (m *Mempool) Clear() {
	m.Lock()
	defer m.Unlock()

	m.txs = make(map[string]*core.RawTransaction)
}
