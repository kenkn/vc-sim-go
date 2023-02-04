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
		workers[i] = models.NewWorker(i, state.UnavailableWorkerState, nil)
	}
	return workers
}

func getInitializedJobs(jobCount int, parallelismNum int) []*models.Job {
	jobs := make([]*models.Job, jobCount*parallelismNum)
	for i := range jobs {
		groupID := i / parallelismNum
		jobs[i] = models.NewJob(i, groupID, state.UnallocatedJobState, nil, false)
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
	dropoutRate, err := strconv.ParseFloat(os.Getenv("DROPOUT_RATE"), 64)
	if err != nil {
		log.Fatal("Error loading dropoutRate")
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

	log.Println(fmt.Sprintf(`
ワーカ数: %d,
ジョブ数: %d,
参加率: %.3f,
離脱率: %.3f,
初期のワーカの参加率: %.3f,
ループ数: %d,
並列数: %d`,
		workerLimit,
		jobLimit,
		joiningRate,
		dropoutRate,
		initialJoiningRate,
		loopCount,
		parallelismNum,
	))
	workers := getInitializedWorkers(workerLimit)
	jobs := getInitializedJobs(jobLimit, parallelismNum)
	simulator := simulation.NewSimulator(workers, jobs, parallelismNum)

	for i := 0; i < loopCount; i++ {
		simulator.SetWorkersState(joiningRate)
		simulator.SetWorkersParticipationRate(dropoutRate, joiningRate)
		cycle := simulator.Simulate()
		fmt.Println(i, "'s cycle : ", cycle)
		simulator.Result.TotalCycle += cycle
	}
	fmt.Println("total cycle : ", simulator.Result.TotalCycle)
}
