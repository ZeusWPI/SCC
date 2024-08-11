package api

import (
	"log"
	"scc/config"
	"scc/screen"
	"slices"
	"time"
)

type zessResponse struct {
	Scans []screen.ZessScan `json:"scans"`
}

var (
	zessURL                = config.GetConfig().Zess.URL
	zessLastOrderTimestamp = time.Time{}
)

func zessRunRequests(app *screen.ScreenApp) {
	headers := []header{jsonHeader}

	for {
		recentScans := &zessResponse{}
		if err := makeGetRequest(zessURL, headers, recentScans); err != nil {
			log.Printf("Error: Unable to get recent scans: %s\n", err)
		}

		slices.SortStableFunc(recentScans.Scans, func(a, b screen.ZessScan) int {
			return a.ScanTime.Compare(b.ScanTime)
		})

		for _, order := range recentScans.Scans {
			if order.ScanTime.After(zessLastOrderTimestamp) {
				app.Zess.Update(&order)
				zessLastOrderTimestamp = order.ScanTime
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
