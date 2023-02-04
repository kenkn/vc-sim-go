package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"vc-sim-go/models"
	"vc-sim-go/state"
	"vc-sim-go/simulation"

	"github.com/joho/godotenv"
)

func getInitializedWorkers(workerCount int) []*models.Worker {
	workers := make([]*models.Worker, workerCount)
	for i := range workers {
		workers[i] = models.NewWorker(i, state.UnavailableWorkerState)
	}
	return workers
}

func getInitializedJobs(jobCount int, parallelismNum int) []*models.Job {
	jobs := make([]*models.Job, jobCount)
	for i := range jobs {
		subjobs := make([]*models.Subjob, parallelismNum)
		for j := range subjobs {
			subjobs[j] = models.NewSubjob(j, state.UnallocatedSubjobState)
		}
		jobs[i] = models.NewJob(i, state.UnallocatedJobState, subjobs)
	}
	return jobs
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	workerLimit, err := strconv.Atoi(os.Getenv("WORKER_LIMIT"))
	if err != nil {
		log.Fatal("Error loading workerLimit")
	}
	jobLimit, err := strconv.Atoi(os.Getenv("JOB_LIMIT"))
	if err != nil {
		log.Fatal("Error loading jobLimit")
	}
	joiningRate, err := strconv.ParseFloat(os.Getenv("JOINING_RATE"), 64)
	if err != nil {
		log.Fatal("Error loading joiningRate")
	}
	secessionRate, err := strconv.ParseFloat(os.Getenv("SECESSION_RATE"), 64)
	if err != nil {
		log.Fatal("Error loading secessionRate")
	}
	initialJoiningRate, err := strconv.ParseFloat(os.Getenv("INITIAL_JOINING_RATE"), 32)
	if err != nil {
		log.Fatal("Error loading initialJoiningRate")
	}
	loopCount, err := strconv.Atoi(os.Getenv("LOOP_COUNT"))
	if err != nil {
		log.Fatal("Error loading loopCount")
	}
	parallelismNum, err := strconv.Atoi(os.Getenv("PARALLELISM_NUM"))
	if err != nil {
		log.Fatal("Error loading parallelismNum")
	}
	redundancy, err := strconv.Atoi(os.Getenv("REDUNDANCY"))
	if err != nil {
		log.Fatal("Error loading redundancy")
	}

	log.Println(fmt.Sprintf(`
ワーカ数: %d,
ジョブ数: %d,
参加率: %.3f,
離脱率: %.3f,
初期のワーカの参加率: %.3f,
ループ数: %d,
並列数: %d,
冗長度: %d`,
		workerLimit,
		jobLimit,
		joiningRate,
		secessionRate,
		initialJoiningRate,
		loopCount,
		parallelismNum,
		redundancy,
	))
	config := simulation.NewConfig(
		workerLimit,
		jobLimit,
		joiningRate,
		secessionRate,
		initialJoiningRate,
		loopCount,
		parallelismNum,
		redundancy,
	)
	for i := 0; i < loopCount; i++ {
		workers := getInitializedWorkers(workerLimit)
		jobs := getInitializedJobs(jobLimit, parallelismNum)
		simulator := simulation.NewSimulator(workers, jobs, *config)
		simulator.SetWorkersState()
		simulator.SetWorkersParticipationRate()
		cycle := simulator.Simulate()
		fmt.Println(i, "'s cycle : ", cycle, "cycle / parallelism :", float64(cycle) / float64(parallelismNum))
		simulator.Result.TotalCycle += cycle
	}
}
