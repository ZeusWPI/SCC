// Package song provides the functions to draw an overview of the song integration
package song

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/internal/pkg/lyrics"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/tui/components/bar"
	"github.com/zeusWPI/scc/tui/components/stopwatch"
	"github.com/zeusWPI/scc/tui/view"
	"go.uber.org/zap"
)

var (
	previousAmount = 5  // Amount of passed lyrics to show
	upcomingAmount = 12 // Amount of upcoming lyrics to show
)

type stat struct {
	title   string
	entries []statEntry
}

type statEntry struct {
	name   string
	amount int
}

type playing struct {
	song     dto.Song
	playing  bool
	lyrics   lyrics.Lyrics
	previous []string // Lyrics already sang
	current  string   // Current lyric
	upcoming []string // Lyrics that are coming up
}

type progression struct {
	stopwatch stopwatch.Model
	bar       bar.Model
}

// Model represents the view model for song
type Model struct {
	db           *db.DB
	current      playing
	progress     progression
	history      stat
	stats        []stat
	statsMonthly []stat
	width        int
	height       int
}

// Msg triggers a song data update
// Required for the view interface
type Msg struct{}

type msgHistory struct {
	history stat
}

type msgStats struct {
	monthly bool
	stats   []stat
}

type msgPlaying struct {
	song   dto.Song
	lyrics lyrics.Lyrics
}

type msgLyrics struct {
	song      dto.Song
	playing   bool
	previous  []string
	current   string
	upcoming  []string
	startNext time.Time
}

// New initializes a new song model
func New(db *db.DB) view.View {
	return &Model{
		db:           db,
		current:      playing{},
		progress:     progression{stopwatch: stopwatch.New(), bar: bar.New(sStatusBar)},
		stats:        make([]stat, 4),
		statsMonthly: make([]stat, 4),
	}
}

// Init starts the song view
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.progress.stopwatch.Init(),
		m.progress.bar.Init(),
	)
}

// Name returns the name of the view
func (m *Model) Name() string {
	return "Songs"
}

// Update updates the song view
func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		entry, ok := msg.Sizes[m.Name()]
		if ok {
			// Update all dependent styles
			m.width = entry.Width
			m.height = entry.Height

			m.updateStyles()
		}

		return m, nil

	case msgPlaying:
		// We're playing a song
		// Initialize the variables
		m.current.song = msg.song
		m.current.playing = true
		m.current.lyrics = msg.lyrics
		m.current.current = ""
		m.current.previous = []string{""}
		m.current.upcoming = []string{""}

		// The song might already been playing for some time
		// Let's go through the lyrics until we get to the current one
		lyric, ok := m.current.lyrics.Current()
		if !ok {
			// Shouldn't happen
			zap.S().Error("song: Unable to get current lyric in initialization phase: ", m.current.song.Title)
			m.current.playing = false
			return m, nil
		}

		startTime := m.current.song.CreatedAt.Add(lyric.Duration) // Start time of the next lyric
		for startTime.Before(time.Now()) {
			// This lyric is already finished, onto the next!
			lyric, ok := m.current.lyrics.Next()
			if !ok {
				// No more lyrics to display, the song is already finished
				m.current.playing = false
				return m, m.progress.stopwatch.Reset()
			}
			startTime = startTime.Add(lyric.Duration)
		}

		// We have the right lyric, let's get the previous and upcoming lyrics
		m.current.current = lyric.Text
		m.current.previous = lyricsToString(m.current.lyrics.Previous(previousAmount))
		m.current.upcoming = lyricsToString(m.current.lyrics.Upcoming(upcomingAmount))

		// Start the update loop
		return m, tea.Batch(
			updateLyrics(m.current, startTime),
			m.progress.stopwatch.Start(time.Since(m.current.song.CreatedAt)),
			m.progress.bar.Start(view.GetWidth(sStatusBar), time.Since(m.current.song.CreatedAt), time.Duration(m.current.song.DurationMS)*time.Millisecond),
		)

	case msgHistory:
		m.history = msg.history

		return m, nil

	case msgStats:
		if msg.monthly {
			// Monthly stats
			m.statsMonthly = msg.stats
			return m, nil
		}

		m.stats = msg.stats
		return m, nil

	case msgLyrics:
		// Check if it's still relevant
		if msg.song.ID != m.current.song.ID {
			// We already switched to a new song
			return m, nil
		}

		m.current.playing = msg.playing
		if !m.current.playing {
			// Song has finished. Reset variables
			return m, m.progress.stopwatch.Reset()
		}

		m.current.previous = msg.previous
		m.current.current = msg.current
		m.current.upcoming = msg.upcoming

		// Start the cmd to update the lyrics
		return m, updateLyrics(m.current, msg.startNext)
	}

	// Maybe a stopwatch message?
	var cmd tea.Cmd
	m.progress.stopwatch, cmd = m.progress.stopwatch.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	// Apparently not, lets try the bar!
	m.progress.bar, cmd = m.progress.bar.Update(msg)

	return m, cmd
}

