.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/timenewbank/go-mit/cmd/swarm
	go build -o bin/mit     github.com/timenewbank/go-mit/cmd/gmit
	go build -o bin/bootnode github.com/timenewbank/go-mit/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/timenewbank/go-mit/swarm/...
