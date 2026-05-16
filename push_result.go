package factoriat

type PushStatus int

const (
	PushStatusSkippedInactive PushStatus = iota
	PushStatusSkippedAlreadyTriggered
	PushStatusEmitted
	PushStatusFailed
)

type PushResult struct {
	Status       PushStatus
	Emitted      bool
	EmittedCount int
	Err          error
}
