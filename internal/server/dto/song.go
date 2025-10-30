package dto

import "github.com/zeusWPI/scc/internal/database/model"

type Song struct {
	SpotifyID string `json:"spotify_id"`
}

func (s *Song) ToModel() *model.Song {
	return &model.Song{
		SpotifyID: s.SpotifyID,
	}
}
