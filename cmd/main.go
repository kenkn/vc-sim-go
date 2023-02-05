package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"vc-sim-go/models"
	"vc-sim-go/simulation"
	"vc-sim-go/state"

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

	flag.IntVar(&workerLimit, "wl", workerLimit, "worker limit")
	flag.IntVar(&jobLimit, "jl", jobLimit, "job limit")
	flag.Float64Var(&joiningRate, "jr", joiningRate, "joining rate")
	flag.Float64Var(&secessionRate, "sr", secessionRate, "secession rate")
	flag.Float64Var(&initialJoiningRate, "ijr", initialJoiningRate, "initial joining rate")
	flag.IntVar(&loopCount, "lc", loopCount, "loop count")
	flag.IntVar(&parallelismNum, "pn", parallelismNum, "parallelism num")
	flag.IntVar(&redundancy, "r", redundancy, "redundancy")
	flag.Parse()

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
	totalCycle := 0.0
	for i := 0; i < loopCount; i++ {
		workers := getInitializedWorkers(workerLimit)
		jobs := getInitializedJobs(jobLimit, parallelismNum)
		simulator := simulation.NewSimulator(workers, jobs, *config)
		simulator.SetWorkersState()
		simulator.SetWorkersParticipationRate()
		cycle := simulator.Simulate()
		fmt.Println(i, "'s cycle : ", cycle, "cycle / parallelism :", float64(cycle)/float64(parallelismNum))
		totalCycle += float64(cycle) / float64(parallelismNum)
	}
	result := []string{
		strconv.Itoa(workerLimit),
		strconv.Itoa(jobLimit),
		strconv.FormatFloat(joiningRate, 'f', 3, 64),
		strconv.FormatFloat(secessionRate, 'f', 3, 64),
		strconv.FormatFloat(initialJoiningRate, 'f', 3, 64),
		strconv.Itoa(loopCount),
		strconv.Itoa(parallelismNum),
		strconv.Itoa(redundancy),
		strconv.FormatFloat(float64(totalCycle)/float64(loopCount), 'f', 3, 64),
	}
	fmt.Println(result)
	f, err := os.OpenFile("result.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	err = w.Write(result)
	if err != nil {
		log.Fatal(err)
	}
	w.Flush()
}
