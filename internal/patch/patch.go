// Package patch generates a minimal set of changes (patch) needed to
// reconcile one config map into another.
package patch

import (
	"fmt"
	"sort"
)

// Op represents a single patch operation.
type Op struct {
	Action string // "set" | "delete"
	Key    string
	Value  string
}

func (o Op) String() string {
	if o.Action == "delete" {
		return fmt.Sprintf("delete %s", o.Key)
	}
	return fmt.Sprintf("set %s=%s", o.Key, o.Value)
}

// Generate returns the list of Ops needed to transform src into dst.
func Generate(src, dst map[string]string) []Op {
	var ops []Op

	for k, dv := range dst {
		if sv, ok := src[k]; !ok || sv != dv {
			ops = append(ops, Op{Action: "set", Key: k, Value: dv})
		}
	}

	for k := range src {
		if _, ok := dst[k]; !ok {
			ops = append(ops, Op{Action: "delete", Key: k})
		}
	}

	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Action != ops[j].Action {
			return ops[i].Action < ops[j].Action
		}
		return ops[i].Key < ops[j].Key
	})

	return ops
}

// Apply executes a slice of Ops against a copy of base and returns the result.
func Apply(base map[string]string, ops []Op) map[string]string {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}
	for _, op := range ops {
		switch op.Action {
		case "set":
			out[op.Key] = op.Value
		case "delete":
			delete(out, op.Key)
		}
	}
	return out
}
