package pkg

import (
	"reflect"
	"testing"

	"github.com/Viking2012/geno/geno"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestConnection(t *testing.T) {
	useConsoleLogger := func(level neo4j.LogLevel) func(config *neo4j.Config) {
		return func(config *neo4j.Config) {
			config.Log = neo4j.ConsoleLogger(level)
		}
	}

	driver, err := neo4j.NewDriver("neo4j://localhost:7687", neo4j.BasicAuth("geno", "genopw", ""), useConsoleLogger(neo4j.INFO))
	if err != nil {
		t.Error(err)
	}
	defer driver.Close()

	err = driver.VerifyConnectivity()
	if err != nil {
		t.Error(err)
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, BoltLogger: neo4j.ConsoleBoltLogger()})
	defer session.Close()

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (n) RETURN n", nil)

		if err != nil {
			return nil, err
		}

		// Next returns false upon error
		for result.Next() {
			_ = result.Record()
		}
		// Err returns the error that caused Next to return false
		if err = result.Err(); err != nil {
			return nil, err
		}

		return nil, err

	})

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetGraphFromJson(t *testing.T) {
	type test struct {
		name       string
		jsonString string
		wantNodes  []geno.Node
		wantRels   []geno.Relationship
	}

	var (
		nodeA geno.Node = geno.Node{Id: 1, Labels: []string{"TypeA"}, Properties: map[string]any{"Prop1": "Value1A", "Prop2": "Value2A"}}
		nodeB geno.Node = geno.Node{Id: 2, Labels: []string{"TypeB"}, Properties: map[string]any{"Prop1": "Value1B", "Prop2": "Value2B"}}
		// nodeC geno.Node         = geno.Node{Id: 2, Labels: []string{"TypeC"}, Properties: map[string]any{"Prop1": "Value1C", "Prop2": "Value2C"}}
		relA geno.Relationship = geno.NewRelationship(4, nodeA, nodeB, "RelTypeA", map[string]any{"RelProp1": "RelValue1A", "RelProp2": "relValue2A"})
		// relB  geno.Relationship = geno.NewRelationship(5, nodeB, nodeC, "RelTypeB", map[string]any{"RelProp1": "RelValue1B", "RelProp2": "relValue2B"})
		// relC  geno.Relationship = geno.NewRelationship(6, nodeA, nodeC, "RelTypeC", map[string]any{"RelProp1": "RelValue1C", "RelProp2": "relValue2C"})
	)

	tests := []test{
		{
			name:       "simple",
			jsonString: `{"nodes":[{"identity":1,"labels": ["TypeA"],"properties":{"Prop1":"Value1A","Prop2":"Value2A"}},{"identity":2,"labels: ["TypeB"],"properties":{"Prop1":"Value1B","Prop2":"Value2B"}}],"rels":[{"identity":0,"start":1,"end":2,"type":"RelA""properties":{"RelProp1":"RelValue1A","RelProp2":"relValue2A"}}]}`,
			wantNodes:  []geno.Node{nodeA, nodeB},
			wantRels:   []geno.Relationship{relA},
		},
	}

	for _, tc := range tests {
		got, err := GetGraphFromJson([]byte(tc.jsonString))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(tc.wantNodes, got.Nodes) {
			t.Errorf("%s: wanted nodes of\n%v\nbut got\n%v\n", tc.name, tc.wantNodes, got.Nodes)
		}
		if !reflect.DeepEqual(tc.wantRels, got.Relationships) {
			t.Errorf("%s: wanted rels  of\n%v\nbut got\n%v\n", tc.name, tc.wantRels, got.Relationships)
		}
	}
}
