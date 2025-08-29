package tap

import (
	"context"
	"encoding/json"
	"fmt"
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

func (o orderResponse) ToModel() []model.Tap {
	taps := make([]model.Tap, 0, len(o.Orders))

	for _, order := range o.Orders {
		var category model.TapCategory = "unknown"
		switch order.ProductCategory {
		case "soft":
			category = model.Soft
		case "mate":
			category = model.Mate
		case "beer":
			category = model.Beer
		case "food":
			category = model.Food
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

	var orders orderResponse

	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("decode tap order response %w", err)
	}

	return orders.ToModel(), nil
}
