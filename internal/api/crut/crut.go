package crut

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"test-task-two/internal/models/request"
	"test-task-two/internal/service/db/pg"
)

type Crut struct {
	repoPg *pg.Service
}

type Params struct {
	RepoPg *pg.Service
}

func NewCrut(params *Params) CurrencyOperations {
	return &Crut{repoPg: params.RepoPg}
}

func (c *Crut) AddAndWithdrawAmount(w http.ResponseWriter, r *http.Request) {
	const op = "router.AddAndWithdrawAmount"
	ctx := context.WithValue(r.Context(), "router", op)

	body := &request.Request{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.Operation == "DEPOSIT" {
		err = c.repoPg.AddAmount(ctx, body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if body.Operation == "WITHDRAW" {
		err = c.repoPg.WithdrawAmount(ctx, body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Данные обновлены")
}

func (c *Crut) GetAmount(w http.ResponseWriter, r *http.Request) {
	const op = "router.GetAmount"
	ctx := context.WithValue(r.Context(), "router", op)

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := c.repoPg.GetAmount(ctx, uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Данные отправлены")
}
