package lyrics

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
)

type Instrumental struct {
	song   model.Song
	lyrics []Lyric
	i      int
}

func newInstrumental(song model.Song) Lyrics {
	return &Instrumental{song: song, lyrics: generateInstrumental(time.Duration(song.DurationMS) * time.Millisecond), i: 0}
}

func (i *Instrumental) GetSong() model.Song {
	return i.song
}

func (i *Instrumental) Previous(amount int) []Lyric {
	lyrics := make([]Lyric, 0, amount)

	for j := 1; j <= amount; j++ {
		if i.i-j-1 < 0 {
			break
		}

		lyrics = append([]Lyric{i.lyrics[i.i-j-1]}, lyrics...)
	}

	return lyrics
}

func (i *Instrumental) Current() (Lyric, bool) {
	if i.i >= len(i.lyrics) {
		return Lyric{}, false
	}

	return i.lyrics[i.i], true
}

func (i *Instrumental) Next() (Lyric, bool) {
	if i.i+1 >= len(i.lyrics) {
		return Lyric{}, false
	}

	i.i++
	return i.lyrics[i.i-1], true
}

func (i *Instrumental) Upcoming(amount int) []Lyric {
	lyrics := make([]Lyric, 0, amount)

	for j := range amount {
		if i.i+j >= len(i.lyrics) {
			break
		}

		lyrics = append(lyrics, i.lyrics[i.i+j])
	}

	return lyrics
}

func (i *Instrumental) Progress() float64 {
	return float64(i.i) / float64(len(i.lyrics))
}

func generateInstrumental(dur time.Duration) []Lyric {
	// Get all instruments with their frequency
	freqs := []instrument{}
	for _, instr := range instruments {
		for range instr.frequency {
			freqs = append(freqs, instr)
		}
	}

	lyrics := []Lyric{}
	// Split up song in segments between 5 and 15 seconds
	currentDur := time.Duration(0)
	for currentDur < dur {
		// Get a random instrument
		instr := freqs[rand.IntN(len(freqs))]

		// Get a random duration
		randomDur := time.Duration(rand.IntN(10)+5) * time.Second
		currentDur += randomDur
		if currentDur >= dur {
			randomDur -= (currentDur - dur)
		}

		// Get the lyrics
		lyrics = append(lyrics, instr.generate(randomDur)...)
	}

	return lyrics
}

// Instruments

type instrument struct {
	frequency int      // Odds of it occuring
	name      string   // Name of the instrument
	sounds    []string // Different ways it could sound
}

// generate creates lyrics for a specific instrument
func (i instrument) generate(dur time.Duration) []Lyric {
	lyrics := []Lyric{}

	// Same logic as in `generateInstrumental` except that the segments are between 1 and 4 seconds

	// Add the start lyric
	text := fmt.Sprintf(openings[rand.IntN(len(openings))], i.name)
	randomDur := time.Duration(1) * time.Second
	lyrics = append(lyrics, Lyric{Text: text, Duration: randomDur})

	currentDur := randomDur
	for currentDur < dur {
		// Get a random text
		textLength := rand.IntN(5) + 1
		var text string
		for range textLength {
			text += i.sounds[rand.IntN(len(i.sounds))] + " "
		}

		// Get a random duration
		randomDur := time.Duration(rand.IntN(3)+1) * time.Second
		currentDur += randomDur
		if currentDur >= dur {
			// Last lyric, add a newline
			text += "\n"
			randomDur -= (currentDur - dur)
		}

		lyrics = append(lyrics, Lyric{Text: text, Duration: randomDur})
	}

	return lyrics
}

