package model

import (
	"time"

	"github.com/zeusWPI/scc/pkg/sqlc"
)

type Event struct {
	ID           int
	Name         string
	Date         time.Time
	AcademicYear string
	Location     string
	Poster       []byte
}

func EventModel(e sqlc.Event) *Event {
	return &Event{
		ID:           int(e.ID),
		Name:         e.Name,
		Date:         e.Date.Time,
		AcademicYear: e.AcademicYear,
		Location:     e.Location,
		Poster:       e.Poster,
	}
}
