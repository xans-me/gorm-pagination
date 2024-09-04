package transaction

import (
	"net/http"
)

// GetTransactions handles the request for paginated transactions.
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	response, err := GetPaginatedTransactions(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, response)
}
