// Package tap provides all tap related logic
package tap

import (
	"context"
	"slices"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Tap struct {
	tap   repository.Tap
	url   string
	beers []string
}

var defaultBeers = []string{
	"Schelfaut",
	"Duvel",
	"Fourchette",
	"Jupiler",
	"Karmeliet",
	"Kriek",
	"Chouffe",
	"Maes",
	"Somersby",
	"Sportzot",
	"Stella",
}

func New(repo repository.Repository) *Tap {
	return &Tap{
		tap:   *repo.NewTap(),
		url:   config.GetDefaultString("backend.tap.url", "https://tap.zeus.gent"),
		beers: config.GetDefaultStringSlice("backend.tap.beers", defaultBeers),
	}
}

// Update gets all new orders from tap
func (t *Tap) Update(ctx context.Context) error {
	// Get latest order
	lastOrder, err := t.tap.GetLast(ctx)
	if err != nil {
		return err
	}
	if lastOrder == nil {
		lastOrder = &model.Tap{OrderID: -1}
	}

	// Get all orders
	allOrders, err := t.getOrders(ctx)
	if err != nil {
		return err
	}

	// Only keep the new orders
	orders := utils.SliceFilter(allOrders, func(o model.Tap) bool { return o.OrderID > lastOrder.OrderID })
	slices.SortFunc(orders, func(a, b model.Tap) int { return a.OrderID - b.OrderID })

	for _, order := range orders {
		if err := t.tap.Create(ctx, &order); err != nil {
			return err
		}
	}

	return nil
}
