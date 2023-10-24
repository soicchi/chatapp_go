package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	err := NewEntityError("user").NotFoundError()
	assert.Equal(t, "user not found", err.Error())
}
