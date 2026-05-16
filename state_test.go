package factoriat

import "testing"

func TestState_AppendAndCount(t *testing.T) {
	st := NewState[int]()

	st.Append(10)
	st.Append(20)

	if got := st.Count(); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}

	data := st.DataSnapshot()
	if len(data) != 2 {
		t.Fatalf("expected data length 2, got %d", len(data))
	}

	if data[0] != 10 || data[1] != 20 {
		t.Fatalf("expected appended factors to be preserved, got %v", data)
	}
}

func TestState_SetAndGet(t *testing.T) {
	st := State[int]{}

	st.Set("sum", 42)

	value, ok := st.Get("sum")
	if !ok {
		t.Fatalf("expected key to exist")
	}

	if value != 42 {
		t.Fatalf("expected value 42, got %v", value)
	}
}

func TestState_Last(t *testing.T) {
	st := NewState[int]()

	if value, ok := st.Last(); ok {
		t.Fatalf("expected empty state to have no last value, got %v", value)
	}

	st.Append(10)
	st.Append(20)

	value, ok := st.Last()
	if !ok {
		t.Fatalf("expected last value to exist")
	}
	if value != 20 {
		t.Fatalf("expected last value 20, got %d", value)
	}
}

func TestState_DataSnapshotCopiesData(t *testing.T) {
	st := NewState[int]()
	st.Append(10)

	data := st.DataSnapshot()
	data[0] = 20

	if got := st.DataSnapshot(); got[0] != 10 {
		t.Fatalf("expected original data to stay unchanged, got %v", got)
	}
}

func TestState_ReplaceDataCopiesAndSyncsCount(t *testing.T) {
	st := NewState[int]()
	data := []int{10, 20}

	st.ReplaceData(data)
	data[0] = 100

	if got := st.Count(); got != 2 {
		t.Fatalf("expected count 2 after replace, got %d", got)
	}
	if got := st.DataSnapshot(); got[0] != 10 || got[1] != 20 {
		t.Fatalf("expected replaced data to be copied, got %v", got)
	}
}

func TestState_MetaSnapshotCopiesMeta(t *testing.T) {
	st := NewState[int]()
	st.Set("sum", 10)

	meta := st.MetaSnapshot()
	meta["sum"] = 20

	if value, _ := st.Get("sum"); value != 10 {
		t.Fatalf("expected original meta to stay unchanged, got %v", value)
	}
}

func TestState_ReplaceMetaCopiesMeta(t *testing.T) {
	st := NewState[int]()
	meta := map[string]any{"sum": 10}

	st.ReplaceMeta(meta)
	meta["sum"] = 20

	if value, _ := st.Get("sum"); value != 10 {
		t.Fatalf("expected replaced meta to be copied, got %v", value)
	}
}

func TestState_TrimLast(t *testing.T) {
	st := NewState[int]()
	st.Append(10)
	st.Append(20)
	st.Append(30)

	st.TrimLast(2)

	if got := st.Count(); got != 2 {
		t.Fatalf("expected count 2 after trim, got %d", got)
	}
	data := st.DataSnapshot()
	if data[0] != 20 || data[1] != 30 {
		t.Fatalf("expected last two factors to remain, got %v", data)
	}
}

func TestState_TrimLastClearsDataForNonPositiveLimit(t *testing.T) {
	st := NewState[int]()
	st.Append(10)
	st.Set("sum", 10)

	st.TrimLast(0)

	if got := st.Count(); got != 0 {
		t.Fatalf("expected count 0 after trim, got %d", got)
	}
	if data := st.DataSnapshot(); len(data) != 0 {
		t.Fatalf("expected data to be cleared, got %v", data)
	}
	if value, ok := st.Get("sum"); !ok || value != 10 {
		t.Fatalf("expected meta to stay unchanged, got %v", value)
	}
}

func TestState_ClearDataKeepsMeta(t *testing.T) {
	st := NewState[int]()
	st.Append(10)
	st.Set("sum", 10)

	st.ClearData()

	if got := st.Count(); got != 0 {
		t.Fatalf("expected count 0 after clear, got %d", got)
	}
	if data := st.DataSnapshot(); len(data) != 0 {
		t.Fatalf("expected data to be cleared, got %v", data)
	}
	if value, ok := st.Get("sum"); !ok || value != 10 {
		t.Fatalf("expected meta to stay unchanged, got %v", value)
	}
}

func TestState_GetMissingKey(t *testing.T) {
	st := State[int]{}

	value, ok := st.Get("missing")
	if ok {
		t.Fatalf("expected missing key, got value %v", value)
	}
}

func TestState_Reset(t *testing.T) {
	st := NewState[int]()
	st.Append(10)
	st.Set("sum", 10)

	st.Reset()

	if got := st.Count(); got != 0 {
		t.Fatalf("expected count 0 after reset, got %d", got)
	}
	if data := st.DataSnapshot(); len(data) != 0 {
		t.Fatalf("expected empty data after reset, got %v", data)
	}
	if _, ok := st.Get("sum"); ok {
		t.Fatalf("expected meta to be cleared after reset")
	}
}

func TestState_SnapshotCopiesState(t *testing.T) {
	st := NewState[int]()
	st.Append(10)
	st.Set("sum", 10)

	snapshot := st.Snapshot()

	snapshot.ReplaceData([]int{20, 30})
	snapshot.Set("sum", 20)

	if data := st.DataSnapshot(); data[0] != 10 {
		t.Fatalf("expected original data to stay unchanged, got %v", data)
	}
	if got := st.Count(); got != 1 {
		t.Fatalf("expected original counter to stay unchanged, got %d", got)
	}
	if value, _ := st.Get("sum"); value != 10 {
		t.Fatalf("expected original meta to stay unchanged, got %v", value)
	}
}
