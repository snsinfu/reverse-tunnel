package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

func main() {
	address := os.Args[1]
	count := int32(0)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&count, 1)
		log.Print("Request ", count)
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(address, nil))
}
