// Package cammie returns the screen containing the cammie messages and other stats
package cammie

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/ui/screen"
	"github.com/zeusWPI/scc/ui/view"
	"github.com/zeusWPI/scc/ui/view/event"
	"github.com/zeusWPI/scc/ui/view/gamification"
	"github.com/zeusWPI/scc/ui/view/message"
	"github.com/zeusWPI/scc/ui/view/tap"
	"github.com/zeusWPI/scc/ui/view/zess"
)

// Cammie represents the cammie screen
type Cammie struct {
	db       *db.DB
	messages view.View
	top      []view.View
	bottom   view.View
	indexTop int
	width    int
	height   int
}

// Message to update the bottomIndex
type msgIndex struct {
	indexBottom int
}

// New creates a new cammie screen
func New(db *db.DB) screen.Screen {
	messages := message.NewModel(db)
	top := event.NewModel(db)
	bottom := []view.View{gamification.NewModel(db), tap.NewModel(db), zess.NewModel(db)}
	return &Cammie{db: db, messages: messages, bottom: top, top: bottom, indexTop: 2, width: 0, height: 0}
}

// Init initializes the cammie screen
func (c *Cammie) Init() tea.Cmd {
	cmds := []tea.Cmd{updateBottomIndex(*c), c.messages.Init(), c.bottom.Init()}
	for _, view := range c.top {
		cmds = append(cmds, view.Init())
	}

	return tea.Batch(cmds...)
}

// Update updates the cammie screen
func (c *Cammie) Update(msg tea.Msg) (screen.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height

		sMsg = sMsg.Width(c.width/2 - sMsg.GetHorizontalFrameSize() - sMsg.GetHorizontalPadding()).Height(c.height - sMsg.GetVerticalFrameSize() - sMsg.GetVerticalPadding())
		sTop = sTop.Width(c.width/2 - sTop.GetHorizontalFrameSize()).Height(c.height/2 - sTop.GetVerticalFrameSize())
		sBottom = sBottom.Width(c.width/2 - sBottom.GetHorizontalFrameSize()).Height(c.height/2 - sBottom.GetVerticalFrameSize())

		return c, updateSize(*c)
	case msgIndex:
		c.indexTop = msg.indexBottom

		return c, updateBottomIndex(*c)
	}

	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd

	c.messages, cmd = c.messages.Update(msg)
	cmds = append(cmds, cmd)

	c.bottom, cmd = c.bottom.Update(msg)
	cmds = append(cmds, cmd)

	for i, view := range c.top {
		c.top[i], cmd = view.Update(msg)
		cmds = append(cmds, cmd)
	}

	return c, tea.Batch(cmds...)
}

// View returns the cammie screen view
func (c *Cammie) View() string {
	if c.width == 0 || c.height == 0 {
		return "Initialzing..."
	}

	// Render messages
	messages := sMsg.Render(c.messages.View())

	// Render top
	// Render tabs
	var topTabs []string
	for i, view := range c.top {
		if i == c.indexTop {
			topTabs = append(topTabs, sActiveTab.Render(view.Name()))
		} else {
			topTabs = append(topTabs, sTabNormal.Render(view.Name()))
		}
	}
	topTab := lipgloss.JoinHorizontal(lipgloss.Bottom, topTabs...)
	topTabsLine := sTabNormal.Render(strings.Repeat(" ", max(0, (c.width/2)-lipgloss.Width(topTab)-10)))
	topTab = lipgloss.JoinHorizontal(lipgloss.Bottom, topTab, topTabsLine)

	// Render top view
	top := lipgloss.JoinVertical(lipgloss.Left, topTab, c.top[c.indexTop].View())
	top = sTop.Render(top)

	// Render bottom
	bottom := sBottom.Render(c.bottom.View())

	// Combine top and bottom
	right := lipgloss.JoinVertical(lipgloss.Left, top, bottom)

	// Combine left and right
	view := lipgloss.JoinHorizontal(lipgloss.Top, messages, right)

	return view
}

// GetUpdateViews returns all the update functions for the cammie screen
func (c *Cammie) GetUpdateViews() []view.UpdateData {
	updates := make([]view.UpdateData, 0)

	updates = append(updates, c.messages.GetUpdateDatas()...)
	updates = append(updates, c.bottom.GetUpdateDatas()...)

	for _, view := range c.top {
		updates = append(updates, view.GetUpdateDatas()...)
	}

	return updates
}

func updateBottomIndex(cammie Cammie) tea.Cmd {
	timeout := time.Duration(config.GetDefaultInt("tui.screen.cammie_interval_change_s", 300) * int(time.Second))
	return tea.Tick(timeout, func(_ time.Time) tea.Msg {
		newIndex := (cammie.indexTop + 1) % len(cammie.top)

		return msgIndex{indexBottom: newIndex}
	})
}

func updateSize(cammie Cammie) tea.Cmd {
	return func() tea.Msg {
		msg := view.MsgSize{
			Sizes: make(map[string]view.Size),
		}

		msgW := sMsg.GetWidth()
		msgH := sMsg.GetHeight()
		msg.Sizes[cammie.messages.Name()] = view.Size{Width: msgW, Height: msgH}

		bottomW := sBottom.GetWidth()
		bottomH := sBottom.GetHeight()
		msg.Sizes[cammie.bottom.Name()] = view.Size{Width: bottomW, Height: bottomH}

		for _, top := range cammie.top {
			topW := sTop.GetWidth()
			topH := sTop.GetHeight()
			msg.Sizes[top.Name()] = view.Size{Width: topW, Height: topH}
		}

		return msg
	}
}
