package model

import "github.com/zeusWPI/scc/pkg/sqlc"

type Gamification struct {
	ID     int
	Name   string
	Score  int
	Avatar []byte
}

func GamificationModel(g sqlc.Gamification) *Gamification {
	return &Gamification{
		ID:     int(g.ID),
		Name:   g.Name,
		Score:  int(g.Score),
		Avatar: g.Avatar,
	}
}
