package store

import (
	"errors"
	"fmt"
)

// Standard error definitions for store operations
var (
	// ErrNotFound is returned when a requested resource doesn't exist
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists is returned when trying to create a resource that already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrResourceInUse is returned when attempting to delete a resource that is referenced by other resources
	ErrResourceInUse = errors.New("resource is in use by other resources")

	// ErrInvalidConfig is returned when a configuration is invalid
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrInternalError is returned when an unexpected internal error occurs
	ErrInternalError = errors.New("internal store error")

	// ErrInvalidID is returned when a resource ID is invalid
	ErrInvalidID = errors.New("invalid resource ID")
)

// DependencyError provides detailed information about dependencies when a resource can't be deleted
type DependencyError struct {
	ResourceType string
	ResourceID   string
	Dependencies []Dependency
}

// Dependency represents a resource that depends on another resource
type Dependency struct {
	ResourceType string `json:"resourceType"` // "router", "service", "middleware"
	ID           string `json:"id"`           // ID of the dependent resource
	Field        string `json:"field"`        // The field that references the dependency
}

func (e *DependencyError) Error() string {
	return fmt.Sprintf("%s '%s' is in use by %d other resources", e.ResourceType, e.ResourceID, len(e.Dependencies))
}

// NewDependencyError creates a new DependencyError with the specified dependencies
func NewDependencyError(resourceType, resourceID string, dependencies []Dependency) *DependencyError {
	return &DependencyError{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Dependencies: dependencies,
	}
}

// ValidationError represents a validation error for a specific resource field
type ValidationError struct {
	ResourceType string
	ResourceID   string
	Field        string
	Message      string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s '%s': %s: %s",
		e.ResourceType, e.ResourceID, e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(resourceType, resourceID, field, message string) *ValidationError {
	return &ValidationError{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Field:        field,
		Message:      message,
	}
}

// IsNotFound returns true if the error is an ErrNotFound error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists returns true if the error is an ErrAlreadyExists error
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsResourceInUse returns true if the error is an ErrResourceInUse error
func IsResourceInUse(err error) bool {
	return errors.Is(err, ErrResourceInUse) || IsDependencyError(err)
}

// IsDependencyError returns true if the error is a DependencyError
func IsDependencyError(err error) bool {
	_, ok := err.(*DependencyError)
	return ok
}

// IsValidationError returns true if the error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// GetDependencies returns the dependencies if the error is a DependencyError, otherwise returns nil
func GetDependencies(err error) []Dependency {
	if depErr, ok := err.(*DependencyError); ok {
		return depErr.Dependencies
	}
	return nil
}
