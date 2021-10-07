.PHONY: run vendor clean build

run: main
	./main

main: main.go
	go build main.go

vendor:
	go get golang.org/x/net/netutil
	go mod vendor

clean:
	rm main

build: main
