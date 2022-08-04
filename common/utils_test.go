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

func Test_JoinPropsToTemplate(t *testing.T) {
	type test struct {
		name     string
		props    map[string]any
		assignor string
		want     []string
	}

	var tests []test = []test{
		{
			name:     "simple",
			props:    map[string]any{"Prop1": nil, "Prop2": nil, "Prop3": nil},
			assignor: ":",
			want:     []string{"Prop1:$Prop1", "Prop2:$Prop2", "Prop3:$Prop3"},
		},
		{
			name:     "with special character",
			props:    map[string]any{"Prop`1": nil, "Prop`2": nil, "Prop`3": nil},
			assignor: ":",
			want:     []string{"Prop`1:$Prop`1", "Prop`2:$Prop`2", "Prop`3:$Prop`3"},
		},
		{
			name:     "with complex seperators and assignors",
			props:    map[string]any{"Prop1": nil, "Prop2": nil, "Prop3": nil},
			assignor: "=",
			want:     []string{"Prop1=$Prop1", "Prop2=$Prop2", "Prop3=$Prop3"},
		},
	}

	for _, tc := range tests {
		got := templatizeProps(tc.props, tc.assignor)
		for i := range tc.want {
			if tc.want[i] != got[i] {
				t.Errorf("%s: wanted %v but got %v", tc.name, tc.want, got)
			}
		}
	}
}
