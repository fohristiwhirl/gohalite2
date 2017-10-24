package main

import (
	"fmt"
	"time"

	hal "./gohalite2"
)

const (
	NAME = "Fohristiwhirl"
	VERSION = 2
)

func main() {

	game := hal.NewGame()

	game.StartLog(fmt.Sprintf("log%d.txt", game.Pid))
	game.Log("--------------------------------------------------------------------------------")
	game.Log("%s %d starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02T15:04:05Z"))

	fmt.Printf("%s %d\n", NAME, VERSION)

	for {
		game.Parse()
		ai(game)
		game.Send()
	}
}

func ai(game *hal.Game) {

	// Idea: if 2 players, how about skipping the planet-getting stage and going for direct assassination?

}