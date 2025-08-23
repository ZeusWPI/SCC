// Package gamification provides all gamification related logic
package gamification

import (
	"context"

	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
)

type Gamification struct {
	repo repository.Repository
	gam  repository.Gamification
	url  string
}

func New(repo repository.Repository) *Gamification {
	return &Gamification{
		repo: repo,
		gam:  *repo.NewGamification(),
		url:  config.GetDefaultString("backend.gamification.url", "https://gamification.zeus.gent"),
	}
}

func (g *Gamification) Update(ctx context.Context) error {
	leaderboard, err := g.getLeaderboard(ctx)
	if err != nil {
		return err
	}

	return g.repo.WithRollback(ctx, func(ctx context.Context) error {
		if err := g.gam.DeleteAll(ctx); err != nil {
			return err
		}

		for _, l := range leaderboard {
			if err := g.gam.Create(ctx, &l); err != nil {
				return err
			}
		}

		return nil
	})
}
