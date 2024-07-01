package util

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWatchDir(t *testing.T) {
	testDir := t.TempDir()
	t.Run("test watch for file creation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		testWg := sync.WaitGroup{}
		fileName := "testFile"
		testWg.Add(1)
		callOnce := true
		onChange := func(actualFilename string) {
			assert.Equal(t, filepath.Join(testDir, fileName), actualFilename)
			testWg.Done()
			assert.True(t, callOnce)
			callOnce = false
		}
		err := WatchPath(ctx, onChange, testDir)
		assert.Nil(t, err)

		createTestFile(t, testDir, fileName)
		testWg.Wait()
	})

	t.Run("test watch for file write", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		testWg := sync.WaitGroup{}
		fileName := "testFile"
		testWg.Add(1)
		callOnce := true
		onChange := func(actualFilename string) {
			assert.Equal(t, filepath.Join(testDir, fileName), actualFilename)
			testWg.Done()
			assert.True(t, callOnce)
			callOnce = false
		}
		err := WatchPath(ctx, onChange, testDir)
		assert.Nil(t, err)

		writeToFile(t, filepath.Join(testDir, fileName))
		testWg.Wait()
	})
}

func createTestFile(t *testing.T, testDir string, fileName string) {
	t.Helper()
	file, err := os.Create(filepath.Join(testDir, fileName))
	assert.Nil(t, err)
	defer file.Close()
}

func writeToFile(t *testing.T, file string) {
	t.Helper()
	err := os.WriteFile(file, []byte("test"), 0644)
	assert.Nil(t, err)
}
