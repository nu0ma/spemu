package executor

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/nu0ma/spemu/pkg/config"
)

type Executor struct {
	client *spanner.Client
}

func New(cfg *config.Config) (*Executor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if cfg.EmulatorHost != "" {
		os.Setenv("SPANNER_EMULATOR_HOST", cfg.EmulatorHost)
	}

	client, err := spanner.NewClient(ctx, cfg.DatabasePath())
	if err != nil {
		return nil, fmt.Errorf("failed to create Spanner client: %w", err)
	}

	return &Executor{client: client}, nil
}

func (e *Executor) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

func (e *Executor) ExecuteStatements(statements []string, verbose bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := e.client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		for i, stmt := range statements {
			if verbose {
				limit := 100
				if len(stmt) < limit {
					limit = len(stmt)
				}
				fmt.Printf("Executing statement %d/%d: %s\n", i+1, len(statements), stmt[:limit]+"...")
			}

			_, err := txn.Update(ctx, spanner.Statement{SQL: stmt})
			if err != nil {
				return fmt.Errorf("failed to execute statement %d: %w\nStatement: %s", i+1, err, stmt)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
