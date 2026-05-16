package factoriat

type State[F any] struct {
	data    []F
	counter int
	meta    map[string]any
}

// NewState Инициализация состояния
func NewState[F any]() State[F] {
	return State[F]{
		data:    make([]F, 0),
		counter: 0,
		meta:    make(map[string]any),
	}
}

// Append сохраняет фактор в буфер состояния и увеличивает счётчик.
func (s *State[F]) Append(factor F) {
	s.data = append(s.data, factor)
	s.counter++
}

// Count возвращает количество факторов, принятых в текущее состояние.
func (s *State[F]) Count() int {
	return s.counter
}

// Last возвращает последний захваченный фактор.
func (s *State[F]) Last() (F, bool) {
	if len(s.data) == 0 {
		var zero F
		return zero, false
	}

	return s.data[len(s.data)-1], true
}

// DataSnapshot возвращает независимую копию буфера факторов.
func (s *State[F]) DataSnapshot() []F {
	data := make([]F, len(s.data))
	copy(data, s.data)
	return data
}

// ReplaceData заменяет буфер факторов и синхронизирует счётчик.
func (s *State[F]) ReplaceData(data []F) {
	s.data = make([]F, len(data))
	copy(s.data, data)
	s.counter = len(s.data)
}

// MetaSnapshot возвращает независимую копию meta-памяти.
func (s *State[F]) MetaSnapshot() map[string]any {
	meta := make(map[string]any, len(s.meta))
	for key, value := range s.meta {
		meta[key] = value
	}
	return meta
}

// ReplaceMeta заменяет meta-память копией переданной map.
func (s *State[F]) ReplaceMeta(meta map[string]any) {
	s.meta = make(map[string]any, len(meta))
	for key, value := range meta {
		s.meta[key] = value
	}
}

// TrimLast оставляет в буфере только последние n факторов.
func (s *State[F]) TrimLast(n int) {
	if n <= 0 {
		s.ClearData()
		return
	}

	if n >= len(s.data) {
		s.counter = len(s.data)
		return
	}

	s.data = s.data[len(s.data)-n:]
	s.counter = len(s.data)
}

// ClearData очищает буфер факторов, не затрагивая Meta.
func (s *State[F]) ClearData() {
	s.data = make([]F, 0)
	s.counter = 0
}

// Set сохраняет произвольное значение в памяти состояния.
func (s *State[F]) Set(key string, value any) {
	s.ensureMeta()
	s.meta[key] = value
}

// Get возвращает значение из памяти состояния.
func (s *State[F]) Get(key string) (any, bool) {
	if s.meta == nil {
		return nil, false
	}

	value, ok := s.meta[key]
	return value, ok
}

// Reset возвращает состояние к пустому устойчивому виду.
func (s *State[F]) Reset() {
	*s = NewState[F]()
}

// Snapshot возвращает независимый снимок текущего состояния.
func (s *State[F]) Snapshot() State[F] {
	snapshot := State[F]{
		data:    make([]F, len(s.data)),
		counter: s.counter,
		meta:    make(map[string]any, len(s.meta)),
	}

	copy(snapshot.data, s.data)
	for key, value := range s.meta {
		snapshot.meta[key] = value
	}

	return snapshot
}

func (s *State[F]) ensureMeta() {
	if s.meta == nil {
		s.meta = make(map[string]any)
	}
}
