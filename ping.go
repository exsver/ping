package ping

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func newConnectionIPv4(timeout time.Duration) (net.PacketConn, error) {
	var connection net.PacketConn
	connection, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		// Possible errors:
		// "listen ip4:icmp :bind: The requested address is not valid in its context"
		// "listen ip4:icmp 0.0.0.0: socket: operation not permitted"   -  need root privilege
		return nil, fmt.Errorf("listen error, %s", err)
	}
	err = connection.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, fmt.Errorf("setDeadline error, %s", err)
	}
	return connection, nil
}

func prepareICMPMessage(target *Target, i *int) *icmp.Message {
	if target.IP.To4() != nil {
		return &icmp.Message{ //Create ICMP message
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   target.ID,
				Seq:  *i,
				Data: target.Options.Data,
			},
		}
	} else {
		return &icmp.Message{ //Create ICMP message
			Type: ipv6.ICMPTypeEchoRequest,
			Code: 0,
			Body: &icmp.Echo{
				ID:   target.ID,
				Seq:  *i,
				Data: target.Options.Data,
			},
		}
	}
}

func findReplyIPv4(connection net.PacketConn, target *Target, i int, testDeadline time.Time) (*icmp.Message, net.IP, error) {
	//IPv4: 32(56) = add 8 bytes of ICMP header data + 20 or 24 bytes of IP header
	//IPv6: 44 = add 4 bytes of ICMPv6 header + 40 bytes of IPv6 header
	receivedBytes := make([]byte, len(target.Options.Data)+56)
	for {
		if testDeadline != *NoDeadline && time.Now().After(testDeadline) {
			return nil, nil, errors.New("Deadline reached. Not enough time to run this test.")
		}
		n, peer, err := connection.ReadFrom(receivedBytes)
		//Possible errors:
		// - "read ip4 0.0.0.0: i/o timeout" = "Request timed out."
		// - "read ip4 0.0.0.0: wsarecvfrom: A message sent on a datagram socket was larger than the internal message buffer or some other network limit, or the buffer used to receive a datagram into was smaller than the datagram itself."
		LogLevel.Message.Printf("Received n:%v peer:%v err:%v ", n, peer, err)

		if n == 0 && peer == nil && err != nil { // Connection timeout
			return nil, nil, errors.New("Request timed out")
		}

		if err != nil { // Error. Something wrong!!!
			continue
		}

		receivedMessage, err := icmp.ParseMessage(1, receivedBytes[:n])
		if err != nil {
			continue
		}

		switch receivedMessage.Type {
		case ipv4.ICMPTypeEcho: // Ignore echo request
			LogLevel.Message.Printf("Ignored n:%v peer:%v err:%v Type:%v Code:%v", n, peer, err, receivedMessage.Type, receivedMessage.Code)
			continue
		case ipv4.ICMPTypeEchoReply:

			/*
				ID receivedBytes[4:6]
				Seq receivedBytes[6:8]
			*/

			body, _ := receivedMessage.Body.(*icmp.Echo)
			if body.ID == target.ID && body.Seq == i && peer.String() == target.IP.String() {
				return receivedMessage, net.ParseIP(peer.String()), err
			}

		case ipv4.ICMPTypeDestinationUnreachable, ipv4.ICMPTypeTimeExceeded, ipv4.ICMPTypeParameterProblem:
			LogLevel.Message.Printf("Received body: %v", receivedMessage.Body)
			/*
				Body:
						original IP header
					TTL :					[16]
					Protocol : 				[17]
					HeaderChecksum : 		[18:20]
					Source : 				[20:24]
					Destination : 			[24:28]
						original ICMP header
					Type :					[28]
					Code : 					[29]
					Checksum : 				[30:32]
					ID : 					[32:34]
					Seq : 					[34:36]
					Data : 					[36:]
			*/
			rDestination := net.IP(receivedBytes[24:28])
			rID := int(binary.BigEndian.Uint16(receivedBytes[32:34]))
			rSeq := int(binary.BigEndian.Uint16(receivedBytes[34:36]))
			LogLevel.Message.Printf("Peer: %v Destination: %v ID: %v Seq: %v", peer, rDestination, rID, rSeq)
			if rID == target.ID && rDestination.Equal(target.IP) && rSeq == i {
				return receivedMessage, net.ParseIP(peer.String()), err
			}
		default:
			continue
		}
	}
}

// PingIPv4 send ICMP ECHO_REQUESTs to target host
func (target *Target) PingIPv4(testDeadline time.Time) (*PingResult, error) {
	var (
		err error
		wb  []byte
	)
	result := &PingResult{TargetID: target.ID, IP: target.IP.String()}

	for i := 1; i <= target.Options.Count; i++ {
		var connection net.PacketConn
		connection, err = newConnectionIPv4(target.Options.Timeout)
		if err != nil {
			break
		}
		defer connection.Close()

		wm := prepareICMPMessage(target, &i)
		wb, err = wm.Marshal(nil) //Marshalling
		if err != nil {
			break
		}
		start := time.Now()
		if _, err := connection.WriteTo(wb, &net.IPAddr{IP: target.IP}); err != nil {
			result.Rtts = append(result.Rtts, Rtt{Err: err})
			// Possible errors:
			// - "network is unreachable"
			LogLevel.Fail.Printf("WriteTo error, %s", err)
			continue
		}
		result.Transmitted++ // increment counter
		receivedMessage, peer, err := findReplyIPv4(connection, target, i, testDeadline)
		if err != nil {
			result.Rtts = append(result.Rtts, Rtt{Err: err})
			if err.Error() == "Deadline reached. Not enough time to run this test." {
				result.Transmitted--
			}
			continue
		}

		elapsedTime := time.Now().Sub(start)
		connection.Close()
		LogLevel.Message.Printf("Find Message type:%v, code:%v, body:%v", receivedMessage.Type, receivedMessage.Code, receivedMessage.Body)
		switch receivedMessage.Type {
		case ipv4.ICMPTypeEchoReply:
			result.Received++
			result.Rtts = append(result.Rtts, Rtt{ReplyTime: elapsedTime})
		case ipv4.ICMPTypeDestinationUnreachable:
			result.Rtts = append(result.Rtts, Rtt{Err: errors.New(peer.String() + " : Destination Unreachable : " + IPv4DestinationUnreachableCode[receivedMessage.Code])})
		case ipv4.ICMPTypeTimeExceeded:
			result.Rtts = append(result.Rtts, Rtt{Err: errors.New(peer.String() + " : Time-to-live exceeded : " + IPv4TimeExceededCode[receivedMessage.Code])})
		case ipv4.ICMPTypeParameterProblem:
			result.Rtts = append(result.Rtts, Rtt{Err: errors.New(peer.String() + " : Parameter Problem")})
		default:
			result.Rtts = append(result.Rtts, Rtt{Err: errors.New(peer.String() + " : Unknown ICMP Type")})
			LogLevel.Message.Printf("Unknown ICMP Type %+v", receivedMessage)
		}
		time.Sleep(target.Options.Interval - elapsedTime)
	}
	return result, err
}
