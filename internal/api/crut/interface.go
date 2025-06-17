package crut

import (
	"net/http"
)

type CurrencyOperations interface {
	AddAndWithdrawAmount(w http.ResponseWriter, r *http.Request)
	GetAmount(w http.ResponseWriter, r *http.Request)
}
