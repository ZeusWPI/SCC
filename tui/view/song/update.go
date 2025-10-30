package song

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/pkg/lyrics"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui/view"
)

// updateCurrentSong checks if there's currently a song playing
func updateCurrentSong(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	song, err := m.repo.GetLastPopulated(ctx)
	if err != nil {
		return nil, err
	}
	if song == nil {
		return nil, nil
	}

	// Check if song is still playing
	if song.PlayedAt.Add(time.Duration(song.DurationMS) * time.Millisecond).Before(time.Now()) {
		// Song is finished
		return nil, nil
	}

	if m.current.playing && song.ID == m.current.song.ID {
		// Song is already set to current
		return nil, nil
	}

	return msgPlaying{song: *song, lyrics: lyrics.New(*song)}, nil
}

// updateHistory updates the recently played list
func updateHistory(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	histories, err := m.repo.GetLast50(ctx)
	if err != nil {
		return nil, err
	}
	histories = utils.SliceGet(histories, statsAmount)

	stat := stat{title: tStatHistory, entries: []statEntry{}}
	for _, h := range histories {
		stat.entries = append(stat.entries, statEntry{name: h.Title, amount: h.PlayCount})
	}

	return msgHistory{history: stat}, nil
}

// Update all monthly stats
func updateMonthlyStats(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	songs, err := m.repo.GetTopSongsMonthly(ctx)
	if err != nil {
		return nil, err
	}
	songs = utils.SliceGet(songs, statsAmount)

	genres, err := m.repo.GetTopGenresMonthly(ctx)
	if err != nil {
		return nil, err
	}
	genres = utils.SliceGet(genres, statsAmount)

	artists, err := m.repo.GetTopArtistsMonthly(ctx)
	if err != nil {
		return nil, err
	}
	artists = utils.SliceGet(artists, statsAmount)

	msg := msgStats{monthly: true, stats: []stat{}}

	// Songs
	s := stat{title: tStatSong, entries: []statEntry{}}
	for _, song := range songs {
		s.entries = append(s.entries, statEntry{name: song.Title, amount: song.PlayCount})
	}
	msg.stats = append(msg.stats, s)

	// Genres
	s = stat{title: tStatGenre, entries: []statEntry{}}
	for _, genre := range genres {
		s.entries = append(s.entries, statEntry{name: genre.Genre, amount: genre.PlayCount})
	}
	msg.stats = append(msg.stats, s)

	// Artists
	s = stat{title: tStatArtist, entries: []statEntry{}}
	for _, artist := range artists {
		s.entries = append(s.entries, statEntry{name: artist.Name, amount: artist.PlayCount})
	}
	msg.stats = append(msg.stats, s)

	return msg, nil
}

// Update all stats
func updateStats(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	songs, err := m.repo.GetTopSongs(ctx)
	if err != nil {
		return nil, err
	}
	songs = utils.SliceGet(songs, statsAmount)

	genres, err := m.repo.GetTopGenres(ctx)
	if err != nil {
		return nil, err
	}
	genres = utils.SliceGet(genres, statsAmount)

	artists, err := m.repo.GetTopArtists(ctx)
	if err != nil {
		return nil, err
	}
	artists = utils.SliceGet(artists, statsAmount)

	// Don't bother checking if anything has changed
	// A single extra refresh won't matter

	msg := msgStats{monthly: false, stats: []stat{}}

	// Songs
	s := stat{title: tStatSong, entries: []statEntry{}}
	for _, song := range songs {
		s.entries = append(s.entries, statEntry{name: song.Title, amount: song.PlayCount})
	}
	msg.stats = append(msg.stats, s)

	// Genres
	s = stat{title: tStatGenre, entries: []statEntry{}}
	for _, genre := range genres {
		s.entries = append(s.entries, statEntry{name: genre.Genre, amount: genre.PlayCount})
	}
	msg.stats = append(msg.stats, s)

	// Artists
	s = stat{title: tStatArtist, entries: []statEntry{}}
	for _, artist := range artists {
		s.entries = append(s.entries, statEntry{name: artist.Name, amount: artist.PlayCount})
	}
	msg.stats = append(msg.stats, s)

	return msg, nil
}

// Update the current lyric
func updateLyrics(song playing, start time.Time) tea.Cmd {
	// How long do we need to wait until we can update the lyric?
	timeout := time.Duration(0)
	now := time.Now()
	if start.After(now) {
		timeout = start.Sub(now)
	}

	return tea.Tick(timeout, func(_ time.Time) tea.Msg {
		// Next lyric
		lyric, ok := song.lyrics.Next()
		if !ok {
			// Song finished
			return msgLyrics{song: song.song, playing: false} // Values in the other fields are not looked at when the song is finished
		}

		previous := song.lyrics.Previous(previousAmount)
		upcoming := song.lyrics.Upcoming(upcomingAmount)

		end := start.Add(lyric.Duration)

		return msgLyrics{
			song:      song.song,
			playing:   true,
			previous:  lyricsToString(previous),
			current:   lyric.Text,
			upcoming:  lyricsToString(upcoming),
			startNext: end,
		}
	})
}
