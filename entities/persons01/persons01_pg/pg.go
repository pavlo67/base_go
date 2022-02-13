package persons01_pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/sqllib/sqllib_pg"

	"github.com/pavlo67/data/entities/crud01"
	"github.com/pavlo67/data/entities/persons01"

	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/selectors"
	"github.com/pavlo67/data/components/vcs"
)

var fields = []string{"firstnames", "middlename", "lastname", "nicknames", "contacts", "info"}

var fieldsToInsert = append(fields, crud01.Description01FieldsToInsert...)
var fieldsToInsertStr = `"` + strings.Join(fieldsToInsert, `","`) + `"`

var fieldsToUpdate = append(fields, crud01.Description01FieldsToUpdate...)

var fieldsToRead = append(fields, crud01.Description01FieldsToRead...)
var fieldsToReadStr = `"` + strings.Join(fieldsToRead, `","`) + `"`

var fieldsToList = append(fieldsToRead, "id")
var fieldsToListStr = `"` + strings.Join(fieldsToList, `","`) + `"`

var _ persons01.Operator = &persons01Pg{}

type persons01Pg struct {
	dbGet *sql.DB // database for data receiving
	dbSet *sql.DB // database for data storing

	table string

	sqlClean, sqlRead, sqlRemove, sqlList, sqlInsert, sqlUpdate string
	stmClean, stmRead, stmRemove, stmList, stmInsert, stmUpdate *sql.Stmt
}

const onNew = "on persons01Pg.New()"

