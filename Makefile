PRODUCTS = rtun rtun-server

.PHONY: all clean lint test

all: $(PRODUCTS)
	@:

clean:
	rm -f $(PRODUCTS)

lint:
	@./lint.sh

test:
	go test ./...

rtun: agent/* agent/cmd/*
	go build -o $@ ./agent/cmd

rtun-server: server/* server/cmd/* server/tcp/* server/udp/* server/service/*
	go build -o $@ ./server/cmd
