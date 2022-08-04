package common

import (
	"sort"
	"strings"
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

func (n *Node) ToCypherMerge(constraints []string, paramPrefix string) (query string, params map[string]any) {
	var (
		q                          strings.Builder = strings.Builder{}
		constrainedProps           map[string]any  = make(map[string]any)
		constrainedPropsTemplate   []string
		unconstrainedProps         map[string]any = make(map[string]any)
		unconstrainedPropsTemplate []string
	)
	sort.Strings(constraints)                        // for more stable testing/query generation
	params = make(map[string]any, len(n.Properties)) // all properties are eventually parameritized

	// segregate constrained props from unconstrained props
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

	constrainedPropsTemplate = templatizeProps(constrainedProps, ":", paramPrefix)
	unconstrainedPropsTemplate = templatizeProps(unconstrainedProps, "=", paramPrefix)

	q.WriteString("MERGE (n:")
	q.WriteString(n.String())
	if len(constrainedProps) > 0 {
		q.WriteString(" {")
		q.WriteString(strings.Join(constrainedPropsTemplate, ", "))
		q.WriteString("}")
	}
	q.WriteString(")")
	if len(unconstrainedProps) > 0 {
		q.WriteString("\nON CREATE SET n.")
		q.WriteString(strings.Join(unconstrainedPropsTemplate, ", n."))
	}
	q.WriteString("\n")

	for key, val := range n.Properties {
		params[paramPrefix+key] = val
	}

	return q.String(), params
}

func (n *Node) ToCypherMatch(constraints []string, paramPrefix string) (query string, params map[string]interface{}) {
	var (
		q                        strings.Builder = strings.Builder{}
		constrainedProps         map[string]any  = make(map[string]any)
		constrainedPropsTemplate []string
	)
	sort.Strings(constraints)                        // for more stable testing/query generation
	params = make(map[string]any, len(n.Properties)) // all properties are eventually parameritized

	// segregate constrained props from unconstrained props
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
		}
	}

	constrainedPropsTemplate = templatizeProps(constrainedProps, ":", paramPrefix)

	q.WriteString("MATCH (n:")
	q.WriteString(n.String())
	if len(constrainedProps) > 0 {
		q.WriteString(" {")
		q.WriteString(strings.Join(constrainedPropsTemplate, ", "))
		q.WriteString("}")
	}
	q.WriteString(")\n")

	for key, val := range constrainedProps {
		params[paramPrefix+key] = val
	}

	return q.String(), params
}

func (n *Node) ToCypherCreate(paramPrefix string) (query string, params map[string]interface{}) {
	var (
		q             strings.Builder = strings.Builder{}
		propKeys      []string        = make([]string, 0, len(n.Properties))
		propsTemplate []string
	)
	params = make(map[string]any, len(n.Properties)) // all properties are eventually parameritized

	for key := range n.Properties {
		propKeys = append(propKeys, key)
	}
	sort.Strings(propKeys)

	propsTemplate = templatizeProps(n.Properties, ":", paramPrefix)

	q.WriteString("CREATE (n:")
	q.WriteString(n.String())
	if len(n.Properties) > 0 {
		q.WriteString(" {")
		q.WriteString(strings.Join(propsTemplate, ", "))
		q.WriteString("}")
	}
	q.WriteString(")\n")

	for key, val := range n.Properties {
		params[paramPrefix+key] = val
	}

	return q.String(), params
}
