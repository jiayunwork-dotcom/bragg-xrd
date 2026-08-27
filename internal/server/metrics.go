package server

import (
	"sync"
	"time"
)

type Metrics struct {
	mu       sync.Mutex
	Requests int64
	Errors   int64
	Started  time.Time
}

func NewMetrics() *Metrics { return &Metrics{Started: time.Now()} }

func (m *Metrics) Observe(status int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Requests++
	if status >= 400 {
		m.Errors++
	}
}

func (m *Metrics) Snapshot() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]interface{}{
		"requests": m.Requests, "errors": m.Errors,
		"uptime": time.Since(m.Started).Seconds(),
	}
}
