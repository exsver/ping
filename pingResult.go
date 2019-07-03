package ping

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Rtt struct {
	ReplyTime time.Duration
	Err       error
}
type PingResult struct {
	IP          string
	TargetID    int
	Transmitted int
	Received    int
	Rtts        []Rtt
}

// Returns ping result as a string (like terminal utility)
// output example:
// "ping 8.8.8.8: 3 packets Transmitted, 3 Received"
func (pingResult *PingResult) String() string {
	return "ping " + pingResult.IP + ": " + strconv.Itoa(pingResult.Transmitted) + " packets Transmitted, " + strconv.Itoa(pingResult.Received) + " Received"
}

func (pingResult *PingResult) RttString() string {
	min, ave, max, standardDeviation := pingResult.CalculateRTT()
	return fmt.Sprintf("min: %v average: %v max: %v stddev: %v", min, ave, max, standardDeviation)
}

func (pingResult *PingResult) RttsString() string {
	result := ""
	for i, value := range pingResult.Rtts {
		if value.Err == nil {
			result = fmt.Sprintf("%s [%v] %v", result, i+1, value.ReplyTime)
		} else {
			result = fmt.Sprintf("%s [%v] %v", result, i+1, value.Err)
		}
	}
	return strings.TrimSpace(result)
}

// Calculete RTT for PingResult Struct
// Returns 0,0,0,0 in this cases:
// - if []Rtt is empty
// - host is Unreachable
func (pingResult *PingResult) CalculateRTT() (min time.Duration, average time.Duration, max time.Duration, standardDeviation time.Duration) {
	if len(pingResult.Rtts) == 0 {
		return
	}
	var total float64  //the sum of rtt
	var stotal float64 //the sum of the squares of the rtt
	received := len(pingResult.Rtts)
	for _, value := range pingResult.Rtts {
		if value.Err == nil {
			total += float64(value.ReplyTime)
			stotal += float64(value.ReplyTime) * float64(value.ReplyTime)
			if min == 0 {
				min = value.ReplyTime
			}
			if value.ReplyTime < min {
				min = value.ReplyTime
			}
			if value.ReplyTime >= max {
				max = value.ReplyTime
			}
		} else {
			received--
		}
	}
	if received > 0 {
		average = time.Duration(int(total) / received)
		// Пошагово вычисление стандартного отклонения (Standard deviation):
		// - вычисляем среднее арифметическое выборки данных
		// - отнимаем это среднее от каждого элемента выборки
		// - все полученные разницы возводим в квадрат
		// - суммируем все полученные квадраты
		// - делим полученную сумму на количество элементов в выборке (или на n-1, если n>30)
		// - вычисляем квадратный корень из полученного частного
		// или (более удобный для вычислений вариант)
		// standardDeviation = sqrt(stotal/received - average * average)   , если n<30
		// standardDeviation = sqrt((stotal/received - average * average)*(received/(received-1))    , если n>30
		standardDeviation = time.Duration(math.Sqrt(stotal/float64(received) - float64(average*average)))
		return
	}
	min = 0
	max = 0
	average = 0
	standardDeviation = 0
	return
}
