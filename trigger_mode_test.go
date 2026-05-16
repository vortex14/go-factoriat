package factoriat

import (
	"sync/atomic"
	"testing"
)

func TestPush_LevelModeEmitsWhileTriggerIsTrue(t *testing.T) {
	var calls atomic.Int32

	f := MustNewFactoriat[int, int](Config[int, int]{
		Stateful: true,
		Capture: func(st *State[int], _ int) error {
			st.Append(1)
			return nil
		},
		Evaluate: func(st *State[int], _ int) (bool, error) {
			return st.Count() >= 2, nil
		},
		Build: func(st *State[int]) ([]int, error) {
			return []int{st.Count()}, nil
		},
		Emit: func(_ int) error {
			calls.Add(1)
			return nil
		},
	})

	f.Push(1)
	f.Push(2)
	f.Push(3)

	if got := calls.Load(); got != 2 {
		t.Fatalf("expected level mode to emit while trigger stays true, got %d calls", got)
	}
}

func TestPush_EdgeModeEmitsOnlyOnFalseToTrueTransition(t *testing.T) {
	var calls atomic.Int32

	f := MustNewFactoriat[int, int](Config[int, int]{
		Stateful:    true,
		TriggerMode: TriggerModeEdge,
		Capture: func(st *State[int], _ int) error {
			st.Append(1)
			return nil
		},
		Evaluate: func(st *State[int], _ int) (bool, error) {
			return st.Count() >= 2, nil
		},
		Build: func(st *State[int]) ([]int, error) {
			return []int{st.Count()}, nil
		},
		Emit: func(_ int) error {
			calls.Add(1)
			return nil
		},
	})

	f.Push(1)
	f.Push(2)
	f.Push(3)

	if got := calls.Load(); got != 1 {
		t.Fatalf("expected edge mode to emit only once while trigger stays true, got %d calls", got)
	}
}

func TestPush_EdgeModeCanEmitAgainAfterDefaultReset(t *testing.T) {
	var calls atomic.Int32

	f := MustNewFactoriat[int, int](Config[int, int]{
		TriggerMode: TriggerModeEdge,
		Capture: func(st *State[int], _ int) error {
			st.Append(1)
			return nil
		},
		Evaluate: func(st *State[int], _ int) (bool, error) {
			return st.Count() >= 2, nil
		},
		Build: func(st *State[int]) ([]int, error) {
			return []int{st.Count()}, nil
		},
		Emit: func(_ int) error {
			calls.Add(1)
			return nil
		},
	})

	f.Push(1)
	f.Push(2)
	f.Push(3)
	f.Push(4)

	if got := calls.Load(); got != 2 {
		t.Fatalf("expected edge mode to emit again after reset, got %d calls", got)
	}
}
