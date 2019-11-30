package alien

type Status string
type Direction string
type Action string

var (
	// status enum
	Free      Status = "free"
	Fight     Status = "fight"
	Destroyed Status = "destroyed"
	// direction enum
	North Direction = "north"
	West  Direction = "west"
	South Direction = "south"
	East  Direction = "east"
	// alien action enum
	Visit Action = "visit"
	Leave Action = "leave"
	Ping  Action = "ping"
)

// City
type City struct {
	ID         int
	Name       string
	Status     Status
	Visitors   []Alien
	OutRoads   []OutRoad
	RadioTower chan Datagram
}

type OutRoad struct {
	Name      string
	Direction Direction
	Road      chan Datagram
}

// Alien
type Alien struct {
	Name   string
	Ears   chan Datagram
	Heart  chan Datagram
	Action Action
}

type PingPong struct {
	Sent bool
	Pong chan Datagram
}

// inter-alien communication payload
type Datagram struct {
	City  City
	Alien Alien
	Ping  PingPong
}

// marshaler

type RawCity struct {
	ID     int
	Name   string
	Status Status
	Paths  map[Direction]string
}
