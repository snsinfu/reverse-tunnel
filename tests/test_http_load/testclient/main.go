package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main() {
	address := os.Args[1]
	url := "http://" + address + "/"

	requests, err := strconv.Atoi(os.Args[2])
	if err != nil {
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < requests; i++ {
		id := i + 1

		go func() {
			wg.Add(1)
			defer wg.Done()

			res, err := http.Get(url)
			if err != nil {
				log.Print("error: ", err)
				return
			}
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Print("error: ", err)
				return
			}

			log.Print(string(data), " ", id)
			return
		}()
	}

	wg.Wait()
}
