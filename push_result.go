package factoriat

type PushStatus int

const (
	PushStatusSkippedInactive PushStatus = iota
	PushStatusSkippedAlreadyTriggered
	PushStatusEmitted
	PushStatusFailed
)

type PushResult struct {
	Status  PushStatus
	Emitted bool
	Err     error
}
