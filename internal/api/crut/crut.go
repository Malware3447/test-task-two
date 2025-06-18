package crut

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"sync"
	"test-task-two/internal/models/request"
	"test-task-two/internal/models/response"
	"test-task-two/internal/service/db/pg"
	"time"
)

type Crut struct {
	repoPg      *pg.Service
	walletQueue map[string]chan *request.Request
	queueMutex  sync.Mutex
	respChan    chan *response.Response
}

type Params struct {
	RepoPg *pg.Service
}

func NewCrut(params *Params) CurrencyOperations {
	return &Crut{
		repoPg:      params.RepoPg,
		walletQueue: make(map[string]chan *request.Request),
		respChan:    make(chan *response.Response, 100), // Глобальный канал для ответов
	}
}

func (c *Crut) processWalletQueue(uuid string) {
	queue := c.walletQueue[uuid]

	for req := range queue {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var err error
		var resp *response.Response
		if req.Operation == "DEPOSIT" {
			resp, err = c.repoPg.AddAmount(ctx, req)
		} else if req.Operation == "WITHDRAW" {
			resp, err = c.repoPg.WithdrawAmount(ctx, req)
		}
		cancel()
		if err != nil {
			log.Printf("Error processing request for wallet %s: %v", uuid, err)
		}
		c.respChan <- resp
	}
}

func (c *Crut) AddAndWithdrawAmount(w http.ResponseWriter, r *http.Request) {
	body := &request.Request{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	c.queueMutex.Lock()
	queue, exists := c.walletQueue[body.Uuid]
	if !exists {
		queue = make(chan *request.Request, 100)
		c.walletQueue[body.Uuid] = queue
		go c.processWalletQueue(body.Uuid)
	}
	c.queueMutex.Unlock()

	select {
	case queue <- body:
		resp := <-c.respChan
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Too many requests for this wallet", http.StatusTooManyRequests)
	}
}

func (c *Crut) GetAmount(w http.ResponseWriter, r *http.Request) {
	const op = "router.GetAmount"
	ctx := context.WithValue(r.Context(), "router", op)

	uuid := chi.URLParam(r, "uuid")
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Данные отправлены")
}
