package syncutil_test

import (
	"TinyKVStore/internal/helpers/syncutil"
	"testing"
	"time"
)

func TestNewEmitterAwaiter(t *testing.T) {
	emitEvent, waitForEvent := syncutil.NewEmitterAwaiter()

	doneCh := make(chan struct{})
	go func() {
		time.Sleep(2 * time.Second)
		emitEvent()
	}()

	go func() {
		defer close(doneCh)
		waitForEvent()
	}()

	select {
	case <-doneCh:
		return
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}