// All instruments to choose from
var instruments = []instrument{
	{frequency: 9, name: "Piano", sounds: []string{"plink", "plonk", "pling", "clink", "clang"}},
	{frequency: 7, name: "Drums", sounds: []string{"boom", "ba-dum", "thwack", "tshh", "bop"}},
	{frequency: 5, name: "Electric Guitar", sounds: []string{"wah", "zzzzzz", "twang", "vrrr", "brrraang"}},
	{frequency: 3, name: "Theremin", sounds: []string{"wooOOOooo", "weeeee", "ooooo", "waaAAaah", "hummmm"}},
	{frequency: 6, name: "Flute", sounds: []string{"toot", "fweee", "trillll", "pip", "peep"}},
	{frequency: 4, name: "Accordion", sounds: []string{"wheeze", "honk", "phwoo", "eep", "squawk"}},
	{frequency: 8, name: "Violin", sounds: []string{"screee", "swish", "zing", "vwee", "mreee"}},
	{frequency: 5, name: "Saxophone", sounds: []string{"saxxy", "bwoop", "dooo", "reebop", "honka"}},
	{frequency: 2, name: "Kazoo", sounds: []string{"bzzzzz", "zwip", "vwoo", "brrr", "zzzzrrt"}},
	{frequency: 9, name: "Trumpet", sounds: []string{"brrraaap", "toot", "doo-doo", "wah-wah", "parp"}},
	{frequency: 3, name: "Cowbell", sounds: []string{"clang", "clong", "ding", "donk", "bonk"}},
	{frequency: 2, name: "Bagpipes", sounds: []string{"drone", "hrooo", "whine", "skree", "rrrrrrr"}},
	{frequency: 6, name: "Triangle", sounds: []string{"ting", "tang", "ding", "dling", "plink"}},
	{frequency: 1, name: "Didgeridoo", sounds: []string{"whooooo", "womp", "drrrrrr", "brrrrr", "hummmmm"}},
	{frequency: 4, name: "Bongos", sounds: []string{"pop", "tap", "dum", "ba-dum", "bop"}},
	{frequency: 7, name: "Harp", sounds: []string{"plink", "tinkle", "zling", "glint", "ding"}},
	{frequency: 5, name: "Maracas", sounds: []string{"sh-sh-sh", "shaka-shaka", "chick", "rattle", "tktktk"}},
	{frequency: 3, name: "Tuba", sounds: []string{"oompah", "bruhm", "whoom", "booo", "phrum"}},
	{frequency: 1, name: "Banjo", sounds: []string{"twang", "plink", "brrrring", "plunk", "doink"}},
	{frequency: 2, name: "Synthesizer", sounds: []string{"beep-boop", "vwee-vwee", "zorp", "waah", "ding"}},
	{frequency: 6, name: "Xylophone", sounds: []string{"ding", "dunk", "plink-plonk", "tok", "tink"}},
	{frequency: 2, name: "Hurdy-Gurdy", sounds: []string{"whirr", "drone", "skreee", "buzz", "rrrrrng"}},
	{frequency: 4, name: "Harmonica", sounds: []string{"wheeze", "toot", "hoo", "blow", "brrrr"}},
	{frequency: 3, name: "Slide Whistle", sounds: []string{"whoooop", "wheeee", "wooo", "boooo", "zwip"}},
	{frequency: 5, name: "Tambourine", sounds: []string{"jingle", "shake-shake", "tshh", "tinkle", "ting-ting"}},
	{frequency: 2, name: "Ocarina", sounds: []string{"woo", "fweee", "doodle", "pip-pip", "toot"}},
	{frequency: 8, name: "Acoustic Guitar", sounds: []string{"strum", "plang", "twang", "zing", "thrum"}},
	{frequency: 1, name: "Sousaphone", sounds: []string{"toot", "boop", "pah-pah", "oompah", "pwaaah"}},
	{frequency: 3, name: "Castanets", sounds: []string{"clack", "click", "clap", "tick", "tack"}},
	{frequency: 7, name: "Synth Drum", sounds: []string{"pshh", "bzzt", "bip", "tsh", "zorp"}},
	{frequency: 2, name: "Bag of Gravel", sounds: []string{"crunch", "scrape", "sh-sh", "clatter", "grrnk"}},
	{frequency: 5, name: "Steel Drum", sounds: []string{"pong", "ding", "donk", "bop", "ting"}},
	{frequency: 4, name: "Mouth Harp", sounds: []string{"boing", "thwong", "zzzt", "doyoyoy", "wobble"}},
	{frequency: 2, name: "Rainstick", sounds: []string{"shhhhh", "rrrrrr", "drip-drop", "fwssh", "ssss"}},
	{frequency: 1, name: "Toy Piano", sounds: []string{"plink", "tink-tink", "chime", "plinkity", "dink"}},
	{frequency: 3, name: "Jaw Harp", sounds: []string{"twang", "boing", "doink", "womp", "zzzrrrt"}},
	{frequency: 4, name: "Bicycle Horn", sounds: []string{"honk", "meeep", "awoooga", "brrrt", "bop-bop"}},
	{frequency: 2, name: "Glass Harp", sounds: []string{"wheee", "zing", "woo", "glint", "oooo"}},
	{frequency: 6, name: "Claves", sounds: []string{"clack", "click", "clonk", "tak", "tok"}},
	{frequency: 3, name: "Rubber Band", sounds: []string{"twang", "ping", "boing", "zing", "snap"}},
	{frequency: 2, name: "Paper Comb", sounds: []string{"buzz", "brrr", "wobble", "zzzt", "drone"}},
	{frequency: 1, name: "Duck Call", sounds: []string{"quack", "wak-wak", "honk", "waak", "weeek"}},
	{frequency: 5, name: "Handbells", sounds: []string{"ding", "dong", "chime", "tinkle", "bong"}},
	{frequency: 4, name: "Foghorn", sounds: []string{"MOOOO", "hoooonk", "BWAAAA", "WOOOO", "brrrmmm"}},
	{frequency: 7, name: "Cello", sounds: []string{"mmmm", "vmmm", "vroom", "dronnn", "zoomm"}},
	{frequency: 6, name: "Clarinet", sounds: []string{"toot", "wooo", "hmmm", "dee-dee", "reeee"}},
	{frequency: 8, name: "Oboe", sounds: []string{"hweee", "hee", "whee", "ooooo", "reee"}},
	{frequency: 5, name: "French Horn", sounds: []string{"vooom", "phoo", "bwoo", "vuuum", "whooo"}},
	{frequency: 6, name: "Bassoon", sounds: []string{"boo", "brrrr", "phrum", "wuuu", "vrrr"}},
	{frequency: 8, name: "Timpani", sounds: []string{"boom", "dum", "rumble", "thud", "pum"}},
	{frequency: 7, name: "Double Bass", sounds: []string{"vrumm", "dumm", "boooom", "grumm", "zzzooom"}},
	{frequency: 9, name: "Trumpet", sounds: []string{"brrrmp", "doo-doo", "toot", "baap", "dah-dah"}},
	{frequency: 6, name: "Trombone", sounds: []string{"wah-wah", "dooo", "wooo", "bwaaah", "vroom"}},
	{frequency: 4, name: "Harp", sounds: []string{"plink", "strum", "zinnnng", "twang", "gliss"}},
	{frequency: 6, name: "Piccolo", sounds: []string{"peep", "tweet", "fweep", "weeet", "pweep"}},
	{frequency: 7, name: "Bass Drum", sounds: []string{"boom", "thud", "pum", "dum", "bomp"}},
	{frequency: 5, name: "Snare Drum", sounds: []string{"rat-a-tat", "tsh", "tktktk", "snap", "crack"}},
	{frequency: 7, name: "Tuba", sounds: []string{"pah-pah", "brumm", "booom", "ooooh", "vrooo"}},
	{frequency: 6, name: "Viola", sounds: []string{"mmmmm", "zoooo", "veee", "whooo", "vrreee"}},
	{frequency: 5, name: "Glockenspiel", sounds: []string{"ding", "tinkle", "ping", "plink", "chime"}},
	{frequency: 7, name: "Organ", sounds: []string{"hummmm", "ooooo", "voooom", "drone", "wooo"}},
	{frequency: 4, name: "Bass Clarinet", sounds: []string{"mmmm", "brooo", "bwooo", "rooo", "vrmmm"}},
	{frequency: 6, name: "English Horn", sounds: []string{"hooo", "wheee", "woooo", "phmmm", "breee"}},
	{frequency: 8, name: "Concert Bass Drum", sounds: []string{"BOOM", "rumble", "dum", "doom", "pum"}},
	{frequency: 5, name: "Cymbals", sounds: []string{"crash", "clang", "clash", "shing", "chhhh"}},
	{frequency: 6, name: "Recorder", sounds: []string{"tweet", "toot", "peep", "reep", "fweee"}},
	{frequency: 5, name: "Baritone Saxophone", sounds: []string{"vrooo", "booo", "bop", "grmmm", "vrooom"}},
	{frequency: 7, name: "Marimba", sounds: []string{"tok", "tonk", "dunk", "dong", "bong"}},
}

