package ai

import (
	"fmt"
	"time"
)

const (
	NAME = "Basic Go Bot"
	VERSION = "1"
)

func Run() {

	game := NewGame()

	game.StartLog(fmt.Sprintf("log%d.txt", game.Pid()))
	game.LogWithoutTurn("--------------------------------------------------------------------------------")
	game.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02T15:04:05Z"))

	fmt.Printf("%s %s\n", NAME, VERSION)

	for {
		game.Parse()
		game.BasicAI()
		game.Send()
	}
}
