package api

import (
	"context"
	"net"
	"time"
)

type DNSMetrics struct {
	LookupTime time.Duration
	IPs        []string
	Error      error
}

func CheckDNS(host string) (DNSMetrics, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	start := time.Now()

	resolver := &net.Resolver{
		PreferGo: true,
	}

	ips, err := resolver.LookupHost(ctx, host)

	return DNSMetrics{
		LookupTime: time.Since(start),
		IPs:        ips,
		Error:      err,
	}, err
}
