package api

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"
)

type HttpMetrics struct {
	DNS     time.Duration
	TCP     time.Duration
	TLS     time.Duration
	TTFB    time.Duration
	Status  int
	Headers http.Header
}

func CheckHttp(url string) (HttpMetrics, error) {
	var (
		start time.Time

		dnsStart, connectStart, tlsStart time.Time
		dnsDone, connectDone, tlsDone    bool

		firstByte time.Time
		gotTTFB   bool
	)

	var headers http.Header = make(http.Header)

	var dnsDur, tcpDur, tlsDur time.Duration

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			if !dnsStart.IsZero() {
				dnsDur = time.Since(dnsStart)
				dnsDone = true
			}
		},

		ConnectStart: func(network, addr string) {
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			if !connectStart.IsZero() {
				tcpDur = time.Since(connectStart)
				connectDone = true
			}
		},

		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if !tlsStart.IsZero() {
				tlsDur = time.Since(tlsStart)
				tlsDone = true
			}
		},

		GotFirstResponseByte: func() {
			firstByte = time.Now()
			gotTTFB = true
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HttpMetrics{}, err
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	start = time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return HttpMetrics{}, err
	}
	defer resp.Body.Close()

	headers = resp.Header

	ttfb := time.Duration(0)
	if gotTTFB {
		ttfb = firstByte.Sub(start)
	}

	_ = dnsDone
	_ = connectDone
	_ = tlsDone

	return HttpMetrics{
		DNS:     dnsDur,
		TCP:     tcpDur,
		TLS:     tlsDur,
		TTFB:    ttfb,
		Status:  resp.StatusCode,
		Headers: headers,
	}, nil
}
