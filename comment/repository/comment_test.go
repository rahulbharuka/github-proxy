package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommentRepo(t *testing.T) {
	h1 := NewCommentRepo()
	h2 := NewCommentRepo()
	assert.Equal(t, h1, h2)
}
