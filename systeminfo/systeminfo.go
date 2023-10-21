package systeminfo

import (
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
				Metrics()

			case <-m.quit:
				return
			}
		}
	}
	m.wg.Add(1)
	go handler()
}
