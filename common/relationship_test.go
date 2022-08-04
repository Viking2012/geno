package common

import (
	"reflect"
	"testing"
)

type test struct {
	name         string
	rel          Relationship
	wantedQuery  string
	wantedParams map[string]any
}

var (
	nodeA Node         = NewNode(1, []string{"TypeA"}, map[string]any{"Prop1": "Value1A", "Prop2": "Value2A", "Prop3": nil})
	nodeB Node         = NewNode(2, []string{"TypeB"}, map[string]any{"Prop1": "Value1B", "Prop2": "Value2B", "Prop3": nil})
	relA  Relationship = NewRelationship(1, nodeA, nodeB, "TypeA", map[string]any{"Prop1": "Value1A", "Prop2": "Value2A"})
)

func Test_ToCypherMerge(t *testing.T) {
	tests := []test{
		{
			name: "simple merge",
			rel:  relA,
			wantedQuery: `MATCH (left:TypeA {Prop1:$leftProp1, Prop2:$leftProp2})
MATCH (right:TypeB {Prop1:$rightProp1, Prop2:$rightProp2})
MERGE (left)-[r:TypeA {Prop1:$relProp1, Prop2:$relProp2}]-(right)`,
			wantedParams: map[string]any{
				"leftProp1":  "Value1A",
				"leftProp2":  "Value2A",
				"rightProp1": "Value1B",
				"rightProp2": "Value2B",
				"relProp1":   "Value1A",
				"relProp2":   "Value2A",
			},
		},
	}

	for _, tc := range tests {
		gotQuery, gotParams := tc.rel.ToCypherMerge()
		if tc.wantedQuery != gotQuery {
			t.Errorf("%s: wanted query \n%s\nbut got \n%s", tc.name, tc.wantedQuery, gotQuery)
		}
		if !reflect.DeepEqual(tc.wantedParams, gotParams) {
			t.Errorf("%s: wanted query \n%v\nbut got %v", tc.name, tc.wantedParams, gotParams)
		}

	}
}

func Test_ToCypherMatch(t *testing.T) {
	tests := []test{
		{
			name: "simple merge",
			rel:  relA,
			wantedQuery: `MATCH (left:TypeA {Prop1:$leftProp1, Prop2:$leftProp2})
MATCH (right:TypeB {Prop1:$rightProp1, Prop2:$rightProp2})
MATCH (left)-[r:TypeA {Prop1:$relProp1, Prop2:$relProp2}]-(right)`,
			wantedParams: map[string]any{
				"leftProp1":  "Value1A",
				"leftProp2":  "Value2A",
				"rightProp1": "Value1B",
				"rightProp2": "Value2B",
				"relProp1":   "Value1A",
				"relProp2":   "Value2A",
			},
		},
	}

	for _, tc := range tests {
		gotQuery, gotParams := tc.rel.ToCypherMatch()
		if tc.wantedQuery != gotQuery {
			t.Errorf("%s: wanted query \n%s\nbut got \n%s", tc.name, tc.wantedQuery, gotQuery)
		}
		if !reflect.DeepEqual(tc.wantedParams, gotParams) {
			t.Errorf("%s: wanted query \n%v\nbut got %v", tc.name, tc.wantedParams, gotParams)
		}

	}
}

func Test_ToCypherCreate(t *testing.T) {
	tests := []test{
		{
			name: "simple merge",
			rel:  relA,
			wantedQuery: `MATCH (left:TypeA {Prop1:$leftProp1, Prop2:$leftProp2})
MATCH (right:TypeB {Prop1:$rightProp1, Prop2:$rightProp2})
CREATE (left)-[r:TypeA {Prop1:$relProp1,Prop2:$relProp2}]-(right)`,
			wantedParams: map[string]any{
				"leftProp1":  "Value1A",
				"leftProp2":  "Value2A",
				"rightProp1": "Value1B",
				"rightProp2": "Value2B",
				"relProp1":   "Value1A",
				"relProp2":   "Value2A",
			},
		},
	}

	for _, tc := range tests {
		gotQuery, gotParams := tc.rel.ToCypherCreate()
		if tc.wantedQuery != gotQuery {
			t.Errorf("%s: wanted query \n%s\nbut got \n%s", tc.name, tc.wantedQuery, gotQuery)
		}
		if !reflect.DeepEqual(tc.wantedParams, gotParams) {
			t.Errorf("%s: wanted query \n%v\nbut got %v", tc.name, tc.wantedParams, gotParams)
		}

	}
}
