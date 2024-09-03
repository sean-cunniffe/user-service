package healthprobes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var probeConfig = ProbeConfig{
	LivenessProbePath:  "/liveness",
	ReadinessProbePath: "/readiness",
}

func TestProbes(t *testing.T) {
	t.Run("test probe setters", func(t *testing.T) {
		probe := NewProbes(probeConfig)
		probe.SetLive()
		assert.True(t, probe.liveness)
		probe.SetUnLive()
		assert.False(t, probe.liveness)
	})

	t.Run("test probe http handler", func(t *testing.T) {
		livenessProbePath := "/liveness"
		readinessProbePath := "/readiness"
		probe := NewProbes(probeConfig)

		testExpectedStatus := func(path string, expectedStatus int) {
			req := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()
			probe.ServeHTTP(w, req)
			assert.Equal(t, expectedStatus, w.Code)
		}

		probe.SetLive()
		testExpectedStatus(livenessProbePath, http.StatusOK)
		probe.SetUnLive()
		testExpectedStatus(livenessProbePath, http.StatusServiceUnavailable)

		probe.SetReady()
		testExpectedStatus(readinessProbePath, http.StatusOK)
		probe.SetUnReady()
		testExpectedStatus(readinessProbePath, http.StatusServiceUnavailable)
	})

	t.Run("test probe http handler with invalid path", func(t *testing.T) {
		probe := NewProbes(probeConfig)
		req := httptest.NewRequest("GET", "/invalid", nil)
		w := httptest.NewRecorder()
		probe.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
