# Apply Facts

Factoriat принимает решение, но не исполняет его. Исполнение происходит во внешнем слое.

```text
facts := Decide(factor)
applyFacts(facts)
```

## Где Применять Facts

Факты применяются там, где доступны side effects:

- application service;
- domain manager;
- handler;
- worker;
- orchestrator.

Пример:

```go
for _, fact := range facts {
	switch f := fact.(type) {
	case HeartbeatCheckShouldReschedule:
		timer.Create(f.PlayerID, f.Deadline)

	case PlayerShouldLeaveTable:
		table.Leave(f.PlayerID)
	}
}
```

## Не Прячьте Важные Flow-Переходы

Если наличие факта управляет переходом в более тяжёлую ветку, держите это явно:

```go
if !factsContainLost(facts) {
	handled, err := applyRecoveryFacts(facts)
	if handled {
		return
	}
}

applyDomainFlow()
```

## Preview И Locked Phase

Если есть preview и locked-проверка, используйте один и тот же факториат:

```go
previewFacts := Decide(previewFactor)

lockedFactor := previewFactor
lockedFactor.Now = time.Now().UTC()
lockedFactor.LastSeen = persisted.LastSeen

lockedFacts := Decide(lockedFactor)
```

Меняются данные, но не меняется формула решения.

