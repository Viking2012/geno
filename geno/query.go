package geno

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type Query struct {
	d *Driver
	c *Constraints
}

func NewQuery(driver *Driver, constraints *Constraints) Query {
	return Query{
		d: driver,
		c: constraints,
	}
}

func (q *Query) MergeNode(database string, n Node) (neo4j.ResultSummary, error) {
	session := q.d.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: database})
	defer session.Close()

	var (
		querySummary neo4j.ResultSummary
		records      []*neo4j.Record
		constraints  []string
	)

	constraints = q.c.GetNodeConstraints(&n)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, txErr := tx.Run(n.ToCypherMerge(constraints, "n"))
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
		return querySummary, err
	}

	return querySummary, nil
}

func (q *Query) MergeRelationship(database string, r Relationship) (neo4j.ResultSummary, error) {
	session := q.d.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: database})
	defer session.Close()

	var (
		querySummary     neo4j.ResultSummary
		records          []*neo4j.Record
		leftConstraints  []string
		rightConstraints []string
		relConstraints   []string
	)

	leftConstraints = q.c.GetNodeConstraints(&r.Start)
	rightConstraints = q.c.GetNodeConstraints(&r.End)
	relConstraints = q.c.GetRelationshipConstraints(&r)
	relConstraints = []string{}

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, txErr := tx.Run(r.ToCypherMerge(leftConstraints, rightConstraints, relConstraints))
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
		return querySummary, err
	}

	return querySummary, nil
}
