// Package song contains the screen displaying the song view
package song

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/tui/screen"
	"github.com/zeusWPI/scc/tui/view"
	"github.com/zeusWPI/scc/tui/view/song"
)

// Song represents the song screen
type Song struct {
	db   *db.DB
	song view.View

	width  int
	height int
}

// New creates a new song screen
func New(db *db.DB) screen.Screen {
	return &Song{db: db, song: song.New(db), width: 0, height: 0}
}

// Init initializes the song screen
func (s *Song) Init() tea.Cmd {
	return s.song.Init()
}

// Update updates the song screen and all its views
func (s *Song) Update(msg tea.Msg) (screen.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height

		sSong = sSong.Width(s.width - view.GetOuterWidth(sSong)).Height(s.height - sSong.GetVerticalFrameSize() - sSong.GetVerticalPadding())

		return s, s.GetSizeMsg
	}

	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

	s.song, cmd = s.song.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

// View returns the song screen view
func (s *Song) View() string {
	if s.width == 0 || s.height == 0 {
		return "Initializing..."
	}

	view := s.song.View()
	view = sSong.Render(view)

	return view
}

// GetUpdateViews returns all the update functions for the song screen
func (s *Song) GetUpdateViews() []view.UpdateData {
	return s.song.GetUpdateDatas()
}

// GetSizeMsg returns a message for the views informing them about their width and height
func (s *Song) GetSizeMsg() tea.Msg {
	sizes := make(map[string]view.Size)

	songW := sSong.GetWidth()
	songH := sSong.GetHeight()
	sizes[s.song.Name()] = view.Size{Width: songW, Height: songH}

	return view.MsgSize{Sizes: sizes}
}
