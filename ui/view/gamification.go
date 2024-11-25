package view

import (
	"context"
	"database/sql"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/config"
)

var width = 20

var (
	base        = lipgloss.NewStyle()
	columnStyle = base.MarginLeft(1)
	nameBase    = base.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(lipgloss.Color("#383838")).Width(width).Align(lipgloss.Center)
	nameStyles  = []lipgloss.Style{
		nameBase.Foreground(lipgloss.Color("#FFD700")),
		nameBase.Foreground(lipgloss.Color("#FF7F00")),
		nameBase.Foreground(lipgloss.Color("#CD7F32")),
		nameBase,
	}
	scoreStyle = base.Width(width).Align(lipgloss.Center)
)

// GamificationModel represents the view model for gamification
type GamificationModel struct {
	db          *db.DB
	leaderboard []gamificationItem
}

type gamificationItem struct {
	image image.Image
	item  dto.Gamification
}

// GamificationMsg contains the data to update the gamification model
type GamificationMsg struct {
	leaderboard []gamificationItem
}

// NewGamificationModel initializes a new gamification model
func NewGamificationModel(db *db.DB) View {
	return &GamificationModel{db: db, leaderboard: []gamificationItem{}}
}

// Init starts the gamification view
func (g *GamificationModel) Init() tea.Cmd {
	return nil
}

// Update updates the gamification view
func (g *GamificationModel) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case GamificationMsg:
		g.leaderboard = msg.leaderboard
	}

	return g, nil
}

// View draws the gamification view
func (g *GamificationModel) View() string {
	columns := make([]string, 0, len(g.leaderboard))

	for i, item := range g.leaderboard {
		user := lipgloss.JoinVertical(lipgloss.Left,
			nameStyles[i%len(nameStyles)].Render(fmt.Sprintf("%d. %s", i+1, item.item.Name)),
			scoreStyle.Render(strconv.Itoa(int(item.item.Score))),
		)

		column := lipgloss.JoinVertical(lipgloss.Left, gamificationToString(width, item.image), user)
		columns = append(columns, columnStyle.Render(column))
	}

	list := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	return list
}

// GetUpdateDatas get all update functions for the gamification view
func (g *GamificationModel) GetUpdateDatas() []UpdateData {
	return []UpdateData{
		{
			Name:     "gamification leaderboard",
			View:     g,
			Update:   updateLeaderboard,
			Interval: config.GetDefaultInt("tui.gamification.interval_s", 3600),
		},
	}
}

func updateLeaderboard(db *db.DB, view View) (tea.Msg, error) {
	g := view.(*GamificationModel)

	gams, err := db.Queries.GetAllGamificationByScore(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	// Check if both leaderboards are equal
	equal := false
	if len(g.leaderboard) == len(gams) {
		equal = true
		for i, l := range g.leaderboard {
			if !l.item.Equal(*dto.GamificationDTO(gams[i])) {
				equal = false
				break
			}
		}
	}

	if equal {
		return nil, nil
	}

	msg := GamificationMsg{leaderboard: []gamificationItem{}}
	for _, gam := range gams {
		if gam.Avatar == "" {
			// No avatar downloaded
			msg.leaderboard = append(msg.leaderboard, gamificationItem{image: nil, item: *dto.GamificationDTO(gam)})
			continue
		}

		file, err := os.Open(filepath.Clean(gam.Avatar))
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = file.Close()
		}()

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, err
		}

		msg.leaderboard = append(msg.leaderboard, gamificationItem{image: img, item: *dto.GamificationDTO(gam)})
	}

	return msg, nil
}

func gamificationToString(width int, img image.Image) string {
	img = imaging.Resize(img, width, 0, imaging.Lanczos)
	b := img.Bounds()
	imageWidth := b.Max.X
	h := b.Max.Y
	str := strings.Builder{}

	for heightCounter := 0; heightCounter < h; heightCounter += 2 {
		for x := imageWidth; x < width; x += 2 {
			str.WriteString(" ")
		}

		for x := 0; x < imageWidth; x++ {
			c1, _ := colorful.MakeColor(img.At(x, heightCounter))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
			color2 := lipgloss.Color(c2.Hex())
			str.WriteString(lipgloss.NewStyle().Foreground(color1).
				Background(color2).Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
}
