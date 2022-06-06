package ping

import (
	"sync"
	"time"
)

// ParallelPing - Run multiple ping tests in parallel.
func ParallelPing(targets []Target, resultChan chan PingResult, testDeadline time.Time, workers int) {
	taskChan := make(chan Target, len(targets))
	for _, task := range targets {
		taskChan <- task
	}
	close(taskChan)
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go worker(&wg, taskChan, resultChan, testDeadline)
	}
	wg.Wait()
	close(resultChan)
}

func worker(wg *sync.WaitGroup, taskChan chan Target, resultChan chan PingResult, testDeadline time.Time) {
	for pinger := range taskChan {
		var pr, _ = pinger.Ping(testDeadline)
		resultChan <- *pr
	}

	wg.Done()
}
