package pkg

import (
	"encoding/json"

	"github.com/Viking2012/geno/geno"
)

type readNode struct {
	Id     int64          `json:"identity"`
	Labels []string       `json:"labels"`
	Props  map[string]any `json:"properties"`
}

type readRelationship struct {
	Id         int64          `json:"identity"`
	Start      int64          `json:"start"`
	End        int64          `json:"end"`
	Label      string         `json:"type"`
	Properties map[string]any `json:"properties"`
}
type readGraph struct {
	Nodes []readNode         `json:"nodes"`
	Rels  []readRelationship `json:"rels"`
}

type Graph struct {
	Nodes         []geno.Node
	Relationships []geno.Relationship
}

func GetGraphFromJson(raw []byte) (g Graph, err error) {
	var js readGraph

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
