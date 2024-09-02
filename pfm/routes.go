package pfm

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/transactions", GetTransactions).Methods("GET")
}
