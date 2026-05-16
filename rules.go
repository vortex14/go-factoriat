package factoriat

// CaptureRule захватывает входной фактор в живую память факториата.
type CaptureRule[F any] func(st *State[F], factor F) error

// EvaluateRule это затвор активации: он решает, активен ли факториат
// при текущей живой памяти и входном факторе.
type EvaluateRule[F any] func(st *State[F], factor F) (bool, error)

// BuildRule строит факт из независимого снимка памяти.
type BuildRule[F any, R any] func(st *State[F]) (R, error)

// StabilizeRule переводит живую память в следующее устойчивое состояние
// после построения факта.
type StabilizeRule[F any] func(st *State[F]) error

// EmitCallback выпускает готовый факт во внешний мир.
type EmitCallback[R any] func(fact R) error
