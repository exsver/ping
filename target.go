package ping

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Target struct represents configuration options for ping host
type Target struct {
	// UID. ID in range 0-65535 (2^16-1).
	ID int
	// IP address of the target host.
	IP      net.IP
	Options TargetOptions
}

type TargetOptions struct {
	// Stop after sending Count ECHO_REQUEST packets. Default is 1.
	Count int
	// Wait interval between sending each packet. Default is defined by const targetDefaultPingInterval.
	Interval time.Duration
	// Time to wait for a response. Default is defined by const targetDefaultPingInterval.
	Timeout time.Duration
	// Specifies the data bytes to be sent in ECHO_REQUEST packet.
	Data []byte
}

// NewTarget returns a new Target struct pointer
func NewTarget(ip net.IP, options TargetOptions) *Target {
	return &Target{
		IP:      ip,
		Options: options,
	}
}

// NewTargetFromString returns a new Target struct pointer with default values (Simple Constructor for Target)
func NewTargetFromString(ipString string) (*Target, error) {
	parsedIP := net.ParseIP(ipString)
	if parsedIP == nil {
		return nil, fmt.Errorf("create NewTargetFromString Error. Can't Parse IP from string: <%s>", ipString)
	}

	return NewTarget(
		parsedIP,
		TargetOptions{
			Count:    defaultConfig.TargetDefaultPingCount,
			Interval: defaultConfig.TargetDefaultPingInterval,
			Timeout:  defaultConfig.TargetDefaultPingTimeout,
			Data:     defaultConfig.TargetDefaultPingData,
		}), nil
}

// Printing an attributes of Target struct(IP, count, interval, timeout) to stdout
// Output example: "IP: 192.168.1.1, count: 4, interval: 1s, timeout: 4s"
// Use for debug purpose.
func (target *Target) String() string {
	return fmt.Sprintf("ip: %s, count: %v, interval: %s, timeout: %s", target.IP, target.Options.Count, target.Options.Interval, target.Options.Timeout)
}

func (target *Target) GenICMPMessage(i int) *icmp.Message {
	if target.IP.To4() != nil {
		return &icmp.Message{ // Create ICMP message
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   target.ID,
				Seq:  i,
				Data: target.Options.Data,
			},
		}
	}

	return &icmp.Message{ // Create ICMP message
		Type: ipv6.ICMPTypeEchoRequest,
		Code: 0,
		Body: &icmp.Echo{
			ID:   target.ID,
			Seq:  i,
			Data: target.Options.Data,
		},
	}
}
