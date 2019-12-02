package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"twotw/alien"
)

var wg sync.WaitGroup
var sm sync.Map

func onError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// read cli args
	argsOnly := os.Args[1:]
	if len(argsOnly) < 1 {
		fmt.Println("usage: ./main.go <number_of_aliens>")
		return
	}
	alienNum, _ := strconv.Atoi(argsOnly[0])

	// open and read cities map
	m, err := ioutil.ReadFile("map.txt")
	if err != nil {
		panic(err)
	}
	// read from alien names list
	rawNames, err := ioutil.ReadFile(".aliennames")
	if err != nil {
		panic(err)
	}
	alienNames := strings.Split(string(rawNames), "\n")

	// Unmarshal file format to the city model - alien.City
	ma := alien.NewMarshaler()
	cities := ma.Unmarshal(m)

	// create a channel for each city, so cities will be able to communicate
	allRoads := make(map[string]chan alien.Datagram)
	for _, c := range cities {
		allRoads[c.Name] = make(chan alien.Datagram)
	}

	// connect the channel of each city to each other (if connected in the map)
	var citySlice []*alien.City
	for _, c := range cities {
		var cityOutRoads []alien.OutRoad
		for direction, name := range c.Paths {
			cityOutRoads = append(cityOutRoads, alien.OutRoad{
				Name:      name,
				Direction: alien.Direction(direction),
				Road:      allRoads[name],
			})
		}
		newCity := alien.NewCity(c.ID, c.Name, cityOutRoads, allRoads[string(c.Name)])
		log.Println("creating a new city", newCity)
		sm.Store(newCity.Name, newCity)
		citySlice = append(citySlice, &newCity)
		go newCity.Run(&sm)
	}
	fmt.Println()

	// create aliens according to the number prompted
	aliens := alien.AlienFactory(alienNames[:alienNum])
	wg.Add(len(aliens))
	for _, a := range aliens {
		time.Sleep(time.Second)
		rand.Seed(time.Now().Unix())
		rn := rand.Intn(len(citySlice))
		go a.Roam(citySlice[rn], &wg)
	}
	wg.Wait()

	// check remaining cities, and marshal into a structured file
	var lastCities []alien.City
	sm.Range(func(name interface{}, city interface{}) bool {
		_ = name.(string)
		cityStruct := city.(alien.City)
		if cityStruct.Status != alien.Fight {
			lastCities = append(lastCities, cityStruct)
		}
		return true
	})
	fmt.Println()
	log.Println("the updated map can be found in out.txt")
	ioutil.WriteFile("out.txt", []byte(ma.Marshal(lastCities)), 0666)
}
