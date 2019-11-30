package alien

import (
	"fmt"
	"strings"
)

type marshaler struct {
}

type Marshaler interface {
	Marshal([]City)
	Unmarshal(m []byte) []RawCity
}

// parses each line into struct
func parseLine(index int, line string) RawCity {
	// split at white spacee
	lineSlice := strings.Split(line, " ")
	// assign name to city
	city := RawCity{
		ID:    index,
		Name:  lineSlice[0],
		Paths: make(map[Direction]string),
	}
	// iterate paths segment, asign path to direction accordingly
	for _, segment := range lineSlice[1:] {
		segmentSplit := strings.Split(segment, "=")

		city.Paths[Direction(segmentSplit[0])] = segmentSplit[1]
	}
	return city
}

func (m marshaler) Unmarshal(file []byte) []RawCity {
	lines := strings.Split(string(file), "\n")
	var cities []RawCity
	for index, line := range lines {
		// testFormat(line)
		city := parseLine(index, line)
		cities = append(cities, city)
	}
	return cities
}

func (m marshaler) Marshal(cities []City) string {
	file := ""
	for _, city := range cities {
		file += fmt.Sprintf("%s ", city.Name)
		for _, p := range city.OutRoads {
			file += fmt.Sprintf("%v=%s", p.Name, p.Direction)
		}
		file += "\n"
	}
	return file
}

func NewMarshaler() marshaler {
	return marshaler{}
}
