package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	ai "./ai"
	hal "./core"
	nav "./navigation"
)

const (
	NAME = "Fohristiwhirl"
	VERSION = "115"
)

func main() {

	config := new(ai.Config)

	flag.BoolVar(&config.Centre, "centre", false, "take the centre first (1v1)")
	flag.BoolVar(&config.Conservative, "conservative", false, "no rushing")
	flag.BoolVar(&config.DockOnly, "dockonly", false, "make initial dockings and stop")
	flag.BoolVar(&config.ForceRush, "forcerush", false, "always rush")
	flag.BoolVar(&config.Imperfect, "imperfect", false, "don't use \"perfect\" GA")
	flag.BoolVar(&config.NoMsg, "nomsg", false, "no angle messages")
	flag.BoolVar(&config.Profile, "profile", false, "run Golang CPU profile")
	flag.BoolVar(&config.Split, "split", false, "split ships at start")
	flag.BoolVar(&config.Timeseed, "timeseed", false, "seed RNG with time")

	flag.Float64Var(&nav.Ignore_Collision_Dist, "icd", 100, "ignore collision distance (nav)")

	flag.IntVar(&config.TestGA, "testga", -1, "test GA on thus turn")

	flag.Parse()

	game := hal.NewGame()

	if config.Profile {
		outfile, _ := os.Create("profile.prof")
		pprof.StartCPUProfile(outfile)
		defer pprof.StopCPUProfile()
	}

	var longest_turn time.Duration
	var longest_turn_number int

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("%v", p)
			game.Log("Quitting: %v", p)
			game.Log("Last known hash: %s", hal.HashFromString(game.RawWorld()))
			game.Log("Current ships...... %3d, %3d, %3d, %3d",
				len(game.ShipsOwnedBy(0)),
				len(game.ShipsOwnedBy(1)),
				len(game.ShipsOwnedBy(2)),
				len(game.ShipsOwnedBy(3)))
			game.Log("Cumulative ships... %3d, %3d, %3d, %3d",
				game.GetCumulativeShipCount(0),
				game.GetCumulativeShipCount(1),
				game.GetCumulativeShipCount(2),
				game.GetCumulativeShipCount(3))
			game.Log("Longest turn (%d) took %v", longest_turn_number, longest_turn)
		}
	}()

	game.StartLog(fmt.Sprintf("log%d.txt", game.Pid()))
	game.LogWithoutTurn("--------------------------------------------------------------------------------")
	game.LogWithoutTurn("%s %s starting up at %s", NAME, VERSION, time.Now().Format("2006-01-02T15:04:05Z"))

	if config.Timeseed {
		seed := time.Now().UTC().UnixNano()
		rand.Seed(seed)
		game.LogWithoutTurn("Seeding own RNG: %v", seed)
	}

	if len(os.Args) < 2 {
		fmt.Printf("%s %s\n", NAME, VERSION)
	} else {
		fmt.Printf("%s %s %s\n", NAME, VERSION, strings.Join(os.Args[1:], " "))
	}

	overmind := ai.NewOvermind(game, config)

	for {
		start_time := time.Now()

		game.Parse()

		if config.Timeseed == false {
			rand.Seed(int64(game.Turn() + game.Width() + game.Pid()))
		}

		if config.TestGA > -1 {								// No moves except on test turn...
			if game.Turn() == config.TestGA {
				overmind.EnterGeneticAlgorithm()
			}
			game.Send(config.NoMsg)
			continue
		}

		overmind.Step()
		game.Send(config.NoMsg)

		if time.Now().Sub(start_time) > longest_turn {
			longest_turn = time.Now().Sub(start_time)
			longest_turn_number = game.Turn()
		}
	}
}
