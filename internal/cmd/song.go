package cmd

import (
	"github.com/zeusWPI/scc/internal/song"
)

func Song() error {
	return song.Init()
}
