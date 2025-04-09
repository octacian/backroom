package cage

import (
	"encoding/json"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/octacian/backroom/api/.gen/backroom/public/model"
	"github.com/octacian/backroom/api/.gen/backroom/public/table"
	"github.com/octacian/backroom/api/db"
)

// Record is a caged entry, identified by UUID, grouped by cage key, and
// containing some data.
// Wraps generated model.Record type.
type Record model.Record

// NewRecord returns a new caged record, identified by a key and storing some data.
func NewRecord(cage string, data db.JSONB) *Record {
	return &Record{
		UUID: db.NewUUID(),
		Cage: cage,
		Data: data,
	}
}

// NewRecordFromString returns a new caged record, identified by a cage key
// and storing some JSON data marshalled to string format.
func NewRecordFromString(key string, data string) (*Record, error) {
	var jsonb db.JSONB
	if err := json.Unmarshal([]byte(data), &jsonb); err != nil {
		return nil, err
	}

	return NewRecord(key, jsonb), nil
}

// CreateRecord creates a new caged record in the database
func CreateRecord(cage *Record) error {
	insert := table.Record.INSERT(table.Record.AllColumns).MODEL(cage)

	_, err := insert.Exec(db.SQLDB)
	if err != nil {
		return err
	}

	return nil
}

// GetRecord retrieves a specific record from the database by its UUID.
func GetRecord(uuid db.UUID) (*Record, error) {
	stmt := table.Record.SELECT(table.Record.AllColumns).
		WHERE(table.Record.UUID.EQ(postgres.UUID(uuid))).
		ORDER_BY(table.Record.UUID.DESC()).
		LIMIT(1)

	var cage Record
	err := stmt.Query(db.SQLDB, &cage)
	if err != nil {
		return nil, err
	}

	return &cage, nil
}

// ListRecordsByCage retrieves all records belonging to a common cage from the database.
func ListRecordsByCage(cage string) ([]*Record, error) {
	stmt := table.Record.SELECT(table.Record.AllColumns).
		WHERE(table.Record.Cage.EQ(postgres.String(cage))).
		ORDER_BY(table.Record.UUID.DESC())

	var cages []*Record
	err := stmt.Query(db.SQLDB, &cages)
	if err != nil {
		return nil, err
	}

	return cages, nil
}

// ListCages retrieves all unique cages from the database.
func ListCages() ([]string, error) {
	stmt := table.Record.SELECT(table.Record.Cage).DISTINCT()

	var cages []string
	err := stmt.Query(db.SQLDB, &cages)
	if err != nil {
		return nil, err
	}

	return cages, nil
}

// DeleteRecord deletes a record from the database by its UUID.
func DeleteRecord(uuid db.UUID) error {
	stmt := table.Record.DELETE().
		WHERE(table.Record.UUID.EQ(postgres.UUID(uuid)))

	_, err := stmt.Exec(db.SQLDB)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCage deletes all records belonging to a common cage from the database.
// Returns the number of deleted records.
func DeleteCage(cage string) (int64, error) {
	stmt := table.Record.DELETE().
		WHERE(table.Record.Cage.EQ(postgres.String(cage)))

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
