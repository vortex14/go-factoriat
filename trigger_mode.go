package factoriat

type TriggerMode int

const (
	// TriggerModeLevel выпускает факт каждый раз, когда Evaluate считает
	// факториат активным.
	TriggerModeLevel TriggerMode = iota

	// TriggerModeEdge выпускает факт только при новом переходе Evaluate
	// из неактивного состояния в активное.
	TriggerModeEdge
)
