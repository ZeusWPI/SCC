package tap

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/utils"
)

type orderResponseItem struct {
	OrderID         int       `json:"order_id"`
	OrderCreatedAt  time.Time `json:"order_created_at"`
	ProductName     string    `json:"product_name"`
	ProductCategory string    `json:"product_category"`
}

type orderResponse struct {
	Orders []orderResponseItem `json:"orders"`
}

func (o orderResponse) ToModel(beers []string) []model.Tap {
	taps := make([]model.Tap, 0, len(o.Orders))

	for _, order := range o.Orders {
		var category model.TapCategory = "unknown"
		switch order.ProductCategory {
		case "food":
			category = model.Food
		case "beverages":
			switch {
			case strings.Contains(order.ProductName, "Mate") || strings.Contains(order.ProductName, "Mio Mio"):
				category = model.Mate
			case slices.ContainsFunc(beers, func(beer string) bool { return strings.Contains(order.ProductName, beer) }):
				category = model.Beer
			default:
				category = model.Soft
			}
		}

		taps = append(taps, model.Tap{
			OrderID:   order.OrderID,
			CreatedAt: order.OrderCreatedAt,
			Name:      order.ProductName,
			Category:  category,
		})
	}

	return taps
}

func (t *Tap) getOrders(ctx context.Context) ([]model.Tap, error) {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    t.url + "/recent",
	})
	if err != nil {
		return nil, fmt.Errorf("get all tap orders %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %s", resp.Status)
	}

	var orders orderResponse

	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("decode tap order response %w", err)
	}

	return orders.ToModel(t.beers), nil
}
