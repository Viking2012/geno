package common

import "testing"

func TestToCypherMerge(t *testing.T) {
	testNode := NewNode(1, []string{"TestLabel"}, map[string]any{"ConstrainedProp": "ConstrainedValue", "UnconstrainedProp": "UnconstrainedValue"})
	want := "MERGE (n {ConstrainedProp:`ConstrainedValue`} SET n:TestLabel AND n.UnconstrainedProp=`UnconstrainedValue`"
	got := testNode.ToCypherMerge([]string{"ConstrainedProp"})
	if want != got {
		t.Errorf("\nWanted: %s\nbut got:%s\n", want, got)
	}
}
