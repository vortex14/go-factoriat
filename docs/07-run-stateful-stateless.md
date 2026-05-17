# Run / Stateful / Stateless

В библиотеке есть два основных способа работы:

```text
Run        -> stateless decision
PushResult -> stateful lifecycle
```

## Stateless: Run

`Run(factor)` подходит, когда один входной `Factor` уже содержит весь контекст.

```go
f, err := factoriat.NewFactoriat[Factor, Fact](factoriat.Config[Factor, Fact]{
	Stateful: false,
	Build: func(st *factoriat.State[Factor]) ([]Fact, error) {
		factor, ok := st.Last()
		if !ok {
			return nil, fmt.Errorf("missing factor")
		}
		return buildFacts(factor), nil
	},
})

facts, err := f.Run(factor)
```

`Run` создаёт временный `State`, вызывает только `Build` и возвращает `[]Fact`.

## Stateful: PushResult

`PushResult(factor)` нужен, когда факториат должен копить память между вызовами:

```go
result := f.PushResult(input)
```

Полный lifecycle:

```text
Capture -> Evaluate -> Build -> Stabilize -> Emit
```

## Когда Что Выбирать

Используйте `Run`, если решение мгновенное, нет накопления, нет `Emit`, а side effects применяются вручную после `Run`.

Используйте `PushResult`, если нужен накопительный state, trigger mode, `Emit` и runtime status.

