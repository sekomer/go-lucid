package p2p

import "go-lucid/p2p"

const (
	TransactionServiceVersion = "1.0.0"
	TransactionServiceName    = "transaction"
)

type TransactionService struct {
	*p2p.BaseService
}

var _ p2p.P2PService = (*TransactionService)(nil)
