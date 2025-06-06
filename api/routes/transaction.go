package routes

import (
	"net/http"

	"go-lucid/api/handlers"
	tx_p2p "go-lucid/p2p/transaction"
)

func RegisterTransactionRoutes(mux *http.ServeMux, transactionService *tx_p2p.TransactionService) {
	handler := handlers.NewTransactionHandler(transactionService)
	mux.HandleFunc("/transaction", handler.GetTransaction)
	mux.HandleFunc("/broadcast", handler.BroadcastTransaction)
}
