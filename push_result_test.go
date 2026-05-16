package factoriat

import (
	"errors"
	"testing"
)

func TestPushResult_Emitted(t *testing.T) {
	var emitted int

	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{10}, nil
		},
		Emit: func(fact int) error {
			emitted = fact
			return nil
		},
	})

	result := f.PushResult(1)

	if result.Status != PushStatusEmitted {
		t.Fatalf("expected emitted status, got %v", result.Status)
	}
	if !result.Emitted {
		t.Fatalf("expected emitted result")
	}
	if result.EmittedCount != 1 {
		t.Fatalf("expected emitted count 1, got %d", result.EmittedCount)
	}
	if emitted != 10 {
		t.Fatalf("expected callback to receive fact 10, got %d", emitted)
	}
}

func TestPushResult_EmitsMultipleFactsInOrder(t *testing.T) {
	var emitted []int

	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{10, 20, 30}, nil
		},
		Emit: func(fact int) error {
			emitted = append(emitted, fact)
			return nil
		},
	})

	result := f.PushResult(1)

	if result.Status != PushStatusEmitted {
		t.Fatalf("expected emitted status, got %v", result.Status)
	}
	if !result.Emitted {
		t.Fatalf("expected emitted result")
	}
	if result.EmittedCount != 3 {
		t.Fatalf("expected emitted count 3, got %d", result.EmittedCount)
	}
	if len(emitted) != 3 || emitted[0] != 10 || emitted[1] != 20 || emitted[2] != 30 {
		t.Fatalf("expected ordered facts [10 20 30], got %v", emitted)
	}
}

func TestPushResult_EmitErrorAfterPartialMultiEmit(t *testing.T) {
	want := errors.New("emit failed")
	var emitted []int

	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{10, 20, 30}, nil
		},
		Emit: func(fact int) error {
			if fact == 20 {
				return want
			}
			emitted = append(emitted, fact)
			return nil
		},
	})

	result := f.PushResult(1)

	if result.Status != PushStatusFailed {
		t.Fatalf("expected failed status, got %v", result.Status)
	}
	if !errors.Is(result.Err, want) {
		t.Fatalf("expected emit error, got %v", result.Err)
	}
	if !result.Emitted {
		t.Fatalf("expected partial emission")
	}
	if result.EmittedCount != 1 {
		t.Fatalf("expected emitted count 1, got %d", result.EmittedCount)
	}
	if len(emitted) != 1 || emitted[0] != 10 {
		t.Fatalf("expected one emitted fact [10], got %v", emitted)
	}
}

func TestPushResult_SkippedInactive(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return false, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	result := f.PushResult(1)

	if result.Status != PushStatusSkippedInactive {
		t.Fatalf("expected inactive status, got %v", result.Status)
	}
	if result.Emitted {
		t.Fatalf("expected no emission")
	}
}

func TestPushResult_SkippedAlreadyTriggered(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Stateful:    true,
		TriggerMode: TriggerModeEdge,
		Capture:     func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	first := f.PushResult(1)
	second := f.PushResult(2)

	if first.Status != PushStatusEmitted {
		t.Fatalf("expected first push to emit, got %v", first.Status)
	}
	if second.Status != PushStatusSkippedAlreadyTriggered {
		t.Fatalf("expected already triggered status, got %v", second.Status)
	}
	if second.Emitted {
		t.Fatalf("expected no emission")
	}
}

func TestSetEmitRejectsNilCallback(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic")
		}
	}()

	f.SetEmit(nil)
}

func TestSetEmitReplacesCallback(t *testing.T) {
	var emitted int

	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{10}, nil
		},
		Emit: func(_ int) error { return nil },
	})
	f.SetEmit(func(fact int) error {
		emitted = fact
		return nil
	})

	f.Push(1)

	if emitted != 10 {
		t.Fatalf("expected replaced callback to receive fact 10, got %d", emitted)
	}
}
