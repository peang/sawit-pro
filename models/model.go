package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tree struct {
	bun.BaseModel `bun:"table:trees"`

	ID        uint64    `bun:"id,pk"`
	EstateID  uint64    `bun:"estate_id,notnull"`
	UUID      string    `bun:"uuid,notnull"`
	X         uint16    `bun:"x,notnull"`
	Y         uint16    `bun:"y,notnull"`
	Height    uint8     `bun:"height,notnull"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`

	Estate *Estate `bun:"rel:belongs-to"`
}

type Estate struct {
	bun.BaseModel `bun:"table:estates"`

	ID               uint64    `bun:"id,pk"`
	UUID             string    `bun:"uuid,notnull"`
	Width            uint16    `bun:"width,notnull"`
	Length           uint16    `bun:"length,notnull"`
	TreeCount        uint8     `bun:"tree_count,notnull"`
	MinTreeHeight    uint8     `bun:"min_tree_height"`
	MaxTreeHeight    uint8     `bun:"max_tree_height"`
	MedianTreeHeight uint8     `bun:"median_tree_height"`
	CreatedAt        time.Time `bun:"created_at"`
	UpdatedAt        time.Time `bun:"updated_at"`
}

func NewEstate(width uint16, length uint16) *Estate {
	return &Estate{
		UUID:      uuid.NewString(),
		Width:     width,
		Length:    length,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewTree(estate *Estate, x uint16, y uint16, height uint8) (*Tree, error) {
	tree := Tree{
		UUID:      uuid.NewString(),
		EstateID:  estate.ID,
		Estate:    estate,
		X:         x,
		Y:         y,
		Height:    height,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := tree.CalculateEstateTreeStats()
	if err != nil {
		return nil, err
	}

	return &tree, nil
}

func (t *Tree) CalculateEstateTreeStats() (err error) {
	if t.X > t.Estate.Length {
		err = errors.New("outside of boundaries")
		return
	}

	if t.Y > t.Estate.Width {
		err = errors.New("outside of boundaries")
		return
	}

	if t.Estate.MinTreeHeight != 0 {
		if t.Height < t.Estate.MinTreeHeight {
			t.Estate.MinTreeHeight = t.Height
		}
	} else {
		t.Estate.MinTreeHeight = t.Height
	}

	if t.Estate.MaxTreeHeight != 0 {
		if t.Height > t.Estate.MaxTreeHeight {
			t.Estate.MaxTreeHeight = t.Height
		}
	} else {
		t.Estate.MaxTreeHeight = t.Height
	}

	var median uint8
	if t.Estate.TreeCount != 0 {
		median = (t.Estate.MedianTreeHeight + t.Height) / 2
	} else {
		median = t.Height
	}
	t.Estate.TreeCount++

	t.Estate.MedianTreeHeight = uint8(median)

	t.Estate.UpdatedAt = time.Now()

	return
}
