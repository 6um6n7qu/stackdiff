package patch_test

import (
	"testing"

	"github.com/stackdiff/internal/patch"
)

func TestGenerate_NoChanges(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1", "B": "2"}
	ops := patch.Generate(src, dst)
	if len(ops) != 0 {
		t.Fatalf("expected 0 ops, got %d", len(ops))
	}
}

func TestGenerate_SetNew(t *testing.T) {
	src := map[string]string{}
	dst := map[string]string{"X": "hello"}
	ops := patch.Generate(src, dst)
	if len(ops) != 1 || ops[0].Action != "set" || ops[0].Key != "X" || ops[0].Value != "hello" {
		t.Fatalf("unexpected ops: %v", ops)
	}
}

func TestGenerate_DeleteMissing(t *testing.T) {
	src := map[string]string{"OLD": "val"}
	dst := map[string]string{}
	ops := patch.Generate(src, dst)
	if len(ops) != 1 || ops[0].Action != "delete" || ops[0].Key != "OLD" {
		t.Fatalf("unexpected ops: %v", ops)
	}
}

func TestGenerate_UpdateChanged(t *testing.T) {
	src := map[string]string{"K": "old"}
	dst := map[string]string{"K": "new"}
	ops := patch.Generate(src, dst)
	if len(ops) != 1 || ops[0].Action != "set" || ops[0].Value != "new" {
		t.Fatalf("unexpected ops: %v", ops)
	}
}

func TestGenerate_Mixed(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{"A": "1", "B": "changed", "D": "4"}
	ops := patch.Generate(src, dst)
	// expect: delete C, set B, set D
	if len(ops) != 3 {
		t.Fatalf("expected 3 ops, got %d: %v", len(ops), ops)
	}
}

func TestApply_RoundTrip(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{"A": "1", "B": "changed", "D": "4"}
	ops := patch.Generate(src, dst)
	result := patch.Apply(src, ops)
	for k, v := range dst {
		if result[k] != v {
			t.Errorf("key %s: want %s got %s", k, v, result[k])
		}
	}
	if len(result) != len(dst) {
		t.Errorf("length mismatch: want %d got %d", len(dst), len(result))
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"X": "original"}
	ops := []patch.Op{{Action: "set", Key: "X", Value: "mutated"}}
	patch.Apply(base, ops)
	if base["X"] != "original" {
		t.Error("Apply mutated the base map")
	}
}

func TestOp_String(t *testing.T) {
	s := patch.Op{Action: "set", Key: "FOO", Value: "bar"}.String()
	if s != "set FOO=bar" {
		t.Errorf("unexpected: %s", s)
	}
	d := patch.Op{Action: "delete", Key: "FOO"}.String()
	if d != "delete FOO" {
		t.Errorf("unexpected: %s", d)
	}
}
