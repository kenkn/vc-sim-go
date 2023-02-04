package models

import (
	"errors"
	"vc-sim-go/state"
)

type Worker struct {
	ID          int
	State       state.WorkerState
	SecessionRate float64
	JoiningRate float64
	AssignedSubjob *Subjob
}

func NewWorker(id int, state state.WorkerState) *Worker {
	return &Worker{
		ID:          id,
		State:       state,
		AssignedSubjob: nil,
	}
}

func (w *Worker) Secession() error {
	if w.State != state.RunningWorkerState && w.State != state.AvailableWorkerState {
		return errors.New("Worker is not available")
	}
	w.State = state.UnavailableWorkerState
	w.AssignedSubjob = nil
	return nil
}

func (w *Worker) Join() error {
	if w.State != state.UnavailableWorkerState {
		return errors.New("Worker is not available")
	}
	w.State = state.AvailableWorkerState
	w.AssignedSubjob = nil
	return nil
}
