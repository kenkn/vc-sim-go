package models

import (
	"fmt"
	"vc-sim-go/state"
)

type Subjob struct {
	ID int
	State state.SubjobState
	AssignedWorker []*Worker
}

func NewSubjob(id int, state state.SubjobState) *Subjob {
	return &Subjob{
		ID: id,
		State: state,
		AssignedWorker: make([]*Worker, 0),
	}
}

func (sj *Subjob) RemoveWorker(i int) error {
	if i >= len(sj.AssignedWorker) || i < 0 {
		return fmt.Errorf("Index is out of range. Index is %d with slice length %d", i, len(sj.AssignedWorker))
	}
	sj.AssignedWorker[i] = sj.AssignedWorker[len(sj.AssignedWorker)-1]
	sj.AssignedWorker = sj.AssignedWorker[:len(sj.AssignedWorker)-1]
	return nil
}
