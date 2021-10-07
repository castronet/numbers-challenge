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
* Channels for communication between go routines

 */

func main() {

    // [+] Variables declaration
	maxConcurrentUsers := 5

    // channels for communications
	system := make(chan interface{}, 1)
	input := make(chan int)

	// mutex declaration to warrant right data management
	mutex := sync.Mutex{}

	// Variables to count numbers
	numbers := map[int]struct{}{}
	totalUniqueNumbers := 0
	duplicatedNumbers := 0
	uniqueNumbers := 0

    // Open file where to write numbers
    outputFile, err := os.OpenFile("./numbers.txt", os.O_RDWR|os.O_CREATE, 0666)
    checkError(err)
    defer outputFile.Close()

    // Truncate the file and seek to the top
    outputFile.Truncate(0)
    outputFile.Seek(0,0)


	// Open port
	listener, err := net.Listen("tcp", ":4000")
	checkError(err)

    // Define maximun concurrent users on our listener
	listener = netutil.LimitListener(listener, maxConcurrentUsers)


	go func() {
		for newNumber := range input {
			// lock
			mutex.Lock()
			if _, ok := numbers[newNumber]; ok {
				duplicatedNumbers++
			} else {
				uniqueNumbers++
			    totalUniqueNumbers++
				numbers[newNumber] = struct{}{}

                // Write to the file each number
                _, err := outputFile.WriteString(fmt.Sprintf("%d\n", newNumber))
//                    strconv.Itoa(number)
                checkError(err)
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

			fmt.Printf("Received %d unique numbers, %d duplicates. Unique total: %d\n", uniqueNumbers, duplicatedNumbers, totalUniqueNumbers)

            // initialize duplicate and unique range variables
			duplicatedNumbers = 0
			uniqueNumbers = 0

			mutex.Unlock()
		}
	}()

	<-system
	close(system)
	listener.Close()

	os.Exit(0)
}

func handleClient(conn net.Conn, input chan int, system chan interface{}) {
    // Open scanner to read data from clients
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()

        // Check if the input has the required lenght
		if len(line) != 9 {
			conn.Close()
            return
		}

        // Convert a string to a int
        // we did it at this point because we expect more correct numbers than errors
		i, err := strconv.Atoi(line)
		if err != nil {
            // We check
			// if here because of cost efficient just one terminate on
			if line == "terminate" {
				system <- struct{}{}
			} else {
				conn.Close()
			}
        } else {
            input <- i
        }
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
