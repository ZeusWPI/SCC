package tap

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type orderResponseItem struct {
	OrderID         int64     `json:"order_id"`
	OrderCreatedAt  time.Time `json:"order_created_at"`
	ProductName     string    `json:"product_name"`
	ProductCategory string    `json:"product_category"`
}

type orderResponse struct {
	Orders []orderResponseItem `json:"orders"`
}

func (t *Tap) getOrders() ([]orderResponseItem, error) {
	zap.S().Info("Tap: Getting orders")

	req := fiber.Get(t.api + "/recent")

	res := new(orderResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return nil, errors.Join(append([]error{errors.New("Tap: Order API request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return nil, errors.New("error getting orders")
	}

	return res.Orders, nil
}
