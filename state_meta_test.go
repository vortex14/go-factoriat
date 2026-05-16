package factoriat

import "testing"

func TestStateMeta_MetaKey(t *testing.T) {
	key := NewMetaKey[int]("sum")

	if key.Name() != "sum" {
		t.Fatalf("expected key name sum, got %q", key.Name())
	}
}

func TestStateMeta_SetMetaAndGetMeta(t *testing.T) {
	st := NewState[int]()
	key := NewMetaKey[int]("sum")

	SetMeta(&st, key, 42)

	value, ok := GetMeta(&st, key)
	if !ok {
		t.Fatalf("expected typed meta value to exist")
	}
	if value != 42 {
		t.Fatalf("expected value 42, got %d", value)
	}
}

func TestStateMeta_GetMetaReturnsFalseForWrongStoredType(t *testing.T) {
	st := NewState[int]()
	key := NewMetaKey[int]("sum")

	st.Set(key.Name(), "42")

	value, ok := GetMeta(&st, key)
	if ok {
		t.Fatalf("expected wrong stored type to return false, got %d", value)
	}
}

func TestStateMeta_MetaOr(t *testing.T) {
	st := NewState[int]()
	key := NewMetaKey[int]("sum")

	if got := MetaOr(&st, key, 10); got != 10 {
		t.Fatalf("expected fallback 10, got %d", got)
	}

	SetMeta(&st, key, 42)

	if got := MetaOr(&st, key, 10); got != 42 {
		t.Fatalf("expected stored value 42, got %d", got)
	}
}

func TestStateMeta_IncMetaInt(t *testing.T) {
	st := NewState[int]()
	key := NewMetaKey[int]("sum")

	if got := IncMetaInt(&st, key, 3); got != 3 {
		t.Fatalf("expected first increment to return 3, got %d", got)
	}
	if got := IncMetaInt(&st, key, 4); got != 7 {
		t.Fatalf("expected second increment to return 7, got %d", got)
	}
	if got := MetaOr(&st, key, 0); got != 7 {
		t.Fatalf("expected stored sum 7, got %d", got)
	}
}
