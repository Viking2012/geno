package geno

import (
	"errors"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

type Constraint struct {
	Label      string
	Properties []string
}

type Constraints struct {
	NodeUniqueness                []Constraint
	NodeKeys                      []Constraint
	NodePropertyExistence         []Constraint
	RelationshipPropertyExistence []Constraint
}

func ConstraintsFromRecords(records []*db.Record) (Constraints, error) {
	var (
		c         Constraints
		uniques   []Constraint = make([]Constraint, 0)
		keys      []Constraint = make([]Constraint, 0)
		nodeProps []Constraint = make([]Constraint, 0)
		relProps  []Constraint = make([]Constraint, 0)
	)

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
			if entityType == "NODE" {
				switch constraintType {
				case "UNIQUENESS":
					uniques = append(uniques, newConstraint)
				case "NODE_KEY":
					keys = append(keys, newConstraint)
				case "NODE_PROPERTY_EXISTENCE":
					nodeProps = append(nodeProps, newConstraint)
				}
			} else {
				relProps = append(relProps, newConstraint)
			}

		}

	}

	c = Constraints{
		NodeUniqueness:                uniques,
		NodeKeys:                      keys,
		NodePropertyExistence:         nodeProps,
		RelationshipPropertyExistence: relProps,
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

func (constraints *Constraints) GetConstraintsFor(n Node) []string {
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
