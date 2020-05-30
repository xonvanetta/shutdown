package shutdown

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChan(t *testing.T) {
	testShutdown(t, Chan())
}

func TestContext(t *testing.T) {
	testShutdown(t, Context().Done())
}

func TestWithContext(t *testing.T) {
	testShutdown(t, WithContext(context.Background()).Done())

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Millisecond)
		cancel()
	}()

	select {
	case <-WithContext(ctx).Done():
	case <-time.After(time.Second):
		assert.Fail(t, "shutdown never happened")
	}
}

func testShutdown(t *testing.T, s <-chan struct{}) {
	process, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)

	go func() {
		time.Sleep(time.Millisecond)
		err := process.Signal(os.Interrupt)
		assert.NoError(t, err)
	}()

	select {
	case <-s:
	case <-time.After(time.Second):
		assert.Fail(t, "shutdown never happened")
	}
}
