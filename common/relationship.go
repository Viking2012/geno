package common

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

func (r *Relationship) ToCypherMerge() (query string, params map[string]interface{}) {
	return "", nil
}

func (r *Relationship) ToCypherMatch() (query string, params map[string]interface{}) {
	return "", nil
}

func (r *Relationship) ToCypherCreate() (query string, params map[string]interface{}) {
	return "", nil
}
