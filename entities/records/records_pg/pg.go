package records_pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pavlo67/data/entities"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/sqllib/sqllib_pg"

	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/entities/records"
)

var fields = []string{"title", "summary", "record_type", "data", "embedded"}

var fieldsToInsert = append(fields, entities.Description01FieldsBasis...)
var fieldsToInsertStr = `"` + strings.Join(fieldsToInsert, `","`) + `"`

var fieldsToUpdate = append(fields, entities.Description01FieldsToUpdate...)

var fieldsToRead = append(fields, entities.Description01FieldsToRead...)
var fieldsToReadStr = `"` + strings.Join(fieldsToRead, `","`) + `"`

var fieldsToList = append(fieldsToRead, "id")
var fieldsToListStr = `"` + strings.Join(fieldsToList, `","`) + `"`

var _ records.Operator = &recordsPg{}

type recordsPg struct {
	dbGet *sql.DB // database for data receiving
	dbSet *sql.DB // database for data storing

	domain string
	table  string

	sqlClean, sqlRead, sqlRemove, sqlList, sqlInsert, sqlUpdate, sqlSetURN string
	stmClean, stmRead, stmRemove, stmList, stmInsert, stmUpdate, stmSetURN *sql.Stmt
}

const onNew = "on recordsPg.New()"

