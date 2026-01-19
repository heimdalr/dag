package dag

import "testing"

func TestVertexStoreAddLookupDelete(t *testing.T) {
	store := newVertexStore(func(v interface{}) interface{} {
		return "h:" + v.(string)
	})

	if err := store.add("id1", "v1"); err != nil {
		t.Fatalf("add id1 failed: %v", err)
	}

	if !store.hasID("id1") {
		t.Fatalf("hasID(id1) = false, want true")
	}

	hash, ok := store.hashByID("id1")
	if !ok {
		t.Fatalf("hashByID(id1) = false, want true")
	}
	if hash != "h:v1" {
		t.Fatalf("hashByID(id1) = %v, want h:v1", hash)
	}
	if !store.hasHash(hash) {
		t.Fatalf("hasHash(hash) = false, want true")
	}

	id, ok := store.idByHash(hash)
	if !ok {
		t.Fatalf("idByHash(hash) = false, want true")
	}
	if id != "id1" {
		t.Fatalf("idByHash(hash) = %s, want id1", id)
	}

	value, ok := store.value("id1")
	if !ok {
		t.Fatalf("value(id1) = false, want true")
	}
	if value != "v1" {
		t.Fatalf("value(id1) = %v, want v1", value)
	}

	store.delete("id1", hash)
	if store.hasID("id1") {
		t.Fatalf("hasID(id1) = true, want false")
	}
	if store.hasHash(hash) {
		t.Fatalf("hasHash(hash) = true, want false")
	}
}

func TestVertexStoreDuplicateAndNil(t *testing.T) {
	store := newVertexStore(nil)

	if err := store.add("id1", "v1"); err != nil {
		t.Fatalf("add id1 failed: %v", err)
	}

	if err := store.add("id2", "v1"); err == nil {
		t.Fatalf("add duplicate value = nil, want VertexDuplicateError")
	} else if _, ok := err.(VertexDuplicateError); !ok {
		t.Fatalf("add duplicate value expected VertexDuplicateError, got %T", err)
	}

	if err := store.add("id1", "v2"); err == nil {
		t.Fatalf("add duplicate id = nil, want IDDuplicateError")
	} else if _, ok := err.(IDDuplicateError); !ok {
		t.Fatalf("add duplicate id expected IDDuplicateError, got %T", err)
	}

	if err := store.add("id3", nil); err == nil {
		t.Fatalf("add nil value = nil, want VertexNilError")
	} else if _, ok := err.(VertexNilError); !ok {
		t.Fatalf("add nil value expected VertexNilError, got %T", err)
	}
}

func TestVertexStoreValuesEachByHash(t *testing.T) {
	store := newVertexStore(nil)
	_ = store.add("id1", "v1")
	_ = store.add("id2", "v2")

	if count := store.count(); count != 2 {
		t.Fatalf("count() = %d, want 2", count)
	}

	values := store.values()
	if len(values) != 2 {
		t.Fatalf("values() = %d, want 2", len(values))
	}
	if values["id1"] != "v1" || values["id2"] != "v2" {
		t.Fatalf("values() = %v, want id1:v1 and id2:v2", values)
	}

	seen := map[string]bool{}
	store.eachByHash(func(_ interface{}, id string, value interface{}) {
		seen[id] = value == values[id]
	})
	if len(seen) != 2 || !seen["id1"] || !seen["id2"] {
		t.Fatalf("eachByHash did not iterate expected entries: %v", seen)
	}
}

func TestVertexStoreSetHashFunc(t *testing.T) {
	store := newVertexStore(func(interface{}) interface{} {
		return "constant"
	})

	if err := store.add("id1", "v1"); err != nil {
		t.Fatalf("add id1 failed: %v", err)
	}

	if err := store.add("id2", "v2"); err == nil {
		t.Fatalf("add hash-colliding value = nil, want VertexDuplicateError")
	} else if _, ok := err.(VertexDuplicateError); !ok {
		t.Fatalf("add hash-colliding value expected VertexDuplicateError, got %T", err)
	}

	store.setHashFunc(nil)

	if err := store.add("id2", "v2"); err != nil {
		t.Fatalf("add id2 after hash reset failed: %v", err)
	}

	hash, ok := store.hashByID("id1")
	if !ok {
		t.Fatalf("hashByID(id1) = false, want true")
	}
	if hash != "v1" {
		t.Fatalf("hashByID(id1) = %v, want v1", hash)
	}
}
