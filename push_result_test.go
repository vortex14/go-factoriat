package factoriat

import "testing"

func TestPushResult_Emitted(t *testing.T) {
	var emitted int

	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) (int, error) {
			return 10, nil
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
	if emitted != 10 {
		t.Fatalf("expected callback to receive fact 10, got %d", emitted)
	}
}

func TestPushResult_SkippedInactive(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return false, nil
		},
		Build: func(_ *State[int]) (int, error) {
			return 1, nil
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
		Build: func(_ *State[int]) (int, error) {
			return 1, nil
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
		Build: func(_ *State[int]) (int, error) {
			return 1, nil
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
		Build: func(_ *State[int]) (int, error) {
			return 10, nil
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
