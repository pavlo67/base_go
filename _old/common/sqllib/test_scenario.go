package sqllib

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T, db *sql.DB) {
	require.NotNil(t, db)

	// prepare table -----------------------------------------------

	sqlDrop := "DROP TABLE IF EXISTS test"
	sqlCreate := "CREATE TABLE test (a TEXT)"

	_, err := db.Exec(sqlDrop)
	require.NoError(t, err)

	_, err = db.Exec(sqlCreate)
	require.NoError(t, err)

	// prepare statements ------------------------------------------

	sqlInsert := "INSERT INTO test (a) VALUES (?)"
	sqlUpdate := "UPDATE test SET a = ? WHERE a = ?"
	sqlDelete := "DELETE FROM test WHERE a = ?"
	sqlSelect := "SELECT a FROM test WHERE a = ?"
	sqlList := SQLList("test", "a", "", nil) // Ranges: &crud.Ranges{OrderBy: []string{"a DESC"}}

	var stmInsert, stmUpdate, stmDelete, stmSelect, stmList *sql.Stmt

	sqlStmts := []SqlStmt{
		{&stmInsert, sqlInsert},
		{&stmUpdate, sqlUpdate},
		{&stmDelete, sqlDelete},
		{&stmSelect, sqlSelect},
		{&stmList, sqlList},
	}

	for _, sqlStmt := range sqlStmts {
		err := Prepare(db, sqlStmt.Sql, sqlStmt.Stmt)
		require.NoError(t, err)
	}

	// insert, update ----------------------------------------------

	_, err = stmInsert.Exec("a1")
	require.NoError(t, err)

	_, err = stmInsert.Exec("a2")
	require.NoError(t, err)

	rows1, err := stmSelect.Query("a2")
	require.NoError(t, err)
	require.NotNil(t, rows1)
	defer rows1.Close()

	_, err = stmUpdate.Exec("a3", "a2")
	require.NoError(t, err)

	// count -------------------------------------------------------

	var num int

	sqlCount1 := SQLCount("test", "a = 'a1'", nil)
	row := db.QueryRow(sqlCount1)
	require.NotNil(t, row)

	err = row.Scan(&num)
	require.NoError(t, err)
	require.Equal(t, 1, num)

	sqlCount2 := SQLCount("test", "a = 'a2'", nil)
	row = db.QueryRow(sqlCount2)
	require.NotNil(t, row)

	err = row.Scan(&num)
	require.NoError(t, err)
	require.Equal(t, 0, num)

	sqlCountAll := SQLCount("test", "", nil)
	row = db.QueryRow(sqlCountAll)
	require.NotNil(t, row)

	err = row.Scan(&num)
	require.NoError(t, err)
	require.Equal(t, 2, num)

	// list --------------------------------------------------------

	rowsList, err := stmList.Query()
	require.NoError(t, err)
	require.NotNil(t, rowsList)
	defer rowsList.Close()

	var items []string
	for rowsList.Next() {
		var item string
		err := rowsList.Scan(&item)
		require.NoError(t, err)
		items = append(items, item)
	}
	err = rowsList.Err()
	require.NoError(t, err)
	require.Equal(t, 2, len(items))

	// count, delete, recount --------------------------------------

	sqlCount3 := SQLCount("test", "a = 'a3'", nil)
	row = db.QueryRow(sqlCount3)
	require.NotNil(t, row)

	err = row.Scan(&num)
	require.NoError(t, err)
	require.Equal(t, 1, num)

	_, err = stmDelete.Exec("a3")
	require.NoError(t, err)

	row = db.QueryRow(sqlCount3)
	require.NotNil(t, row)

	err = row.Scan(&num)
	require.NoError(t, err)
	require.Equal(t, 0, num)

}
