package factoriat

import (
	"errors"
	"testing"
)

func TestPushResult_CaptureError(t *testing.T) {
	want := errors.New("capture failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error {
			return want
		},
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	result := f.PushResult(1)

	if result.Status != PushStatusFailed {
		t.Fatalf("expected failed status, got %v", result.Status)
	}
	if !errors.Is(result.Err, want) {
		t.Fatalf("expected capture error, got %v", result.Err)
	}
	if result.Emitted {
		t.Fatalf("expected no emission")
	}
}

func TestPushResult_EvaluateError(t *testing.T) {
	want := errors.New("evaluate failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return false, want
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	result := f.PushResult(1)

	if result.Status != PushStatusFailed {
		t.Fatalf("expected failed status, got %v", result.Status)
	}
	if !errors.Is(result.Err, want) {
		t.Fatalf("expected evaluate error, got %v", result.Err)
	}
}

func TestPushResult_BuildError(t *testing.T) {
	want := errors.New("build failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return nil, want
		},
		Emit: func(_ int) error { return nil },
	})

	result := f.PushResult(1)

	if result.Status != PushStatusFailed {
		t.Fatalf("expected failed status, got %v", result.Status)
	}
	if !errors.Is(result.Err, want) {
		t.Fatalf("expected build error, got %v", result.Err)
	}
}

func TestPushResult_StabilizeError(t *testing.T) {
	want := errors.New("stabilize failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Stabilize: func(_ *State[int]) error {
			return want
		},
		Emit: func(_ int) error { return nil },
	})

	result := f.PushResult(1)

	if result.Status != PushStatusFailed {
		t.Fatalf("expected failed status, got %v", result.Status)
	}
	if !errors.Is(result.Err, want) {
		t.Fatalf("expected stabilize error, got %v", result.Err)
	}
}

func TestPushResult_EmitError(t *testing.T) {
	want := errors.New("emit failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error {
			return want
		},
	})

	result := f.PushResult(1)

	if result.Status != PushStatusFailed {
		t.Fatalf("expected failed status, got %v", result.Status)
	}
	if !errors.Is(result.Err, want) {
		t.Fatalf("expected emit error, got %v", result.Err)
	}
	if result.Emitted {
		t.Fatalf("expected no successful emission")
	}
}
