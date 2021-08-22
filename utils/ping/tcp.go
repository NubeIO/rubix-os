package utils

import (
	"fmt"
	"net"
	"time"
)


//PingPort will ping an address and port, example 192.168.15.10:1515
func PingPort(network string, timeout time.Duration) (message string, fail bool) {
	var seqNumber uint64 = 0
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", network, timeout)
	endTime := time.Now()
	if err != nil {
		st := startTime.Format("[2006-01-02T15:04:05]:")
		msg := fmt.Sprintf("%s:%s", st, "connection failed")
		return msg, true
	} else {
		defer conn.Close()
		et := float64(endTime.Sub(startTime)) / float64(time.Millisecond)
		st := startTime.Format("[2006-01-02T15:04:05]:")
		addr := fmt.Sprintf(" addr=%s seq=%d time=%4.2fms", conn.RemoteAddr().String(), seqNumber, et)
		msg := fmt.Sprintf("%s %s", st, addr)
		return msg, false
	}
}