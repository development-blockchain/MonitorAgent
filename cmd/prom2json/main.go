package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
	"net/http"
	"os"
	"strings"
	"time"
)

type Metrics struct {
	Labels map[string]string
}

type SysMonitor struct {
	Name string `json:"name"`
}

func makeTransport() *http.Transport {
	// Start with the DefaultTransport for sane defaults.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// Conservatively disable HTTP keep-alives as this program will only
	// ever need a single HTTP request.
	transport.DisableKeepAlives = true
	// Timeout early if the server doesn't even return the headers.
	transport.ResponseHeaderTimeout = time.Minute
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport.TLSClientConfig = tlsConfig
	return transport
}

var (
	url = flag.String("url", "http://127.0.0.1:9000/metrics", "metrics url")
)

func main() {
	flag.Parse()
	mfChan := make(chan *dto.MetricFamily, 1024)
	prom2json.FetchMetricFamilies(*url, mfChan, makeTransport())
	result := []*prom2json.Family{}
	for mf := range mfChan {
		if strings.HasPrefix(mf.GetName(), "monitor") == true {
			result = append(result, prom2json.NewFamily(mf))
		}
	}
	jsonText, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error marshaling JSON:", err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(jsonText); err != nil {
		fmt.Fprintln(os.Stderr, "error writing to stdout:", err)
		os.Exit(1)
	}
	fmt.Println()
}
