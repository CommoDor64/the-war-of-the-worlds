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

# What could have been done better:

1. **Complex and convoluted code** - The low degree of spereation between the client (Alien) and Network/Node/Server (City) makes the codebase a bit spaghetti.
2. **Data Race** - Because of lots of message transfer without a suitable protocol, there are some incidents of data race, for example in the case of creating a larger number of aliens than cities. This can be fixed but requires some model change so I left it as it is.
3. **More flexibilty in config** - If I had the time I would add config to the city and alien (now the internal behavior is hardcoded), preferably as idiomatic golang functional configurations.
4. **Tests** - I did not want to make half ass work with the tests so I decided to submit without.

# Notes

In an ideal situation, I would just create the Client (Alien) and Network/Node/Server (City) in a complete seperate manner, which means with established communication protocol and better structure overall.  
Since this project was implemented using 100% Go with no dependencies, it was challenging but also quite interesting, as the challenges were quite unusual.

The entire project took around 10 hours, where it includes thinking about a model, trying simple ones and getting to the end result.

# Project Structure and Run

**/alien** - the alien package defines the entire case specific progam  
**/alien/alien.go** - defiens alien's behavior and data model  
**/aliien/city** - defines city's behavior and data model  
**/alien/repository.go** - help functions  
**/alien/marshaler.go** - marshaler to convert from plain text to the model and from model to plain text  
**/alien/mode.go** - defiens all the models, for city, alien and messages between them

`clone https://github.com/CommoDor64/the-war-of-the-worlds.git`  
`cd the-war-of-the-worlds`  
`go mode install`  
`go run main.go <number_of_aliens>`

The output will be stored in out.txt
