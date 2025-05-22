package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"coral.daniel-guo.com/internal/aws"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

type Pool struct {
	pool *pgxpool.Pool
}

func NewPool(env string) (*Pool, error) {
	config, err := loadDBConfig(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load db config: %w", err)
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.DBName)
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

func loadDBConfig(env string) (*DBConfig, error) {
	secretName := fmt.Sprintf("hub-insights-rds-cluster-readonly-%s", env)
	secretData, err := aws.GetSecret(secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret from :%s : %w", env, err)
	}

	var dbConfig DBConfig
	if err := json.Unmarshal([]byte(secretData), &dbConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret data: %w", err)
	}

	return &dbConfig, nil
}