// View draws the song view
func (m *Model) View() string {
	if m.current.playing {
		return m.viewPlaying()
	}

	return m.viewNotPlaying()
}

// GetUpdateDatas gets all update functions for the song view
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "update current song",
			View:     m,
			Update:   updateCurrentSong,
			Interval: config.GetDefaultInt("tui.view.song.interval_current_s", 5),
		},
		{
			Name:     "update history",
			View:     m,
			Update:   updateHistory,
			Interval: config.GetDefaultInt("tui.view.song.interval_history_s", 5),
		},
		{
			Name:     "monthly stats",
			View:     m,
			Update:   updateMonthlyStats,
			Interval: config.GetDefaultInt("tui.view.song.interval_monthly_stats_s", 300),
		},
		{
			Name:     "all time stats",
			View:     m,
			Update:   updateStats,
			Interval: config.GetDefaultInt("tui.view.song.interval_stats_s", 3600),
		},
	}
}

// updateCurrentSong checks if there's currently a song playing
func updateCurrentSong(view view.View) (tea.Msg, error) {
	m := view.(*Model)

	songs, err := m.db.Queries.GetLastSongFull(context.Background())
	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
		return nil, err
	}
	if len(songs) == 0 {
		return nil, nil
	}

	// Check if song is still playing
	if songs[0].CreatedAt.Time.Add(time.Duration(songs[0].DurationMs) * time.Millisecond).Before(time.Now()) {
		// Song is finished
		return nil, nil
	}

	if m.current.playing && songs[0].ID == m.current.song.ID {
		// Song is already set to current
		return nil, nil
	}

	// Convert sqlc song to a dto song
	song := *dto.SongDTOHistory(songs)

	return msgPlaying{song: song, lyrics: lyrics.New(song)}, nil
}

// updateHistory updates the recently played list
func updateHistory(view view.View) (tea.Msg, error) {
	m := view.(*Model)

	history, err := m.db.Queries.GetSongHistory(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	stat := stat{title: tStatHistory, entries: []statEntry{}}
	for _, h := range history {
		stat.entries = append(stat.entries, statEntry{name: h.Title, amount: int(h.PlayCount)})
	}

	return msgHistory{history: stat}, nil
}

// Update all monthly stats
func updateMonthlyStats(view view.View) (tea.Msg, error) {
	m := view.(*Model)

	songs, err := m.db.Queries.GetTopMonthlySongs(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	genres, err := m.db.Queries.GetTopMonthlyGenres(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	artists, err := m.db.Queries.GetTopMonthlyArtists(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	msg := msgStats{monthly: true, stats: []stat{}}

	// Songs
	s := stat{title: tStatSong, entries: []statEntry{}}
	for _, song := range songs {
		s.entries = append(s.entries, statEntry{name: song.Title, amount: int(song.PlayCount)})
	}
	msg.stats = append(msg.stats, s)

	// Genres
	s = stat{title: tStatGenre, entries: []statEntry{}}
	for _, genre := range genres {
		s.entries = append(s.entries, statEntry{name: genre.GenreName, amount: int(genre.TotalPlays)})
	}
	msg.stats = append(msg.stats, s)

	// Artists
	s = stat{title: tStatArtist, entries: []statEntry{}}
	for _, artist := range artists {
		s.entries = append(s.entries, statEntry{name: artist.ArtistName, amount: int(artist.TotalPlays)})
	}
	msg.stats = append(msg.stats, s)

	return msg, nil
}

// Update all stats
func updateStats(view view.View) (tea.Msg, error) {
	m := view.(*Model)

	songs, err := m.db.Queries.GetTopSongs(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	genres, err := m.db.Queries.GetTopGenres(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	artists, err := m.db.Queries.GetTopArtists(context.Background())
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	// Don't bother checking if anything has changed
	// A single extra refresh won't matter

	msg := msgStats{monthly: false, stats: []stat{}}

	// Songs
	s := stat{title: tStatSong, entries: []statEntry{}}
	for _, song := range songs {
		s.entries = append(s.entries, statEntry{name: song.Title, amount: int(song.PlayCount)})
	}
	msg.stats = append(msg.stats, s)

	// Genres
	s = stat{title: tStatGenre, entries: []statEntry{}}
	for _, genre := range genres {
		s.entries = append(s.entries, statEntry{name: genre.GenreName, amount: int(genre.TotalPlays)})
	}
	msg.stats = append(msg.stats, s)

	// Artists
	s = stat{title: tStatArtist, entries: []statEntry{}}
	for _, artist := range artists {
		s.entries = append(s.entries, statEntry{name: artist.ArtistName, amount: int(artist.TotalPlays)})
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
