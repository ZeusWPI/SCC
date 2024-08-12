package buzzer

import (
	"log"
	"scc/config"

	"github.com/a-h/beeper"
	"github.com/stianeikeland/go-rpio"
)

var buzzerOptions = map[string]func(rpio.Pin){
	"default": playMusic,
}

func PlayBuzzer() {
	err := rpio.Open()
	if err != nil {
		log.Printf("Error: Unable to open pin: %s", err)
		return
	}
	defer rpio.Close()

	pin := rpio.Pin(config.GetConfig().Buzzer.Pin)

	buzzerSong := config.GetConfig().Buzzer.Song
	val, ok := buzzerOptions[buzzerSong]
	if !ok {
		log.Printf("Error: Selected buzzer song: %s does not exist\n", buzzerSong)
		return
	}

	val(pin)
}

func playMusic(pin rpio.Pin) {
	bpm := 300
	music := beeper.NewMusic(pin, bpm)

	music.Note("A5", beeper.Quaver)
	music.Note("B5", beeper.Quaver)
	music.Note("D5", beeper.Quaver)
	music.Note("B5", beeper.Quaver)
	music.Note("E5", beeper.Crotchet)
	music.Note("E5", beeper.Crotchet)
	music.Note("D5", beeper.Quaver)
	music.Note("C#5", beeper.Quaver)
	music.Note("B4", beeper.Quaver)
}
