package synctest

import (
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// 测试共享值的可见性 go test -run TestSharedValue -count=1000
func TestSharedValue1(t *testing.T) {
	var shared atomic.Int64

	go func() {
		shared.Store(1)
		time.Sleep(1 * time.Microsecond)
		shared.Store(2)
	}()

	time.Sleep(10 * time.Microsecond)
	if shared.Load() != 2 {
		t.Errorf("shared = %d, want 2", shared.Load())
	}
}

func TestSharedValue2(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var shared atomic.Int64

		go func() {
			shared.Store(1)
			time.Sleep(1 * time.Microsecond)
			shared.Store(2)
		}()

		time.Sleep(5 * time.Microsecond)
		if shared.Load() != 2 {
			t.Errorf("shared = %d, want 2", shared.Load())
		}
	})
}
