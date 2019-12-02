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
