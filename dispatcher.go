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
	nbRequests			int
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
	d.nbRequests = 0
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
					counter := len(d.WorkQueue)
					log.SetPrefix("[Dispatcher] ")
					log.Printf(
						"nbWorkerQueue - nbWorkQueue - nbWorkers: %v - %v - %v",
						len(d.WorkerQueue), counter, len(d.workers),
					)

					if counter > 10 {
						workersToCreate := int(math.Ceil(float64(counter / 10)))
						log.SetPrefix("[Dispatcher] ")
						log.Printf("Going to add %v workers", workersToCreate)
						for i := 0; i < workersToCreate; i++ {
							d.createWorker(len(d.workers)+1, d.WorkerQueue)
						}
						//log.Printf("%#v", len(workers))
					} else if counter < 5 && len(d.workers) > nworkers {
						// Stop 2 workers a time but do not drop below the initial number of workers
						d.removeWorker(len(d.workers))
						if len(d.workers) > nworkers {
							d.removeWorker(len(d.workers))
						}
					}
				}()
			case work := <-d.WorkQueue:
				// Count the requests received	
				d.nbRequests++
				log.SetPrefix("[Dispatcher] ")
				log.Println("Received work request")
				worker := <-d.WorkerQueue
				worker <- work
			}
		}
	}()
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
