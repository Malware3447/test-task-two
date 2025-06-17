package api

import (
	"net/http"
	"test-task-two/internal/api/crut"
)

type Service struct {
	api crut.CurrencyOperations
}

type Params struct {
	Api crut.CurrencyOperations
}

func NewService(params *Params) *Service {
	return &Service{api: params.Api}
}

func (s *Service) AddAndWithdrawAmount(w http.ResponseWriter, r *http.Request) {
	s.api.AddAndWithdrawAmount(w, r)
}

func (s *Service) GetAmount(w http.ResponseWriter, r *http.Request) {
	s.api.GetAmount(w, r)
}
