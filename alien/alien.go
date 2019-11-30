// package alien

// import (
// 	"fmt"
// 	"math/rand"
// 	"reflect"
// 	"time"
// 	"twotw/common"
// )

// type Alien struct {
// 	Name     string
// 	Brain    Brain
// 	Ears     chan common.Datagram
// 	Mouth    []chan common.Datagram
// 	Heart    chan bool
// 	Location uint
// }

// type Brain struct {
// 	Cities        []common.City
// 	AlienLocation map[string]uint
// 	Leader        int
// }

// type alien interface {
// }

// // NewAlien sets initial data for the alien
// func NewAlien(name string, cities []common.City, ears chan common.Datagram, mouth []chan common.Datagram) Alien {
// 	rand.Seed(time.Now().Unix())
// 	rn := uint(rand.Intn(len(cities)))
// 	fmt.Println("alien", name, "start at", cities[rn])
// 	return Alien{
// 		Name:     name,
// 		Brain:    Brain{Cities: cities, AlienLocation: make(map[string]uint)},
// 		Ears:     ears,
// 		Mouth:    mouth,
// 		Heart:    make(chan bool, 1),
// 		Location: rn,
// 	}
// }

// // shout transmit location information to other aliens
// func (alien *Alien) shout(quit chan bool) {
// 	time.Sleep(time.Second)
// 	datagram := common.Datagram{
// 		Name:     alien.Name,
// 		Location: alien.Location,
// 	}
// 	fmt.Println("alien ", alien.Name, "sent:", datagram)
// 	for _, c := range alien.Mouth {
// 		if c == alien.Ears {
// 			continue
// 		}
// 		c <- datagram
// 	}

// }

// // listen receives location information from other aliens
// func (alien *Alien) listen(quit chan bool) {
// 	for {
// 		select {
// 		case datagram, ok := <-alien.Ears:
// 			if ok {
// 				alien.Brain.AlienLocation[datagram.Name] = datagram.Location

// 				go alien.shout(quit)
// 				alien.move(quit)
// 			}
// 		case <-quit:
// 			return
// 		}
// 		fmt.Print()
// 	}
// }

// // move sets a new location of the alien, randomly
// func (alien *Alien) move(quit chan bool) {
// 	// check before moving, whether another alien is on the planet
// 	for _, v := range alien.Brain.AlienLocation {
// 		if v == alien.Location {
// 			fmt.Println("alien", alien.Name, "dies on", alien.Location)
// 			quit <- true
// 			return
// 		}
// 	}
// 	time.Sleep(time.Millisecond * 100)
// 	// current city struct of alien
// 	currentLocation := alien.Brain.Cities[alien.Location]
// 	// create a slice from direction keys for easy random selection
// 	directionSlice := reflect.ValueOf(currentLocation.Paths).MapKeys()
// 	// fmt.Println("panic", len(directionSlice))
// 	direction := directionSlice[uint(rand.Intn(len(directionSlice)))].String()

// 	for _, city := range alien.Brain.Cities {
// 		if currentLocation.Paths[common.Direction(direction)] == city.Name {
// 			alien.Location = uint(city.ID)
// 			fmt.Println("alien", alien.Name, "moves to", alien.Brain.Cities[alien.Location])
// 			break
// 		}
// 	}
// }

// func (alien *Alien) think(quit chan bool) {

// }

// // Unleash starts an alien
// func (alien *Alien) Roam() {
// 	quit := make(chan bool)
// 	go alien.listen(quit)
// 	go alien.shout(quit)
// 	<-quit
// 	fmt.Println("done", alien.Brain.Cities)

// 	return
// }

package alien

import (
	"fmt"
	"math/rand"
	"time"
)

func NewAlien(name string) Alien {
	return Alien{
		Name:   name,
		Ears:   make(chan Datagram),
		Heart:  make(chan Datagram),
		Action: Visit,
	}
}

func AlienFactory(names []string) []Alien {
	var aliens []Alien
	for _, name := range names {
		fmt.Println("creats", name)
		aliens = append(aliens, NewAlien(name))
	}
	return aliens
}
func (a Alien) die(quit chan bool) {
	for {
		select {
		case m, ok := <-a.Heart:
			if ok {
				fmt.Println("got", m, "alien", a.Name)
				if m.City.Status == Fight {
					quit <- true
					return
				}
			}
		}
	}
}

func (a Alien) move(city City) {
	radioTower := city.RadioTower
	for j := 0; j < 10000; j++ {
		time.Sleep(time.Second * 2)
		a.Action = Visit
		radioTower <- Datagram{
			Alien: Alien{
				Name:   a.Name,
				Ears:   a.Ears,
				Heart:  a.Heart,
				Action: Visit,
			},
		}
		select {
		case m, ok := <-a.Ears:
			if ok {
				fmt.Println("got", m, "alien", a.Name)
				// if m.City.Status == Fight {
				// 	fmt.Println("alien", a.Name, "dies on", city.Name)
				// 	return
				// }
			}
		}
		// spend time in city
		time.Sleep(time.Second * 5)
		// leave city
		radioTower <- Datagram{Alien: Alien{
			Name:   a.Name,
			Ears:   a.Ears,
			Heart:  a.Heart,
			Action: Leave,
		}}
		datagram := <-a.Ears
		outRoadsLen := len(datagram.City.OutRoads)
		if outRoadsLen <= 0 {
			fmt.Println("alien", a.Name, "is stuck on", city.Name)
			return
		}
		rn := rand.Intn(outRoadsLen)
		radioTower = datagram.City.OutRoads[rn].Road
	}
}
func (a Alien) Roam(city City) {
	quit := make(chan bool)
	go a.die(quit)
	go a.move(city)
	<-quit
	fmt.Println("Roam ended")
	return
}
