package gohalite2

import (
	"fmt"
)

// ---------------------------------------

type Ship struct {
	Id					int
	X					float64
	Y					float64
	HP					int
	Docked				int			// Is this really a bool?
	DockedPlanet		int
	DockingProgress		int

	Order				string
	Birth				int			// Turn this ship was first seen
}

func (self *Ship) Thrust(speed, angle int) {
	if self == nil {return}
	self.Order = fmt.Sprintf("t %d %d %d", self.Id, speed, angle)
}

func (self *Ship) Dock(planet int) {
	if self == nil {return}
	self.Order = fmt.Sprintf("d %d %d", self.Id, planet)
}

func (self *Ship) Undock(planet int) {
	if self == nil {return}
	self.Order = fmt.Sprintf("u %d", self.Id)
}

func (self *Ship) ClearOrder() {
	if self == nil {return}
	self.Order = ""
}