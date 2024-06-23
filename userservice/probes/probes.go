package healthprobes

import (
	"net/http"
	"sync"
)

// Probes is a struct that holds the health probe status
type Probes struct {
	mu                 sync.Mutex
	livenessProbePath  string
	readinessProbePath string
	liveness           bool
	readiness          bool
}

type ProbeConfig struct {
	LivenessProbePath  string `json:"livenessProbePath" yaml:"livenessProbePath"`
	ReadinessProbePath string `json:"readinessProbePath" yaml:"readinessProbePath"`
}

func NewProbes(config ProbeConfig) *Probes {
	return &Probes{
		liveness:           true,
		readiness:          false,
		livenessProbePath:  config.LivenessProbePath,
		readinessProbePath: config.ReadinessProbePath,
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
	switch r.URL.Path {
	case p.livenessProbePath:
		health = p.liveness
	case p.readinessProbePath:
		health = p.readiness
	default:
		w.WriteHeader(http.StatusNotFound)
		return
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
