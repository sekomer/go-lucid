package mempool

import (
	"go-lucid/core/transaction"
	"sync"
)

var (
	once     sync.Once
	instance *Mempool
)

type Mempool struct {
	sync.RWMutex
	txs map[string]*transaction.RawTransaction
}

func GetMempool() *Mempool {
	once.Do(func() {
		instance = &Mempool{
			txs: make(map[string]*transaction.RawTransaction),
		}
	})
	return instance
}

func GetTestMempool() *Mempool {
	return &Mempool{
		txs: make(map[string]*transaction.RawTransaction),
	}
}

func (m *Mempool) AddTx(tx *transaction.RawTransaction) error {
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

func (m *Mempool) GetTx(hash string) *transaction.RawTransaction {
	m.RLock()
	defer m.RUnlock()

	return m.txs[hash]
}

func (m *Mempool) GetTxs() []*transaction.RawTransaction {
	m.RLock()
	defer m.RUnlock()

	txs := make([]*transaction.RawTransaction, 0, len(m.txs))
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

	m.txs = make(map[string]*transaction.RawTransaction)
}
