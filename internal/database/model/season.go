package model

import (
	"github.com/zeusWPI/scc/pkg/date"
	"github.com/zeusWPI/scc/pkg/sqlc"
)

type Season struct {
	ID      int
	Name    string
	Start   date.Date
	End     date.Date
	Current bool
}

func SeasonModel(s sqlc.Season) *Season {
	return &Season{
		ID:      int(s.ID),
		Name:    s.Name,
		Start:   date.Date(s.Start.Time),
		End:     date.Date(s.End.Time),
		Current: s.Current,
	}
}
