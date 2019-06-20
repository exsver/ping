package main

import (
	".."
	"time"
)

func main() {
	// Configure Targets
	target1, _ := ping.NewTargetFromString("8.8.8.8")
	target1.Options.Interval = time.Second * 1
	target1.Options.Count = 4
	// Id in range 0-65535 (2^16-1)
	target1.ID = 1

	target2, _ := ping.NewTargetFromString("8.8.4.4")
	target2.Options.Interval = time.Second * 1
	target2.Options.Count = 4
	target2.ID = 2

	// Create slice of Targets
	var targets []ping.Target
	// Add Targets to slice
	targets = append(targets, *target1, *target2)
	// Create chan for results
	resultChan := make(chan ping.PingResult, len(targets))
	// Run test
	go ping.ParallelPing(targets, resultChan, time.Now().Add(time.Minute), len(targets))
	// Read and print results
	for pr := range resultChan {
		println(pr.String(), pr.RttsString())
	}
}