package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type Pool struct {
	pool *pgxpool.Pool
}

func NewPool(config *DBConfig) (*Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.DBName)
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Pool{pool: dbPool}, nil
}

func (p *Pool) Close() {
	p.pool.Close()
}

func (p *Pool) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, query, args...)
}

func (p *Pool) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, query, args...)
}

func (p *Pool) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	result := p.pool.QueryRow(ctx, "SELECT 1")
	if err := result.Scan(new(int)); err != nil {
		return 0, fmt.Errorf("connection test failed: %w", err)
	}

	commandTag, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag.RowsAffected(), nil
}
