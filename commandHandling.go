package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	zmq "github.com/pebbe/zmq4"
)

// Command Ideas
// broadcast ______ #
// (broadcasts ____ to all subs, # of times)
// list
// (lists all connected subs)
// kick ______
// disconnects _____ sub

// As the name implies, this starts up a goroutine that just
// listens and replies back to the console.
// It also runs the handleInput func...
func listenAndReply() {
	fmt.Println("Listening...")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("BAD INPUT")
		} else {
			// Some input sanitation.
			input = strings.TrimSuffix(input, "\n")
			handleInput(input)
			fmt.Println(input)
		}
	}
}

// Turns console input strings into command calls.
func handleInput(input string) {
	// Split the string to individual words.
	strArr := strings.Split(input, " ")

	// If the first word is 'broadcast'...
	if strArr[0] == "broadcast" {
		// Ensure you have enough args for broadcast.
		if len(strArr) >= 3 {
			fmt.Println("Attempting Broadcast!")
			count, _ := strconv.Atoi(strArr[2])
			broadcastMsg(strArr[1], count)
		}
	}
}

// Broadcast message "msg" to all subscribers, "times" times.
func broadcastMsg(msg string, times int) {
	if !pubStarted {
		pubStarted = true
		// TODO: Hacky fix
		time.Sleep(time.Second * 1)

		// ZMQ4 initialization
		pubContext, _ := zmq.NewContext()
		pub, _ := pubContext.NewSocket(zmq.PUB)
		pub.Bind("tcp://*:5556")

		log.Println("Sending ", msg, "!")

		// Ellipses are nice, but take up a lot of room.
		ellipses := 0
		for i := 0; i < times; i++ {
			ellipses += 1
			if ellipses >= ellipsesMax {
				ellipses = 0
				fmt.Print(".")
			}
			pub.Send(msg, 0)
		}

		fmt.Println("")
		pubStarted = false

		// Close the socket.
		defer pub.Close()
	}
}
