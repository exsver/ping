package ping

import (
	"time"
)

type Config struct {
	// Specifies the default value of Counter field for Target constructors
	TargetDefaultPingCount int
	// Wait interval between sending each packet.
	TargetDefaultPingInterval time.Duration
	TargetMaxPingInterval     time.Duration
	// Time to wait for a response.
	// The default value of this timeout is:
	//  - 2 seconds on Cisco routers;
	//  - 10 seconds on Linux systems;
	//  - 4 seconds on Windows.
	TargetDefaultPingTimeout time.Duration
	TargetMinPingTimeout     time.Duration
	TargetMaxPingTimeout     time.Duration
	// Specifies the maximum number of data bytes in echo requests.
	// For interfaces with mtu 1500 (IPv4) set TargetMaxPingDataLenth = 1458
	TargetMaxPingDataLenth int
	TargetDefaultPingData  []byte
}

var defaultConfig = Config{
	TargetDefaultPingCount:    1,
	TargetDefaultPingInterval: time.Millisecond * 100,
	TargetMaxPingInterval:     time.Second * 10,
	TargetDefaultPingTimeout:  time.Second * 2,
	TargetMinPingTimeout:      time.Millisecond * 10,
	TargetMaxPingTimeout:      time.Second * 10,
	TargetMaxPingDataLenth:    1458,
	TargetDefaultPingData:     []byte("abcdefghijklmnopqrstuvwabcdefghi"), //windows like (32 bytes + 8 bytes of ICMP header data + 20 bytes of IP header + Ethernet frame data)
}

var NoDeadline = new(time.Time)