func New(dbGet, dbSet *sql.DB, domain, table string) (records.Operator, db.Cleaner, error) {
	if dbGet == nil {
		return nil, nil, errors.New(onNew + ": no dbGet")
	}
	if dbSet == nil {
		dbSet = dbGet
	}

	if domain = strings.TrimSpace(domain); domain == "" {
		return nil, nil, errors.New(onNew + ": no domain defined")
	}
	if table = strings.TrimSpace(table); table == "" {
		return nil, nil, errors.New(onNew + ": no table name defined")
	}

	recordsOp := recordsPg{
		dbGet: dbGet,
		dbSet: dbSet,
		table: table,

		domain: domain,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",
		sqlUpdate: "UPDATE " + table + " SET " + sqllib_pg.WildcardsForUpdate(fieldsToUpdate) + " WHERE id = $" + strconv.Itoa(len(fieldsToUpdate)+1) + " AND history = $" + strconv.Itoa(len(fieldsToUpdate)+2),
		sqlSetURN: "UPDATE " + table + " SET URN = $1 WHERE URN = '' AND id = $2",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1",
		sqlList:   "SELECT " + fieldsToListStr + " FROM " + table + ` ORDER BY id`,
		sqlRemove: "DELETE FROM " + table + " WHERE id = $1",
		sqlClean:  "TRUNCATE " + table,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&recordsOp.stmList, recordsOp.sqlList},
		{&recordsOp.stmRead, recordsOp.sqlRead},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(dbGet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	sqlStmtsSet := []sqllib.SqlStmt{
		{&recordsOp.stmInsert, recordsOp.sqlInsert},
		{&recordsOp.stmUpdate, recordsOp.sqlUpdate},
		{&recordsOp.stmSetURN, recordsOp.sqlSetURN},
		{&recordsOp.stmRemove, recordsOp.sqlRemove},
		{&recordsOp.stmClean, recordsOp.sqlClean},
	}

	for _, sqlStmt := range sqlStmtsSet {
		if err := sqllib.Prepare(dbSet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &recordsOp, &recordsOp, nil

}

// operator ----------------------------------------------------------------------------------------------------------------

var _ records.Operator = &recordsPg{}

const onSetURN = "on recordsPg.setURN()"

func (recordsOp recordsPg) setURN(id records.ID) (ns.URN, error) {

	if strings.TrimSpace(string(id)) == "" {
		return "", fmt.Errorf(onSetURN + ": empty id to set urn")
	}

	urn := ns.CreateURN(recordsOp.domain, string(records.CRUD), string(id))

	values := []interface{}{urn, id}

	if _, err := recordsOp.stmSetURN.Exec(values...); err != nil {
		return "", errors.Wrapf(err, onSetURN+": "+sqllib.CantExec, recordsOp.sqlSetURN, values)
	}

	return urn, nil
}

func (recordsOp recordsPg) Update(ri records.Record, ID records.ID, _ auth.Actor) error { // , vcs.History
}

const onSave = "on recordsPg.Add()"

func (recordsOp recordsPg) Add(ri records.Record, _ auth.Actor) (records.ID, ns.URN, error) { // , vcs.History

	// "title", "summary", "record_type", "data", "embedded"

	var embeddedBytes []byte
	var err error

	if len(ri.Additions) > 0 {
		if embeddedBytes, err = json.Marshal(ri.Additions); err != nil {
			return "", "", errors.Wrapf(err, "can't marshal .Contacts (%#v)", ri.Additions)
		}
	}

	onInsert := ri.ID == ""

	var descriptionValues []interface{}
	var historyOriginalStr string

	descriptionValues, ri.History, historyOriginalStr, err = ri.Description.FoldToSavePg(onInsert)
	if err != nil {
		return "", "", errors.Wrap(err, onSave)
	}

	if onInsert {
		values := append([]interface{}{ri.Title, ri.Summary, ri.Type, ri.Data, embeddedBytes}, descriptionValues...)

		var idInt64 int64
		if err := recordsOp.stmInsert.QueryRow(values...).Scan(&idInt64); err != nil {
			return "", "", errors.Wrapf(err, onSave+": "+sqllib.CantExec, recordsOp.sqlInsert, values)
		}

		ri.ID = (common.IDNum(idInt64)).Key()

		if ri.URN == "" {
			if ri.URN, err = recordsOp.setURN(ri.ID); err != nil {
				return "", "", errors.Wrap(err, onSave)
			}
		}

	} else {

		values := append(
			append([]interface{}{ri.Title, ri.Summary, ri.Type, ri.Data, embeddedBytes}, descriptionValues...),
			ri.ID, historyOriginalStr,
		)
		if res, err := recordsOp.stmUpdate.Exec(values...); err != nil {
			return "", "", errors.Wrapf(err, onSave+": "+sqllib.CantExec, recordsOp.sqlUpdate, values)

		} else {
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return "", "", errors.Wrapf(err, onSave+": "+sqllib.CantGetRowsAffected, recordsOp.sqlUpdate, values)
			} else if rowsAffected < 1 {
				return "", "", fmt.Errorf(onSave+": res.RowsAffected() < 1 on "+sqllib.CantExec, recordsOp.sqlUpdate, values)
			}
		}
	}

	return ri.ID, ri.URN, nil // , ri.History
}

const onRead = "on recordsPg.Read()"

func (recordsOp recordsPg) Read(id records.ID, _ auth.Actor) (*records.Item, error) {
	values := []interface{}{id}
	ri := records.Item{ID: id}

	// "title", "summary", "record_type", "data", "embedded"
	// "urn", "tags", "relations_map", "owner_nss", "viewer_nss", "history"

	var embeddedBytes, urnBytes, relationsMapBytes, historyBytes []byte

	if err := recordsOp.stmRead.QueryRow(values...).Scan(
		&ri.Title, &ri.Summary, &ri.Type, &ri.Data, &embeddedBytes,
		&urnBytes, pq.Array(&ri.Description.Tags), &relationsMapBytes, &ri.Description.OwnerNSS, &ri.Description.ViewerNSS, &historyBytes,
		&ri.Description.UpdatedAt, &ri.Description.CreatedAt); err == sql.ErrNoRows {
		return nil, errors.Wrapf(common.ErrNotFound, onRead+": "+sqllib.CantScanQueryRow, recordsOp.sqlRead, values)
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+": "+sqllib.CantScanQueryRow, recordsOp.sqlRead, values)
	}

	if len(embeddedBytes) > 0 {
		if err := json.Unmarshal(embeddedBytes, &ri.Additions); err != nil {
			return nil, errors.Wrapf(err, onRead+": can't unmarshal .Additions (%s)", embeddedBytes)
		}
	}

	if err := ri.Description.UnfoldReaded(urnBytes, relationsMapBytes, historyBytes); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return &ri, nil
}

const onList = "on recordsPg.List()"

func (recordsOp recordsPg) List(*entities.Term, auth.Actor) ([]records.Item, error) {

	// TODO!!! selector

	var values []interface{}
	rows, err := recordsOp.stmList.Query(values...)

	var items []records.Item

	if err == sql.ErrNoRows {
		return items, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, recordsOp.sqlList, values)
	}
	defer rows.Close()

	for rows.Next() {
		var idInt64 int64
		var ri records.Item
		var embeddedBytes, urnBytes, relationsMapBytes, historyBytes []byte

		if err := rows.Scan(
			&ri.Title, &ri.Summary, &ri.Type, &ri.Data, &embeddedBytes,
			&urnBytes, pq.Array(&ri.Description.Tags), &relationsMapBytes, &ri.Description.OwnerNSS, &ri.Description.ViewerNSS, &historyBytes,
			&ri.Description.UpdatedAt, &ri.Description.CreatedAt, &idInt64); err != nil {
			return nil, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, recordsOp.sqlList, values)
		}

		if len(embeddedBytes) > 0 {
			if err := json.Unmarshal(embeddedBytes, &ri.Additions); err != nil {
				return nil, errors.Wrapf(err, onRead+": can't unmarshal .Additions (%s)", embeddedBytes)
			}
		}

		if err := ri.Description.UnfoldReaded(urnBytes, relationsMapBytes, historyBytes); err != nil {
			return nil, errors.Wrap(err, onList)
		}

		ri.ID = (common.IDNum(idInt64)).Key()

		items = append(items, ri)
	}

	if err = rows.Err(); err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, recordsOp.sqlList, values)
	}

	return items, nil
}

const onRemove = "on recordsPg.Remove()"

func (recordsOp recordsPg) Remove(id records.ID, _ auth.Actor) error {
	values := []interface{}{id}

	if _, err := recordsOp.stmRemove.Exec(values...); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, recordsOp.sqlRemove, values)
	}

	return nil
}

const onClose = "on recordsPg.Close()"

func (recordsOp recordsPg) Close() (err error) {
	if err = recordsOp.dbGet.Close(); err != nil {
		return errors.Wrap(err, onClose+": can't .dbGet.Close()")
	}

	if recordsOp.dbSet != recordsOp.dbGet {
		if err = recordsOp.dbSet.Close(); err != nil {
			return errors.Wrap(err, onClose+": can't .dbSet.Close()")
		}
	}

	return nil
}

// cleaner -----------------------------------------------------------------------------------------------------------------

var _ db.Cleaner = &recordsPg{}

const onClean = "on recordsPg.Clean()"

func (recordsOp recordsPg) Clean() error {
	if env := os.Getenv("ENV"); env != "test" {
		return fmt.Errorf("wrong ENV environment value (%s), must be 'test'", env)
	}

	if _, err := recordsOp.stmClean.Exec(); err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, recordsOp.sqlClean, nil)
	}

	return nil
}
