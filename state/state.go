package state

type WorkerState int

const (
	UnavailableWorkerState WorkerState = iota
	AvailableWorkerState
	RunningWorkerState
)

type JobState int

const (
	UnallocatedJobState JobState = iota
	ProcessingJobState
	FinishedJobState
)
