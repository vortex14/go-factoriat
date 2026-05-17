package factoriat

import "testing"

func validConfig() Config[int, int] {
	return Config[int, int]{
		Stateful: true,
		Capture:  func(_ *State[int], _ int) error { return nil },
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Emit: func(_ int) error { return nil },
	}
}

func TestConfig_ValidateRequiresOnlyBuildWhenStateless(t *testing.T) {
	cfg := Config[int, int]{
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected stateless config with only build to be valid, got %v", err)
	}
}

func TestConfig_ValidateRequiresCaptureWhenStateful(t *testing.T) {
	cfg := validConfig()
	cfg.Capture = nil

	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestConfig_ValidateRequiresEvaluateWhenStateful(t *testing.T) {
	cfg := validConfig()
	cfg.Evaluate = nil

	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestConfig_ValidateRequiresBuild(t *testing.T) {
	cfg := validConfig()
	cfg.Build = nil

	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestConfig_ValidateRequiresEmitWhenStateful(t *testing.T) {
	cfg := validConfig()
	cfg.Emit = nil

	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestConfig_ValidateAllowsErrorReturningLifecycle(t *testing.T) {
	cfg := Config[int, int]{
		Stateful: true,
		Capture: func(_ *State[int], _ int) error {
			return nil
		},
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(_ *State[int]) ([]int, error) {
			return []int{1}, nil
		},
		Stabilize: func(_ *State[int]) error {
			return nil
		},
		Emit: func(_ int) error {
			return nil
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected config to be valid, got %v", err)
	}
}

func TestConfig_ValidateRejectsUnknownTriggerMode(t *testing.T) {
	cfg := validConfig()
	cfg.TriggerMode = TriggerMode(100)

	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestNewFactoriatReturnsValidationError(t *testing.T) {
	if _, err := NewFactoriat[int, int](Config[int, int]{}); err == nil {
		t.Fatalf("expected validation error")
	}
}
