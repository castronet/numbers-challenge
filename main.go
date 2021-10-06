package main

// go get -u golang.org/x/net/netutil
import (
    "net"
    "os"
    "fmt"
//    "strconv"
    "golang.org/x/net/netutil"
)

/* Considerations
* TPC/IP v4
* Manage concurrent users with netutil // option to create a semaphore to do that


*/


func main() {

    maxConcurrentUsers := 5;

    // Open port
    listener, err := net.Listen("tcp", ":4000")
    checkError(err)
    defer listener.Close()

    listener = netutil.LimitListener(listener, maxConcurrentUsers)

    // Client bucle
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }

        // Handle clients in a gorouting
        go handleClient(conn)
    }

}

func handleClient(conn net.Conn) {
    // close connection on exit
    defer conn.Close()

    // Valid inputs:
    // 9 decimal number
    // 123456789
    // terminate

    var buf [10]byte
    for {
        // read upto 10 bytes
        n, err := conn.Read(buf[0:9])
        if err != nil {
            return
        }

        // write the n bytes read
        _, err2 := conn.Write(buf[0:n])
        if err2 != nil {
            return
        }

//         write the concurrentUsers
//        _, err3 := conn.Write([]byte(strconv.Itoa(concurrentUsers)))
//        if err3 != nil {
//            return
//        }
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
