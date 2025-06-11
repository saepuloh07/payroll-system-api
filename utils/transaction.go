package utils

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type WrapTransactionFunc = func(ctx context.Context, tx *sql.Tx) error

func WrapTransaction(ctx context.Context, db *sqlx.DB, wrapFunc WrapTransactionFunc) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = wrapFunc(ctx, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
