package geno

import (
	"errors"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

type EntityType string
type ConstraintType string

const (
	IS_NODE                         EntityType     = "NODE"
	IS_RELATIONSHIP                 EntityType     = "RELATIONSHIP"
	UNKNOWN_ENTITY                  EntityType     = "_"
	NODE_UNIQUE_CONSTRAINT          ConstraintType = "UNIQUENESS"
	NODE_KEY_CONSTRAINT             ConstraintType = "NODE_KEY"
	NODE_PROPERTY_EXISTS_CONSTRAINT ConstraintType = "NODE_PROPERTY_EXISTENCE"
	REL_UNIQUE_CONSTRAINT           ConstraintType = "RELATIONSHIP_UNIQUENESS"
	REL_KEY_CONSTRAINT              ConstraintType = "RELATIONSHIP_KEY"
	REL_PROPERTY_EXISTS_CONSTRAINT  ConstraintType = "RELATIONSHIP_PROPERTY_EXISTENCE"
)

type Constraint struct {
	Label      string
	Properties []string
}

type Constraints struct {
	NodeUniqueness                []Constraint
	NodeKeys                      []Constraint
	NodePropertyExistence         []Constraint
	RelationshipUniqueness        []Constraint
	RelationshipKeys              []Constraint
	RelationshipPropertyExistence []Constraint
}

func (c *Constraints) AddConstraint(entityType EntityType, constraintType ConstraintType, newConstraint Constraint) error {
	switch entityType {
	case IS_NODE:
		switch constraintType {
		case NODE_UNIQUE_CONSTRAINT:
			c.NodeUniqueness = append(c.NodeUniqueness, newConstraint)
		case NODE_KEY_CONSTRAINT:
			c.NodeKeys = append(c.NodeKeys, newConstraint)
		case NODE_PROPERTY_EXISTS_CONSTRAINT:
			c.NodePropertyExistence = append(c.NodePropertyExistence, newConstraint)
		default:
			return fmt.Errorf("node constraint type %s could not be added as a constraint", constraintType)
		}
	case IS_RELATIONSHIP:
		switch constraintType {
		case REL_UNIQUE_CONSTRAINT:
			c.RelationshipUniqueness = append(c.RelationshipUniqueness, newConstraint)
		case REL_KEY_CONSTRAINT:
			c.RelationshipKeys = append(c.RelationshipKeys, newConstraint)
		case REL_PROPERTY_EXISTS_CONSTRAINT:
			c.RelationshipPropertyExistence = append(c.RelationshipPropertyExistence, newConstraint)
		default:
			return fmt.Errorf("relationship constraint type %s could not be added as a constraint", constraintType)
		}
	case UNKNOWN_ENTITY:
		return fmt.Errorf("entity type %s could not be added as a constraint", entityType)
	}
	return nil
}

func ConstraintsFromRecords(records []*db.Record) (Constraints, error) {
	var c Constraints

	for _, record := range records {
		var (
			rawField       interface{}
			listField      []interface{}
			labelOrTypes   []string
			entityType     string
			properties     []string
			constraintType string
			ok             bool
		)
		rawField, found := record.Get("labelsOrTypes")
		if !found {
			return c, errors.New("a constraint was found which did not have a label or type")
		}
		listField, ok = rawField.([]interface{})
		if !ok {
			return c, errors.New("a constraint was found which did not have a properly formatted label or type")
		}
		labelOrTypes = make([]string, len(listField))
		if len(labelOrTypes) > 1 {
			return c, errors.New("a constraint was applied to more than one node label or type")
		}
		for i, l := range listField {
			s, ok := l.(string)
			if !ok {
				return c, errors.New("a constraint was found which did not have a properly formatted label or type")
			}
			labelOrTypes[i] = s
		}

		rawField, found = record.Get("entityType")
		if !found {
			return c, errors.New("a constraint was found which did not have an entity type")
		}
		entityType, ok = rawField.(string)
		if !ok {
			return c, errors.New("a constraint was found which did not have a properly formatted entity type")
		}

		rawField, found = record.Get("properties")
		if !found {
			return c, errors.New("a constraint was found which did not have properties")
		}
		listField, ok = rawField.([]interface{})
		if !ok {
			return c, errors.New("a constraint was found which did not have a properly formatted label or type")
		}
		properties = make([]string, len(listField))
		for i, p := range listField {
			s, ok := p.(string)
			if !ok {
				return c, errors.New("a constraint was found which did not have a properly formatted list of properties")
			}
			properties[i] = s
		}

		rawField, found = record.Get("type")
		if !found {
			return c, errors.New("a constraint was found which did not have a type")
		}
		constraintType, ok = rawField.(string)
		if !ok {
			return c, errors.New("a constraint was found which did not have a properly formatted type")
		}

		for _, label := range labelOrTypes {
			var newConstraint Constraint = Constraint{Label: label, Properties: properties}
			c.AddConstraint(EntityType(entityType), ConstraintType(constraintType), newConstraint)
		}

	}

	return c, nil
}

func (c *Constraints) String() string {
	ret := strings.Builder{}

	ret.WriteString("NODE UNIQUENESS")
	for _, u := range c.NodeUniqueness {
		ret.WriteString("\n\t")
		ret.WriteString(u.Label)
		ret.WriteString(": ")
		for i, p := range u.Properties {
			if i > 0 {
				ret.WriteString(", ")
			}
			ret.WriteString(p)
		}
	}
	ret.WriteString("\nNODE KEYS")
	for _, u := range c.NodeKeys {
		ret.WriteString("\n\t")
		ret.WriteString(u.Label)
		ret.WriteString(": ")
		for i, p := range u.Properties {
			if i > 0 {
				ret.WriteString(", ")
			}
			ret.WriteString(p)
		}
	}
	ret.WriteString("\nNODE REQUIRED PROPERTIES")
	for _, u := range c.NodePropertyExistence {
		ret.WriteString("\n\t")
		ret.WriteString(u.Label)
		ret.WriteString(": ")
		for i, p := range u.Properties {
			if i > 0 {
				ret.WriteString(", ")
			}
			ret.WriteString(p)
		}
	}
	ret.WriteString("REL  REQUIRED PROPERTIES")
	for _, u := range c.RelationshipPropertyExistence {
		ret.WriteString("\n\t")
		ret.WriteString(u.Label)
		ret.WriteString(": ")
		for i, p := range u.Properties {
			if i > 0 {
				ret.WriteString(", ")
			}
			ret.WriteString(p)
		}
	}

	return ret.String()
}

func (constraints *Constraints) GetNodeConstraints(n *Node) []string {
	var (
		reducer        map[string]byte = make(map[string]byte)
		allConstraints []string
	)

	for _, label := range n.Labels {
		for _, c := range constraints.NodeKeys {
			if c.Label == label {
				for _, p := range c.Properties {
					reducer[p] = 1
				}
			}
		}
		for _, c := range constraints.NodePropertyExistence {
			if c.Label == label {
				for _, p := range c.Properties {
					reducer[p] = 1
				}
			}
		}
		for _, c := range constraints.NodeUniqueness {
			if c.Label == label {
				for _, p := range c.Properties {
					reducer[p] = 1
				}
			}
		}
	}

	for key := range reducer {
		allConstraints = append(allConstraints, key)
	}

	return allConstraints
}

func (constraints *Constraints) GetRelationshipConstraints(r *Relationship) []string {
	var (
		reducer        map[string]byte = make(map[string]byte)
		allConstraints []string
	)

	for _, c := range constraints.RelationshipKeys {
		if c.Label == r.Label {
			for _, p := range c.Properties {
				reducer[p] = 1
			}
		}
	}
	for _, c := range constraints.RelationshipPropertyExistence {
		if c.Label == r.Label {
			for _, p := range c.Properties {
				reducer[p] = 1
			}
		}
	}
	for _, c := range constraints.RelationshipUniqueness {
		if c.Label == r.Label {
			for _, p := range c.Properties {
				reducer[p] = 1
			}
		}
	}

	for key := range reducer {
		allConstraints = append(allConstraints, key)
	}

	return allConstraints
}

func ReadConstraints(driver neo4j.Driver, database string) (Constraints, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: database})
	defer session.Close()

	var querySummary neo4j.ResultSummary
	var records []*neo4j.Record

	_, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, txErr := tx.Run("SHOW CONSTRAINTS", nil)
		if txErr != nil {
			return nil, txErr
		}

		for result.Next() {
			record := result.Record()
			records = append(records, record)
		}

		// return result.Consume()
		summary, summaryErr := result.Consume()
		querySummary = summary
		return summary, summaryErr
	})
	if err != nil {
		return Constraints{}, err
	}

	fmt.Print("Has updates: ")
	fmt.Println(querySummary.Counters().ContainsUpdates())

	c, err := ConstraintsFromRecords(records)
	if err != nil {
		return Constraints{}, err
	}

	return c, nil
}
