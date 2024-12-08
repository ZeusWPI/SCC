// Package gamification provides all gamification related logic
package gamification

import (
	"context"
	"errors"

	"github.com/zeusWPI/scc/internal/pkg/db"
)

// Gamification represents a gamification instance
type Gamification struct {
	db *db.DB
}

// New creates a new gamification instance
func New(db *db.DB) *Gamification {
	return &Gamification{db: db}
}

// Update gets the current leaderboard from gamification
func (g *Gamification) Update() error {
	if _, err := g.db.Queries.DeleteGamificationAll(context.Background()); err != nil {
		return err
	}

	leaderboard, err := g.getLeaderboard()
	if err != nil {
		return err
	}

	var errs []error
	for _, item := range leaderboard {
		if _, err := g.db.Queries.CreateGamification(context.Background(), item.CreateParams()); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
