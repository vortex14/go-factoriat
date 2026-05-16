package factoriat_test

import (
	"fmt"

	"github.com/vortex14/go-factoriat"
)

func Example_threshold() {
	metaTotal := factoriat.NewMetaKey[int]("total")

	type LargePurchaseFact struct {
		Total int
	}

	var emitted []int

	f := factoriat.MustNewFactoriat[int, LargePurchaseFact](factoriat.Config[int, LargePurchaseFact]{
		Capture: func(st *factoriat.State[int], amount int) error {
			factoriat.IncMetaInt(st, metaTotal, amount)
			return nil
		},
		Evaluate: func(st *factoriat.State[int], _ int) (bool, error) {
			return factoriat.MetaOr(st, metaTotal, 0) >= 1000, nil
		},
		Build: func(st *factoriat.State[int]) ([]LargePurchaseFact, error) {
			return []LargePurchaseFact{{Total: factoriat.MetaOr(st, metaTotal, 0)}}, nil
		},
		Stabilize: func(st *factoriat.State[int]) error {
			st.Reset()
			return nil
		},
		Emit: func(fact LargePurchaseFact) error {
			emitted = append(emitted, fact.Total)
			return nil
		},
	})

	f.Push(300)
	f.Push(400)
	f.Push(400)

	fmt.Println(emitted)

	// Output:
	// [1100]
}

func Example_slidingWindow() {
	type LoginEvent struct {
		UserID string
		OK     bool
	}

	type BruteForceFact struct {
		UserID string
		Count  int
	}

	var emitted []BruteForceFact

	f := factoriat.MustNewFactoriat[LoginEvent, BruteForceFact](factoriat.Config[LoginEvent, BruteForceFact]{
		Capture: func(st *factoriat.State[LoginEvent], event LoginEvent) error {
			if event.OK {
				st.Reset()
				return nil
			}
			st.Append(event)
			st.TrimLast(5)
			return nil
		},
		Evaluate: func(st *factoriat.State[LoginEvent], _ LoginEvent) (bool, error) {
			return st.Count() == 5, nil
		},
		Build: func(st *factoriat.State[LoginEvent]) ([]BruteForceFact, error) {
			events := st.DataSnapshot()
			return []BruteForceFact{{
				UserID: events[0].UserID,
				Count:  st.Count(),
			}}, nil
		},
		Stabilize: func(st *factoriat.State[LoginEvent]) error {
			st.Reset()
			return nil
		},
		Emit: func(fact BruteForceFact) error {
			emitted = append(emitted, fact)
			return nil
		},
	})

	for range 5 {
		f.Push(LoginEvent{UserID: "u1"})
	}

	fmt.Printf("%s %d\n", emitted[0].UserID, emitted[0].Count)

	// Output:
	// u1 5
}

func Example_edgeTrigger() {
	metaScore := factoriat.NewMetaKey[int]("score")

	type RiskFact struct {
		Score int
	}

	var emitted []int

	f := factoriat.MustNewFactoriat[int, RiskFact](factoriat.Config[int, RiskFact]{
		TriggerMode: factoriat.TriggerModeEdge,
		Stateful:    true,
		Capture: func(st *factoriat.State[int], score int) error {
			factoriat.SetMeta(st, metaScore, score)
			return nil
		},
		Evaluate: func(st *factoriat.State[int], _ int) (bool, error) {
			return factoriat.MetaOr(st, metaScore, 0) >= 80, nil
		},
		Build: func(st *factoriat.State[int]) ([]RiskFact, error) {
			return []RiskFact{{Score: factoriat.MetaOr(st, metaScore, 0)}}, nil
		},
		Emit: func(fact RiskFact) error {
			emitted = append(emitted, fact.Score)
			return nil
		},
	})

	f.Push(70)
	f.Push(80)
	f.Push(90)
	f.Push(70)
	f.Push(81)

	fmt.Println(emitted)

	// Output:
	// [80 81]
}

func Example_pipeline() {
	type LoginEvent struct {
		UserID string
		OK     bool
	}

	type SuspiciousLoginFact struct {
		UserID string
	}

	type SecurityAlertFact struct {
		UserID string
		Level  string
	}

	var emitted []SecurityAlertFact

	alerts := factoriat.MustNewFactoriat[SuspiciousLoginFact, SecurityAlertFact](factoriat.Config[SuspiciousLoginFact, SecurityAlertFact]{
		Capture: func(st *factoriat.State[SuspiciousLoginFact], fact SuspiciousLoginFact) error {
			st.Append(fact)
			return nil
		},
		Evaluate: func(st *factoriat.State[SuspiciousLoginFact], _ SuspiciousLoginFact) (bool, error) {
			return st.Count() >= 1, nil
		},
		Build: func(st *factoriat.State[SuspiciousLoginFact]) ([]SecurityAlertFact, error) {
			fact, _ := st.Last()
			return []SecurityAlertFact{{UserID: fact.UserID, Level: "high"}}, nil
		},
		Emit: func(fact SecurityAlertFact) error {
			emitted = append(emitted, fact)
			return nil
		},
	})

	logins := factoriat.MustNewFactoriat[LoginEvent, SuspiciousLoginFact](factoriat.Config[LoginEvent, SuspiciousLoginFact]{
		Capture: func(st *factoriat.State[LoginEvent], event LoginEvent) error {
			if !event.OK {
				st.Append(event)
			}
			return nil
		},
		Evaluate: func(st *factoriat.State[LoginEvent], _ LoginEvent) (bool, error) {
			return st.Count() >= 3, nil
		},
		Build: func(st *factoriat.State[LoginEvent]) ([]SuspiciousLoginFact, error) {
			event, _ := st.Last()
			return []SuspiciousLoginFact{{UserID: event.UserID}}, nil
		},
		Stabilize: func(st *factoriat.State[LoginEvent]) error {
			st.Reset()
			return nil
		},
		Emit: func(fact SuspiciousLoginFact) error {
			alerts.Push(fact)
			return nil
		},
	})

	logins.Push(LoginEvent{UserID: "u1"})
	logins.Push(LoginEvent{UserID: "u1"})
	logins.Push(LoginEvent{UserID: "u1"})

	fmt.Println(emitted)

	// Output:
	// [{u1 high}]
}
