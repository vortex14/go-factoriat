package factoriat

// MetaKey описывает типизированный ключ meta-памяти состояния.
type MetaKey[T any] struct {
	name string
}

// NewMetaKey создаёт типизированный ключ meta-памяти.
func NewMetaKey[T any](name string) MetaKey[T] {
	return MetaKey[T]{name: name}
}

// Name возвращает строковое имя ключа.
func (k MetaKey[T]) Name() string {
	return k.name
}

// SetMeta сохраняет значение по типизированному ключу meta-памяти.
func SetMeta[F any, T any](st *State[F], key MetaKey[T], value T) {
	st.Set(key.name, value)
}

// GetMeta возвращает значение по типизированному ключу meta-памяти.
func GetMeta[F any, T any](st *State[F], key MetaKey[T]) (T, bool) {
	value, ok := st.Get(key.name)
	if !ok {
		var zero T
		return zero, false
	}

	typed, ok := value.(T)
	return typed, ok
}

// MetaOr возвращает значение по типизированному ключу или fallback.
func MetaOr[F any, T any](st *State[F], key MetaKey[T], fallback T) T {
	value, ok := GetMeta(st, key)
	if !ok {
		return fallback
	}

	return value
}

// IncMetaInt увеличивает int-значение по типизированному ключу.
func IncMetaInt[F any](st *State[F], key MetaKey[int], delta int) int {
	value := MetaOr(st, key, 0)
	value += delta
	SetMeta(st, key, value)
	return value
}
