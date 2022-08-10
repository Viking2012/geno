package pkg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Viking2012/geno/geno"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type jsonNode struct {
	Id     int64          `json:"identity"`
	Labels []string       `json:"labels"`
	Props  map[string]any `json:"properties"`
}

type jsonRelationship struct {
	Id         int64          `json:"identity"`
	Start      int64          `json:"start"`
	End        int64          `json:"end"`
	Label      string         `json:"type"`
	Properties map[string]any `json:"properties"`
}
type jsonFile struct {
	Nodes []jsonNode         `json:"nodes"`
	Rels  []jsonRelationship `json:"rels"`
}

type Graph struct {
	Nodes         []geno.Node
	Relationships []geno.Relationship
}

func getGraphFromJson(raw []byte) (g Graph, err error) {
	var js jsonFile

	if err := json.Unmarshal(raw, &js); err != nil {
		return g, err
	}

	g.Nodes = make([]geno.Node, len(js.Nodes))
	for i := range js.Nodes {
		rawN := js.Nodes[i]
		g.Nodes[i] = geno.NewNode(rawN.Id, rawN.Labels, rawN.Props)
	}

	g.Relationships = make([]geno.Relationship, len(js.Rels))
	for i := range js.Rels {
		rawR := js.Rels[i]
		start, err := findNodeById(rawR.Start, g.Nodes)
		if err != nil {
			return g, err
		}
		end, err := findNodeById(rawR.End, g.Nodes)
		if err != nil {
			return g, err
		}
		g.Relationships[i] = geno.NewRelationship(rawR.Id, start, end, rawR.Label, rawR.Properties)
	}
	return g, nil
}

func findNodeById(id int64, nodes []geno.Node) (geno.Node, error) {
	for i := range nodes {
		if id == nodes[i].Id {
			return nodes[i], nil
		}
	}
	return geno.EmptyNode, fmt.Errorf("node with id %d could not be found", id)
}

func ImportJson(uri, database, username, password, filepath string) error {
	// raw, err := os.ReadFile(path.Join("test_data", "test.json"))
	var graph Graph

	raw, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	graph, err = getGraphFromJson(raw)
	if err != nil {
		return err
	}

	driver, err := geno.NewDriver("neo4j://"+uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return err
	}
	defer driver.Close()

	var nodeConstraints geno.Constraints
	nodeConstraints, err = driver.GetConstraints(database)
	fmt.Println("CONSTRAINTS")
	fmt.Println(nodeConstraints)
	if err != nil {
		return err
	}

	var query geno.Query = geno.NewQuery(&driver, &nodeConstraints)

	for _, node := range graph.Nodes {
		summary, err := query.MergeNode(database, node)
		if err != nil {
			return err
		}
		fmt.Printf("merged %d node(s) with label %s\n", summary.Counters().NodesCreated(), node)
	}

	for _, rel := range graph.Relationships {
		summary, err := query.MergeRelationship(database, rel)
		if err != nil {
			return err
		}
		fmt.Printf("merged %d relationship(s) of type %s between %s and %s\n", summary.Counters().RelationshipsCreated(), rel.Label, rel.Start, rel.End)
	}

	return nil

}
