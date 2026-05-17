# Build Rules

`Build` — это место, где `Factor` превращается в `[]Fact`.

Для stateless-сценариев `Build` часто является главным rules builder:

```go
func buildFacts(f Factor) []Fact {
	if f.inStartupGrace() {
		return []Fact{HeartbeatCheckShouldReschedule{}}
	}

	if f.sitOutStillAwaiting() {
		return []Fact{SitOutHeartbeatStillAwaiting{}}
	}

	return []Fact{SeatedPlayerHeartbeatLost{}}
}
```

## Пишите Как Decision Table

Хороший `buildFacts` читается сверху вниз как таблица решений:

```go
if conditionA {
	return factA
}

if conditionB {
	return factB
}

return defaultFact
```

В нём не должно быть длинных boolean-выражений. Сложность должна жить в предикатах.

## Ранний Return Или Накопление

Используйте ранний `return`, когда правила взаимоисключающие. Используйте накопление, когда решений может быть несколько:

```go
facts := make([]Fact, 0, 4)

if f.nonFoldedSitOut() {
	facts = append(facts, SitOutPlayerShouldMarkLeavePending{})
}

facts = append(facts, HeartbeatMonitorShouldStop{})

return facts
```

## Build Должен Быть Чистым

`Build` не должен писать в live-state, вызывать repository, создавать timer, логировать business flow или менять domain objects.

## Ошибки В Build

Возвращайте ошибку, если невозможно построить корректное решение:

```go
Build: func(st *factoriat.State[Factor]) ([]Fact, error) {
	factor, ok := st.Last()
	if !ok {
		return nil, fmt.Errorf("missing factor")
	}

	return buildFacts(factor), nil
}
```

