package workers

type Manager struct {
	// A buffered channel that we can send work requests on.
	WorkQueue chan WorkRequest
	// A buffered channel that we are going to put the workers' work channels into.
	WorkerQueue chan chan WorkRequest
}

func NewManager() *Manager {
	m := Manager{make(chan WorkRequest, 10000), make(chan chan WorkRequest, 10000)}
	return &m;
}
