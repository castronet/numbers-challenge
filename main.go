package main

// go get -u golang.org/x/net/netutil
import (
	"bufio"
	"fmt"
	"golang.org/x/net/netutil"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

/* Considerations
* TPC/IP v4
* Manage concurrent users with netutil // option to create a semaphore to do that


 */

func main() {

	maxConcurrentUsers := 5

	// Open port
	listener, err := net.Listen("tcp", ":4000")
	checkError(err)
	//    defer listener.Close()

	listener = netutil.LimitListener(listener, maxConcurrentUsers)

	system := make(chan interface{}, 1)
	input := make(chan int)
	//	connections := make(chan conn)

	// my mutex
	mutex := sync.Mutex{}

	// a ordenar - inicialicazion
	numbers := map[int]struct{}{}
	totalNumbers := 0
	duplicatedNumbers := 0
	uniqueNumbers := 0
	// open file - trucante

	go func() {
		for newNumber := range input {
			// lock
			mutex.Lock()
			totalNumbers++
			if _, ok := numbers[newNumber]; ok {
				duplicatedNumbers++
			} else {
				uniqueNumbers++
				numbers[newNumber] = struct{}{}

				// write to file
				// file write newNumber
			}

			//unlock
			mutex.Unlock()
		}
	}()

	go func() {
		// Client bucle
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			// Handle clients in a gorouting
			go handleClient(conn, input, system)
		}
	}()

	go func() {
		tick := time.Tick(10 * time.Second)
		for range tick {
			mutex.Lock()
			fmt.Printf("Received %d unique numbers, %d duplicates. Unique total: %d\n", uniqueNumbers, duplicatedNumbers, totalNumbers)
			duplicatedNumbers = 0
			uniqueNumbers = 0
			mutex.Unlock()
		}
	}()

	<-system
	close(system)
	// todo el codigo de cleanup
	listener.Close()
	// can I close all go routines?

	os.Exit(1)
}

func handleClient(conn net.Conn, input chan int, system chan interface{}) {

	// Valid inputs:
	// 9 decimal number
	// 123456789
	// terminate
	fmt.Println("Debug: handling connection")
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) > 9 { // check if the lenth must be 9 or 10
			conn.Close()
		}

		i, err := strconv.Atoi(line)
		if err != nil {
			// if here because of cost efficient just one terminate on
			if line == "terminate" {
				system <- struct{}{}
			} else {
				conn.Close()
			}
		}

		input <- i

		/*
		   case "terminate":
		       system <- struct{}{}

		   default:
		       if len(line) > 9 { // check if the lenth must be 9 or 10
		           system <- struct{}{}
		       }
		*/

		/*
		   switch line {
		       case "terminate":
		           system <- struct{}{}

		       default:
		           if len(line) > 9 { // check if the lenth must be 9 or 10
		               system <- struct{}{}
		           }

		           i, err := strconv.Atoi(line)
		           input <- i
		   }
		*/
//		fmt.Println("New line", line)

		/*
		   // check if the line is "terminate"
		   if err != nil {
		       system <- struct{}{}
		   }
		*/

	}

	//  inputLength := 10

	//  buffer := make([]byte, inputLength)

	//		fmt.Println("test test close close handler")
	// close connection on exit
	//    connections <- conn
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
