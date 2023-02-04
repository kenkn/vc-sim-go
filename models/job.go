package models

import (
	"vc-sim-go/state"
)

type Job struct {
	ID      int
	State   state.JobState
	Subjobs []*Subjob
}

func NewJob(id int, state state.JobState, subjobs []*Subjob) *Job {
	return &Job{
		ID: id,
		State: state,
		Subjobs: subjobs,
	}
}

func (j *Job) Failed() {
	for _, subjob := range j.Subjobs {
		for _, aw := range subjob.AssignedWorker {
			if aw.State == state.RunningWorkerState {
				aw.State = state.AvailableWorkerState
			}
		}
		subjob.State = state.UnallocatedSubjobState
		subjob.AssignedWorker = make([]*Worker, 0)
	}
	j.State = state.UnallocatedJobState
}
