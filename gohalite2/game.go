package gohalite2

import (
	"fmt"
	"strings"
)

type Game struct {
	Pid					int					// Our own ID
	Width				int
	Height				int

	PlayerCount			int					// Does this change if a player dies?

	// These 3 things below can contain references to dead objects.
	// But the Planet and Player structs themselves do not.

	Players				map[int]*Player
	Planets				map[int]*Planet
	Ships				map[int]*Ship

	logfile             *Logfile
	token_parser		*TokenParser
}

func (self *Game) GetMe() *Player {
	return self.Players[self.Pid]
}

func (self *Game) Init(logfilename string, log_enabled bool) {

	self.token_parser = NewTokenParser()

	self.Pid = self.token_parser.Int()
	self.Width = self.token_parser.Int()
	self.Height = self.token_parser.Int()

	self.Players = make(map[int]*Player)
	self.Planets = make(map[int]*Planet)
	self.Ships = make(map[int]*Ship)

	self.logfile = NewLog(logfilename, log_enabled)

	self.Parse()
}

func (self *Game) Parse() {

	// Set all objects to have 0 HP, on the assumption that they are only sent to us if they have
	// 1 or more HP. Thus this is a default value if we receive no info.

	for _, planet := range(self.Planets) {
		planet.HP = 0
	}

	for _, ship := range(self.Ships) {
		ship.HP = 0
		ship.ClearOrder()
	}

	// Player parsing.............................................................................

	self.PlayerCount = self.token_parser.Int()

	for p := 0; p < self.PlayerCount; p++ {

		// Get or create the player in memory...

		pid := self.token_parser.Int()
		player, ok := self.Players[pid]

		if ok == false {
			player = new(Player)
			player.Id = pid
			self.Players[pid] = player
		}

		ship_count := self.token_parser.Int()

		player.Ships = nil		// Clear the player's ship slice, it will be recreated now...

		for s := 0; s < ship_count; s++ {

			// Get or create the ship in memory...

			sid := self.token_parser.Int()
			ship, ok := self.Ships[sid]

			if ok == false {
				ship = new(Ship)
				ship.Id = sid
				self.Ships[sid] = ship
			}

			ship.X = self.token_parser.Float()
			ship.Y = self.token_parser.Float()
			ship.HP = self.token_parser.Int()
			self.token_parser.Float()							// Skip deprecated "speedx"
			self.token_parser.Float()							// Skip deprecated "speedy"
			ship.Docked = self.token_parser.Int()
			ship.DockedPlanet = self.token_parser.Int()
			ship.DockingProgress = self.token_parser.Int()
			ship.Cooldown = self.token_parser.Int()

			player.Ships = append(player.Ships, ship)
		}
	}

	// Planet parsing.............................................................................

	planet_count := self.token_parser.Int()

	for p := 0; p < planet_count; p++ {

		plid := self.token_parser.Int()

		planet, ok := self.Planets[plid]
		if ok == false {
			planet = new(Planet)
			planet.Id = plid
			self.Planets[plid] = planet
		}

		planet.X = self.token_parser.Float()
		planet.Y = self.token_parser.Float()
		planet.HP = self.token_parser.Int()
		planet.Radius = self.token_parser.Float()
		planet.DockingSpots = self.token_parser.Int()
		planet.CurrentProduction = self.token_parser.Int()
		self.token_parser.Int()									// Skip deprecated "remaining production"
		planet.Owned = self.token_parser.Int()					// This should probably be converted to bool
		planet.Owner = self.token_parser.Int()

		if planet.Owned == 0 {
			planet.Owner = -1
		}

		// Docked ships are given to us as their IDs...

		planet.DockedShips = nil	// Clear the planet's ship slice, it will be recreated now...

		docked_ship_count := self.token_parser.Int()

		for s := 0; s < docked_ship_count; s++ {
			sid := self.token_parser.Int()
			planet.DockedShips = append(planet.DockedShips, self.Ships[sid])
		}
	}
}

func (self *Game) Send() {
	me := self.GetMe()

	var commands []string

	for _, ship := range(me.Ships) {
		if ship.Order != "" {
			commands = append(commands, ship.Order)
		}
	}
	fmt.Printf(strings.Join(commands, " "))
	fmt.Printf("\n")
}
