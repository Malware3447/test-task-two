package pg

import (
	"context"
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

func (r *RepositoryPg) AddAmount(ctx context.Context, ReqModel *request.Request) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const q = `
	UPDATE currency
	SET amount = amount + $2
	WHERE uuid = $1
`
	_, err = tx.Exec(ctx, q, ReqModel.Uuid, ReqModel.Amount)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *RepositoryPg) WithdrawAmount(ctx context.Context, ReqModel *request.Request) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const q = `
	UPDATE currency
	SET amount = amount - $2
	WHERE uuid = $1
`
	_, err = tx.Exec(ctx, q, ReqModel.Uuid, ReqModel.Amount)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
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
