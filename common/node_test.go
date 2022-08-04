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

type cypherTest struct {
	name         string
	node         Node
	constraints  []string
	paramPrefix  string
	wantedQuery  string
	wantedParams map[string]any
}

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
	tests := []cypherTest{
		{
			name:        "simple test",
			node:        testNode,
			constraints: []string{"ConstrainedProp1", "ConstrainedProp2"},
			wantedQuery: `MERGE (n:TestLabel {ConstrainedProp1:$ConstrainedProp1, ConstrainedProp2:$ConstrainedProp2})
ON CREATE SET n.UnconstrainedProp1=$UnconstrainedProp1, n.UnconstrainedProp2=$UnconstrainedProp2
`,
			wantedParams: map[string]any{
				"ConstrainedProp1":   "ConstrainedValue1",
				"ConstrainedProp2":   "ConstrainedValue2",
				"UnconstrainedProp1": "UnconstrainedValue1",
				"UnconstrainedProp2": "UnconstrainedValue2",
			},
		},
	}

	for _, tc := range tests {
		gotQuery, gotParams := testNode.ToCypherMerge(tc.constraints, tc.paramPrefix)
		if tc.wantedQuery != gotQuery {
			t.Errorf("%s: wanted query \n%s\nbut got \n%s", tc.name, tc.wantedQuery, gotQuery)
		}
		if !reflect.DeepEqual(tc.wantedParams, gotParams) {
			t.Errorf("%s: wanted query \n%v\nbut got %v", tc.name, tc.wantedParams, gotParams)
		}
	}
}

func TestToCypherMatch(t *testing.T) {
	tests := []cypherTest{
		{
			name:        "simple test",
			node:        testNode,
			constraints: []string{"ConstrainedProp1", "ConstrainedProp2"},
			wantedQuery: "MATCH (n:TestLabel {ConstrainedProp1:$ConstrainedProp1, ConstrainedProp2:$ConstrainedProp2})\n",
			wantedParams: map[string]any{
				"ConstrainedProp1": "ConstrainedValue1",
				"ConstrainedProp2": "ConstrainedValue2",
			},
		},
	}

	for _, tc := range tests {
		gotQuery, gotParams := testNode.ToCypherMatch(tc.constraints, tc.paramPrefix)
		if tc.wantedQuery != gotQuery {
			t.Errorf("%s: wanted query \n%s\nbut got \n%s", tc.name, tc.wantedQuery, gotQuery)
		}
		if !reflect.DeepEqual(tc.wantedParams, gotParams) {
			t.Errorf("%s: wanted query \n%v\nbut got %v", tc.name, tc.wantedParams, gotParams)
		}
	}
}

func TestToCypherCreate(t *testing.T) {
	tests := []cypherTest{
		{
			name:        "simple test",
			node:        testNode,
			wantedQuery: "CREATE (n:TestLabel {ConstrainedProp1:$ConstrainedProp1, ConstrainedProp2:$ConstrainedProp2, UnconstrainedProp1:$UnconstrainedProp1, UnconstrainedProp2:$UnconstrainedProp2})\n",
			wantedParams: map[string]any{
				"ConstrainedProp1":   "ConstrainedValue1",
				"ConstrainedProp2":   "ConstrainedValue2",
				"UnconstrainedProp1": "UnconstrainedValue1",
				"UnconstrainedProp2": "UnconstrainedValue2",
			},
		},
	}

	for _, tc := range tests {
		gotQuery, gotParams := testNode.ToCypherCreate(tc.paramPrefix)
		if tc.wantedQuery != gotQuery {
			t.Errorf("%s: wanted query \n%s\nbut got \n%s", tc.name, tc.wantedQuery, gotQuery)
		}
		if !reflect.DeepEqual(tc.wantedParams, gotParams) {
			t.Errorf("%s: wanted query \n%v\nbut got %v", tc.name, tc.wantedParams, gotParams)
		}
	}
}
