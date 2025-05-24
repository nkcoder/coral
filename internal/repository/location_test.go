package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"coral.daniel-guo.com/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockPool implements the PoolInterface for testing
type MockPool struct {
	mock.Mock
}

func (m *MockPool) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(pgx.Row)
}

func (m *MockPool) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

func (m *MockPool) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0).(int64), mockArgs.Error(1)
}

func (m *MockPool) Close() {
	m.Called()
}

// MockRow implements pgx.Row for testing
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	mockArgs := m.Called(dest[0], dest[1], dest[2])
	return mockArgs.Error(0)
}

func TestNewLocationRepository(t *testing.T) {
	t.Run("should create repository with pool", func(t *testing.T) {
		pool := &MockPool{}
		repo := NewLocationRepository(pool)

		assert.NotNil(t, repo)
		assert.Equal(t, pool, repo.db)
	})
}

func TestLocationRepository_FindByName(t *testing.T) {
	tests := []struct {
		name           string
		locationName   string
		expectedName   string // the name that should be passed to the query
		setupMock      func(*MockPool, *MockRow, string)
		expectedResult *model.Location
		expectedError  string
	}{
		{
			name:         "successful find with email",
			locationName: "Test Club",
			expectedName: "Test Club",
			setupMock: func(pool *MockPool, row *MockRow, expectedName string) {
				pool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.MatchedBy(func(args []interface{}) bool {
					return len(args) == 1 && args[0] == expectedName
				})).Return(row)
				row.On("Scan", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					// Simulate successful scan
					id := args.Get(0).(*string)
					name := args.Get(1).(*string)
					email := args.Get(2).(*sql.NullString)

					*id = "123"
					*name = "Test Club"
					*email = sql.NullString{String: "test@example.com", Valid: true}
				}).Return(nil)
			},
			expectedResult: &model.Location{
				ID:    "123",
				Name:  "Test Club",
				Email: "test@example.com",
			},
		},
		{
			name:         "successful find without email",
			locationName: "Test Club",
			expectedName: "Test Club",
			setupMock: func(pool *MockPool, row *MockRow, expectedName string) {
				pool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.MatchedBy(func(args []interface{}) bool {
					return len(args) == 1 && args[0] == expectedName
				})).Return(row)
				row.On("Scan", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					id := args.Get(0).(*string)
					name := args.Get(1).(*string)
					email := args.Get(2).(*sql.NullString)

					*id = "123"
					*name = "Test Club"
					*email = sql.NullString{Valid: false}
				}).Return(nil)
			},
			expectedResult: &model.Location{
				ID:    "123",
				Name:  "Test Club",
				Email: "",
			},
		},
		{
			name:         "location not found",
			locationName: "Nonexistent Club",
			expectedName: "Nonexistent Club",
			setupMock: func(pool *MockPool, row *MockRow, expectedName string) {
				pool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.MatchedBy(func(args []interface{}) bool {
					return len(args) == 1 && args[0] == expectedName
				})).Return(row)
				row.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedResult: nil,
		},
		{
			name:         "database error",
			locationName: "Test Club",
			expectedName: "Test Club",
			setupMock: func(pool *MockPool, row *MockRow, expectedName string) {
				pool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.MatchedBy(func(args []interface{}) bool {
					return len(args) == 1 && args[0] == expectedName
				})).Return(row)
				row.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "error querying location by name",
		},
		{
			name:         "should trim whitespace from name",
			locationName: "  Test Club  ",
			expectedName: "Test Club", // Should be trimmed
			setupMock: func(pool *MockPool, row *MockRow, expectedName string) {
				// Should call with trimmed name
				pool.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.MatchedBy(func(args []interface{}) bool {
					return len(args) == 1 && args[0] == expectedName
				})).Return(row)
				row.On("Scan", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					id := args.Get(0).(*string)
					name := args.Get(1).(*string)
					email := args.Get(2).(*sql.NullString)

					*id = "123"
					*name = "Test Club"
					*email = sql.NullString{String: "test@example.com", Valid: true}
				}).Return(nil)
			},
			expectedResult: &model.Location{
				ID:    "123",
				Name:  "Test Club",
				Email: "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := &MockPool{}
			row := &MockRow{}
			repo := NewLocationRepository(pool)

			tt.setupMock(pool, row, tt.expectedName)

			result, err := repo.FindByName(tt.locationName)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			pool.AssertExpectations(t)
			row.AssertExpectations(t)
		})
	}
}
