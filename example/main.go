package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	workers "github.com/go-workers"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
	dispatcherWork workers.Dispatcher
	dispatcherHello workers.Dispatcher
)

func main() {
	// Parse the command-line flags.
	flag.Parse()

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	dispatcherWork = workers.Dispatcher{}
	dispatcherWork.Start(*NWorkers, workerWorkHandler)

	dispatcherHello = workers.Dispatcher{}
	dispatcherHello.Start(*NWorkers, workerHelloHandler)

	// Register our collector as an HTTP handler function.
	log.Println("Registering the hanlders")
	http.HandleFunc("/work", hanlderWork)
	http.HandleFunc("/hello", hanlderHello)

	// Start the HTTP server!
	log.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
