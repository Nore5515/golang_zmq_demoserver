package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"
)

var incremental = true
var ellipsesMax = 10000
var pubStarted = false

// Main Func; This runs first.
func main() {

	// Start the server-command listener.
	go listenAndReply() // commandHandling.go

	// Gather and organize CLI arguments.
	handleArguments()

	// ZMQ4 initialization.
	zctx, _ := zmq.NewContext()
	s, _ := zctx.NewSocket(zmq.REP)
	s.Bind("tcp://*:5555")

	// Main Loop
	for {
		// Wait for and display next request from client
		msg, _ := s.Recv(0)
		log.Printf("Received %s\n", msg)

		// Convert recieved message to int.
		// If it is an integer (aka Atoi returns > 0), try to start publishing that many times.
		msgInt, _ := strconv.Atoi(msg)
		if msgInt > 0 {
			s.Send("OK-NUMBER", 0)
			tryStartPub(msgInt, s)
		}

		// Non-Int message handling.
		switch msg {
		case "start":
			s.Send("OK", 0)
			tryStartPub(100000, s)
		case "inf":
			s.Send("OK-INF", 0)
			tryStartPub(1000000, s) // Eh it's not infinite but it's close enough
		default:
			s.Send("NO", 0)
		}
	}
}

// See if we are currently publishing data!
// If so, hold off on starting anymore data publishing.
// Otherwise, publish away!
func tryStartPub(amount int, s *zmq.Socket) {
	if !pubStarted {
		pubStarted = true
		go startPub(amount)
	} else {
		fmt.Println("OCCUPIED")
		s.Send("SOCKET OCCUPIED", 0)
	}
}

// Handle CLI arguments.
func handleArguments() {
	// Get all arguments in CLI call
	argsWithoutProg := os.Args[1:]
	fmt.Println(argsWithoutProg)
	incremental = sliceContains(argsWithoutProg, "incremental")
}

// Start publishing data!
func startPub(amount int) {
	// Wait one second
	// TODO: Hacky fix
	time.Sleep(time.Second * 1)

	// ZMQ4 Socket initialization.
	pubContext, _ := zmq.NewContext()
	pub, _ := pubContext.NewSocket(zmq.PUB)
	pub.Bind("tcp://*:5556")

	// This value is only used for an asthetic ellipse filling in feature.
	// Completely non-essential, but I think it looks nice.
	ellipses := 0

	log.Println("Sending ", amount, "!")
	fmt.Println()
	for i := 0; i < amount; i++ {
		// Ellipses are non-essential but pretty.
		ellipses += 1
		if ellipses >= ellipsesMax {
			ellipses = 0
			fmt.Print(".")
		}

		// Initialize x and y.
		x, y := "", ""
		// If argument for incremental is passed, generate a Hex Pair from current iteration.
		if incremental {
			x, y = getHexPair(i, i)
		} else {
			// Otherwise, generate a random pairing.
			x, y = genRandomPairing()
		}
		// Put the x and y into a single string.
		pairString := x + ", " + y
		// Send the single string on up!
		pub.Send(pairString, 0)
	}
	fmt.Println()
	log.Println("\nDone Sending!")
	pubStarted = false

	// Close the socket.
	defer pub.Close()
}
