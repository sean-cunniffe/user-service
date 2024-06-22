package healthprobes

import (
	"net/http"
	"sync"
)

// Probes is a struct that holds the health probe status
type Probes struct {
	mu        sync.Mutex
	liveness  bool
	readiness bool
}

func newProbes() *Probes {
	return &Probes{
		liveness:  true,
		readiness: false,
	}
}

// SetReady sets the health probe to ready
func (p *Probes) SetLive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.liveness = true
}

// SetUnReady sets the health probe to unready
func (p *Probes) SetUnLive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.liveness = false
}

func (p *Probes) SetReady() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.readiness = true
}

func (p *Probes) SetUnReady() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.readiness = false
}

// ServeHTTP serves the health probe
func (p *Probes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	health := false
	if r.URL.Path != "/liveness" {
		health = p.liveness
	}
	if r.URL.Path != "/readiness" {
		health = p.readiness
	}

	if health {
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
