package common

import (
	"reflect"
	"testing"
)

var (
	testLabels []string       = []string{"TestLabel"}
	testProps  map[string]any = map[string]any{
		"ConstrainedProp1":   "ConstrainedValue1",
		"ConstrainedProp2":   "ConstrainedValue2",
		"UnconstrainedProp1": "UnconstrainedValue1",
		"UnconstrainedProp2": "UnconstrainedValue2",
	}
	testNode Node = Node{
		Id:         1,
		Labels:     testLabels,
		Properties: testProps,
	}
)

func TestNewNode(t *testing.T) {
	type test struct {
		name  string
		id    int64
		labs  []string
		props map[string]any
		want  Node
	}

	tests := []test{
		{
			name: "simple node",
			id:   1,
			labs: []string{"TestLabel"},
			props: map[string]any{
				"ConstrainedProp1":   "ConstrainedValue1",
				"ConstrainedProp2":   "ConstrainedValue2",
				"UnconstrainedProp1": "UnconstrainedValue1",
				"UnconstrainedProp2": "UnconstrainedValue2",
			},
			want: Node{Id: 1, Labels: testLabels, Properties: testProps},
		},
		{
			name: "multilabel node",
			id:   1,
			labs: []string{"TestLabel", "TestLabel2"},
			props: map[string]any{
				"ConstrainedProp1":   "ConstrainedValue1",
				"ConstrainedProp2":   "ConstrainedValue2",
				"UnconstrainedProp1": "UnconstrainedValue1",
				"UnconstrainedProp2": "UnconstrainedValue2",
			},
			want: Node{Id: 1, Labels: []string{"TestLabel", "TestLabel2"}, Properties: testProps},
		},
	}

	for _, tc := range tests {
		got := NewNode(tc.id, tc.labs, tc.props)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("%s: expected:%v but got:%v", tc.name, tc.want, got)
		}
	}
}

func TestToCypherMerge(t *testing.T) {
	want := "MERGE (n {ConstrainedProp1:'ConstrainedValue1', ConstrainedProp2:'ConstrainedValue2'}) SET n:TestLabel, n.UnconstrainedProp1='UnconstrainedValue1', n.UnconstrainedProp2='UnconstrainedValue2'"
	got := testNode.ToCypherMerge([]string{"ConstrainedProp1", "ConstrainedProp2"})
	if want != got {
		t.Errorf("\nWanted: %s\nbut got:%s\n", want, got)
	}
}

func TestToCypherMatch(t *testing.T) {
	want := "MATCH (n:TestLabel {ConstrainedProp1:'ConstrainedValue1', ConstrainedProp2:'ConstrainedValue2'})"
	got := testNode.ToCypherMatch([]string{"ConstrainedProp1", "ConstrainedProp2"})
	if want != got {
		t.Errorf("\nWanted: %s\nbut got:%s\n", want, got)
	}
}

func TestToCypherCreate(t *testing.T) {
	want := "CREATE (n:TestLabel {ConstrainedProp1:'ConstrainedValue1', ConstrainedProp2:'ConstrainedValue2', UnconstrainedProp1:'UnconstrainedValue1', UnconstrainedProp2:'UnconstrainedValue2'})"
	got := testNode.ToCypherCreate()
	if want != got {
		t.Errorf("\nWanted: %s\nbut got:%s\n", want, got)
	}
}
