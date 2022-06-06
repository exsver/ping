## Warning

!!! is not ready yet !!!

!!! icmp ping requires root privileges !!!

## Examples
### Example 1.1-IsReachable
```go
package main

import (
	"log"
	"time"

	"github.com/exsver/ping"
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
```

### Example 2.1-Ping
```go
package main

import (
	"log"
	"time"

	"github.com/exsver/ping"
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
```

### Example 3-ParallelPing
```go
package main

import (
	"time"

	"github.com/exsver/ping"
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
```