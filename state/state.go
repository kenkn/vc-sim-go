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

type SubjobState int

const (
	UnallocatedSubjobState SubjobState = iota
	ProcessingSubjobState
	FinishedSubjobState
)
