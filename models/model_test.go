package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEstate(t *testing.T) {
	width := uint16(10)
	length := uint16(20)

	estate := NewEstate(width, length)

	assert.NotNil(t, estate)

	assert.NotEmpty(t, estate.UUID)

	assert.Equal(t, width, estate.Width)
	assert.Equal(t, length, estate.Length)
}
