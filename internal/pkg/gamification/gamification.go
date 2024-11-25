// Package gamification provides all gamification related logic
package gamification

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
)

// Gamification represents a gamification instance
type Gamification struct {
	db  *db.DB
	api string
}

// New creates a new gamification instance
func New(db *db.DB) *Gamification {
	api := config.GetDefaultString("gamification.api", "https://gamification.zeus.gent")

	return &Gamification{db: db, api: api}
}

// Update gets the current leaderboard from gamification
func (g *Gamification) Update() error {
	leaderboard, err := g.db.Queries.GetAllGamification(context.Background())
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	gamLeaderboard, err := g.getLeaderboard()
	if err != nil {
		return err
	}

	// Delete old
	var errs []error
	var errsOS []error // OS specific errors.
	for _, l := range leaderboard {
		// Remove picture
		if err := os.Remove(l.Avatar); err != nil && err != os.ErrNotExist {
			errsOS = append(errsOS, err)
		}

		// Remove DB entry
		if _, err = g.db.Queries.DeleteGamification(context.Background(), l.ID); err != nil {
			errs = append(errs, err)
		}
	}

	// Don't quit if the only error(s) are os related
	if errs != nil {
		return errors.Join(append(errs, errsOS...)...)
	}

	// Insert new ones
	for _, gamL := range *gamLeaderboard {
		location, err := downloadAvatar(*gamL)
		if err != nil {
			errs = append(errs, err)
		}
		gamL.Avatar = location

		if _, err = g.db.Queries.CreateGamification(context.Background(), gamL.CreateParams()); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
