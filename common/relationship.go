package common

import (
	"sort"
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
	Label      string
	Properties map[string]any
}

// String prints all types of a Relationship in neo4j format
func (r *Relationship) String() string { return r.Label }

func NewRelationship(id int64, from, to Node, label string, props map[string]any) Relationship {
	return Relationship{
		Id:         id,
		Start:      from,
		End:        to,
		Label:      label,
		Properties: props,
	}
}

func (r *Relationship) ToCypherMerge(leftNodeConstraints, rightNodeConstraints, relConstraints []string) (query string, params map[string]interface{}) {
	var (
		leftMatchQuery             string
		leftMatchParams            map[string]any
		rightMatchQuery            string
		rightMatchParams           map[string]any
		q                          strings.Builder = strings.Builder{}
		constrainedProps           map[string]any  = make(map[string]any)
		constrainedPropsTemplate   []string
		unconstrainedProps         map[string]any = make(map[string]any)
		unconstrainedPropsTemplate []string
		paramPrefix                string = "rel"
	)
	sort.Strings(relConstraints)
	params = make(map[string]interface{}) // check the length of the parameters at the end and panic if not lengths of r.Properties + leftMatchParams + rightMatchParams

	leftMatchQuery, leftMatchParams = r.Start.ToCypherMatch(leftNodeConstraints, "left")
	rightMatchQuery, rightMatchParams = r.End.ToCypherMatch(rightNodeConstraints, "right")

	// segregate constrained props from unconstrained props
	for key, val := range r.Properties {
		var isConstrained bool = false
		for _, constrainedKey := range relConstraints {
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

	q.WriteString(leftMatchQuery)
	q.WriteString(rightMatchQuery)
	q.WriteString("MERGE (left)-[r:")
	q.WriteString(r.String())
	if len(constrainedProps) > 0 {
		q.WriteString(" {")
		q.WriteString(strings.Join(constrainedPropsTemplate, ", "))
		q.WriteString("}")
	}
	q.WriteString("]-(right)")
	if len(unconstrainedProps) > 0 {
		q.WriteString("\nON CREATE SET n.")
		q.WriteString(strings.Join(unconstrainedPropsTemplate, ", n."))
	}
	q.WriteString("\n")

	for key, val := range r.Properties {
		params[paramPrefix+key] = val
	}
	for key, val := range leftMatchParams {
		params[key] = val //do not use the paramPrefix here, as it is statically set to "rel" and templatize should have already set them
	}
	for key, val := range rightMatchParams {
		params[key] = val //do not use the paramPrefix here, as it is statically set to "rel" and templatize should have already set them
	}

	if len(params) != (len(r.Properties) + len(rightMatchParams) + len(leftMatchParams)) {
		panic("Relationship templating failed because of mismatch in query parameters")
	}

	return q.String(), params
}

func (r *Relationship) ToCypherMatch(leftNodeConstraints, rightNodeConstraints, relConstraints []string) (query string, params map[string]interface{}) {
	var (
		leftMatchQuery           string
		leftMatchParams          map[string]any
		rightMatchQuery          string
		rightMatchParams         map[string]any
		q                        strings.Builder = strings.Builder{}
		constrainedProps         map[string]any  = make(map[string]any)
		constrainedPropsTemplate []string
		paramPrefix              string = "rel"
	)
	sort.Strings(relConstraints)
	params = make(map[string]interface{}) // check the length of the parameters at the end and panic if not lengths of r.Properties + leftMatchParams + rightMatchParams

	leftMatchQuery, leftMatchParams = r.Start.ToCypherMatch(leftNodeConstraints, "left")
	rightMatchQuery, rightMatchParams = r.End.ToCypherMatch(rightNodeConstraints, "right")

	// segregate constrained props from unconstrained props
	for key, val := range r.Properties {
		var isConstrained bool = false
		for _, constrainedKey := range relConstraints {
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

	q.WriteString(leftMatchQuery)
	q.WriteString(rightMatchQuery)
	q.WriteString("MERGE (left)-[r:")
	q.WriteString(r.String())
	if len(constrainedProps) > 0 {
		q.WriteString(" {")
		q.WriteString(strings.Join(constrainedPropsTemplate, ", "))
		q.WriteString("}")
	}
	q.WriteString("]-(right)")
	q.WriteString("\n")

	for key, val := range r.Properties {
		params[paramPrefix+key] = val
	}
	for key, val := range leftMatchParams {
		params[key] = val //do not use the paramPrefix here, as it is statically set to "rel" and templatize should have already set them
	}
	for key, val := range rightMatchParams {
		params[key] = val //do not use the paramPrefix here, as it is statically set to "rel" and templatize should have already set them
	}

	if len(params) != (len(r.Properties) + len(rightMatchParams) + len(leftMatchParams)) {
		panic("Relationship templating failed because of mismatch in query parameters")
	}

	return q.String(), params
}

func (r *Relationship) ToCypherCreate(leftNodeConstraints, rightNodeConstraints []string) (query string, params map[string]interface{}) {
	var (
		leftMatchQuery   string
		leftMatchParams  map[string]any
		rightMatchQuery  string
		rightMatchParams map[string]any
		q                strings.Builder = strings.Builder{}
		propKeys         []string        = make([]string, 0, len(r.Properties))
		relPropsTemplate []string
		paramPrefix      string = "rel"
	)

	for key := range r.Properties {
		propKeys = append(propKeys, key)
	}
	sort.Strings(propKeys)
	params = make(map[string]interface{}) // check the length of the parameters at the end and panic if not lengths of r.Properties + leftMatchParams + rightMatchParams

	leftMatchQuery, leftMatchParams = r.Start.ToCypherMatch(leftNodeConstraints, "left")
	rightMatchQuery, rightMatchParams = r.End.ToCypherMatch(rightNodeConstraints, "right")

	relPropsTemplate = templatizeProps(r.Properties, ":", paramPrefix)

	q.WriteString(leftMatchQuery)
	q.WriteString(rightMatchQuery)
	q.WriteString("CREATE (left)-[r:")
	q.WriteString(r.String())
	if len(relPropsTemplate) > 0 {
		q.WriteString(" {")
		q.WriteString(strings.Join(relPropsTemplate, ", "))
		q.WriteString("}")
	}
	q.WriteString("]-(right)")
	q.WriteString("\n")

	for key, val := range r.Properties {
		params[paramPrefix+key] = val
	}
	for key, val := range leftMatchParams {
		params[key] = val //do not use the paramPrefix here, as it is statically set to "rel" and templatize should have already set them
	}
	for key, val := range rightMatchParams {
		params[key] = val //do not use the paramPrefix here, as it is statically set to "rel" and templatize should have already set them
	}

	if len(params) != (len(r.Properties) + len(rightMatchParams) + len(leftMatchParams)) {
		panic("Relationship templating failed because of mismatch in query parameters")
	}

	return q.String(), params
}
