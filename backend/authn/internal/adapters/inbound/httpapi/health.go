package httpapi

import (
	"log/slog"
	"sync"

	"github.com/faber-numeris/beholder/backend/authn/internal/adapters/outbound"
)

type CheckFn func() bool

type HealthChecker struct {
	mu     sync.RWMutex
	checks []namedCheck
}

type namedCheck struct {
	name string
	fn   func() bool
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

func (hc *HealthChecker) Register(name string, fn func() bool) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks = append(hc.checks, namedCheck{name: name, fn: fn})
}

func (hc *HealthChecker) RegisterAdapter(name string, a outbound.Adapter) {
	hc.Register(name, a.Ping)
}

func (hc *HealthChecker) IsReady() bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	for _, c := range hc.checks {
		if !c.fn() {
			slog.Warn("Health check failed", "component", c.name)
			return false
		}
	}
	return true
}
