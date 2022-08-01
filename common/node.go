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

// DOTIDD allows Node to have better labels in dot files
func (n Node) DOTID() string {
	if len(n.Labels) == 0 {
		return fmt.Sprintf("%d", n.ID())
	}
	return fmt.Sprintf("%s, %d", n.String(), n.ID())
}

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
