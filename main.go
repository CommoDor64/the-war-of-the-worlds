package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
	"twotw/alien"
)

func onError(err error) {
	if err != nil {
		panic(err)
	}
}

// must fix
// func testFormat(line string) {
// 	// name-length was set to 85 to cover the city Taumatawhakatangi­hangakoauauotamatea­turipukakapikimaunga­horonukupokaiwhen­uakitanatahu
// 	r := regexp.MustCompile(`[A-Za-z0-9_-]{0,85} (north\=[A-Za-z0-9_-]+)? (west\=[A-Za-z0-9_-]+)? (south\=[A-Za-z0-9_-]+)? (east\=[A-Za-z0-9_-]+)?`)
// 	submatch := r.FindStringSubmatch(line)
// 	fmt.Println(submatch)
// }

func main() {
	// read cli args
	argsOnly := os.Args[1:]
	if len(argsOnly) < 1 {
		fmt.Println("usage: ./main.go <number_of_aliens>")
		return
	}
	_, _ = strconv.Atoi(argsOnly[0])

	// open and read cities map
	m, err := ioutil.ReadFile("map.txt")
	if err != nil {
		panic(err)
	}
	ma := alien.NewMarshaler()
	cities := ma.Unmarshal(m)

	allRoads := make(map[string]chan alien.Datagram)
	// create a channel for each city
	for _, c := range cities {
		allRoads[c.Name] = make(chan alien.Datagram, 10)
	}

	// connect the channel of each city to each other (if connected in the map)
	var citySlice []alien.City
	for _, c := range cities {
		// paths := make(map[city.Direction]chan city.Datagram, aliensNum)
		var cityOutRoads []alien.OutRoad
		for direction, name := range c.Paths {
			cityOutRoads = append(cityOutRoads, alien.OutRoad{
				Name:      name,
				Direction: alien.Direction(direction),
				Road:      allRoads[name],
			})
		}
		newCity := alien.NewCity(c.ID, c.Name, cityOutRoads, allRoads[string(c.Name)])
		fmt.Println(newCity)
		citySlice = append(citySlice, newCity)
		go newCity.Run()
	}
	fmt.Println()

	// create aliens
	aliens := alien.AlienFactory([]string{"joe", "christopher"})

	for _, a := range aliens {
		time.Sleep(time.Second * 2)
		rand.Seed(time.Now().Unix())
		rn := rand.Intn(len(citySlice))
		// fmt.Println(a.Name, "on", citySlice[rn])
		go a.Roam(citySlice[rn])
	}
	fmt.Println("test")
	for {
		fmt.Println(citySlice)
		time.Sleep(time.Second * 5)
	}
}
