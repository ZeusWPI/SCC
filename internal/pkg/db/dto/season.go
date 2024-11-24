package dto

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Season is the DTO for the season
type Season struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name" validate:"required"`
	Start   time.Time `json:"start" validate:"required"`
	End     time.Time `json:"end" validate:"required"`
	Current bool      `json:"is_current" validate:"required"`
}

// SeasonDTO converts a sqlc.Season to a Season
func SeasonDTO(season sqlc.Season) *Season {
	return &Season{
		ID:      season.ID,
		Name:    season.Name,
		Start:   season.Start,
		End:     season.End,
		Current: season.Current,
	}
}

// SeasonCmp compares two seasons
// Returns an int so it can be used in compare functions
func SeasonCmp(s1, s2 *Season) int {
	if s1.ID == s2.ID && s1.Name == s2.Name && s1.Start == s2.Start && s1.End == s2.End && s1.Current == s2.Current {
		return 0
	}

	return 1
}

// CreateParams converts a Season to sqlc.CreateSeasonParams
func (s *Season) CreateParams() sqlc.CreateSeasonParams {
	return sqlc.CreateSeasonParams{
		Name:    s.Name,
		Start:   s.Start,
		End:     s.End,
		Current: s.Current,
	}
}

// UpdateParams converts a Season to sqlc.UpdateSeasonParams
func (s *Season) UpdateParams() sqlc.UpdateSeasonParams {
	return sqlc.UpdateSeasonParams{
		ID:      s.ID,
		Name:    s.Name,
		Start:   s.Start,
		End:     s.End,
		Current: s.Current,
	}
}
