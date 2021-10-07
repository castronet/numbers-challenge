# Numbers challenge

This document will explain how to run the application and defend why I choose to do some functions as I did.

### Requirements

The application needs

* Golang v1.17
* Recommended: GNU Make

### Developed on

This challenge was developed and tested on a macOS Big Sur (11.6)

## How to run

If requirements are satisfied, to be able to run the application we can execute the following commands

* To install the required modules and create a vendor folder to be able to port the app to elsewhere
```
$ make vendor
```

* To build the application
```
$ make build
```

this step will generate a binary called `main`


You can run the application executing `./main` or
```
$ make run
```

When you execute the application, it will be on foreground waiting for clients to connect.

## Decisions taken

### IPv4

The application only listens for IPv4 protocol version, unfortunately IPv6 is not as implanted as some of us would want.

### Concurrent users

There are many possibilities to implement the maximum concurrent users behavior. I choose a package called `netutil` for
simplicity.

I could always code my own implementation with a maxConcurrentUsers-buffered channel and also send a message to other clients telling that the server is at maximum
capacity.

### Using go routines

My opinion is that using go routines is the key of golang and it lets us to work concurrently with different functions
on different threads.

### Go routines as anonymous function

I would like to create different functions for each go routine but the code is so small to don't be difficult to read
and I tried to order it as good as I could to let the code understandable.

On the other hand if you asked me to think in performance, create lots of nested functions will add a non desired overhead.

### handleClient function

`handleClient` function is an exception, that is not an anonymous function because that function has more logic and I prefer to move it a part.
And I created this point to tell that in Golang is not a best practice to have nested "if" but again,
for performance and be able to shortcut a client connection as fast as possible I decided to do the algorithm that way.


### Algorithm and complexity

I decided to split the application in three big parts:

* Client connection handler
* Received numbers logic
* Standard output message

I think that I already talked about client handling.

About received numbers I looked for the bests way to save the current received numbers with the fastest access to the
position and at same time the fastest way to save the number. Is for that reason that I chose a map.

Theoretically for hash tables working with `int` access, insert and delete items the cost in general is O(1), I assume
that this is also true for Go's maps.


Finally I treated the output message in a different function because it depends on a timer and for clarity I decided to
let it on it own scope.


Of course, we need a mutex while saving, reading and re-initializing values to avoid multiple threads working over the
same piece of memory.


### Output file

I decided to write to the file each time that I receive a unique number as the easiest option.


### Data loss

I tried to use buffers and mutex everywhere needed to be resilient with regard to data loss.

### Negative numbers


I didn't really know what to do with negative numbers, so I decided to discard them because the sign at the beginning is not a digit, then is an invalid line.



## How to improve

### Build method

The first point that comes to my mind about how to improve all the challenge is deliver it as a Docker.

I couldn't do that because Docker is not native for Mac and I don't have Docker for Mac (if I'm not wrong it changes its license few weeks ago).

On my defense I want to say that port 4000 is over 1024 and I can execute the application also as user space (it doesn't require to be executed as privileged user)


### Go routines

Go routines are ok but I could create functions and don't let them as anonymous functions. I know that it could be a
huge improvement to have a main function smaller.


### Saved numbers struct and type

To be more elegant I could use a struct like

```
type Challenge struct {
  numbers map[int]struct{}
  totalUniqueNumbers,duplicatedNumbers int
  uniqueNumbers int
}
```

to have all numbers well identified.

I also want to say that I used uniqueNumbers as an `int` type an improvement could be to use it as `bitInt` or `uint` type to be able to record more numbers with the same memory.


### Writing to a file

To improve that part I could create some cache/buffer with all new numbers between each stdout prompt. That pretty sure
is faster because it implies use less system calls (and maybe) less string formatting/parsing.

On the other hand writing numbers when I receive them it assures me that I write all numbers if the application receives a
"terminate" command or the application exits abnormally.


### Reading numbers from clients

I think that this part could also be more and more efficient, for example reading byte per byte from the client. That
lets us know if the byte is not a number, neither a part of "terminate" string, and then cut the connection because the input is not valid.

I expect to don't forget any edge case!

### Terminate application

Another part that could be improved is how to handle the program exit/termination. I'm not sure if when you exits the
main function all the clients are disconnected cleanly.

I think that is possible to create and manage a list of current active connections and `Close` them before exit the application.


## Conclusions

I enjoyed a lot developing that application and moreover I also learned a lot.

Is my fault but I want to say that I didn't have much time to develop the test.

I know that it could be done better but in overall I'm happy with the readability and the simplicity of the code without forget about the performance.

Thank you for this opportunity.
