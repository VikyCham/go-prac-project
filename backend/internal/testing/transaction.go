package testing

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// TxFn represents a function that executes within a transaction
type TxFn func(tx pgx.Tx) error

// With Transaction runs a function within a transaction and rolls it back forward
func WithTransaction(ctx context.Context, db *TestDB, fn TxFn) error {
	// begin transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}

	// ensure rollback happens if commit dosen't occur
	defer tx.Rollback(ctx)

	// run the function within the transaction
	if err := fn(tx); err != nil {
		return err
	}

	// transaction was successful, commit it
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

// WithRollbackTransaction runs a function within a transaction and always rolls it back
// Useful for tests where you want to execute operations but never persist them
func WithRollbackTransaction(ctx context.Context, db *TestDB, fn TxFn) error {
	// begin transaction

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}

	// Always rollback at the end
	defer tx.Rollback(ctx)

	// run the function within the transaction
	return fn(tx)
}
