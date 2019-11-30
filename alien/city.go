package alien

import "fmt"

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

// func CityFactory(citySliceRaw []) []City {
// 	var citySlice City
// 	for _, c := range cities {
// 		paths := make(map[city.Direction]chan city.Datagram, aliensNum)
// 		for k, v := range c.Paths {
// 			paths[city.Direction(k)] = cityChannel[v]
// 		}
// 		newCity := city.NewCity(c.ID, c.Name, paths, cityChannel[c.Name])
// 		fmt.Println(newCity)
// 		citySlice = append(citySlice, newCity)
// 		go newCity.Run()
// 	}
// }
func (c City) receive(quit chan bool) {
	for {
		select {
		case datagram := <-c.RadioTower:
			// if 2 aliens visit at the same time, death
			fmt.Println("Alien", datagram.Alien.Name, "says", datagram.Alien.Action)
			// if aliens visits, add to visitor list, if leaves, removes
			switch {
			// aliens enters the city
			case datagram.Alien.Action == Visit:
				c.Visitors = append(c.Visitors, datagram.Alien)
				fmt.Println("name:", datagram.Alien.Name, "channel:", datagram.Alien.Ears)
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
			}
		}
		if len(c.Visitors) > 1 {
			c.Status = Fight
			// send die message to aliens
			// fmt.Println(c.Visitors)
			resDatagram := Datagram{
				City: c,
			}
			fmt.Println("Death", c.Visitors, "on", c.Name)

			// notify aliens
			for _, visitor := range c.Visitors {
				fmt.Println("sends", resDatagram.City.Status, "to", visitor)
				visitor.Heart <- resDatagram
			}

			//notify other cities
			for _, outRoad := range c.OutRoads {
				fmt.Println("sending to", outRoad)
				outRoad.Road <- resDatagram
			}

			quit <- true
			return
		}

	}
}

// Run opens up a city for bussiness
func (c *City) Run() {
	// go c.teleport()
	quit := make(chan bool)
	go c.receive(quit)
	<-quit
	fmt.Print("Run ended")
	return
}
