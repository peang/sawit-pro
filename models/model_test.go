package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func TestNewTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEstate := &Estate{
		ID:               1,
		Width:            100,
		Length:           100,
		MinTreeHeight:    0,
		MaxTreeHeight:    0,
		MedianTreeHeight: 0,
		TreeCount:        0,
	}

	newTree, err := NewTree(mockEstate, 10, 10, 20)
	assert.NoError(t, err)
	assert.NotNil(t, newTree)
	assert.IsType(t, "", newTree.UUID, "newTree.UUID is not a string")

	assert.Equal(t, uint8(20), newTree.Height)
	assert.Equal(t, uint16(10), newTree.X)
	assert.Equal(t, uint16(10), newTree.Y)
	assert.Equal(t, uint64(1), newTree.EstateID)

	assert.Equal(t, uint8(20), newTree.Estate.MinTreeHeight)
	assert.Equal(t, uint8(20), newTree.Estate.MaxTreeHeight)
	assert.Equal(t, uint8(20), newTree.Estate.MedianTreeHeight)
	assert.Equal(t, uint8(1), newTree.Estate.TreeCount)
}

func TestCalculateEstateTreeStats(t *testing.T) {
	mockEstate := &Estate{
		ID:               1,
		Width:            100,
		Length:           100,
		MinTreeHeight:    0,
		MaxTreeHeight:    0,
		MedianTreeHeight: 0,
		TreeCount:        0,
	}

	tree := &Tree{
		UUID:   "mockUUID",
		Estate: mockEstate,
		Height: 20,
	}

	err := tree.CalculateEstateTreeStats()
	assert.NoError(t, err)

	assert.Equal(t, uint8(20), tree.Estate.MinTreeHeight)
	assert.Equal(t, uint8(20), tree.Estate.MaxTreeHeight)
	assert.Equal(t, uint8(20), tree.Estate.MedianTreeHeight)
	assert.Equal(t, uint8(1), tree.Estate.TreeCount)
}

func TestCalculateEstateTreeStats_OutsideBoundaries(t *testing.T) {
	mockEstate := &Estate{
		ID:               1,
		Width:            5,
		Length:           5,
		MinTreeHeight:    0,
		MaxTreeHeight:    0,
		MedianTreeHeight: 0,
		TreeCount:        0,
	}

	tree := &Tree{
		UUID:   "mockUUID",
		Estate: mockEstate,
		X:      10,
		Y:      10,
		Height: 20,
	}

	err := tree.CalculateEstateTreeStats()

	assert.Error(t, err)
	assert.Equal(t, "outside of boundaries", err.Error())
}
