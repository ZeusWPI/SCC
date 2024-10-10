package buzzer

import (
	"log"
	"os/exec"
	"scc/config"
)

var buzzerOptions = map[string]func(){
	"default": playMusic,
}

func PlayBuzzer() {
	buzzerSong := config.GetConfig().Buzzer.Song
	fun, ok := buzzerOptions[buzzerSong]
	if !ok {
		log.Printf("Error: Selected buzzer song: %s does not exist\n", buzzerSong)
		return
	}
	fun()
}

func playMusic() {
	// See 'man beep'
	cmd := exec.Command(
		"beep",
		"-n", "-f880", "-l100", "-d0",
		"-n", "-f988", "-l100", "-d0",
		"-n", "-f588", "-l100", "-d0",
		"-n", "-f989", "-l100", "-d0",
		"-n", "-f660", "-l200", "-d0",
		"-n", "-f660", "-l200", "-d0",
		"-n", "-f588", "-l100", "-d0",
		"-n", "-f555", "-l100", "-d0",
		"-n", "-f495", "-l100", "-d0",
	)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running command 'beep': %s\n", err)
	}
}
