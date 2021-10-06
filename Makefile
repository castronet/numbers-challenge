run: main
	./main

main: main.go
	go build main.go

modules:
	go get golang.org/x/net/netutil
