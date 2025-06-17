package api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"test-task-two/internal/service/api"
)

type Router struct {
	router  *chi.Mux
	apiServ *api.Service
}

func NewRouter(apiServ *api.Service) *Router {
	return &Router{
		router:  nil,
		apiServ: apiServ,
	}
}

func (r *Router) Init(ctx context.Context) {
	const op = "router.Init"
	ctx = context.WithValue(ctx, "router", op)

	r.router = chi.NewRouter()

	r.router.Route("api/v1", func(router chi.Router) {
		router.Post("/wallet/", r.apiServ.AddAndWithdrawAmount)
		router.Get("/wallets/{uuid}", r.apiServ.GetAmount)
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", 8080), r.router); err != nil {
			panic(fmt.Sprintf("%v: %v", op, err))
		}
	}()
}
