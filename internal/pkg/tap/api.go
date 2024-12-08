package tap

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type orderResponseItem struct {
	OrderID         int32     `json:"order_id"`
	OrderCreatedAt  time.Time `json:"order_created_at"`
	ProductName     string    `json:"product_name"`
	ProductCategory string    `json:"product_category"`
}

type orderResponse struct {
	Orders []orderResponseItem `json:"orders"`
}

func (t *Tap) getOrders() ([]orderResponseItem, error) {
	zap.S().Info("Tap: Getting orders")

	api := config.GetDefaultString("backend.tap.api", "https://tap.zeus.gent")
	req := fiber.Get(api + "/recent")

	res := new(orderResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return nil, errors.Join(append(errs, errors.New("Tap: Order API request failed"))...)
	}
	if status != fiber.StatusOK {
		return nil, errors.New("error getting orders")
	}

	return res.Orders, nil
}
