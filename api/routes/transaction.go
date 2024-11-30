package routes

import (
	"go-lucid/api/handlers"
	"net/http"
)

func RegisterTransactionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/transaction", handlers.GetTransaction)
	mux.HandleFunc("/broadcast", handlers.BroadcastTransaction)
}
