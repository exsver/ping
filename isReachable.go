package ping

import (
	"golang.org/x/net/ipv4"
	"net"
	"time"
)

// IsReachableIPv4 returns true if target is reachable
func (target *Target) IsReachableIPv4(testDeadline time.Time) (bool, error) {
	for i := 1; i <= 2; i++ {
		connection, err := newConnectionIPv4(target.Options.Timeout)
		if err != nil {
			return false, err
		}
		defer connection.Close()

		wm := target.GenICMPMessage(i)
		wb, err := wm.Marshal(nil) //Marshalling
		if err != nil {
			LogLevel.Fail.Printf("Marshaling error, %s", err)
			return false, err
		}

		_, err = connection.WriteTo(wb, &net.IPAddr{IP: target.IP}) // Write to connection
		if err != nil {
			return false, err
		}
		receivedMessage, _, err := findReplyIPv4(connection, target, i, testDeadline) // skip peer
		connection.Close()
		if receivedMessage != nil {
			LogLevel.Message.Printf("type:%v code:%v body:%v", receivedMessage.Type, receivedMessage.Code, receivedMessage.Body)
		}

		if err == nil {
			switch receivedMessage.Type {
			case ipv4.ICMPTypeEchoReply:
				return true, nil
			default:
				return false, nil
			}
		} else {
			continue
		}
	}

	return false, nil
}
