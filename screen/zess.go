package screen

import (
	"scc/config"
	"scc/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type ZessScan struct {
	ID       int       `json:"id"`
	ScanTime time.Time `json:"scan_time"`
}

type Zess struct {
	ScreenApp *ScreenApp
	view      *tview.Flex
	chart     *tvxwidgets.Plot
}

var (
	scans = [][]float64{
		make([]float64, 0, config.GetConfig().Zess.DayAmount),
	}
	day = -1
)

func NewZess(screenApp *ScreenApp) *Zess {
	zess := Zess{
		ScreenApp: screenApp,
		view:      tview.NewFlex(),
		chart:     tvxwidgets.NewPlot(),
	}

	zess.view.SetBorder(true).SetTitle(" Zess ")

	zess.chart.SetBorder(false)
	zess.chart.SetLineColor([]tcell.Color{tcell.ColorOrange})
	zess.chart.SetAxesLabelColor(tcell.ColorYellow)
	zess.chart.SetAxesColor(tcell.ColorYellow)
	zess.chart.SetMarker(tvxwidgets.PlotMarkerBraille)
	zess.chart.SetDrawYAxisLabelFloat(false)
	zess.chart.SetData(scans)

	zess.view.AddItem(zess.chart, 0, 1, false)

	return &zess
}

func (zess *Zess) Run() {
}

func (zess *Zess) Update(scan *ZessScan) {
	if day == -1 {
		scans[0] = append(scans[0], 0)
		day = scan.ScanTime.YearDay()
	}

	scanDay := scan.ScanTime.YearDay()

	if scanDay == day {
		// Same day, increase the amount
		scans[0][len(scans[0])-1]++
	} else {
		// New day

		// Add the offset of days
		dayDifference := utils.GetDayDifference(day, scan.ScanTime) - 1
		for i := 0; i < dayDifference; i++ {
			scans[0] = utils.AddSliceElement(scans[0], 0)
		}

		scans[0] = utils.AddSliceElement(scans[0], 1)
		day = scanDay
	}

	zess.ScreenApp.execute(func() {
		zess.chart.SetData(scans)
	})
}
