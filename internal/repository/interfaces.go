// Package repository provides data access functionality
package repository

import "coral.daniel-guo.com/internal/model"

// LocationRepositoryInterface defines the interface for location repository operations
type LocationRepositoryInterface interface {
	FindByName(name string) (*model.Location, error)
}
