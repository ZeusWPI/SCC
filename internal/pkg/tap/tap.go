// Package tap provides all tap related logic
package tap

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/util"
)

// Tap represents a tap instance
type Tap struct {
	db    *db.DB
	api   string
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

// New creates a new tap instance
func New(db *db.DB) *Tap {
	api := config.GetDefaultString("tap.api", "https://tap.zeus.gent")
	beers := config.GetDefaultStringSlice("tap.beers", defaultBeers)

	return &Tap{db: db, api: api, beers: beers}
}

// Update gets all new orders from tap
func (t *Tap) Update() error {
	// Get latest order
	lastOrder, err := t.db.Queries.GetLastOrderByOrderID(context.Background())
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		lastOrder = sqlc.Tap{OrderID: -1}
	}

	// Get all orders
	allOrders, err := t.getOrders()
	if err != nil {
		return err
	}

	// Only keep the new orders
	orders := util.SliceFilter(allOrders, func(o orderResponseItem) bool { return o.OrderID > lastOrder.OrderID })

	if len(orders) == 0 {
		return nil
	}

	// Adjust categories
	t.adjustCategories(orders)

	// Insert all new orders
	errs := make([]error, 0)
	for _, order := range orders {
		_, err := t.db.Queries.CreateTap(context.Background(), sqlc.CreateTapParams{
			OrderID:        order.OrderID,
			OrderCreatedAt: order.OrderCreatedAt,
			Name:           order.ProductName,
			Category:       order.ProductCategory,
		})
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// adjustCategories changes the categories of the orders to the custom ones
func (t *Tap) adjustCategories(orders []orderResponseItem) {
	for i := range orders {
		order := &orders[i] // Take a pointer to the struct to modify it directly
		switch order.ProductCategory {
		case "food":
			order.ProductCategory = "Food"
		case "other":
			order.ProductCategory = "Other"
		case "beverages":
			// Atm only beverages get special categories
			if strings.Contains(order.ProductName, "Mate") {
				order.ProductCategory = "Mate"
			} else if slices.ContainsFunc(t.beers, func(beer string) bool { return strings.Contains(order.ProductName, beer) }) {
				order.ProductCategory = "Beer"
			} else {
				order.ProductCategory = "Soft"
			}
		}
	}
}
