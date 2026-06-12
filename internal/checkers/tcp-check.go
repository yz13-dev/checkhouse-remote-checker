package api

import (
	"net"
	"time"
)

type TCPMetrics struct {
	ConnectTime time.Duration
	Error       error
}

func CheckTCP(host string, port string) (TCPMetrics, error) {
	address := net.JoinHostPort(host, port)

	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return TCPMetrics{
			Error: err,
		}, err
	}
	defer conn.Close()

	return TCPMetrics{
		ConnectTime: time.Since(start),
	}, err
}
