package records01_pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	crud012 "github.com/pavlo67/data/entities/crud01"

	"github.com/pavlo67/data/components/vcs"

	"github.com/pavlo67/data/components/selectors"

	"github.com/lib/pq"
	"github.com/pavlo67/data/components/crud"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/sqllib/sqllib_pg"

	"github.com/pavlo67/data/entities/records01"
)

var fields = []string{"title", "summary", "record_type", "data", "embedded"}

var fieldsToInsert = append(fields, crud012.Description01FieldsToInsert...)
var fieldsToInsertStr = `"` + strings.Join(fieldsToInsert, `","`) + `"`

var fieldsToUpdate = append(fields, crud012.Description01FieldsToUpdate...)

var fieldsToRead = append(fields, crud012.Description01FieldsToRead...)
var fieldsToReadStr = `"` + strings.Join(fieldsToRead, `","`) + `"`

var fieldsToList = append(fieldsToRead, "id")
var fieldsToListStr = `"` + strings.Join(fieldsToList, `","`) + `"`

var _ records01.Operator = &records01Pg{}

type records01Pg struct {
	dbGet *sql.DB // database for data receiving
	dbSet *sql.DB // database for data storing

	table string

	sqlClean, sqlRead, sqlRemove, sqlList, sqlInsert, sqlUpdate string
	stmClean, stmRead, stmRemove, stmList, stmInsert, stmUpdate *sql.Stmt
}

const onNew = "on records01Pg.New()"

