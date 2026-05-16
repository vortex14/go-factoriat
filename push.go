package factoriat

func (f *Factoriat[F, R]) Push(factor F) {
	_ = f.PushResult(factor)
}

func (f *Factoriat[F, R]) PushResult(factor F) PushResult {
	var fact R
	var emit EmitCallback[R]
	var result PushResult

	func() {
		f.mu.Lock()
		defer f.mu.Unlock()

		record, err := f.stateRepository.LoadState(f.stateKey)
		if err != nil {
			result.Status = PushStatusFailed
			result.Err = err
			return
		}

		state := record.State
		triggered := record.Triggered

		save := func(status PushStatus, emitted bool, err error) {
			if saveErr := f.stateRepository.SaveState(f.stateKey, StateRecord[F]{
				State:     state,
				Triggered: triggered,
			}); saveErr != nil && err == nil {
				status = PushStatusFailed
				emitted = false
				err = saveErr
			}

			f.state = state.Snapshot()
			f.triggered = triggered
			result.Status = status
			result.Emitted = emitted
			result.Err = err
		}

		// 1. захватываем входной фактор в состояние
		if err := f.capture(&state, factor); err != nil {
			save(PushStatusFailed, false, err)
			return
		}

		active, err := f.evaluate(&state, factor)
		if err != nil {
			save(PushStatusFailed, false, err)
			return
		}
		if !active {
			triggered = false
			save(PushStatusSkippedInactive, false, nil)
			return
		}

		if f.triggerMode == TriggerModeEdge && triggered {
			save(PushStatusSkippedAlreadyTriggered, false, nil)
			return
		}

		snapshot := state.Snapshot()
		fact, err = f.build(&snapshot)
		if err != nil {
			save(PushStatusFailed, false, err)
			return
		}
		emit = f.emit

		// 4. стабилизируем состояние после порождения факта
		if f.stabilize != nil {
			if err := f.stabilize(&state); err != nil {
				save(PushStatusFailed, false, err)
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

		save(PushStatusEmitted, true, nil)
	}()

	if !result.Emitted {
		return result
	}

	// 5. выпускаем факт наружу после завершения внутреннего цикла
	if err := emit(fact); err != nil {
		return PushResult{
			Status: PushStatusFailed,
			Err:    err,
		}
	}

	return result
}
