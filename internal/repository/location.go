// Package repository provides data access functionality
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"coral.daniel-guo.com/internal/model"
)

// LocationRepository provides data access for locations
type LocationRepository struct {
	db PoolInterface
}

// NewLocationRepository creates a new location repository
func NewLocationRepository(db PoolInterface) *LocationRepository {
	return &LocationRepository{db: db}
}

// FindByName looks up a location by its name
func (r *LocationRepository) FindByName(name string) (*model.Location, error) {
	ctx := context.Background()
	trimmedName := strings.TrimSpace(name)

	query := `
		SELECT id, name, email
		FROM location
		WHERE TRIM(name) = $1
	`

	var location model.Location
	var email sql.NullString

	err := r.db.QueryRow(ctx, query, trimmedName).Scan(&location.ID, &location.Name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying location by name: %w", err)
	}

	if email.Valid {
		location.Email = email.String
	}

	return &location, nil
}
