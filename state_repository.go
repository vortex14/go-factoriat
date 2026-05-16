package factoriat

import (
	"fmt"
	"sync"
)

const defaultStateKey = "default"

// StateRecord описывает persistable runtime-состояние факториата.
type StateRecord[F any] struct {
	State     State[F]
	Triggered bool
}

// Snapshot возвращает независимую копию записи состояния.
func (r StateRecord[F]) Snapshot() StateRecord[F] {
	return StateRecord[F]{
		State:     r.State.Snapshot(),
		Triggered: r.Triggered,
	}
}

// StateRepository хранит runtime-состояние факториата.
//
// Репозиторий должен сохранять не только State, но и Triggered, потому что
// TriggerModeEdge зависит от памяти о предыдущем активном импульсе.
type StateRepository[F any] interface {
	LoadState(key string) (StateRecord[F], error)
	SaveState(key string, record StateRecord[F]) error
}

// InMemoryStateRepository хранит состояния факториатов в памяти процесса.
type InMemoryStateRepository[F any] struct {
	mu      sync.Mutex
	records map[string]StateRecord[F]
}

// NewInMemoryStateRepository создаёт in-memory репозиторий состояний.
func NewInMemoryStateRepository[F any]() *InMemoryStateRepository[F] {
	return &InMemoryStateRepository[F]{
		records: make(map[string]StateRecord[F]),
	}
}

// LoadState загружает состояние по ключу или возвращает пустое состояние.
func (r *InMemoryStateRepository[F]) LoadState(key string) (StateRecord[F], error) {
	if key == "" {
		return StateRecord[F]{}, fmt.Errorf("factoriat state repository: state key is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.records[key]
	if !ok {
		return StateRecord[F]{State: NewState[F]()}, nil
	}

	return record.Snapshot(), nil
}

// SaveState сохраняет состояние по ключу.
func (r *InMemoryStateRepository[F]) SaveState(key string, record StateRecord[F]) error {
	if key == "" {
		return fmt.Errorf("factoriat state repository: state key is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.records[key] = record.Snapshot()
	return nil
}
