# Predicates

Предикаты — это именованные условия. Они превращают сложную boolean-логику в доменный язык.

## 1. Атомарные Предикаты

Атомарный предикат отвечает на один простой вопрос:

```go
func (f Factor) hasSeat() bool {
	return f.HasSeat
}

func (f Factor) isFolded() bool {
	return f.Folded
}
```

## 2. Базовые Доменные Предикаты

Базовый доменный предикат объединяет несколько простых признаков:

```go
func (f Factor) isSeatedPlayer() bool {
	return f.IsPlayer && f.HasSeat
}

func (f Factor) nonFoldedSitOut() bool {
	return f.HasSeat && !f.Folded && f.PlayerState == player.SitOut
}
```

## 3. Временные Предикаты

Если решение зависит от времени, время должно быть частью `Factor`:

```go
func (f Factor) elapsed() time.Duration {
	elapsed := f.Now.Sub(f.LastSeen)
	if elapsed < 0 {
		return 0
	}
	return elapsed
}
```

## 4. Композиционные Предикаты

Композиционный предикат собирает несколько условий в доменное решение:

```go
func (f Factor) canWaitingPlayer() bool {
	return f.isSeatedPlayer() &&
		f.HasPositiveBalance &&
		!f.LastSeen.IsZero() &&
		!f.heartbeatDeadlinePassed()
}
```

## 5. Decision Predicates

Decision predicate напрямую соответствует будущему факту:

```go
func (f Factor) shouldRescheduleHeartbeat() bool {
	return f.inStartupGrace() || f.heartbeatBelowThreshold()
}
```

## 6. Предикаты Над Фактами

Иногда нужно проверить уже созданный набор фактов:

```go
func heartbeatFactsContainLost(facts []Fact) bool {
	for _, fact := range facts {
		switch fact.(type) {
		case SeatedPlayerHeartbeatLost, SpectatorHeartbeatLost:
			return true
		}
	}
	return false
}
```

Такие предикаты лучше держать во внешнем слое, если они управляют execution flow.

## Пирамида Предикатов

```text
Atomic predicates
    ↓
Domain predicates
    ↓
Time predicates
    ↓
Composed predicates
    ↓
Decision predicates
    ↓
Facts
```

Главное правило: предикаты должны быть pure. Они не пишут в БД, не логируют, не создают timers и не мутируют состояние.

