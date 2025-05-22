// Package repository provides data access functionality
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"coral.daniel-guo.com/internal/logger"
	"coral.daniel-guo.com/internal/secrets"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBConfig contains database connection configuration
type DBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

// Pool provides a wrapper around pgxpool for database operations
type Pool struct {
	pool *pgxpool.Pool
}

// PoolConfig holds configuration for database pool creation
type PoolConfig struct {
	Environment    string
	SecretsManager *secrets.Manager
}

// NewPool creates a new database connection pool
func NewPool(cfg PoolConfig) (*Pool, error) {
	config, err := loadDBConfig(cfg.Environment, cfg.SecretsManager)
	if err != nil {
		return nil, fmt.Errorf("failed to load db config: %w", err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Username, config.Password, config.Host, config.Port, config.DBName)

	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to database at %s:%d/%s", config.Host, config.Port, config.DBName)
	return &Pool{pool: dbPool}, nil
}

// Close closes the database connection pool
func (p *Pool) Close() {
	p.pool.Close()
}

// QueryRow executes a query that is expected to return a single row
func (p *Pool) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, query, args...)
}

// Query executes a query that returns rows
func (p *Pool) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, query, args...)
}

// Exec executes a query that doesn't return rows
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

// loadDBConfig loads database configuration from AWS Secrets Manager
func loadDBConfig(env string, secretsManager *secrets.Manager) (*DBConfig, error) {
	secretName := fmt.Sprintf("hub-insights-rds-cluster-readonly-%s", env)
	logger.Info("Loading database configuration from secret: %s", secretName)

	secretData, err := secretsManager.GetSecret(secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret from %s: %w", env, err)
	}

	var dbConfig DBConfig
	if err := json.Unmarshal([]byte(secretData), &dbConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret data: %w", err)
	}

	return &dbConfig, nil
}
