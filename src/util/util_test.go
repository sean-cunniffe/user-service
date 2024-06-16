package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotNil(t *testing.T) {
	t.Run("test does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			NotNil("not nil")
			NotNil(3)
			NotNil(time.Hour)
		})
	})

	t.Run("test does panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NotNil(nil)
		})
	})
}
