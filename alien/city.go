package alien

import (
	"fmt"
	"log"
	"sync"
)

// NewCity creates a city struct, with inbound and outbound channels
func NewCity(id int, name string, outRoads []OutRoad, radioTower chan Datagram) City {
	return City{
		ID:         id,
		Name:       name,
		Status:     Free,
		Visitors:   []Alien{},
		OutRoads:   outRoads,
		RadioTower: radioTower,
	}
}

func (c City) receive(sm *sync.Map, quit chan bool) {
	for {
		select {
		case datagram := <-c.RadioTower:
			// if 2 aliens visit at the same time, death
			// if aliens visits, add to visitor list, if leaves, removes
			switch {
			// aliens enters the city
			case datagram.Alien.Action == Visit:
				c.Visitors = append(c.Visitors, datagram.Alien)
				log.Printf("alien %s landed on %s\n", datagram.Alien.Name, c.Name)
				datagram.Alien.Ears <- Datagram{
					City: City{
						ID:       c.ID,
						Name:     c.Name,
						OutRoads: c.OutRoads,
						Status:   c.Status,
					},
				}
			// aliens leaves the city
			case datagram.Alien.Action == Leave:
				c.Visitors = removeAlien(c.Visitors, datagram.Alien)
				log.Printf("alien %s left %s and gets %v\n", datagram.Alien.Name, c.Name, c.OutRoads)
				datagram.Alien.Ears <- Datagram{
					City: City{
						ID:       c.ID,
						Name:     c.Name,
						OutRoads: c.OutRoads,
						Status:   c.Status,
					},
				}
			// neighbour city explodes
			case datagram.City.Status == Fight:
				var newOutRoads []OutRoad
				for _, outRoad := range c.OutRoads {
					if datagram.City.Name == outRoad.Name {
						continue
					}
					newOutRoads = append(newOutRoads, outRoad)
				}
				c.OutRoads = newOutRoads
				fmt.Println(c.OutRoads)
				sm.Store(c.Name, c)
			}
		}
		if len(c.Visitors) > 1 {
			c.Status = Fight

			// prepare die message to aliens
			resDatagram := Datagram{
				City: c,
			}

			// notify aliens
			log.Printf("%s and %s die on %s\n", c.Visitors[0].Name, c.Visitors[1].Name, c.Name)
			for _, visitor := range c.Visitors {
				visitor.Heart <- resDatagram
			}

			//notify other cities

			for _, outRoad := range c.OutRoads {
				outRoad.Road <- resDatagram
			}
			sm.Store(c.Name, c)
			quit <- true
			return
		}

	}
}

// Run opens up a city for bussiness
func (c *City) Run(sm *sync.Map) {
	quit := make(chan bool)
	go c.receive(sm, quit)
	<-quit
	return
}
