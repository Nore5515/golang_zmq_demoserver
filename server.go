package main  
 
import (
	"fmt" 
	"time"
	"math/rand"
	"log"
	zmq "github.com/pebbe/zmq4"
)

const incremental = true

func main() {  
	fmt.Println("Hello World")  

	zctx, _ := zmq.NewContext()
	s, _ := zctx.NewSocket(zmq.REP)
	s.Bind("tcp://*:5555")

	pubStarted := false

	for {
			// Wait for next request from client
			msg, _ := s.Recv(0)
			log.Printf("Received %s\n", msg)

			if msg == "start" || msg == "1" {
				s.Send("OK",0)
				if (!pubStarted){
					pubStarted = true
					startPub(&pubStarted)
				}

			} else if msg == "inf" || msg == "runInf" {
				s.Send("OK-INF", 0)
				if (!pubStarted){
					pubStarted = true
					go infinitePub(&pubStarted)
				}
			} else {
				s.Send("NO",0)
			}
	}
}

func genRandomPairing() (string, string){
	x := rand.Int31()
	y := rand.Int31()

	hexX := fmt.Sprintf("%x", x)
	hexY := fmt.Sprintf("%x", y)

	return hexX, hexY
}

func getHexPair(x int, y int) (string, string){
	hexX := fmt.Sprintf("%x", x)
	hexY := fmt.Sprintf("%x", y)

	return hexX, hexY
}

func startPub(ptr *bool){
	time.Sleep(time.Second * 1)
	pubContext, _ := zmq.NewContext()
	pub, _ := pubContext.NewSocket(zmq.PUB)
	pub.Bind("tcp://*:5556")
	for i := 0; i < 100; i++{
		log.Println("SENDING")
		pub.Send("DATA", 0)
	}
	log.Println("\n\nDone Sending!")
	*ptr = false;
	// pubStarted = false;
	defer pub.Close()
}


func infinitePub(ptr *bool){
	time.Sleep(time.Second * 1)
	pubContext, _ := zmq.NewContext()
	pub, _ := pubContext.NewSocket(zmq.PUB)
	pub.Bind("tcp://*:5556")
	count := 0
	for {
		log.Println("SENDING")
		x, y := "", ""
		if (incremental){
			x, y = getHexPair(count, count)
			count += 1
		} else{
			x, y = genRandomPairing()
		}
		pairString := x + ", " + y
		pub.Send(pairString, 0)
	}
	// log.Println("\n\nDone Sending!")
	// *ptr = false;
	// // pubStarted = false;
	// defer pub.Close()
}