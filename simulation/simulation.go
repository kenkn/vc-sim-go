package simulation

import (
	"crypto/rand"
	"log"
	"math/big"
	"vc-sim-go/models"
	"vc-sim-go/state"
)

type Result struct {
	TotalCycle int
}

type Simulator struct {
	Workers        []*models.Worker
	Jobs           []*models.Job
	ParallelismNum int
	Result         Result
}

func NewSimulator(workers []*models.Worker, jobs []*models.Job, parallelismNum int) *Simulator {
	return &Simulator{
		Workers:        workers,
		Jobs:           jobs,
		ParallelismNum: parallelismNum,
		Result: Result{
			TotalCycle: 0,
		},
	}
}

func (s *Simulator) SetWorkersState(joiningRate float64) {
	for i := range s.Workers {
		if float64(i) < float64(len(s.Workers))*joiningRate {
			s.Workers[i].State = state.AvailableWorkerState
			continue
		}
		s.Workers[i].State = state.UnavailableWorkerState
	}
}

func (s *Simulator) SetWorkersParticipationRate(dropoutRate float64, joiningRate float64) {
	for i := range s.Workers {
		s.Workers[i].DropoutRate = dropoutRate
		s.Workers[i].JoiningRate = joiningRate
	}
}

func (s *Simulator) areAllJobsFinished() bool {
	for i := range s.Jobs {
		if s.Jobs[i].State != state.FinishedJobState {
			return false
		}
	}
	return true
}

func (s *Simulator) Simulate() int {
	cycle := 0
	for !s.areAllJobsFinished() {
		for i := 0; i < s.ParallelismNum; i++ {
			s.assignJobs()
			s.participationEvent()
			s.dropoffJobs()
		}
		cycle++
	}
	return cycle
}

func (s *Simulator) assignJobs() {
	subjobNum := len(s.Jobs)
	for i := 0; i < subjobNum; i++ {
		if s.Jobs[i].State == state.UnallocatedJobState {
			for j := range s.Workers {
				if s.Workers[j].State == state.AvailableWorkerState {
					if s.Workers[j].AssignedJob != nil || s.Jobs[i].AssignedWorker != nil {
						log.Fatal("Worker or Job is already assigned")
					}
					s.Workers[j].State = state.RunningWorkerState
					s.Jobs[i].State = state.ProcessingJobState
					s.Workers[j].AssignedJob = s.Jobs[i]
					s.Jobs[i].AssignedWorker = s.Workers[j]
					break
				}
			}
		}
	}
}

func (s *Simulator) participationEvent() {
	for i := range s.Workers {
		if s.Workers[i].State == state.RunningWorkerState || s.Workers[i].State == state.AvailableWorkerState {
			n, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				log.Fatal(err)
			}
			if n.Int64() < int64(s.Workers[i].DropoutRate*100) {
				err := s.Workers[i].Dropout()
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if s.Workers[i].State == state.UnavailableWorkerState {
			n, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				log.Fatal(err)
			}
			if n.Int64() < int64(s.Workers[i].JoiningRate*100) {
				s.Workers[i].Join()
			}
		}
	}
}

func (s *Simulator) dropoffJobs() {
	for i := range s.Workers {
		if s.Workers[i].State == state.RunningWorkerState {
			job := s.Workers[i].AssignedJob
			s.Workers[i].State = state.AvailableWorkerState
			s.Workers[i].AssignedJob.State = state.FinishedJobState
			s.Workers[i].AssignedJob = nil
			job.AssignedWorker = nil
		}
	}
}
