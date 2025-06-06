package handlers

import (
	"encoding/json"
	"go-lucid/core/transaction"
	"go-lucid/mempool"
	tx_p2p "go-lucid/p2p/transaction"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	TxService *tx_p2p.TransactionService
}

func NewTransactionHandler(txService *tx_p2p.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		TxService: txService,
	}
}

// BroadcastTransaction allows users to broadcast a transaction into mempool and into the network
func (h *TransactionHandler) BroadcastTransaction(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo add dtos for request and get tx

	random_hash := []byte(strconv.Itoa(rand.Int()))
	tx := transaction.RawTransaction{
		Hash: random_hash,
	}

	if rand.Int()%2 == 0 {
		errorResponse := map[string]any{
			"reason":  "validation_failed",
			"message": "Transaction validation failed",
			"code":    400,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	err = mempool.GetMempool().AddTx(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Transaction broadcasted successfully"))
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	// Logic to get a transaction
}
