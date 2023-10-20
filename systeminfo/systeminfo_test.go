package systeminfo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSystemInfo(T *testing.T) {
	info := SystemInfo()
	d, _ := json.Marshal(info)
	fmt.Printf("%s\n", string(d))
}

func TestPrometheus(t *testing.T) {
	txt := Prometheus(SystemInfo())
	fmt.Printf("%s\n", txt)
}
