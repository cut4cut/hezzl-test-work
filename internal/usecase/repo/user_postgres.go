package repo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/cut4cut/hezzl-test-work/internal/entity"
	"github.com/cut4cut/hezzl-test-work/pkg/postgres"
	"github.com/georgysavva/scany/pgxscan"
)

type UserRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func newSortExpr(pag entity.Pagination) (sortName string, sortCreated string) {
	sortName, sortCreated = "name ASC", "created_dt ASC"

	if pag.DescName {
		sortName = "name DESC"
	}
	if pag.DescCreated {
		sortCreated = "created_dt DESC"
	}

	return
}

func (r *UserRepo) Create(ctx context.Context, name string) (int64, error) {
	id := int64(-1)
	sql, _, err := r.Builder.
		Insert("dim_users").
		Columns("id, created_dt, name").
		Values(
			sq.Expr("DEFAULT"),
			sq.Expr("DEFAULT"),
			name,
		).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		return id, fmt.Errorf("userRepo - create - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, name).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("userRepo - create - tx.QueryRow: %w", err)
	}

	return id, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) (int64, error) {
	sql, _, err := r.Builder.
		Delete("dim_users").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		return id, fmt.Errorf("userRepo - delete - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, id).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("userRepo - delete - tx.QueryRow: %w", err)
	}

	return id, nil
}

func (r *UserRepo) GetList(ctx context.Context, pag entity.Pagination) ([]entity.User, error) {
	offset := pag.Page * pag.ItemsPerPage
	sortExprPrice, sortExprCreated := newSortExpr(pag)

	sql, _, err := r.Builder.
		Select("id, created_dt, name").
		From("dim_users").
		OrderBy(sortExprPrice).
		OrderBy(sortExprCreated).
		Limit(pag.ItemsPerPage).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("userRepo - getList - r.Builder: %w", err)
	}

	users := make([]entity.User, int(pag.ItemsPerPage))
	if err := pgxscan.Select(
		ctx, r.Pool, &users, sql,
	); err != nil {
		return nil, fmt.Errorf("userRepo - getList - pgxscan.Select: %w", err)
	}

	return users, nil
}
