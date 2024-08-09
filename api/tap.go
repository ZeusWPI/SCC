package api

import (
	"encoding/json"
	"log"
	"net/http"
	"scc/config"
	"scc/screen"
	"time"
)

type tapReponse struct {
	Orders []screen.TapOrder `json:"orders"`
}

var (
	tapURL             = config.GetConfig().Tap.URL
	timestampLayout    = config.GetConfig().Tap.TimestampLayout
	lastOrderTimestamp = time.Time{}
)

func runTapRequests(app *screen.ScreenApp) {
	for true {
		recentOrders, err := tapGetRecentOrders()
		if err != nil {
			log.Printf("Error: Unable to get recent order: %s\n", err)
		}
		for _, order := range recentOrders.Orders {
			timestamp, err := time.Parse(timestampLayout, order.OrderCreatedAt)
			if err != nil {
				log.Printf("Error: Unable to parse timestamp: %s\n", err)
			}

			if order.ProductCategory == "beverages" && timestamp.After(lastOrderTimestamp) {
				app.Tap.Update(&order)
				lastOrderTimestamp = timestamp
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func tapGetRecentOrders() (*tapReponse, error) {
	req, err := http.NewRequest("GET", tapURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &tapReponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}
