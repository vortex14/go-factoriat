# Тестирование

Главная польза Factoriat — решения можно тестировать без инфраструктуры.

## Тест Factor -> Facts

Самый важный тест:

```go
func TestFactoriat_SitOutStillAwaiting(t *testing.T) {
	factor := Factor{
		IsSitOut: true,
		Now:      now,
		LastSeen: now.Add(-20 * time.Second),
		Window:   10 * time.Second,
	}

	facts, err := DecideFacts(factor)
	require.NoError(t, err)

	require.IsType(t, SitOutHeartbeatStillAwaiting{}, facts[0])
}
```

Такой тест не требует БД, timer service, locks, network или domain mutation.

## Table-Driven Tests

Для decision table удобно использовать table-driven style:

```go
tests := []struct {
	name string
	factor Factor
	want []string
}{
	{
		name: "startup grace",
		factor: Factor{Now: now, UpTime: now.Add(-10 * time.Second)},
		want: []string{"HeartbeatCheckShouldReschedule"},
	},
}
```

## Проверка Имён Фактов

Если у фактов есть `Name()`, можно проверять результат списком:

```go
func factNames(facts []Fact) []string {
	names := make([]string, 0, len(facts))
	for _, fact := range facts {
		names = append(names, fact.Name())
	}
	return names
}
```

## Тестируйте Границы

Для временных правил проверяйте ровно до порога, на пороге, после порога, отрицательный elapsed и zero window.

