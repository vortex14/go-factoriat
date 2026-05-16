package factoriat

import "fmt"

type Config[F any, R any] struct {
	Stateful bool

	TriggerMode TriggerMode

	// StateKey identifies this Factoriat state inside StateRepository.
	StateKey string
	// StateRepository stores runtime state. If nil, an in-memory repository is used.
	StateRepository StateRepository[F]

	Capture  CaptureRule[F]
	Evaluate EvaluateRule[F]
	Build    BuildRule[F, R]

	Stabilize StabilizeRule[F]

	Emit EmitCallback[R]
}

func (cfg Config[F, R]) Validate() error {
	if cfg.Capture == nil {
		return fmt.Errorf("factoriat config: capture is required")
	}

	if cfg.Evaluate == nil {
		return fmt.Errorf("factoriat config: evaluate is required")
	}

	if cfg.Build == nil {
		return fmt.Errorf("factoriat config: build is required")
	}

	if cfg.Emit == nil {
		return fmt.Errorf("factoriat config: emit callback is required")
	}

	if cfg.StateRepository != nil && cfg.StateKey == "" {
		return fmt.Errorf("factoriat config: state key is required when state repository is configured")
	}

	switch cfg.TriggerMode {
	case TriggerModeLevel, TriggerModeEdge:
		return nil
	default:
		return fmt.Errorf("factoriat config: unknown trigger mode %d", cfg.TriggerMode)
	}
}
