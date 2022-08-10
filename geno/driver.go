package geno

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Driver struct {
	neo4j.Driver
}

func NewDriver(uri string, auth neo4j.AuthToken) (Driver, error) {
	driver, err := neo4j.NewDriver(uri, auth)
	if err != nil {
		return Driver{}, err
	}
	return Driver{driver}, nil
}

func (d *Driver) GetConstraints(database string) (Constraints, error) {
	session := d.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: database})
	defer session.Close()

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

		return result.Consume()
	})
	if err != nil {
		return Constraints{}, err
	}

	c, err := ConstraintsFromRecords(records)
	if err != nil {
		return Constraints{}, err
	}

	return c, nil
}
