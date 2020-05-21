package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDBHandler(t *testing.T) {
	h1 := NewDBHandler()
	h2 := NewDBHandler()
	assert.Equal(t, h1, h2)
}
