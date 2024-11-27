package event

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
)

var layout = "Monday 02 January, 15:04 2006"

func (e *Event) getEvents() ([]dto.Event, error) {
	var events []dto.Event
	var errs []error
	c := colly.NewCollector()

	c.OnHTML(".event-tile", func(el *colly.HTMLElement) {
		event := dto.Event{}

		// Name
		event.Name = el.ChildText(".is-size-4-mobile")

		// Date & Location
		dateLoc := el.DOM.Find(".event-time-loc").Contents()
		dateLocStr := strings.Split(dateLoc.Text(), "  ")

		if len(dateLocStr) != 2 {
			errs = append(errs, fmt.Errorf("Event: Unable to scrape date and location %s", dateLocStr))
			return
		}

		// Location
		event.Location = strings.TrimSpace(dateLocStr[1])

		// Date
		date := strings.TrimSpace(dateLocStr[0])

		yearString := el.Attr("href")
		yearParts := strings.Split(yearString, "/")
		if len(yearParts) == 5 {
			rangeParts := strings.Split(yearParts[2], "-")
			if len(rangeParts) == 2 {
				dateWithYear := fmt.Sprintf("%s 20%s", date, rangeParts[0])
				parsedDate, err := time.Parse(layout, dateWithYear)
				if err == nil {
					event.AcademicYear = yearParts[2]
					event.Date = parsedDate
				}
			}
		}
		// Check for error
		if event.Date.IsZero() {
			errs = append(errs, fmt.Errorf("Event: Unable to parse date %s %s", date, yearString))
			return
		}

		events = append(events, event)
	})

	err := c.Visit(e.api)
	if err != nil {
		return nil, err
	}

	c.Wait()

	return events, errors.Join(errs...)
}
