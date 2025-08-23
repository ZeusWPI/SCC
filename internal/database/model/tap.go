package model

import (
	"time"

	"github.com/zeusWPI/scc/pkg/sqlc"
)

type TapCategory string

const (
	Soft    TapCategory = "soft"
	Mate    TapCategory = "mate"
	Beer    TapCategory = "beer"
	Food    TapCategory = "food"
	Unknown TapCategory = "unknown"
)

type Tap struct {
	ID        int
	OrderID   int
	Name      string
	Category  TapCategory
	CreatedAt time.Time
}

func TapModel(t sqlc.Tap) *Tap {
	return &Tap{
		ID:        int(t.ID),
		OrderID:   int(t.OrderID),
		Name:      t.Name,
		Category:  TapCategory(t.Category),
		CreatedAt: t.CreatedAt.Time,
	}
}

type TapCount map[TapCategory]int
