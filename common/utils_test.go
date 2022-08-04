package common

import (
	"testing"
)

func Test_interfaceToFloat(t *testing.T) {
	var want float64 = 1
	got := interfaceToFloat(int(1))
	if want != got {
		t.Error("wanted 1 as a float, but didn't get it")
	}
}

func Test_templatizeProps(t *testing.T) {
	type test struct {
		name        string
		props       map[string]any
		assignor    string
		paramPrefix string
		want        []string
	}

	var tests []test = []test{
		{
			name:        "simple",
			props:       map[string]any{"Prop1": nil, "Prop2": nil, "Prop3": nil},
			assignor:    ":",
			paramPrefix: "",
			want:        []string{"Prop1:$Prop1", "Prop2:$Prop2", "Prop3:$Prop3"},
		},
		{
			name:        "with special character",
			props:       map[string]any{"Prop`1": nil, "Prop`2": nil, "Prop`3": nil},
			assignor:    ":",
			paramPrefix: "",
			want:        []string{"Prop`1:$Prop`1", "Prop`2:$Prop`2", "Prop`3:$Prop`3"},
		},
		{
			name:        "with complex seperators and assignors",
			props:       map[string]any{"Prop1": nil, "Prop2": nil, "Prop3": nil},
			assignor:    "=",
			paramPrefix: "",
			want:        []string{"Prop1=$Prop1", "Prop2=$Prop2", "Prop3=$Prop3"},
		},
		{
			name:        "simple with prefix",
			props:       map[string]any{"Prop1": nil, "Prop2": nil, "Prop3": nil},
			assignor:    ":",
			paramPrefix: "left",
			want:        []string{"Prop1:$leftProp1", "Prop2:$leftProp2", "Prop3:$leftProp3"},
		},
	}

	for _, tc := range tests {
		got := templatizeProps(tc.props, tc.assignor, tc.paramPrefix)
		for i := range tc.want {
			if tc.want[i] != got[i] {
				t.Errorf("%s: wanted %v but got %v", tc.name, tc.want, got)
			}
		}
	}
}
