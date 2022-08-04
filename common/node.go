package common

import (
	"fmt"
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

func (n *Node) ToCypherMerge(constraints []string) (query string, params map[string]interface{}) {
	var (
		constrainedProps   map[string]any = make(map[string]any)
		constrainedKeys    []string       // set and sorted later for more stable testing/query generation
		unconstrainedProps map[string]any = make(map[string]any)
		unconstrainedKeys  []string       // set and sorted later for more stable testing/query generation
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
	// saving and sorting the relevant keys ensures that props are listed in alphabetical order
	// this make no difference in imterpreting the generated cypher command, but it critical for testing
	constrainedKeys = make([]string, 0, len(constrainedProps))
	for key := range constrainedProps {
		constrainedKeys = append(constrainedKeys, key)
	}
	sort.Strings(constrainedKeys)
	unconstrainedKeys = make([]string, 0, len(unconstrainedProps))
	for key := range unconstrainedProps {
		unconstrainedKeys = append(unconstrainedKeys, key)
	}
	sort.Strings(unconstrainedKeys)

	// initiate the bulding of the query
	query := strings.Builder{}
	query.WriteString("MERGE (n {") // first comes the merge statement

	// Constrained properties - keep track of the number writter
	writtenContraints := 0
	for _, key := range constrainedKeys {
		val := constrainedProps[key]
		if writtenContraints > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("%s:'%v'", key, val))
		writtenContraints++
	}

	query.WriteString("}) SET ") // then the properties we will set on the merge (total overwrite of existing properties)

	// Uncontrained properties - keep track of the number written
	writtenSetTerms := 0
	if len(n.Labels) != 0 {
		query.WriteString(fmt.Sprintf("n:%s", strings.Join(n.Labels, ":")))
		writtenSetTerms++
	}
	for _, key := range unconstrainedKeys {
		val := unconstrainedProps[key]
		if writtenSetTerms > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("n.%s='%v'", key, val))
		writtenSetTerms++
	}

	return query.String()
}

func (n *Node) ToCypherMatch(constraints []string) (query string, params map[string]interface{}) {
	var (
		constrainedProps map[string]any = make(map[string]any)
		constrainedKeys  []string       // set and sort later
	)

	for key, val := range n.Properties {
		for _, constrainedKey := range constraints {
			if key == constrainedKey {
				constrainedProps[key] = val
				break
			}
		}
	}
	constrainedKeys = make([]string, 0, len(constrainedProps))
	for key := range constrainedProps {
		constrainedKeys = append(constrainedKeys, key)
	}
	sort.Strings(constrainedKeys)

	query := strings.Builder{}
	query.WriteString(fmt.Sprintf("MATCH (n:%s {", strings.Join(n.Labels, ":")))
	writtenContraints := 0
	for _, key := range constrainedKeys {
		val := constrainedProps[key]
		if writtenContraints > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("%s:'%v'", key, val))
		writtenContraints++
	}
	query.WriteString("})")
	return query.String()
}

func (n *Node) ToCypherCreate() (query string, params map[string]interface{}) {
	var propKeys []string = make([]string, 0, len(n.Properties))
	for key := range n.Properties {
		propKeys = append(propKeys, key)
	}
	sort.Strings(propKeys)

	query := strings.Builder{}
	query.WriteString(fmt.Sprintf("CREATE (n:%s {", strings.Join(n.Labels, ":")))
	writtenContraints := 0

	for _, key := range propKeys {
		val := n.Properties[key]
		if writtenContraints > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("%s:'%v'", key, val))
		writtenContraints++
	}
	query.WriteString("})")
	return query.String()
}
