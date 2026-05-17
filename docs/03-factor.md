# Factor

`Factor` — это входной объект, который описывает состояние мира в момент принятия решения.

Хороший `Factor` не говорит, что делать. Он только описывает то, что известно:

```go
type Factor struct {
	TableID  uuid.UUID
	PlayerID uuid.UUID

	HasSeat            bool
	Folded             bool
	HasPositiveBalance bool
	PlayerState        player.State

	Now       time.Time
	LastSeen  time.Time
	Window    time.Duration
	MaxMissed int
}
```

## Что Класть В Factor

Кладите в `Factor` исходные признаки:

- текущее состояние игрока;
- наличие seat;
- folded/all-in/sitout flags;
- текущее время `Now`;
- последнее событие `LastSeen`;
- настройки окна или порога;
- состояние table/street;
- внешне загруженные данные.

## Что Не Класть В Factor

Не кладите в `Factor` готовые решения:

```go
ShouldLeave bool
ShouldDelete bool
ShouldCreateTimer bool
```

Такие поля превращают факториат в формальность. Лучше передавать исходные данные, а решение вычислять предикатами.

## Время Как Данные

Не используйте `time.Now()` внутри предикатов. Передавайте время в `Factor`:

```go
type Factor struct {
	Now      time.Time
	LastSeen time.Time
}
```

Так решение становится детерминированным и легко тестируется.

## Factor Как Контракт

`Factor` — это контракт между orchestration layer и decision layer.

Service или Manager собирает данные:

```go
factor := Factor{
	PlayerID: player.UserID,
	HasSeat:  seat != nil,
	Folded:   player.IsFolded(),
	Now:      time.Now().UTC(),
	LastSeen: player.HeartbeatLastSeen,
}
```

Factoriat принимает решение:

```go
facts, err := DecideFacts(factor)
```

Важно: факториат не должен сам ходить в repository, брать lock или читать runtime state.

