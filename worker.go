package workers

import (
	"log"
	"runtime"
	"strconv"
	"time"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest, handler func(worker *Worker, work WorkRequest)) *Worker {
	// Create, and return the worker.
	w := Worker{id, make(chan WorkRequest), workerQueue, make(chan bool, 1), handler}
	return &w
}

type WorkRequest struct {
	Name  string
	Delay time.Duration
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	Handler     func(w *Worker, work WorkRequest)
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	receivedQuitSignal := false
	go func() {
		for {
			if receivedQuitSignal == false {
				// Add ourselves into the worker queue if not going to stop
				w.WorkerQueue <- w.Work
			}

			select {
			case work := <-w.Work:
				// Receive a work request
				log.SetPrefix("[Worker " + strconv.Itoa(w.ID) + "] ")
				log.Printf("Received work request")

				w.Handler(w, work)

				if receivedQuitSignal != false {
					log.SetPrefix("[Worker " + strconv.Itoa(w.ID) + "] ")
					log.Printf("Terminating\n")
					return
				}
			case <-w.QuitChan:
				// Flag the worker to stop after next request
				// This is in order to remove the worker itself from WorkerQueue
				log.SetPrefix("[Worker " + strconv.Itoa(w.ID) + "] ")
				log.Printf("Going to stop after processing another work request\n")
				receivedQuitSignal = true
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		runtime.SetFinalizer(w, finalizer)
		//fmt.Printf("%#v", w)
		//log.Printf("worker %d Going to send signal quit\n", w.ID)
		w.QuitChan <- true
	}()
}

func finalizer(w *Worker) {
	log.Printf("Destroy worker %v\n", w.ID)
}
