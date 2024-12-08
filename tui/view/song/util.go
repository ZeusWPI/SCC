package song

import (
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/internal/pkg/lyrics"
)

func equalTopSongs(s1 []topStat, s2 []sqlc.GetTopSongsRow) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, s := range s1 {
		if s.name != s2[i].Title || s.amount != int(s2[i].PlayCount) {
			return false
		}
	}

	return true
}

func topStatSqlcSong(songs []sqlc.GetTopSongsRow) []topStat {
	topstats := make([]topStat, 0, len(songs))
	for _, s := range songs {
		topstats = append(topstats, topStat{name: s.Title, amount: int(s.PlayCount)})
	}
	return topstats
}

func equalTopGenres(s1 []topStat, s2 []sqlc.GetTopGenresRow) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, s := range s1 {
		if s.name != s2[i].GenreName || s.amount != int(s2[i].TotalPlays) {
			return false
		}
	}

	return true
}

func topStatSqlcGenre(songs []sqlc.GetTopGenresRow) []topStat {
	topstats := make([]topStat, 0, len(songs))
	for _, s := range songs {
		topstats = append(topstats, topStat{name: s.GenreName, amount: int(s.TotalPlays)})
	}
	return topstats
}

func equalTopArtists(s1 []topStat, s2 []sqlc.GetTopArtistsRow) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, s := range s1 {
		if s.name != s2[i].ArtistName || s.amount != int(s2[i].TotalPlays) {
			return false
		}
	}

	return true
}

func topStatSqlcArtist(songs []sqlc.GetTopArtistsRow) []topStat {
	topstats := make([]topStat, 0, len(songs))
	for _, s := range songs {
		topstats = append(topstats, topStat{name: s.ArtistName, amount: int(s.TotalPlays)})
	}
	return topstats
}

func lyricsToString(lyrics []lyrics.Lyric) []string {
	text := make([]string, 0, len(lyrics))
	for _, lyric := range lyrics {
		text = append(text, lyric.Text)
	}
	return text
}
