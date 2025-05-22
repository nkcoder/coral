package domain

import (
	"context"
	"database/sql"
	"strings"

	"coral.daniel-guo.com/internal/db"
)

type LocationRepository struct {
	db *db.Pool
}

func NewLocationRepository(db *db.Pool) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) FindByName(name string) (*Location, error) {
	ctx := context.Background()
	trimmedName := strings.TrimSpace(name)

	query := `
		SELECT id, name, email
		FROM location
		WHERE TRIM(name) = $1
	`

	var location Location
	var email sql.NullString

	err := r.db.QueryRow(ctx, query, trimmedName).Scan(&location.ID, &location.Name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if email.Valid {
		location.Email = email.String
	}

	return &location, nil
}
