package factoriat

import "sync"

type Factoriat[F any, R any] struct {
	mu sync.Mutex

	state State[F]

	stateful    bool
	triggerMode TriggerMode
	triggered   bool

	stateKey        string
	stateRepository StateRepository[F]

	capture   CaptureRule[F]
	evaluate  EvaluateRule[F]
	build     BuildRule[F, R]
	stabilize StabilizeRule[F]

	emit EmitCallback[R]
}

func NewFactoriat[F any, R any](cfg Config[F, R]) (*Factoriat[F, R], error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return newFactoriatUnchecked(cfg), nil
}

func MustNewFactoriat[F any, R any](cfg Config[F, R]) *Factoriat[F, R] {
	f, err := NewFactoriat(cfg)
	if err != nil {
		panic(err)
	}

	return f
}

func newFactoriatUnchecked[F any, R any](cfg Config[F, R]) *Factoriat[F, R] {
	stateKey := cfg.StateKey
	if stateKey == "" {
		stateKey = defaultStateKey
	}

	stateRepository := cfg.StateRepository
	if stateRepository == nil {
		stateRepository = NewInMemoryStateRepository[F]()
	}

	return &Factoriat[F, R]{
		state:           NewState[F](),
		stateful:        cfg.Stateful,
		triggerMode:     cfg.TriggerMode,
		stateKey:        stateKey,
		stateRepository: stateRepository,
		capture:         cfg.Capture,
		evaluate:        cfg.Evaluate,
		build:           cfg.Build,
		stabilize:       cfg.Stabilize,
		emit:            cfg.Emit,
	}
}
