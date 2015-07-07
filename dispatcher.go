package workers

import (
	"log"
	"math"
	"time"
)

type Dispatcher struct {
	Id                 string
	// Map used to stored created workers
	workers            map[int]Worker
	workRequestCounter int
	workerHandler      func(worker *Worker, work WorkRequest)
	// A buffered channel that we can send work requests on.
	WorkQueue chan WorkRequest
	// A buffered channel that we are going to put the workers' work channels into.
	WorkerQueue chan chan WorkRequest
}

// Start dispatcher for the given manager with nworkers as the initial number of
// workers
func (d *Dispatcher) Start(nworkers int, workerHandler func(worker *Worker, work WorkRequest)) {
	d.workerHandler = workerHandler
	d.workRequestCounter = 0
	d.workers = make(map[int]Worker)
	d.WorkQueue = make(chan WorkRequest, 10000)
	d.WorkerQueue = make(chan chan WorkRequest, 10000)

	tickerCheck := time.NewTicker(5 * time.Second)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		d.createWorker(i+1, d.WorkerQueue)
	}

	go func() {
		for {
			select {
			case <-tickerCheck.C:
				go func() {
					log.SetPrefix("[Dispatcher] ")
					log.Printf("Available slots in WorkerQueue: %v", len(d.WorkerQueue))
					log.Printf("Number of workers registered: %v", len(d.workers))

					if d.workRequestCounter > 10 {
						workersToCreate := int(math.Ceil(float64(d.workRequestCounter / 10)))
						log.SetPrefix("[Dispatcher] ")
						log.Printf("Going to add %v workers", workersToCreate)
						for i := 0; i < workersToCreate; i++ {
							d.createWorker(len(d.workers)+1, d.WorkerQueue)
						}
						//log.Printf("%#v", len(workers))
					} else if d.workRequestCounter < 5 && len(d.workers) > nworkers {
						//log.SetPrefix("[Dispatcher] ")
						//log.Printf("Going to stop worker %v\n", len(d.workers))

						// Stop 2 workers a time but do not drop below the initial number of workers
						d.removeWorker(len(d.workers))
						if len(d.workers) > nworkers {
							d.removeWorker(len(d.workers))
						}
					}
				}()
			case work := <-d.WorkQueue:
				go func() {
					log.SetPrefix("[Dispatcher] ")
					log.Println("Received work requeust")
					d.workRequestCount("+")
					worker := <-d.WorkerQueue
					worker <- work
					log.SetPrefix("[Dispatcher] ")
					log.Println("Dispatched work request")
					d.workRequestCount("-")
				}()
			}
		}
	}()
}

func (d *Dispatcher) workRequestCount(op string) {
	switch op {
	case "+":
		d.workRequestCounter++
	case "-":
		d.workRequestCounter--
	}
	log.Printf("WorkRequests in queue: %v\n", d.workRequestCounter)
	//log.Printf("%#v\n", workers)
}

func (d *Dispatcher) createWorker(id int, workerQueue chan chan WorkRequest) {
	log.SetPrefix("[Dispatcher] ")
	log.Println("Starting worker", id)
	worker := NewWorker(id, workerQueue, d.workerHandler)
	worker.Start()
	d.workers[id] = *worker
}

func (d *Dispatcher) removeWorker(id int) {
	w := d.workers[id]
	w.Stop()
	delete(d.workers, id)
}
