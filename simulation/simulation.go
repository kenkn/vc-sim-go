package simulation

import (
	"log"
	"math/rand"
	"time"
	"vc-sim-go/models"
	"vc-sim-go/state"
)

type Result struct {
	TotalCycle int
}

type Simulator struct {
	Workers                   []*models.Worker
	Jobs                      []*models.Job
	Config                    Config
	AvailableWorkerStateCount int
	FinishedJobStateCount     int
	Result                    Result
}

func NewSimulator(workers []*models.Worker, jobs []*models.Job, config Config) *Simulator {
	return &Simulator{
		Workers:                   workers,
		Jobs:                      jobs,
		Config:                    config,
		AvailableWorkerStateCount: 0,
		FinishedJobStateCount:     0,
		Result: Result{
			TotalCycle: 0,
		},
	}
}

func (s *Simulator) SetWorkersState() {
	for i, w := range s.Workers {
		if float64(i) < float64(len(s.Workers))*s.Config.InitialJoiningRate {
			s.AvailableWorkerStateCount++
			w.State = state.AvailableWorkerState
			continue
		}
		w.State = state.UnavailableWorkerState
	}
}

func (s *Simulator) SetWorkersParticipationRate() {
	for i := range s.Workers {
		s.Workers[i].SecessionRate = s.Config.SecessionRate
		s.Workers[i].JoiningRate = s.Config.JoiningRate
	}
}

func (s *Simulator) areAllJobsFinished() bool {
	if s.FinishedJobStateCount == len(s.Jobs) {
		return true
	}
	return false
}

func (s *Simulator) Simulate() int {
	cycle := 0
	for !s.areAllJobsFinished() {
		s.assignJobs()
		s.workerSecessionEvent()
		s.finishJobs()
		s.workerJoinEvent()
		cycle++
	}
	return cycle
}

func (s *Simulator) countAvailableWorker() int {
	c := 0
	for _, worker := range s.Workers {
		if worker.State == state.AvailableWorkerState {
			c++
		}
	}
	return c
}

func (s *Simulator) assignJobs() {
label:
	for _, job := range s.Jobs {
		if (s.Config.ParallelismNum * s.Config.Redundancy) > s.AvailableWorkerStateCount {
			return
		}
		if job.State != state.UnallocatedJobState {
			continue
		}
		for _, subjob := range job.Subjobs {
			if subjob.State != state.UnallocatedSubjobState {
				continue
			}
			for i := 0; i < s.Config.Redundancy; i++ {
				for _, worker := range s.Workers {
					if worker.State != state.AvailableWorkerState {
						continue
					}
					s.AvailableWorkerStateCount--
					worker.State = state.RunningWorkerState
					worker.AssignedSubjob = subjob
					subjob.AssignedWorker = append(subjob.AssignedWorker, worker)
					subjob.State = state.ProcessingSubjobState
					break
				}
			}
			if subjob.State == state.UnallocatedSubjobState {
				break label
			}
		}
		job.State = state.ProcessingJobState
	}
}

func (s *Simulator) workerSecessionEvent() {
	for _, worker := range s.Workers {
		if worker.State == state.UnavailableWorkerState {
			continue
		}
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(100)
		if n < int(worker.SecessionRate*100) {
			if worker.State != state.RunningWorkerState {
				continue
			}
			for i, aw := range worker.AssignedSubjob.AssignedWorker {
				if aw.ID == worker.ID {
					worker.AssignedSubjob.RemoveWorker(i)
					break
				}
			}
			if worker.State == state.AvailableWorkerState {
				s.AvailableWorkerStateCount--
			}
			err := worker.Secession()
			if err != nil {
				log.Fatal(err)
			}

		}
	}
	for _, job := range s.Jobs {
		for _, subjob := range job.Subjobs {
			if len(subjob.AssignedWorker) == 0 && subjob.State == state.ProcessingSubjobState {
				beAvailableCount := job.Failed()
				s.AvailableWorkerStateCount += beAvailableCount
			}
		}
	}
}

func (s *Simulator) finishJobs() {
	for _, job := range s.Jobs {
		if job.State != state.ProcessingJobState {
			continue
		}
		for _, subjob := range job.Subjobs {
			if subjob.State != state.ProcessingSubjobState {
				continue
			}
			for _, aw := range subjob.AssignedWorker {
				s.AvailableWorkerStateCount++
				aw.State = state.AvailableWorkerState
			}
			subjob.State = state.FinishedSubjobState
		}
		s.FinishedJobStateCount++
		job.State = state.FinishedJobState
	}
}

func (s *Simulator) workerJoinEvent() {
	for _, worker := range s.Workers {
		if worker.State != state.UnavailableWorkerState {
			continue
		}
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(100)
		if n < int(worker.JoiningRate*100) {
			s.AvailableWorkerStateCount++
			err := worker.Join()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
