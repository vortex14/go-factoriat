# Базовая Модель: Factor -> Predicates -> Facts -> Apply

Базовая модель Factoriat выглядит так:

```text
Factor -> Predicates -> Rules -> Facts -> Apply
```

Для stateless-сценария поток ещё проще:

```text
Factor -> Build -> []Fact
```

Для stateful-сценария используется полный lifecycle:

```text
Factor -> Capture -> Evaluate -> Build -> Stabilize -> Emit
```

## Factor

`Factor` — это входной снимок состояния. Он содержит данные, нужные для решения:

```go
type Factor struct {
	PlayerID string
	HasSeat  bool
	Folded   bool
	State    player.State
	Now      time.Time
	LastSeen time.Time
}
```

В `Factor` не нужно класть готовые решения вроде `ShouldLeave` или `CanDelete`. Лучше передавать исходные данные, а решение пусть принимает факториат.

## Predicates

Предикаты дают имена условиям:

```go
func (f Factor) nonFoldedSitOut() bool {
	return f.HasSeat && !f.Folded && f.State == player.SitOut
}
```

После этого rules builder читает доменный смысл, а не boolean-проводку.

## Facts

`Fact` — это решение. Он говорит, что должно произойти, но сам ничего не делает:

```go
type SitOutPlayerShouldMarkLeavePending struct {
	PlayerID string
}
```

## Apply

Факты применяются снаружи:

```go
for _, fact := range facts {
	switch f := fact.(type) {
	case SitOutPlayerShouldMarkLeavePending:
		player.MarkLeavePending(f.PlayerID)
	}
}
```

Так факториат остаётся чистым, а side effects выполняются в нужном application или domain layer.

