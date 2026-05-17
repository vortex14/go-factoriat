package factoriat

import (
	"errors"
	"testing"
)

func TestRun_ReturnsBuiltFacts(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Build: func(st *State[int]) ([]int, error) {
			last, _ := st.Last()
			return []int{last * 2, last * 3}, nil
		},
	})

	facts, err := f.Run(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(facts) != 2 || facts[0] != 20 || facts[1] != 30 {
		t.Fatalf("expected facts [20 30], got %v", facts)
	}
}

func TestRun_DoesNotCallCaptureOrEvaluate(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error {
			t.Fatalf("run must not capture factor")
			return nil
		},
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			t.Fatalf("run must not evaluate factor")
			return false, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	facts, err := f.Run(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(facts) != 1 || facts[0] != 1 {
		t.Fatalf("expected facts [1], got %v", facts)
	}
}

func TestRun_DoesNotReadOrSaveStateRepository(t *testing.T) {
	repoErr := errors.New("state repository used")
	f := MustNewFactoriat[int, int](Config[int, int]{
		StateKey:        "broken",
		StateRepository: failingStateRepository[int]{loadErr: repoErr, saveErr: repoErr},
		Capture:         func(_ *State[int], _ int) error { return nil },
		Evaluate:        func(_ *State[int], _ int) (bool, error) { return false, nil },
		Build: func(st *State[int]) ([]int, error) {
			last, _ := st.Last()
			return []int{last}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	facts, err := f.Run(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(facts) != 1 || facts[0] != 10 {
		t.Fatalf("expected facts [10], got %v", facts)
	}
}

func TestRun_DoesNotReuseOrPersistLiveState(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Stateful: true,
		Capture:  func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) { return false, nil },
		Build: func(st *State[int]) ([]int, error) {
			last, _ := st.Last()
			return []int{last}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	first, err := f.Run(10)
	if err != nil {
		t.Fatalf("unexpected first error: %v", err)
	}
	second, err := f.Run(20)
	if err != nil {
		t.Fatalf("unexpected second error: %v", err)
	}

	if len(first) != 1 || first[0] != 10 {
		t.Fatalf("expected first facts [10], got %v", first)
	}
	if len(second) != 1 || second[0] != 20 {
		t.Fatalf("expected second facts [20], got %v", second)
	}
	if f.state.Count() != 0 {
		t.Fatalf("expected live state to stay empty, got count %d", f.state.Count())
	}
}

func TestRun_ReturnsBuildError(t *testing.T) {
	want := errors.New("build failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		Build: func(_ *State[int]) ([]int, error) {
			return nil, want
		},
	})

	_, err := f.Run(10)
	if !errors.Is(err, want) {
		t.Fatalf("expected build error, got %v", err)
	}
}

func TestPushResult_ReturnsFailedWhenStatelessConfigHasNoLifecycleRules(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
	})

	result := f.PushResult(10)

	if result.Status != PushStatusFailed {
		t.Fatalf("expected failed status, got %v", result.Status)
	}
	if !errors.Is(result.Err, errCaptureRequired) {
		t.Fatalf("expected capture required error, got %v", result.Err)
	}
}
