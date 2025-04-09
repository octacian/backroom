package cage

import (
	"encoding/json"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/octacian/backroom/api/.gen/backroom/public/model"
	"github.com/octacian/backroom/api/.gen/backroom/public/table"
	"github.com/octacian/backroom/api/db"
)

// Cage stores
type Cage model.Cage

// NewRecord returns a new caged record, identified by a key and storing some data.
func NewRecord(key string, data db.JSONB) *Cage {
	return &Cage{
		UUID: db.NewUUID(),
		Key:  key,
		Data: data,
	}
}

// NewRecordFromString returns a new caged record, identified by a key and
// storing some data in string format.
func NewRecordFromString(key string, data string) (*Cage, error) {
	var jsonb db.JSONB
	if err := json.Unmarshal([]byte(data), &jsonb); err != nil {
		return nil, err
	}

	return NewRecord(key, jsonb), nil
}

// CreateRecord creates a new caged record in the database
func CreateRecord(cage *Cage) error {
	insert := table.Cage.INSERT(table.Cage.AllColumns).MODEL(cage)

	_, err := insert.Exec(db.SQLDB)
	if err != nil {
		return err
	}

	return nil
}

// GetRecord retrieves a specific caged record from the database by its UUID.
func GetRecord(uuid db.UUID) (*Cage, error) {
	stmt := table.Cage.SELECT(table.Cage.AllColumns).
		WHERE(table.Cage.UUID.EQ(postgres.UUID(uuid))).
		ORDER_BY(table.Cage.UUID.DESC()).
		LIMIT(1)

	var cage Cage
	err := stmt.Query(db.SQLDB, &cage)
	if err != nil {
		return nil, err
	}

	return &cage, nil
}

// ListRecordsByKey retrieves all caged records with a common key from the database.
func ListRecordsByKey(key string) ([]*Cage, error) {
	stmt := table.Cage.SELECT(table.Cage.AllColumns).
		WHERE(table.Cage.Key.EQ(postgres.String(key))).
		ORDER_BY(table.Cage.UUID.DESC())

	var cages []*Cage
	err := stmt.Query(db.SQLDB, &cages)
	if err != nil {
		return nil, err
	}

	return cages, nil
}

// ListCageKeys retrieves all unique cage keys from the database.
func ListCageKeys() ([]string, error) {
	stmt := table.Cage.SELECT(table.Cage.Key).DISTINCT()

	var keys []string
	err := stmt.Query(db.SQLDB, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// DeleteRecord deletes a caged record from the database by its UUID.
func DeleteRecord(uuid db.UUID) error {
	stmt := table.Cage.DELETE().
		WHERE(table.Cage.UUID.EQ(postgres.UUID(uuid)))

	_, err := stmt.Exec(db.SQLDB)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCage deletes all caged records with a common key from the database.
// Returns the number of deleted records.
func DeleteCage(key string) (int64, error) {
	stmt := table.Cage.DELETE().
		WHERE(table.Cage.Key.EQ(postgres.String(key)))

	res, err := stmt.Exec(db.SQLDB)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
