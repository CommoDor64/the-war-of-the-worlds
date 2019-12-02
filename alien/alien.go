package alien

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
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
		log.Println("creats", name)
		aliens = append(aliens, NewAlien(name))
	}
	fmt.Println()
	return aliens
}
func (a Alien) die(quit chan bool) {
	for {
		select {
		case m, ok := <-a.Heart:
			if ok {
				if m.City.Status == Fight {
					quit <- true
					return
				}
			}
		}
	}
}

func (a Alien) move(city *City, quit chan bool) {
	radioTower := city.RadioTower
	for j := 0; j < 20; j++ {
		time.Sleep(time.Second)
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
		case _, ok := <-a.Ears:
			if ok {
			}
		}
		// spend time in city
		time.Sleep(time.Second)
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
			quit <- true
			return
		}
		rn := rand.Intn(outRoadsLen)
		radioTower = datagram.City.OutRoads[rn].Road
	}
	quit <- true
}
func (a Alien) Roam(city *City, wg *sync.WaitGroup) {
	quit := make(chan bool, 1)
	go a.die(quit)
	go a.move(city, quit)
	<-quit
	wg.Done()
	return
}
