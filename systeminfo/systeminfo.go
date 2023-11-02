package systeminfo

import (
	"github.com/develope/MonitorAgent/log"
	"sync"
	"time"
)

var (
	interval = time.Second * 5
)

type SystemMonitor struct {
	quit chan struct{}
	wg   sync.WaitGroup
}

func (m *SystemMonitor) Start() {
	m.monitor()
}

func (m *SystemMonitor) Stop() {
	close(m.quit)
	m.wg.Wait()
}

func (m *SystemMonitor) monitor() {
	handler := func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		defer m.wg.Done()
		for {
			select {
			case <-ticker.C:
				l := log.Entry()
				Metrics()
				l.Info("finish metrics one time")

			case <-m.quit:
				return
			}
		}
	}
	m.wg.Add(1)
	go handler()
}
