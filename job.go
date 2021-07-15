package dispatcher

// Job represents the job to be run
type Job struct {
	DoJob func()
}

// A buffered channel that we can send work requests on.
var JobQueue chan *Job

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) *Worker {
	JobQueue = make(chan *Job)
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				job.Run()

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	w.quit <- true
}

// NewTask initializes a new task based on a given work
// function.
func NewJob(f func()) *Job {
	return &Job{DoJob: f}
}

// Run runs a Task and does appropriate accounting via a
// given sync.WorkGroup.
func (j *Job) Run() {
	j.DoJob()
}
