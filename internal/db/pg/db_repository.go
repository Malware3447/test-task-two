package pg

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"test-task-two/internal/models/request"
	"test-task-two/internal/models/response"
)

type RepositoryPg struct {
	db *pgxpool.Pool
}

type Params struct {
	Db *pgxpool.Pool
}

func NewRepository(params *Params) Repository {
	return &RepositoryPg{db: params.Db}
}

func (r *RepositoryPg) AddAmount(ctx context.Context, ReqModel *request.Request) (*response.Response, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Error: "internal error"}
		return errorResp, err
	}
	defer tx.Rollback(ctx)

	const q = `
	UPDATE currency
	SET amount = amount + $2
	WHERE uuid = $1
`
	_, err = tx.Exec(ctx, q, ReqModel.Uuid, ReqModel.Amount)
	if err != nil {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Error: "internal error"}
		return errorResp, err
	}

	var currentAmount int
	err = tx.QueryRow(ctx, "SELECT amount FROM currency WHERE uuid = $1 FOR UPDATE", ReqModel.Uuid).Scan(&currentAmount)
	if err != nil {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Amount: currentAmount, Error: "internal error"}
		return errorResp, err
	}

	errorResp := &response.Response{Uuid: ReqModel.Uuid, Amount: currentAmount, Error: "OK"}

	return errorResp, tx.Commit(ctx)
}

func (r *RepositoryPg) WithdrawAmount(ctx context.Context, ReqModel *request.Request) (*response.Response, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Error: "internal error"}
		return errorResp, err
	}
	defer tx.Rollback(ctx)

	var currentAmount int
	err = tx.QueryRow(ctx, "SELECT amount FROM currency WHERE uuid = $1 FOR UPDATE", ReqModel.Uuid).Scan(&currentAmount)
	if err != nil {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Amount: currentAmount, Error: "internal error"}
		return errorResp, err
	}

	if currentAmount < ReqModel.Amount {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Amount: currentAmount, Error: "Недостаточно средств для снятия"}
		return errorResp, errors.New("Недостаточно средств для снятия")
	}

	const q = `
	UPDATE currency
	SET amount = amount - $2
	WHERE uuid = $1
`
	_, err = tx.Exec(ctx, q, ReqModel.Uuid, ReqModel.Amount)
	if err != nil {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Amount: currentAmount, Error: "internal error"}
		return errorResp, err
	}

	err = tx.QueryRow(ctx, "SELECT amount FROM currency WHERE uuid = $1 FOR UPDATE", ReqModel.Uuid).Scan(&currentAmount)
	if err != nil {
		errorResp := &response.Response{Uuid: ReqModel.Uuid, Amount: currentAmount, Error: "internal error"}
		return errorResp, err
	}

	errorResp := &response.Response{Uuid: ReqModel.Uuid, Amount: currentAmount, Error: "OK"}

	return errorResp, tx.Commit(ctx)
}

func (r *RepositoryPg) GetAmount(ctx context.Context, uuid string) (*response.Response, error) {
	const q = `
    SELECT uuid, amount
    FROM currency
    WHERE uuid = $1
    `

	var resp response.Response
	err := r.db.QueryRow(ctx, q, uuid).Scan(&resp.Uuid, &resp.Amount)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
