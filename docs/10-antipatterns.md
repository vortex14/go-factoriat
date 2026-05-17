# Антипаттерны

## Side Effects Внутри Factoriat

Плохо:

```go
func buildFacts(f Factor) []Fact {
	repo.Delete(f.PlayerID)
	return nil
}
```

Factoriat должен только возвращать факты.

## Готовые Решения В Factor

Плохо:

```go
type Factor struct {
	ShouldLeave bool
}
```

Лучше:

```go
type Factor struct {
	HasSeat bool
	Folded  bool
	State   player.State
}
```

## Слишком Большие Предикаты

Плохо:

```go
func (f Factor) shouldLeaveAndFinalizeAndReleaseSeat() bool
```

Лучше разбить:

```go
isSeatedPlayer()
nonFoldedSitOut()
canFinalizeLeave()
```

## Скрытый Execution Flow

Плохо, когда apply-метод скрыто решает, идти ли дальше в domain flow.

Лучше держать важный переход явно:

```go
if !factsContainCriticalDecision(facts) {
	handled, err := applyRecoveryFacts(facts)
	if handled {
		return
	}
}
```

## Инфраструктурные Имена Фактов

Плохо:

```go
CallPostgresDelete
SendKafkaMessage
CreateGrpcTimer
```

Лучше:

```go
PlayerShouldLeaveTable
HeartbeatCheckShouldReschedule
SeatShouldBeReleased
```

## time.Now() В Предикатах

Плохо:

```go
func (f Factor) expired() bool {
	return time.Since(f.LastSeen) > f.Timeout
}
```

Лучше:

```go
func (f Factor) expired() bool {
	return f.Now.Sub(f.LastSeen) > f.Timeout
}
```

