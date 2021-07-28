PRODUCTS = rtun rtun-server

.PHONY: all clean lint test deps

all: $(PRODUCTS)
	@:

clean:
	rm -f $(PRODUCTS)

lint:
	@./lint.sh

test:
	go test ./...

deps:
	go get -d ./...

rtun: deps agent/* agent/cmd/*
	go build -o $@ ./agent/cmd

rtun-server: deps server/* server/cmd/* server/tcp/* server/udp/* server/service/*
	go build -o $@ ./server/cmd
