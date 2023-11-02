package common

import (
	"bytes"
	"net/http"
	"os"
)

var (
	nodename = ""
)

func SetHostName(name string) {
	nodename = name
}

func GetHostName() string {
	if nodename == "" {
		hostname, _ := os.Hostname()
		return hostname
	}
	return nodename
}

func GetExternal() (string, error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	ip := buf.String()
	return ip, nil
}
