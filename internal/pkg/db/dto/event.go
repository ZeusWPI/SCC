package dto

import (
	"bytes"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Event represents the DTO object for event
type Event struct {
	ID           int32
	Name         string
	Date         time.Time
	AcademicYear string
	Location     string
	Poster       []byte
}

// EventDTO converts a sqlc Event object to a DTO Event
func EventDTO(e sqlc.Event) *Event {
	return &Event{
		ID:           e.ID,
		Name:         e.Name,
		Date:         e.Date.Time,
		AcademicYear: e.AcademicYear,
		Location:     e.Location,
		Poster:       e.Poster,
	}
}

// Equal compares 2 events
func (e *Event) Equal(e2 Event) bool {
	return e.Name == e2.Name && e.Date.Equal(e2.Date) && e.AcademicYear == e2.AcademicYear && e.Location == e2.Location && bytes.Equal(e.Poster, e2.Poster)
}

// CreateParams converts a Event DTO to a sqlc CreateEventParams object
func (e *Event) CreateParams() sqlc.CreateEventParams {
	return sqlc.CreateEventParams{
		Name:         e.Name,
		Date:         pgtype.Timestamptz{Time: e.Date, Valid: true},
		AcademicYear: e.AcademicYear,
		Location:     e.Location,
		Poster:       e.Poster,
	}
}
