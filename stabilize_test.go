package factoriat

import "testing"

func TestPush_StabilizeRunsAfterBuild(t *testing.T) {
	var built int
	var emitted int

	f := MustNewFactoriat[int, int](Config[int, int]{
		Capture: func(st *State[int], _ int) error {
			st.Append(1)
			return nil
		},
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(st *State[int]) ([]int, error) {
			built = st.Count()
			return []int{st.Count()}, nil
		},
		Stabilize: func(st *State[int]) error {
			st.ClearData()
			return nil
		},
		Emit: func(fact int) error {
			emitted = fact
			return nil
		},
	})

	f.Push(1)

	if built != 1 {
		t.Fatalf("expected Build to see state before Stabilize, got %d", built)
	}
	if emitted != 1 {
		t.Fatalf("expected emitted fact to come from Build result, got %d", emitted)
	}
	if got := f.state.Count(); got != 0 {
		t.Fatalf("expected Stabilize to update state, got counter %d", got)
	}
}
