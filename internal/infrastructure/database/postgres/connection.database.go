package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/L1mus/backend-ewallet/internal/infrastructure/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func DBConnection(cfg config.DBConfig) (*sql.DB, error) {
	// open connection
	db, err := sql.Open("pgx", cfg.DataSourceName())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	//pool connection tuning
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// verifies a connection to the database is still alive ?
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}
