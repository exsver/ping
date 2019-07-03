package main

import (
	"log"
	"time"

	".."
)

func main() {
	target, err := ping.NewTargetFromString("8.8.8.8") // Create New target
	if err != nil {
		log.Fatal(err)
	}
	result, err := target.IsReachableIPv4(time.Now().Add(time.Minute)) // Run test
	if err == nil {
		println(result) // Print result. IsReachable return true if target is reachable, or false if not.
	} else {
		println(err.Error()) // or error
	}
}
