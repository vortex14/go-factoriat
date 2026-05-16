package factoriat

import "fmt"

// Emit возвращает текущий callback выпуска факта.
func (f *Factoriat[F, R]) Emit() EmitCallback[R] {
	return f.emit
}

// SetEmit задаёт новый callback выпуска факта.
func (f *Factoriat[F, R]) SetEmit(cb EmitCallback[R]) {
	if cb == nil {
		panic(fmt.Errorf("factoriat: emit callback is required"))
	}

	f.emit = cb
}
