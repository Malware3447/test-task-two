package pg

import (
	"context"
	"test-task-two/internal/db/pg"
	"test-task-two/internal/models/request"
	"test-task-two/internal/models/response"
)

type Service struct {
	repo pg.Repository
}

func NewService(repo pg.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddAmount(ctx context.Context, reqModel *request.Request) error {
	return s.repo.AddAmount(ctx, reqModel)
}

func (s *Service) WithdrawAmount(ctx context.Context, reqModel *request.Request) error {
	return s.repo.WithdrawAmount(ctx, reqModel)
}

func (s *Service) GetAmount(ctx context.Context, uuid string) (*response.Response, error) {
	return s.repo.GetAmount(ctx, uuid)
}
