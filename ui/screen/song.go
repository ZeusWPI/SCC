package screen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/ui/view"
	"github.com/zeusWPI/scc/ui/view/song"
)

// Song represents the song screen
type Song struct {
	db   *db.DB
	song view.View
}

// NewSong creates a new song screen
func NewSong(db *db.DB) Screen {
	return &Song{db: db, song: song.NewModel(db)}
}

// Init initializes the song screen
func (s *Song) Init() tea.Cmd {
	return s.song.Init()
}

// Update updates the song screen
func (s *Song) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	default:
		song, cmd := s.song.Update(msg)
		s.song = song

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return s, tea.Batch(cmds...)
}

// View returns the song screen view
func (s *Song) View() string {
	return s.song.View()
}

// GetUpdateViews returns all the update functions for the song screen
func (s *Song) GetUpdateViews() []view.UpdateData {
	return s.song.GetUpdateDatas()
}
