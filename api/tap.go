package api

import (
	"log"
	"scc/config"
	"scc/screen"
	"slices"
	"time"
)

type tapReponse struct {
	Orders []screen.TapOrder `json:"orders"`
}

var (
	tapURL                = config.GetConfig().Tap.URL
	tapLastOrderTimestamp = time.Time{}
)

func tapRunRequests(app *screen.ScreenApp) {
	headers := []header{jsonHeader}

	for {
		recentOrders := &tapReponse{}
		if err := makeGetRequest(tapURL, headers, recentOrders); err != nil {
			log.Printf("Error: Unable to get recent orders: %s\n", err)
		}

		slices.SortStableFunc(recentOrders.Orders, func(a, b screen.TapOrder) int {
			return a.OrderCreatedAt.Compare(b.OrderCreatedAt)
		})

		for _, order := range recentOrders.Orders {
			if order.OrderCreatedAt.After(tapLastOrderTimestamp) {
				app.Tap.Update(&order)
				tapLastOrderTimestamp = order.OrderCreatedAt
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
