package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Tree struct {
	bun.BaseModel `bun:"table:trees"`

	ID        uint64    `bun:"id,pk"`
	Estate    Estate    `bun:"rel:belongs-to"`
	X         uint16    `bun:"x,notnull"`
	Y         uint16    `bun:"y,notnull"`
	Height    uint8     `bun:"height,notnull"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
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
		UUID:   uuid.NewString(),
		Width:  width,
		Length: length,
	}
}
