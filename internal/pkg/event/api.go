package event

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"go.uber.org/zap"
)

// TODO: Look at https://github.com/PuerkitoBio/goquery

var layout = "Monday 02 January, 15:04 2006"

func (e *Event) getEvents() ([]dto.Event, error) {
	zap.S().Info("Events: Getting all events")

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

	err := c.Visit(e.website)
	if err != nil {
		return nil, err
	}

	c.Wait()

	return events, errors.Join(errs...)
}

func (e *Event) getPoster(event *dto.Event) error {
	zap.S().Info("Events: Getting poster for ", event.Name)

	yearParts := strings.Split(event.AcademicYear, "-")
	if len(yearParts) != 2 {
		return fmt.Errorf("Event: Academic year not properly formatted %s", event.AcademicYear)
	}

	yearStart, err := strconv.Atoi(yearParts[0])
	if err != nil {
		return fmt.Errorf("Event: Unable to convert academic year to int %v", yearParts)
	}
	yearEnd, err := strconv.Atoi(yearParts[1])
	if err != nil {
		return fmt.Errorf("Event: Unable to convert academic year to int %v", yearParts)
	}

	year := fmt.Sprintf("20%d-20%d", yearStart, yearEnd)

	url := fmt.Sprintf("%s/%s/%s/scc.png", e.websitePoster, year, event.Name)

	req := fiber.Get(url)
	status, body, errs := req.Bytes()
	if len(errs) != 0 {
		return errors.Join(append(errs, errors.New("Event: Download poster request failed"))...)
	}
	if status != fiber.StatusOK {
		// No poster for event
		return nil
	}

	event.Poster = body

	return nil
}
