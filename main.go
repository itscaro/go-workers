package workers

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"log"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

func main() {
	// Parse the command-line flags.
	flag.Parse()

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	dispatcherWork := Dispatcher{}
	managerWork := NewManager()
	dispatcherWork.StartDispatcher(*NWorkers, managerWork)

	dispatcherHello := Dispatcher{}
	managerHello := NewManager()
	dispatcherHello.StartDispatcher(*NWorkers, managerHello)

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", managerWork.Work)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
