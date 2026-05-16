package factoriat

import (
	"sync/atomic"
	"testing"
	"time"
)

func newAlwaysTriggerFactoriat() *Factoriat[int, int] {
	return MustNewFactoriat[int, int](Config[int, int]{
		Stateful: false,
		Capture: func(_ *State[int], _ int) error {
			return nil
		},
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	})
}

func TestPush_DeadlockOnSelfLoopCallback(t *testing.T) {
	f := newAlwaysTriggerFactoriat()

	var calls atomic.Int32
	f.SetEmit(func(_ int) error {
		if calls.Add(1) == 1 {
			f.Push(1)
		}
		return nil
	})

	done := make(chan struct{})
	go func() {
		f.Push(1)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("expected Push to finish without blocking on self-loop callback")
	}

	if got := calls.Load(); got != 2 {
		t.Fatalf("expected callback to be called twice, got %d", got)
	}
}

func TestPush_DeadlockOnCyclicCallbacks(t *testing.T) {
	a := newAlwaysTriggerFactoriat()
	b := newAlwaysTriggerFactoriat()

	var aCalls atomic.Int32
	var bCalls atomic.Int32

	a.SetEmit(func(_ int) error {
		if aCalls.Add(1) == 1 {
			b.Push(1)
		}
		return nil
	})

	b.SetEmit(func(_ int) error {
		if bCalls.Add(1) == 1 {
			a.Push(1)
		}
		return nil
	})

	done := make(chan struct{})
	go func() {
		a.Push(1)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("expected Push to finish without blocking on cyclic callbacks")
	}

	if got := aCalls.Load(); got != 2 {
		t.Fatalf("expected A callback to be called twice, got %d", got)
	}
	if got := bCalls.Load(); got != 1 {
		t.Fatalf("expected B callback to be called once, got %d", got)
	}
}
