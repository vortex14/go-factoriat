package factoriat

import "errors"

var (
	errCaptureRequired  = errors.New("factoriat push: capture is required")
	errEvaluateRequired = errors.New("factoriat push: evaluate is required")
	errEmitRequired     = errors.New("factoriat push: emit callback is required")
)

func (f *Factoriat[F, R]) Push(factor F) {
	_ = f.PushResult(factor)
}

func (f *Factoriat[F, R]) Run(factor F) ([]R, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	state := State[F]{
		data:    []F{factor},
		counter: 1,
	}

	return f.build(&state)
}

func (f *Factoriat[F, R]) PushResult(factor F) PushResult {
	var facts []R
	var emit EmitCallback[R]
	var result PushResult

	func() {
		f.mu.Lock()
		defer f.mu.Unlock()

		if f.capture == nil {
			result.Status = PushStatusFailed
			result.Err = errCaptureRequired
			return
		}

		if f.evaluate == nil {
			result.Status = PushStatusFailed
			result.Err = errEvaluateRequired
			return
		}

		if f.emit == nil {
			result.Status = PushStatusFailed
			result.Err = errEmitRequired
			return
		}

		record, err := f.stateRepository.LoadState(f.stateKey)
		if err != nil {
			result.Status = PushStatusFailed
			result.Err = err
			return
		}

		state := record.State
		triggered := record.Triggered

		save := func(status PushStatus, emittedCount int, err error) {
			if saveErr := f.stateRepository.SaveState(f.stateKey, StateRecord[F]{
				State:     state,
				Triggered: triggered,
			}); saveErr != nil && err == nil {
				status = PushStatusFailed
				emittedCount = 0
				err = saveErr
			}

			f.state = state.Snapshot()
			f.triggered = triggered
			result.Status = status
			result.Emitted = emittedCount > 0
			result.EmittedCount = emittedCount
			result.Err = err
		}

		// 1. захватываем входной фактор в состояние
		if err := f.capture(&state, factor); err != nil {
			save(PushStatusFailed, 0, err)
			return
		}

		active, err := f.evaluate(&state, factor)
		if err != nil {
			save(PushStatusFailed, 0, err)
			return
		}
		if !active {
			triggered = false
			save(PushStatusSkippedInactive, 0, nil)
			return
		}

		if f.triggerMode == TriggerModeEdge && triggered {
			save(PushStatusSkippedAlreadyTriggered, 0, nil)
			return
		}

		snapshot := state.Snapshot()
		facts, err = f.build(&snapshot)
		if err != nil {
			save(PushStatusFailed, 0, err)
			return
		}
		emit = f.emit

		// 4. стабилизируем состояние после порождения факта
		if f.stabilize != nil {
			if err := f.stabilize(&state); err != nil {
				save(PushStatusFailed, 0, err)
				return
			}
			triggered = false
		} else if !f.stateful {
			// по умолчанию, если не stateful — сбрасываем состояние
			state = NewState[F]()
			triggered = false
		} else {
			triggered = true
		}

		save(PushStatusEmitted, len(facts), nil)
	}()

	if !result.Emitted {
		return result
	}

	// 5. выпускаем факт наружу после завершения внутреннего цикла
	for i, fact := range facts {
		if err := emit(fact); err != nil {
			return PushResult{
				Status:       PushStatusFailed,
				Emitted:      i > 0,
				EmittedCount: i,
				Err:          err,
			}
		}
	}

	return result
}
