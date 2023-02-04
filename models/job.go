package models

import (
	"vc-sim-go/state"
)

type Job struct {
	ID             int
	GroupID        int
	State          state.JobState
	// TODO Worker との二重結合を解消し、親子関係をはっきりさせる
	AssignedWorker *Worker
	IsAssigned     bool
}

func NewJob(id int, groupID int, state state.JobState, assignedWorker *Worker, isAssigned bool) *Job {
	return &Job{
		ID: id,
		GroupID: groupID,
		State: state,
		AssignedWorker: assignedWorker,
		IsAssigned: isAssigned,
	}
}
