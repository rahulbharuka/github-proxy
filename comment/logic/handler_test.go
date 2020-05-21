package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	h1 := GetHandler()
	h2 := GetHandler()
	assert.Equal(t, h1, h2)
}
