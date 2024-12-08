package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/date"
)

// Season is the DTO for the season
type Season struct {
	ID      int32     `json:"id"`
	Name    string    `json:"name" validate:"required"`
	Start   date.Date `json:"start" validate:"required"`
	End     date.Date `json:"end" validate:"required"`
	Current bool      `json:"is_current"` // FIXME: This should have `required`. However when added the validation fails even though it's present
}

// SeasonDTO converts a sqlc.Season to a Season
func SeasonDTO(season sqlc.Season) *Season {
	return &Season{
		ID:      season.ID,
		Name:    season.Name,
		Start:   date.Date(season.Start.Time),
		End:     date.Date(season.End.Time),
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
		Start:   pgtype.Timestamp{Time: s.Start.ToTime(), Valid: true},
		End:     pgtype.Timestamp{Time: s.End.ToTime(), Valid: true},
		Current: s.Current,
	}
}

// UpdateParams converts a Season to sqlc.UpdateSeasonParams
func (s *Season) UpdateParams() sqlc.UpdateSeasonParams {
	return sqlc.UpdateSeasonParams{
		ID:      s.ID,
		Name:    s.Name,
		End:     pgtype.Timestamp{Time: s.End.ToTime(), Valid: true},
		Start:   pgtype.Timestamp{Time: s.Start.ToTime(), Valid: true},
		Current: s.Current,
	}
}
