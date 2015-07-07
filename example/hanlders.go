package main

import (
	"workers"
	"log"
	"net/http"
	"time"
	"fmt"
)

func workerWorkHandler(worker *workers.Worker, work workers.WorkRequest) {
	time.Sleep(work.Delay)

	switch work.Name {
		default:
			fmt.Printf("Hello, %s!\n", work.Name)
	}	
}

func workerHelloHandler(worker *workers.Worker, work workers.WorkRequest) {

}

func hanlderHello(w http.ResponseWriter, r *http.Request) {
}

func hanlderWork(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the delay.
	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check to make sure the delay is anywhere from 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Now, we retrieve the person's name from the request.
	name := r.FormValue("name")

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	// Now, we take the delay, and the person's name, and make a WorkRequest out of them.
	work := workers.WorkRequest{Name: name, Delay: delay}

	// Push the work onto the queue.
	dispatcherWork.WorkQueue <- work
	log.SetPrefix("[Handler] ")
	log.Println("Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}
