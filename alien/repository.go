package alien

// convinient function to remove alien from visitors list
func removeAlien(visitors []Alien, visitorToRemove Alien) []Alien {
	var newAliens []Alien
	for _, a := range visitors {
		if a.Name == visitorToRemove.Name {
			continue
		}
		newAliens = append(newAliens, a)
	}
	return newAliens
}

// selects a random city for next step of alien
// func SelectRandomCity(city City, cities map[Direction]chan Datagram) chan Datagram {
// 	var citySlice []chan Datagram
// 	rand.Seed(time.Now().Unix())
// 	rn := rand.Intn(len(cities))
// 	if rn >= len(cities) {
// 		return city.RadioTower
// 	}
// 	for _, c := range cities {
// 		citySlice = append(citySlice, c)
// 	}
// 	return citySlice[rn]
// }

// pingNeighbourCities checks if city path exists.
// func pingNeighbourCities(outRoads []OutRoad) []OutRoad {
// 	newRoads := []OutRoad{}
// 	quit := make(chan bool)
// 	go func() {
// 		for _, road := range outRoads {
// 			channel := make(chan bool, 10)

// 			datagram := Datagram{
// 				Ping: PingPong{
// 					Sent: true,
// 					Pong: channel,
// 				},
// 			}
// 			// sends ping to city
// 			road.Road <- datagram
// 			// wait for pong from city
// 			select {
// 			case alive := <-channel:
// 				fmt.Println(alive)
// 				if !alive {
// 					return
// 				}
// 				newRoads = append(newRoads, road)

// 			}
// 		}
// 		quit <- true
// 	}()
// 	<-quit
// 	return newRoads
// }