func New(dbGet, dbSet *sql.DB, table string) (persons01.Operator, db.Cleaner, error) {
	if dbGet == nil {
		return nil, nil, errors.New(onNew + ": no dbGet")
	}
	if dbSet == nil {
		dbSet = dbGet
	}

	if table = strings.TrimSpace(table); table == "" {
		return nil, nil, errors.New(onNew + ": no table name defined")
	}

	persons01Op := persons01Pg{
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
		{&persons01Op.stmList, persons01Op.sqlList},
		{&persons01Op.stmRead, persons01Op.sqlRead},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(dbGet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	sqlStmtsSet := []sqllib.SqlStmt{
		{&persons01Op.stmInsert, persons01Op.sqlInsert},
		{&persons01Op.stmUpdate, persons01Op.sqlUpdate},
		{&persons01Op.stmRemove, persons01Op.sqlRemove},
		{&persons01Op.stmClean, persons01Op.sqlClean},
	}

	for _, sqlStmt := range sqlStmtsSet {
		if err := sqllib.Prepare(dbSet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &persons01Op, &persons01Op, nil

}

// operator ----------------------------------------------------------------------------------------------------------------

var _ persons01.Operator = &persons01Pg{}

const onSave = "on persons01Pg.Save()"

func (persons01Op persons01Pg) Save(pi persons01.Item, _ auth.Actor) (persons01.ID, vcs.History, error) {

	// "firstnames", "middlename", "lastname", "nicknames", "contacts", "info"

	var contactsBytes, infoBytes []byte
	var err error

	if len(pi.Contacts) > 0 {
		if contactsBytes, err = json.Marshal(pi.Contacts); err != nil {
			return "", nil, errors.Wrapf(err, onSave+": can't marshal .Contacts (%#v)", pi.Contacts)
		}
	}
	if len(pi.Info) > 0 {
		if infoBytes, err = json.Marshal(pi.Info); err != nil {
			return "", nil, errors.Wrapf(err, "can't marshal .Info (%#v)", pi.Info)
		}
	}

	onInsert := pi.ID == ""

	descriptionValues, historyChanged, historyOriginalStr, err := pi.Description.FoldToSavePg(onInsert)
	if err != nil {
		return "", nil, errors.Wrap(err, onSave)
	}

	values := append(
		[]interface{}{pq.Array(pi.Firstnames), pi.Middlename, pi.Lastname, pq.Array(pi.Nicknames), contactsBytes, infoBytes},
		descriptionValues...)

	if onInsert {

		var idInt64 int64

		if err := persons01Op.stmInsert.QueryRow(values...).Scan(&idInt64); err != nil {
			return "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantExec, persons01Op.sqlInsert, values)
		}

		pi.ID = crud.NewIDInt64(idInt64)

	} else {
		values = append(values, pi.ID, historyOriginalStr)
		if res, err := persons01Op.stmUpdate.Exec(values...); err != nil {
			return "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantExec, persons01Op.sqlUpdate, values)
		} else {
			rowsAffected, err := res.RowsAffected()

			if err != nil {
				return "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantGetRowsAffected, persons01Op.sqlUpdate, values)
			} else if rowsAffected < 1 {
				return "", nil, fmt.Errorf(onSave+": res.RowsAffected() < 1 on "+sqllib.CantExec, persons01Op.sqlUpdate, values)
			}
		}
	}

	return pi.ID, historyChanged, nil
}

const onRead = "on persons01Pg.Read()"

func (persons01Op persons01Pg) Read(id persons01.ID, _ auth.Actor) (*persons01.Item, error) {

	values := []interface{}{id}
	pi := persons01.Item{ID: id}

	// "firstnames", "middlename", "lastname", "nicknames", "contacts", "info"
	// "urn", "tags", "relations_map", "owner_nss", "viewer_nss", "history"

	var urnBytes, contactBytes, infoBytes, relationsMapBytes, historyBytes []byte

	if err := persons01Op.stmRead.QueryRow(values...).Scan(
		pq.Array(&pi.Firstnames), &pi.Middlename, &pi.Lastname, pq.Array(&pi.Nicknames), &contactBytes, &infoBytes,
		&urnBytes, pq.Array(&pi.Description.Tags), &relationsMapBytes, &pi.Description.OwnerNSS, &pi.Description.ViewerNSS, &historyBytes,
		&pi.Description.UpdatedAt, &pi.Description.CreatedAt); err == sql.ErrNoRows {
		return nil, errors.Wrapf(common.ErrNotFound, onRead+": "+sqllib.CantScanQueryRow, persons01Op.sqlRead, values)
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+": "+sqllib.CantScanQueryRow, persons01Op.sqlRead, values)
	}

	if len(contactBytes) > 0 {
		if err := json.Unmarshal(contactBytes, &pi.Contacts); err != nil {
			return nil, errors.Wrapf(err, onRead+": can't unmarshal .Contacts (%s)", contactBytes)
		}
	}
	if len(infoBytes) > 0 {
		if err := json.Unmarshal(infoBytes, &pi.Info); err != nil {
			return nil, errors.Wrapf(err, onRead+": can't unmarshal .Info (%s)", infoBytes)
		}
	}

	if err := pi.Description.UnfoldReaded(urnBytes, relationsMapBytes, historyBytes); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return &pi, nil
}

const onList = "on persons01Pg.List()"

func (persons01Op persons01Pg) List(*selectors.Term, auth.Actor) ([]persons01.Item, error) {

	// TODO!!! selector

	var values []interface{}
	rows, err := persons01Op.stmList.Query(values...)

	var items []persons01.Item

	if err == sql.ErrNoRows {
		return items, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, persons01Op.sqlList, values)
	}
	defer rows.Close()

	for rows.Next() {
		var idInt64 int64
		var pi persons01.Item
		var urnBytes, contactBytes, infoBytes, relationsMapBytes, historyBytes []byte

		if err := rows.Scan(pq.Array(&pi.Firstnames), &pi.Middlename, &pi.Lastname, pq.Array(&pi.Nicknames), &contactBytes, &infoBytes,
			&urnBytes, pq.Array(&pi.Description.Tags), &relationsMapBytes, &pi.Description.OwnerNSS, &pi.Description.ViewerNSS, &historyBytes,
			&pi.Description.UpdatedAt, &pi.Description.CreatedAt, &idInt64); err != nil {
			return nil, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, persons01Op.sqlList, values)
		}

		if len(contactBytes) > 0 {
			if err := json.Unmarshal(contactBytes, &pi.Contacts); err != nil {
				return nil, errors.Wrapf(err, onList+": can't unmarshal .Contacts (%s)", contactBytes)
			}
		}
		if len(infoBytes) > 0 {
			if err := json.Unmarshal(infoBytes, &pi.Info); err != nil {
				return nil, errors.Wrapf(err, onList+": can't unmarshal .Info (%s)", infoBytes)
			}
		}

		if err := pi.Description.UnfoldReaded(urnBytes, relationsMapBytes, historyBytes); err != nil {
			return nil, errors.Wrap(err, onList)
		}

		pi.ID = crud.NewIDInt64(idInt64)

		items = append(items, pi)
	}

	if err = rows.Err(); err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, persons01Op.sqlList, values)
	}

	return items, nil
}

const onRemove = "on persons01Pg.Remove()"

func (persons01Op persons01Pg) Remove(id persons01.ID, _ auth.Actor) error {
	values := []interface{}{id}

	if _, err := persons01Op.stmRemove.Exec(values...); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, persons01Op.sqlRemove, values)
	}

	return nil
}

const onClose = "on persons01Pg.Close()"

func (persons01Op persons01Pg) Close() (err error) {
	if err = persons01Op.dbGet.Close(); err != nil {
		return errors.Wrap(err, onClose+": can't .dbGet.Close()")
	}

	if persons01Op.dbSet != persons01Op.dbGet {
		if err = persons01Op.dbSet.Close(); err != nil {
			return errors.Wrap(err, onClose+": can't .dbSet.Close()")
		}
	}

	return nil
}

// cleaner -----------------------------------------------------------------------------------------------------------------

var _ db.Cleaner = &persons01Pg{}

const onClean = "on persons01Pg.Clean()"

func (persons01Op persons01Pg) Clean() error {
	if env := os.Getenv("ENV"); env != "test" {
		return fmt.Errorf("wrong ENV environment value (%s), must be 'test'", env)
	}

	if _, err := persons01Op.stmClean.Exec(); err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, persons01Op.sqlClean, nil)
	}

	return nil
}
