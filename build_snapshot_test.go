package factoriat

import "testing"

func TestPush_BuildReceivesStateSnapshot(t *testing.T) {
	f := MustNewFactoriat[int, int](Config[int, int]{
		Stateful: true,
		Capture: func(st *State[int], factor int) error {
			st.Append(factor)
			st.Set("sum", factor)
			return nil
		},
		Evaluate: func(_ *State[int], _ int) (bool, error) {
			return true, nil
		},
		Build: func(st *State[int]) ([]int, error) {
			st.ReplaceData([]int{100})
			st.Set("sum", 100)
			return []int{st.Count()}, nil
		},
		Emit: func(_ int) error { return nil },
	})

	f.Push(10)

	if data := f.state.DataSnapshot(); data[0] != 10 {
		t.Fatalf("expected Build data mutation to stay in snapshot, got %v", data)
	}
	if got := f.state.Count(); got != 1 {
		t.Fatalf("expected Build counter mutation to stay in snapshot, got %d", got)
	}
	if value, _ := f.state.Get("sum"); value != 10 {
		t.Fatalf("expected Build meta mutation to stay in snapshot, got %v", value)
	}
}
