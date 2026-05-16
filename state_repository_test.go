package factoriat

import (
	"errors"
	"testing"
)

func TestInMemoryStateRepository_LoadMissingReturnsEmptyState(t *testing.T) {
	repo := NewInMemoryStateRepository[int]()

	record, err := repo.LoadState("counter")
	if err != nil {
		t.Fatalf("expected missing state to load as empty, got %v", err)
	}
	if got := record.State.Count(); got != 0 {
		t.Fatalf("expected empty state, got count %d", got)
	}
	if record.Triggered {
		t.Fatalf("expected missing record to be inactive")
	}
}

func TestInMemoryStateRepository_SavesSnapshot(t *testing.T) {
	repo := NewInMemoryStateRepository[int]()
	st := NewState[int]()
	st.Append(10)

	if err := repo.SaveState("counter", StateRecord[int]{State: st, Triggered: true}); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	st.Append(20)

	record, err := repo.LoadState("counter")
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}
	if got := record.State.Count(); got != 1 {
		t.Fatalf("expected saved snapshot count 1, got %d", got)
	}
	if !record.Triggered {
		t.Fatalf("expected triggered flag to be saved")
	}
}

func TestConfig_ValidateRequiresStateKeyWhenRepositoryConfigured(t *testing.T) {
	cfg := validConfig()
	cfg.StateRepository = NewInMemoryStateRepository[int]()

	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestPush_UsesConfiguredStateRepository(t *testing.T) {
	repo := NewInMemoryStateRepository[int]()
	key := "shared-counter"

	newCounter := func() *Factoriat[int, int] {
		return MustNewFactoriat[int, int](Config[int, int]{
			Stateful:        true,
			StateKey:        key,
			StateRepository: repo,
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
				return nil
			},
		})
	}

	first := newCounter()
	second := newCounter()

	if result := first.PushResult(1); result.Status != PushStatusSkippedInactive {
		t.Fatalf("expected first push to be inactive, got %v", result.Status)
	}
	if result := second.PushResult(2); result.Status != PushStatusEmitted {
		t.Fatalf("expected second push to use shared state and emit, got %v", result.Status)
	}
}

func TestPush_PersistsTriggeredFlagInStateRepository(t *testing.T) {
	repo := NewInMemoryStateRepository[int]()
	key := "edge"
	var emitted int

	newEdge := func() *Factoriat[int, int] {
		return MustNewFactoriat[int, int](Config[int, int]{
			Stateful:        true,
			TriggerMode:     TriggerModeEdge,
			StateKey:        key,
			StateRepository: repo,
			Capture: func(st *State[int], _ int) error {
				st.Append(1)
				return nil
			},
			Evaluate: func(st *State[int], _ int) (bool, error) {
				return st.Count() >= 1, nil
			},
			Build: func(st *State[int]) ([]int, error) {
				return []int{st.Count()}, nil
			},
			Emit: func(_ int) error {
				emitted++
				return nil
			},
		})
	}

	first := newEdge()
	second := newEdge()

	if result := first.PushResult(1); result.Status != PushStatusEmitted {
		t.Fatalf("expected first push to emit, got %v", result.Status)
	}
	if result := second.PushResult(2); result.Status != PushStatusSkippedAlreadyTriggered {
		t.Fatalf("expected second push to see persisted triggered flag, got %v", result.Status)
	}
	if emitted != 1 {
		t.Fatalf("expected one emission, got %d", emitted)
	}
}

type failingStateRepository[F any] struct {
	loadErr error
	saveErr error
}

func (r failingStateRepository[F]) LoadState(_ string) (StateRecord[F], error) {
	if r.loadErr != nil {
		return StateRecord[F]{}, r.loadErr
	}
	return StateRecord[F]{State: NewState[F]()}, nil
}

func (r failingStateRepository[F]) SaveState(_ string, _ StateRecord[F]) error {
	return r.saveErr
}

func TestPush_ReturnsFailedWhenStateRepositoryLoadFails(t *testing.T) {
	loadErr := errors.New("load failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		StateKey:        "broken",
		StateRepository: failingStateRepository[int]{loadErr: loadErr},
		Capture:         func(_ *State[int], _ int) error { return nil },
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
	if !errors.Is(result.Err, loadErr) {
		t.Fatalf("expected load error, got %v", result.Err)
	}
}

func TestPush_ReturnsFailedWhenStateRepositorySaveFails(t *testing.T) {
	saveErr := errors.New("save failed")
	f := MustNewFactoriat[int, int](Config[int, int]{
		StateKey:        "broken",
		StateRepository: failingStateRepository[int]{saveErr: saveErr},
		Capture:         func(_ *State[int], _ int) error { return nil },
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
	if !errors.Is(result.Err, saveErr) {
		t.Fatalf("expected save error, got %v", result.Err)
	}
}
