package common

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/graph/encoding"
)

var (
	EmptyNode  Node   = Node{}
	EmptyNodes []Node = []Node{}
)

// Node is a representation of a neo4j driver Node
type Node struct {
	Id         int64
	Labels     []string
	Properties map[string]any
}

// ID allows Node to satisfy the interface requirements of a gonum graph.Node
func (n Node) ID() int64 {
	return n.Id
}

// String prints all labels of a Node in a neo4j format
func (n Node) String() string { return strings.Join(n.Labels[:], ":") }

func NewNode(id int64, labels []string, props map[string]any) (n Node) {
	return Node{
		Id:         id,
		Labels:     labels,
		Properties: props,
	}
}

func (n *Node) Attributes() (attrs []encoding.Attribute) {
	for k, v := range n.Properties {
		attrs = append(attrs, encoding.Attribute{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	return attrs
}

func (n *Node) ToCypherMerge(constraints []string) string {
	var (
		constrainedProps   map[string]any = make(map[string]any)
		unconstrainedProps map[string]any = make(map[string]any)
	)

	for key, val := range n.Properties {
		var isConstrained bool = false
		for _, constrainedKey := range constraints {
			if key == constrainedKey {
				isConstrained = true
				break
			}
		}
		if isConstrained {
			constrainedProps[key] = val
		} else {
			unconstrainedProps[key] = val
		}
	}

	query := strings.Builder{}
	query.WriteString("MERGE (n {")
	writtenContraints := 0
	for key, val := range constrainedProps {
		if writtenContraints > 0 {
			query.WriteString(",")
		}
		query.WriteString(fmt.Sprintf("%s:`%v`", key, val))
		writtenContraints++
	}
	query.WriteString("} SET ")
	if len(n.Labels) != 0 {
		query.WriteString(fmt.Sprintf("n:%s", strings.Join(n.Labels, ":")))
		if len(unconstrainedProps) != 0 {
			query.WriteString(" AND ")
		}
	}
	writtenProps := 0
	for key, val := range unconstrainedProps {
		if writtenProps > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString(fmt.Sprintf("n.%s = %v", key, val))
	}

	return query.String()
}
