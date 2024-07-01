package configuration

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	defaultConfigFile = "test_config/test_config.yaml"
	t.Run("can read default config", func(t *testing.T) {
		config := ReadEnvConfig()
		if reflect.DeepEqual(config, Config{}) {
			t.Error("config is empty")
		}
	})

	t.Run("panic file doesnt exist", func(t *testing.T) {
		defaultConfigFile = "nofile.yaml"
		assert.Panics(t, func() {
			ReadEnvConfig()
		})
	})

	t.Run("panic invalid file", func(t *testing.T) {
		defaultConfigFile = "test_config/invalid_config.yaml"
		assert.Panics(t, func() {
			ReadEnvConfig()
		})
	})
}
