# the-war-of-the-worlds

1898 H. G. Wells Book being modernized as a simple simulation

# Toughts, Anthology and Analogy

Initially I wanted to give a distributed system anaology to each component in the user story

**City** - A stateful node communicating via some channel (go channel/socket/grpc/json-rpc)  
**Alien** - A client running a process which performes random transactions  
**Aliens visit** - transaction to a specific node  
**Aliens fight** - a transaction which causes fail-stop of the node

## Ideas I had

1. Represent each City as a server, a node, communicating via http/grpc/json-rpc and maintaining the state with some Sqlite or Badger.
2. Represent each City as a node in Tendermint, which is probably the best approach. I wanted to simulate a fail-stop of each node when being destroyed, I wasn't sure how to achieve it, so I passed on it.
3. Represent each city as container with a server instance and run it as a kubernetes pod, this approach is quite useful as k8s runs the Raft algorithm interenaly to manage the distributed state.

in hindsight, going with either one of the approached could have been easier, and much more robust.

## Chosen approach

I ended up with a rather simple model in my simulation

**City** - A stateful _goroutine_ communicating via some gochannel, it can send messages to aliens and other cities ( only neighbour cities)  
**Alien** - A _goroutine_ communicating only with cities  
**Aliens visit** - A message with status - "visit" being sent to a city  
**Aliens fight** - A result of two aliens visiting a city at the same time, a city notices that and pings both aliens as well as other cities regarding the fight.

## Flow

1. Cities and aliens are being generated according to input and names list
2. Cities and alien are being run randomly
3. An alien visits a city by sending a message on a dedicated channel, two scenarios might occure  
   a. The city is free and nothing really happens
   b. Another alien is in the city, in this case the aliens goroutines will halt, the city will let other cities (only neighbour cities) know that they should ignore it in the future, and halt as well
4. After some time if all aliens are dead or after some amount of aliens city hops (default 10,000), the program will stop and will print out to a file updated map.

# Run

`clone https://github.com/CommoDor64/the-war-of-the-worlds.git`  
`cd the-war-of-the-worlds`  
`go mode install`  
`go run main.go <number_of_aliens>`
