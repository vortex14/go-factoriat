# Facts

`Fact` — это результат решения факториата. Он описывает, что должно произойти, но сам этого не выполняет.

```go
type Fact interface {
	Name() string
}
```

Примеры фактов:

```go
type HeartbeatCheckShouldReschedule struct {
	PlayerID uuid.UUID
	Deadline time.Time
}

type PlayerShouldLeaveTable struct {
	PlayerID uuid.UUID
}
```

## Как Называть Facts

Имя факта должно звучать как доменное решение:

```go
PlayerShouldLeaveTable
SeatShouldBeReleased
SitOutHeartbeatStillAwaiting
```

Не называйте факты именами инфраструктурных действий:

```go
CallPlayerRepoDelete
SendGrpcTimer
UpdatePostgresRow
```

Факт не знает, как он будет применён.

## Fact Не Является Side Effect

Факт не должен писать в БД, создавать timer, вызывать network, мутировать table/player или брать lock. Он только переносит решение наружу.

## Почему Полезен `Name()`

`Name()` делает тесты и логи проще:

```go
func (PlayerShouldLeaveTable) Name() string {
	return "PlayerShouldLeaveTable"
}
```

## Один Или Несколько Facts

Факториaт может вернуть один факт:

```go
return []Fact{SpectatorHeartbeatLost{}}
```

Или несколько фактов:

```go
return []Fact{
	HeartbeatCheckShouldReschedule{},
	SeatedPlayerHeartbeatLost{},
}
```

Несколько фактов полезны, когда одно решение требует нескольких последующих действий на внешнем слое.

