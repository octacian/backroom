//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Record = newRecordTable("public", "record", "")

type recordTable struct {
	postgres.Table

	// Columns
	UUID postgres.ColumnString
	Cage postgres.ColumnString
	Data postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
	DefaultColumns postgres.ColumnList
}

type RecordTable struct {
	recordTable

	EXCLUDED recordTable
}

// AS creates new RecordTable with assigned alias
func (a RecordTable) AS(alias string) *RecordTable {
	return newRecordTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new RecordTable with assigned schema name
func (a RecordTable) FromSchema(schemaName string) *RecordTable {
	return newRecordTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new RecordTable with assigned table prefix
func (a RecordTable) WithPrefix(prefix string) *RecordTable {
	return newRecordTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new RecordTable with assigned table suffix
func (a RecordTable) WithSuffix(suffix string) *RecordTable {
	return newRecordTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newRecordTable(schemaName, tableName, alias string) *RecordTable {
	return &RecordTable{
		recordTable: newRecordTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newRecordTableImpl("", "excluded", ""),
	}
}

func newRecordTableImpl(schemaName, tableName, alias string) recordTable {
	var (
		UUIDColumn     = postgres.StringColumn("uuid")
		CageColumn     = postgres.StringColumn("cage")
		DataColumn     = postgres.StringColumn("data")
		allColumns     = postgres.ColumnList{UUIDColumn, CageColumn, DataColumn}
		mutableColumns = postgres.ColumnList{CageColumn, DataColumn}
		defaultColumns = postgres.ColumnList{}
	)

	return recordTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UUID: UUIDColumn,
		Cage: CageColumn,
		Data: DataColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
