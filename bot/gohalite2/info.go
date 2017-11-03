package gohalite2

import (
	"sort"
)

// ----------------------------------------------

func (self *Game) GetShip(sid int) Ship {
	return self.shipMap[sid]
}

func (self *Game) GetPlanet(plid int) Planet {
	return self.planetMap[plid]
}

// ----------------------------------------------

func (self *Game) AllShips() []Ship {
	ret := make([]Ship, len(self.all_ships_cache))
	copy(ret, self.all_ships_cache)
	return ret
}

func (self *Game) AllPlanets() []Planet {
	ret := make([]Planet, len(self.all_planets_cache))
	copy(ret, self.all_planets_cache)
	return ret
}

func (self *Game) AllImmobile() []Entity {						// Returns all planets and all docked ships
	ret := make([]Entity, len(self.all_immobile_cache))
	copy(ret, self.all_immobile_cache)
	return ret
}

// ----------------------------------------------

func (self *Game) ShipsOwnedBy(pid int) []Ship {
	ret := make([]Ship, len(self.playershipMap[pid]))
	copy(ret, self.playershipMap[pid])
	return ret
}

func (self *Game) MyShips() []Ship {
	return self.ShipsOwnedBy(self.pid)
}

func (self *Game) EnemyShips() []Ship {
	ret := make([]Ship, len(self.enemy_ships_cache))
	copy(ret, self.enemy_ships_cache)
	return ret
}

// ----------------------------------------------

func (self *Game) MyNewShipIDs() []int {			// My ships born this turn.
	var ret []int
	for sid, _ := range self.shipMap {
		ship := self.shipMap[sid]
		if ship.Birth == self.turn && ship.Owner == self.pid {
			ret = append(ret, ship.Id)
		}
	}
	sort.Slice(ret, func(a, b int) bool {
		return ret[a] < ret[b]
	})
	return ret
}

func (self *Game) ShipsDockedAt(planet Planet) []Ship {
	ret := make ([]Ship, len(self.dockMap[planet.Id]))
	copy(ret, self.dockMap[planet.Id])
	return ret
}

func (self *Game) LastTurnMoveBy(ship Ship) MoveInfo {
	return self.lastmoveMap[ship.Id]
}

func (self *Game) ClosestPlanet(e Entity) Planet {

	var best_dist float64 = 9999999
	var ret Planet

	for _, planet := range self.AllPlanets() {
		dist := e.ApproachDist(planet)
		if dist < best_dist {
			best_dist = dist
			ret = planet
		}
	}

	return ret
}

func (self *Game) RawWorld() string {
	return self.raw
}

func (self *Game) GetCumulativeShipCount(pid int) int {
	return self.cumulativeShips[pid]
}
