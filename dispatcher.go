package workers

import (
	"fmt"
	"log"
	"time"
	"math"
	//"sync"
)

type Dispatcher struct {
	Id string
	workers map[int]Worker
	nbCurrentWorkers int
	workRequestCounter int
}

// Start dispatcher for the given manager with nworkers as the initial number of
// workers
func (d *Dispatcher) Start(nworkers int, manager *Manager) {
	d.workers = make(map[int]Worker)
	d.nbCurrentWorkers = nworkers
	d.workRequestCounter = 0
	tickerCheck := time.NewTicker(5 * time.Second)
	
	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		d.createWorker(i + 1, manager.WorkerQueue)
	}
	
	go func() {
		for {
			select {
			case <-tickerCheck.C:
				go func() {
					log.SetPrefix("[Dispatcher] ")
					log.Printf("Available slots in WorkerQueue: %v", len(manager.WorkerQueue))
					log.Printf("Number of workers registered: %v", len(d.workers))

					if d.workRequestCounter > 10 {
						workersToCreate := int(math.Ceil(float64(d.workRequestCounter / 10)))
						log.SetPrefix("[Dispatcher] ")
						log.Printf("Going to add %v workers", workersToCreate)
						for i := 0; i < workersToCreate; i++ {
							d.createWorker(d.nbCurrentWorkers + 1, manager.WorkerQueue)
						}
						//log.Printf("%#v", len(workers))
					} else if d.workRequestCounter < 5 && d.nbCurrentWorkers > nworkers {
						//log.SetPrefix("[Dispatcher] ")
						//log.Printf("Going to stop worker %v\n", nbCurrentWorkers)
						
						// Stop 2 workers a time but do not drop below the initial number of workers
						d.removeWorker(d.nbCurrentWorkers)
						if d.nbCurrentWorkers > nworkers {
							d.removeWorker(d.nbCurrentWorkers)
						}
					}
				}()
			case work := <-manager.WorkQueue:
				go func() {
					log.SetPrefix("[Dispatcher] ")
					log.Println("Received work requeust")
					d.workRequestCount("+")
					worker := <-manager.WorkerQueue
					worker <- work
					log.SetPrefix("[Dispatcher] ")					
					log.Println("Dispatching work request")
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
	d.nbCurrentWorkers = id
	//fmt.SetPrefix("[Dispatcher] ")
	fmt.Println("Starting worker", id)
	worker := NewWorker(id, workerQueue)
	worker.Start()
	d.workers[id] = *worker
}

func (d *Dispatcher) removeWorker(id int) {
	w := d.workers[id]
	w.Stop()
	delete(d.workers, id)
	d.nbCurrentWorkers--
}