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

rtun: agent/* config/*
	go build -o $@ ./agent

rtun-server: server/* config/*
	go build -o $@ ./server