var openings = []string{
	"The sound of %s fills the air",
	"Everyone listens as %s takes over",
	"A melody rises, played by %s",
	"The stage belongs to %s now",
	"You can hear %s in the distance",
	"All eyes are on %s as it begins",
	"The music swells, led by %s",
	"A soft hum emerges from %s",
	"Powerful notes erupt from %s",
	"The rhythm shifts, thanks to %s",
	"From the corner, %s adds its voice",
	"The harmony is completed by %s",
	"Suddenly, %s makes its presence known",
	"In the mix, %s finds its place",
	"A delicate tune floats out of %s",
	"The energy builds, driven by %s",
	"A resonant sound comes from %s",
	"The silence is broken by %s",
	"An unmistakable sound flows from %s",
	"Everything changes when %s joins in",
	"The audience is captivated by %s",
	"The backdrop hums with the sound of %s",
	"A new tone emerges, thanks to %s",
	"The piece takes flight with %s",
	"A rich sound emanates from %s",
	"The music deepens as %s plays",
	"Out of nowhere, %s begins to play",
	"The atmosphere transforms with %s",
	"The melody comes alive with %s",
	"A wave of sound builds around %s",
	"The air is electrified by %s",
	"The composition breathes through %s",
	"A bright tone emerges from %s",
	"The song's heartbeat is driven by %s",
	"In the chaos, %s finds its voice",
	"The layers of sound are enriched by %s",
	"A subtle rhythm flows from %s",
	"The crowd stirs as %s joins the fray",
	"The lead shifts to %s for a moment",
	"The balance is perfected by %s",
	"The soul of the piece resonates with %s",
	"A cascade of notes falls from %s",
	"The performance peaks with %s",
	"Each note feels alive with %s playing",
	"The essence of the tune shines through %s",
	"A haunting sound drifts from %s",
	"The soundscape expands with %s",
	"The magic unfolds around %s",
	"The rhythm breathes new life through %s",
	"From the shadows, %s contributes a tone",
	"The journey continues with %s",
	"A bold entrance by %s turns heads",
	"The crescendo builds, led by %s",
	"The quiet is punctuated by %s",
	"The song finds its pulse in %s",
	"The atmosphere shimmers with %s",
	"A tender phrase is born from %s",
	"The mood shifts under the spell of %s",
	"%s brings a new layer to the melody",
	"%s fills the space with its sound",
	"%s adds depth to the composition",
	"%s carries the tune to new heights",
	"%s weaves through the harmony effortlessly",
	"%s resonates with a rich and vibrant tone",
	"%s shapes the rhythm with precision",
	"%s colors the soundscape beautifully",
	"%s takes the lead with bold notes",
	"%s softens the mood with its melody",
	"%s breathes life into the music",
	"%s anchors the harmony with steady tones",
	"%s dances through the melody with ease",
	"%s punctuates the silence with clarity",
	"%s soars above the other instruments",
	"%s enriches the atmosphere with its presence",
	"%s blends seamlessly into the symphony",
	"%s echoes the spirit of the piece",
	"%s shines as the centerpiece of the sound",
	"%s threads its voice into the composition",
	"%s carries the weight of the rhythm",
	"%s bursts forth with dynamic energy",
	"%s hums softly, anchoring the melody",
	"%s paints vivid colors with its notes",
	"%s rises and falls with graceful precision",
	"%s whispers a delicate phrase into the mix",
	"%s transforms the tune with its entrance",
	"%s gives the piece a fresh perspective",
	"%s stirs emotions with every note",
	"%s intertwines with the harmony effortlessly",
	"%s drives the pulse of the music forward",
}
