package song

import "github.com/zeusWPI/scc/pkg/lyrics"

func lyricsToString(lyrics []lyrics.Lyric) []string {
	text := make([]string, 0, len(lyrics))
	for _, lyric := range lyrics {
		text = append(text, lyric.Text)
	}
	return text
}
