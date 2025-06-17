package pg

import (
	"context"
	"test-task-two/internal/models/request"
	"test-task-two/internal/models/response"
)

type Repository interface {
	AddAmount(ctx context.Context, ReqModel *request.Request) error
	WithdrawAmount(ctx context.Context, ReqModel *request.Request) error
	GetAmount(ctx context.Context, uuid string) (*response.Response, error)
}
