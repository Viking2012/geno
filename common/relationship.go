package common

import (
	"strings"
)

var (
	EmptyRelationship  Relationship   = Relationship{}
	EmptyRelationships []Relationship = []Relationship{}
)

// Relationship is a representation of a neo4j driver Relationship
type Relationship struct {
	Id         int64
	Start      Node
	End        Node
	Types      []string
	Properties map[string]any
}

// String prints all types of a Relationship in neo4j format
func (r *Relationship) String() string { return strings.Join(r.Types[:], ":") }

func NewRelationship(id int64, from, to Node, types []string, props map[string]any) Relationship {
	return Relationship{
		Id:         id,
		Start:      from,
		End:        to,
		Types:      types,
		Properties: props,
	}
}
