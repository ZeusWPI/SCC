package lyrics

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/dto"
)

var re = regexp.MustCompile(`^\[(\d{2}):(\d{2})\.(\d{2})\]`)

// LRC represents synced lyrics
type LRC struct {
	song   dto.Song
	lyrics []Lyric
	i      int
}

func newLRC(song dto.Song) Lyrics {
	return &LRC{song: song, lyrics: parseLRC(song.Lyrics, time.Duration(song.DurationMS)), i: 0}
}

// GetSong returns the song associated to the lyrics
func (l *LRC) GetSong() dto.Song {
	return l.song
}

// Previous provides the previous `amount` of lyrics without affecting the current lyric
func (l *LRC) Previous(amount int) []Lyric {
	lyrics := make([]Lyric, 0, amount)

	for i := 1; i <= amount; i++ {
		if l.i-i-1 < 0 {
			break
		}

		lyrics = append([]Lyric{l.lyrics[l.i-i-1]}, lyrics...)
	}

	return lyrics
}

// Current provides the current lyric if any.
// If the song is finished the boolean is set to false
func (l *LRC) Current() (Lyric, bool) {
	if l.i >= len(l.lyrics) {
		return Lyric{}, false
	}

	return l.lyrics[l.i], true
}

// Next provides the next lyric if any.
// If the song is finished the boolean is set to false
func (l *LRC) Next() (Lyric, bool) {
	if l.i+1 >= len(l.lyrics) {
		return Lyric{}, false
	}

	l.i++
	return l.lyrics[l.i-1], true
}

// Upcoming provides the next `amount` lyrics without affecting the current lyric
func (l *LRC) Upcoming(amount int) []Lyric {
	lyrics := make([]Lyric, 0, amount)

	for i := 0; i < amount; i++ {
		if l.i+i >= len(l.lyrics) {
			break
		}

		lyrics = append(lyrics, l.lyrics[l.i+i])
	}

	return lyrics
}

// Progress shows the fraction of lyrics that have been used.
func (l *LRC) Progress() float64 {
	return float64(l.i) / float64(len(l.lyrics))
}

func parseLRC(text string, totalDuration time.Duration) []Lyric {
	lines := strings.Split(text, "\n")

	if len(lines) == 0 {
		return []Lyric{}
	}

	lyrics := make([]Lyric, 0, len(lines)+1) // + 1 for a start empty lyric
	var previousTimestamp time.Duration

	// Add first lyric (no text)
	lyrics = append(lyrics, Lyric{Text: ""})
	previousTimestamp = time.Duration(0)

	for i, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		// Duration part
		timeParts := re.FindStringSubmatch(parts[0])
		minutes, _ := strconv.Atoi(timeParts[1])
		seconds, _ := strconv.Atoi(timeParts[2])
		hundredths, _ := strconv.Atoi(timeParts[3])
		timestamp := time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second +
			time.Duration(hundredths)*10*time.Millisecond

		// Actual lyric
		lyric := parts[1]

		lyrics = append(lyrics, Lyric{Text: lyric})
		lyrics[i].Duration = timestamp - previousTimestamp
		previousTimestamp = timestamp
	}

	// Set duration of last lyric
	lyrics[len(lyrics)-1].Duration = totalDuration - previousTimestamp

	return lyrics
}
