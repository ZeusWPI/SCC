package dto

import (
	"bytes"

	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Gamification represents the DTO object for gamification
type Gamification struct {
	ID     int64  `json:"id"`
	Name   string `json:"github_name"`
	Score  int64  `json:"score"`
	Avatar []byte `json:"avatar"`
}

// GamificationDTO converts a sqlc Gamification object to a DTO gamification
func GamificationDTO(gam sqlc.Gamification) *Gamification {
	return &Gamification{
		ID:     gam.ID,
		Name:   gam.Name,
		Score:  gam.Score,
		Avatar: gam.Avatar,
	}
}

// Equal compares 2 Gamification objects for equality
func (g *Gamification) Equal(g2 Gamification) bool {
	return g.Name == g2.Name && g.Score == g2.Score && bytes.Equal(g.Avatar, g2.Avatar)
}

// CreateParams converts a Gamification DTO to a sqlc CreateGamificationParams object
func (g *Gamification) CreateParams() sqlc.CreateGamificationParams {
	return sqlc.CreateGamificationParams{
		Name:   g.Name,
		Score:  g.Score,
		Avatar: g.Avatar,
	}
}

// UpdateScoreParams converts a Gamification DTO to a sqlc UpdateScoreParams object
func (g *Gamification) UpdateScoreParams() sqlc.UpdateGamificationScoreParams {
	return sqlc.UpdateGamificationScoreParams{
		ID:    g.ID,
		Score: g.Score,
	}
}
