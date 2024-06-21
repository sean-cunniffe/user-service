package healthprobes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthProbe_ServeHTTP(t *testing.T) {
	p := newProbes()
	t.Run("test healthy probe", func(t *testing.T) {
		p.SetReady()
		req, err := http.NewRequest("GET", "/healthz", nil)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(p.ServeHTTP)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Equal(t, "OK", rr.Body.String())
	})

	t.Run("test unhealthy probe", func(t *testing.T) {
		p.SetUnReady()
		req, err := http.NewRequest("GET", "/healthz", nil)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(p.ServeHTTP)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

		assert.Equal(t, "Service Unavailable", rr.Body.String())
	})
}
