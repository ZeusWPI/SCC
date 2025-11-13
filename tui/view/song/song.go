// Package song provides the functions to draw an overview of the song integration
package song

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/lyrics"
	"github.com/zeusWPI/scc/tui/components/bar"
	"github.com/zeusWPI/scc/tui/components/stopwatch"
	"github.com/zeusWPI/scc/tui/view"
	"go.uber.org/zap"
)

const (
	previousAmount = 2
	upcomingAmount = 4
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
	song     model.Song
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

type Model struct {
	repo     repository.Song
	current  playing
	progress progression

	history           stat
	stats             []stat
	statsMonthly      []stat
	statAmount        int
	statAmountPlaying int

	width  int
	height int
}

// Interface compliance
var _ view.View = (*Model)(nil)

// Msg contains the data to update the gamification model
type Msg struct{}

type msgHistory struct {
	history stat
}

type msgStats struct {
	monthly bool
	stats   []stat
}

type msgPlaying struct {
	song   model.Song
	lyrics lyrics.Lyrics
}

type msgLyrics struct {
	song      model.Song
	playing   bool
	previous  []string
	current   string
	upcoming  []string
	startNext time.Time
}

// New initializes a new song model
func New(repo repository.Repository) view.View {
	return &Model{
		repo: *repo.NewSong(),
		progress: progression{
			stopwatch: stopwatch.New(),
			bar:       bar.New(sStatusBar),
		},
		stats:             make([]stat, 4),
		statsMonthly:      make([]stat, 4),
		statAmount:        config.GetDefaultInt("tui.view.song.stat_amount", 3),
		statAmountPlaying: config.GetDefaultInt("tui.view.song.stat_amount_playing", 5),
		width:             0,
		height:            0,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.progress.stopwatch.Init(),
		m.progress.bar.Init(),
	)
}

func (m *Model) Name() string {
	return "Song"
}

func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		if entry, ok := msg.Sizes[m.Name()]; ok {
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

		startTime := m.current.song.PlayedAt.Add(lyric.Duration) // Start time of the next lyric
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
			m.progress.stopwatch.Start(time.Since(m.current.song.PlayedAt)),
			m.progress.bar.Start(view.GetWidth(sStatusBar), time.Since(m.current.song.PlayedAt), time.Duration(m.current.song.DurationMS)*time.Millisecond),
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
	stopwatchNew, cmd := m.progress.stopwatch.Update(msg)
	m.progress.stopwatch = stopwatchNew.(stopwatch.Model)
	if cmd != nil {
		return m, cmd
	}

	// Apparently not, lets try the bar!
	barNew, cmd := m.progress.bar.Update(msg)
	m.progress.bar = barNew.(bar.Model)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
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
