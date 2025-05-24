// Package repository provides data access functionality
package repository

import (
	"context"

	"coral.daniel-guo.com/internal/model"
	"github.com/jackc/pgx/v5"
)

// PoolInterface defines the interface for database pool operations
type PoolInterface interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	Close()
}

// LocationRepositoryInterface defines the interface for location repository operations
type LocationRepositoryInterface interface {
	FindByName(name string) (*model.Location, error)
}
