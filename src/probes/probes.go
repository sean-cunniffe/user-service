package healthprobes

import (
	"net/http"
	"sync"
)

// Probes is a struct that holds the health probe status
type Probes struct {
	mu      sync.Mutex
	healthy bool
}

func newProbes() *Probes {
	return &Probes{
		healthy: true,
	}
}

// SetReady sets the health probe to ready
func (p *Probes) SetReady() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.healthy = true
}

// SetUnReady sets the health probe to unready
func (p *Probes) SetUnReady() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.healthy = false
}

// ServeHTTP serves the health probe
func (p *Probes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.healthy {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			panic(err) // will never reach here
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err := w.Write([]byte("Service Unavailable"))
		if err != nil {
			panic(err) // will never reach here
		}
	}
}
