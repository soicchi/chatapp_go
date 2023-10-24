package errors

import (
	"fmt"
)

type EntityError struct {
	EntityType string
}

func NewEntityError(entityType string) *EntityError {
	return &EntityError{EntityType: entityType}
}

func (e *EntityError) NotFoundError() error {
	return fmt.Errorf("%s not found", e.EntityType)
}