func New(dbGet, dbSet *sql.DB, table string) (records01.Operator, db.Cleaner, error) {
	if dbGet == nil {
		return nil, nil, errors.New(onNew + ": no dbGet")
	}
	if dbSet == nil {
		dbSet = dbGet
	}

	if table = strings.TrimSpace(table); table == "" {
		return nil, nil, errors.New(onNew + ": no table name defined")
	}

	records01Op := records01Pg{
		dbGet: dbGet,
		dbSet: dbSet,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",
		sqlUpdate: "UPDATE " + table + " SET " + sqllib_pg.WildcardsForUpdate(fieldsToUpdate) + " WHERE id = $" + strconv.Itoa(len(fieldsToUpdate)+1) + " AND history = $" + strconv.Itoa(len(fieldsToUpdate)+2),
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1",
		sqlList:   "SELECT " + fieldsToListStr + " FROM " + table + ` ORDER BY id`,
		sqlRemove: "DELETE FROM " + table + " WHERE id = $1",
		sqlClean:  "TRUNCATE " + table,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&records01Op.stmList, records01Op.sqlList},
		{&records01Op.stmRead, records01Op.sqlRead},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(dbGet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	sqlStmtsSet := []sqllib.SqlStmt{
		{&records01Op.stmInsert, records01Op.sqlInsert},
		{&records01Op.stmUpdate, records01Op.sqlUpdate},
		{&records01Op.stmRemove, records01Op.sqlRemove},
		{&records01Op.stmClean, records01Op.sqlClean},
	}

	for _, sqlStmt := range sqlStmtsSet {
		if err := sqllib.Prepare(dbSet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &records01Op, &records01Op, nil

}

// operator ----------------------------------------------------------------------------------------------------------------

var _ records01.Operator = &records01Pg{}

const onSave = "on records01Pg.Save()"

func (records01Op records01Pg) Save(ri records01.Item, _ auth.Actor) (records01.ID, vcs.History, error) {

	// "title", "summary", "record_type", "data", "embedded"

	var embeddedBytes []byte
	var err error

	if len(ri.Embedded) > 0 {
		if embeddedBytes, err = json.Marshal(ri.Embedded); err != nil {
			return "", nil, errors.Wrapf(err, onSave+": can't marshal .Contacts (%#v)", ri.Embedded)
		}
	}

	onInsert := ri.ID == ""
	descriptionValues, historyChanged, historyOriginalStr, err := ri.Description.FoldToSavePg(onInsert)
	if err != nil {
		return "", nil, errors.Wrap(err, onSave)
	}

	values := append([]interface{}{ri.Title, ri.Summary, ri.Type, ri.Data, embeddedBytes}, descriptionValues...)

	if onInsert {
		var idInt64 int64

		if err := records01Op.stmInsert.QueryRow(values...).Scan(&idInt64); err != nil {
			return "", historyChanged, errors.Wrapf(err, onSave+": "+sqllib.CantExec, records01Op.sqlInsert, values)
		}

		ri.ID = crud.NewIDInt64(idInt64)

	} else {
		values = append(values, ri.ID, historyOriginalStr)
		res, err := records01Op.stmUpdate.Exec(values...)
		if err != nil {
			return "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantExec, records01Op.sqlUpdate, values)
		} else {
			rowsAffected, err := res.RowsAffected()

			if err != nil {
				return "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantGetRowsAffected, records01Op.sqlUpdate, values)
			} else if rowsAffected < 1 {
				return "", nil, fmt.Errorf(onSave+": res.RowsAffected() < 1 on "+sqllib.CantExec, records01Op.sqlUpdate, values)
			}
		}
	}

	return ri.ID, historyChanged, nil
}

const onRead = "on records01Pg.Read()"

func (records01Op records01Pg) Read(id records01.ID, _ auth.Actor) (*records01.Item, error) {

	values := []interface{}{id}
	ri := records01.Item{ID: id}

	// "title", "summary", "record_type", "data", "embedded"
	// "urn", "tags", "relations_map", "owner_nss", "viewer_nss", "history"

	var embeddedBytes, urnBytes, relationsMapBytes, historyBytes []byte

	if err := records01Op.stmRead.QueryRow(values...).Scan(
		&ri.Title, &ri.Summary, &ri.Type, &ri.Data, &embeddedBytes,
		&urnBytes, pq.Array(&ri.Description.Tags), &relationsMapBytes, &ri.Description.OwnerNSS, &ri.Description.ViewerNSS, &historyBytes,
		&ri.Description.UpdatedAt, &ri.Description.CreatedAt); err == sql.ErrNoRows {
		return nil, errors.Wrapf(common.ErrNotFound, onRead+": "+sqllib.CantScanQueryRow, records01Op.sqlRead, values)
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+": "+sqllib.CantScanQueryRow, records01Op.sqlRead, values)
	}

	if len(embeddedBytes) > 0 {
		if err := json.Unmarshal(embeddedBytes, &ri.Embedded); err != nil {
			return nil, errors.Wrapf(err, onRead+": can't unmarshal .Embedded (%s)", embeddedBytes)
		}
	}

	if err := ri.Description.UnfoldReaded(urnBytes, relationsMapBytes, historyBytes); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return &ri, nil
}

const onList = "on records01Pg.List()"

func (records01Op records01Pg) List(*selectors.Term, auth.Actor) ([]records01.Item, error) {

	// TODO!!! selector

	var values []interface{}
	rows, err := records01Op.stmList.Query(values...)

	var items []records01.Item

	if err == sql.ErrNoRows {
		return items, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, records01Op.sqlList, values)
	}
	defer rows.Close()

	for rows.Next() {
		var idInt64 int64
		var ri records01.Item
		var embeddedBytes, urnBytes, relationsMapBytes, historyBytes []byte

		if err := rows.Scan(
			&ri.Title, &ri.Summary, &ri.Type, &ri.Data, &embeddedBytes,
			&urnBytes, pq.Array(&ri.Description.Tags), &relationsMapBytes, &ri.Description.OwnerNSS, &ri.Description.ViewerNSS, &historyBytes,
			&ri.Description.UpdatedAt, &ri.Description.CreatedAt, &idInt64); err != nil {
			return nil, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, records01Op.sqlList, values)
		}

		if len(embeddedBytes) > 0 {
			if err := json.Unmarshal(embeddedBytes, &ri.Embedded); err != nil {
				return nil, errors.Wrapf(err, onRead+": can't unmarshal .Embedded (%s)", embeddedBytes)
			}
		}

		if err := ri.Description.UnfoldReaded(urnBytes, relationsMapBytes, historyBytes); err != nil {
			return nil, errors.Wrap(err, onList)
		}

		ri.ID = crud.NewIDInt64(idInt64)

		items = append(items, ri)
	}

	if err = rows.Err(); err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, records01Op.sqlList, values)
	}

	return items, nil
}

const onRemove = "on records01Pg.Remove()"

func (records01Op records01Pg) Remove(id records01.ID, _ auth.Actor) error {
	values := []interface{}{id}

	if _, err := records01Op.stmRemove.Exec(values...); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, records01Op.sqlRemove, values)
	}

	return nil
}

const onClose = "on records01Pg.Close()"

func (records01Op records01Pg) Close() (err error) {
	if err = records01Op.dbGet.Close(); err != nil {
		return errors.Wrap(err, onClose+": can't .dbGet.Close()")
	}

	if records01Op.dbSet != records01Op.dbGet {
		if err = records01Op.dbSet.Close(); err != nil {
			return errors.Wrap(err, onClose+": can't .dbSet.Close()")
		}
	}

	return nil
}

// cleaner -----------------------------------------------------------------------------------------------------------------

var _ db.Cleaner = &records01Pg{}

const onClean = "on records01Pg.Clean()"

func (records01Op records01Pg) Clean() error {
	if env := os.Getenv("ENV"); env != "test" {
		return fmt.Errorf("wrong ENV environment value (%s), must be 'test'", env)
	}

	if _, err := records01Op.stmClean.Exec(); err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, records01Op.sqlClean, nil)
	}

	return nil
}
