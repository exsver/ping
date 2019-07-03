package main

import (
	"log"
	"time"

	".."
)

func main() {
	target, err := ping.NewTargetFromString("8.8.8.8")
	if err != nil {
		log.Fatal(err)
	}
	target.Options.Count = 4                                    // Configure ping to send 4 ICMP ECHO_REQUESTs
	target.Options.Interval = time.Second * 1                   // with interval in 1 second
	pingResult, err := target.Ping(time.Now().Add(time.Minute)) // Run test
	if err == nil {
		println(pingResult.String())     // Print results
		println(pingResult.RttString())  // Print results
		println(pingResult.RttsString()) // Print results
	} else {
		println(err.Error())
	}
}
