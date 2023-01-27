package persons_pg

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

	"github.com/pavlo67/data/entities/persons"
	"github.com/pavlo67/data/entities/records"

	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/components/vcs"
)

var fields = []string{"firstnames", "middlename", "lastname", "nicknames", "contacts", "info"}

var fieldsToInsert = append(fields, crud.Description01FieldsBasis...)
var fieldsToInsertStr = `"` + strings.Join(fieldsToInsert, `","`) + `"`

var fieldsToUpdate = append(fields, crud.Description01FieldsToUpdate...)

var fieldsToRead = append(fields, crud.Description01FieldsToRead...)
var fieldsToReadStr = `"` + strings.Join(fieldsToRead, `","`) + `"`

var fieldsToList = append(fieldsToRead, "id")
var fieldsToListStr = `"` + strings.Join(fieldsToList, `","`) + `"`

var _ persons.Operator = &personsPg{}

type personsPg struct {
	dbGet *sql.DB // database for data receiving
	dbSet *sql.DB // database for data storing

	domain string
	table  string

	sqlClean, sqlRead, sqlRemove, sqlList, sqlInsert, sqlUpdate, sqlSetURN string
	stmClean, stmRead, stmRemove, stmList, stmInsert, stmUpdate, stmSetURN *sql.Stmt
}

const onNew = "on personsPg.New()"

func New(dbGet, dbSet *sql.DB, domain, table string) (persons.Operator, db.Cleaner, error) {
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

	personsOp := personsPg{
		dbGet: dbGet,
		dbSet: dbSet,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",
		sqlUpdate: "UPDATE " + table + " SET " + sqllib_pg.WildcardsForUpdate(fieldsToUpdate) + " WHERE id = $" + strconv.Itoa(len(fieldsToUpdate)+1) + " AND history = $" + strconv.Itoa(len(fieldsToUpdate)+2),
		sqlSetURN: "UPDATE " + table + " SET URN = $1 WHERE URN = '' AND id = $2",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1",
		sqlList:   "SELECT " + fieldsToListStr + " FROM " + table + ` ORDER BY id`,
		sqlRemove: "DELETE FROM " + table + " WHERE id = $1",
		sqlClean:  "TRUNCATE " + table,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&personsOp.stmList, personsOp.sqlList},
		{&personsOp.stmRead, personsOp.sqlRead},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(dbGet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	sqlStmtsSet := []sqllib.SqlStmt{
		{&personsOp.stmInsert, personsOp.sqlInsert},
		{&personsOp.stmUpdate, personsOp.sqlUpdate},
		{&personsOp.stmSetURN, personsOp.sqlSetURN},
		{&personsOp.stmRemove, personsOp.sqlRemove},
		{&personsOp.stmClean, personsOp.sqlClean},
	}

	for _, sqlStmt := range sqlStmtsSet {
		if err := sqllib.Prepare(dbSet, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &personsOp, &personsOp, nil
}

// operator ----------------------------------------------------------------------------------------------------------------

var _ persons.Operator = &personsPg{}

const onSetURN = "on personsPg.SetURN()"

func (personsOp personsPg) SetURN(id records.ID) (ns.URN, error) {

	if strings.TrimSpace(string(id)) == "" {
		return "", fmt.Errorf(onSetURN + ": empty id to set urn")
	}

	urn := ns.CreateURN(personsOp.domain, string(persons.CRUD), string(id))

	values := []interface{}{urn, id}

	if _, err := personsOp.stmSetURN.Exec(values...); err != nil {
		return "", errors.Wrapf(err, onSetURN+": "+sqllib.CantExec, personsOp.sqlSetURN, values)
	}

	return urn, nil
}

const onSave = "on personsPg.Save()"

func (personsOp personsPg) Save(pi persons.Item, _ auth.Actor) (persons.ID, ns.URN, vcs.History, error) {

	// "firstnames", "middlename", "lastname", "nicknames", "contacts", "info"

	var contactsBytes, infoBytes []byte
	var err error

	if len(pi.Contacts) > 0 {
		if contactsBytes, err = json.Marshal(pi.Contacts); err != nil {
			return "", "", nil, errors.Wrapf(err, "can't marshal .Contacts (%#v)", pi.Contacts)
		}
	}
	if len(pi.Info) > 0 {
		if infoBytes, err = json.Marshal(pi.Info); err != nil {
			return "", "", nil, errors.Wrapf(err, "can't marshal .Info (%#v)", pi.Info)
		}
	}

	onInsert := pi.ID == ""

	var descriptionValues []interface{}
	var historyOriginalStr string

	descriptionValues, pi.History, historyOriginalStr, err = pi.Description.FoldToSavePg(onInsert)
	if err != nil {
		return "", "", nil, errors.Wrap(err, onSave)
	}

	values := append(
		[]interface{}{pq.Array(pi.Firstnames), pi.Middlename, pi.Lastname, pq.Array(pi.Nicknames), contactsBytes, infoBytes},
		descriptionValues...)

	if onInsert {
		var idInt64 int64

		if err := personsOp.stmInsert.QueryRow(values...).Scan(&idInt64); err != nil {
			return "", "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantExec, personsOp.sqlInsert, values)
		}

		pi.ID = (common.IDNum(idInt64)).Key()

		if pi.URN == "" {
			if pi.URN, err = personsOp.SetURN(pi.ID); err != nil {
				return "", "", nil, errors.Wrap(err, onSave)
			}
		}

	} else {

		values = append(values, pi.ID, historyOriginalStr)
		if res, err := personsOp.stmUpdate.Exec(values...); err != nil {
			return "", "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantExec, personsOp.sqlUpdate, values)
		} else {
			rowsAffected, err := res.RowsAffected()
			if err != nil {
				return "", "", nil, errors.Wrapf(err, onSave+": "+sqllib.CantGetRowsAffected, personsOp.sqlUpdate, values)
			} else if rowsAffected < 1 {
				return "", "", nil, fmt.Errorf(onSave+": res.RowsAffected() < 1 on "+sqllib.CantExec, personsOp.sqlUpdate, values)
			}
		}
	}

	return pi.ID, pi.URN, pi.History, nil
}

const onRead = "on personsPg.Read()"

func (personsOp personsPg) Read(id persons.ID, _ auth.Actor) (*persons.Item, error) {

	values := []interface{}{id}
	pi := persons.Item{ID: id}

	// "firstnames", "middlename", "lastname", "nicknames", "contacts", "info"
	// "urn", "tags", "relations_map", "owner_nss", "viewer_nss", "history"

	var urnBytes, contactBytes, infoBytes, relationsMapBytes, historyBytes []byte

	if err := personsOp.stmRead.QueryRow(values...).Scan(
		pq.Array(&pi.Firstnames), &pi.Middlename, &pi.Lastname, pq.Array(&pi.Nicknames), &contactBytes, &infoBytes,
		&urnBytes, pq.Array(&pi.Description.Tags), &relationsMapBytes, &pi.Description.OwnerNSS, &pi.Description.ViewerNSS, &historyBytes,
		&pi.Description.UpdatedAt, &pi.Description.CreatedAt); err == sql.ErrNoRows {
		return nil, errors.Wrapf(common.ErrNotFound, onRead+": "+sqllib.CantScanQueryRow, personsOp.sqlRead, values)
	} else if err != nil {
		return nil, errors.Wrapf(err, onRead+": "+sqllib.CantScanQueryRow, personsOp.sqlRead, values)
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

const onList = "on personsPg.List()"

func (personsOp personsPg) List(*crud.Term, auth.Actor) ([]persons.Item, error) {

	// TODO!!! selector

	var values []interface{}
	rows, err := personsOp.stmList.Query(values...)

	var items []persons.Item

	if err == sql.ErrNoRows {
		return items, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, personsOp.sqlList, values)
	}
	defer rows.Close()

	for rows.Next() {
		var idInt64 int64
		var pi persons.Item
		var urnBytes, contactBytes, infoBytes, relationsMapBytes, historyBytes []byte

		if err := rows.Scan(pq.Array(&pi.Firstnames), &pi.Middlename, &pi.Lastname, pq.Array(&pi.Nicknames), &contactBytes, &infoBytes,
			&urnBytes, pq.Array(&pi.Description.Tags), &relationsMapBytes, &pi.Description.OwnerNSS, &pi.Description.ViewerNSS, &historyBytes,
			&pi.Description.UpdatedAt, &pi.Description.CreatedAt, &idInt64); err != nil {
			return nil, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, personsOp.sqlList, values)
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

		pi.ID = (common.IDNum(idInt64)).Key()

		items = append(items, pi)
	}

	if err = rows.Err(); err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, personsOp.sqlList, values)
	}

	return items, nil
}

const onRemove = "on personsPg.Remove()"

func (personsOp personsPg) Remove(id persons.ID, _ auth.Actor) error {
	values := []interface{}{id}

	if _, err := personsOp.stmRemove.Exec(values...); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, personsOp.sqlRemove, values)
	}

	return nil
}

const onClose = "on personsPg.Close()"

func (personsOp personsPg) Close() (err error) {
	if err = personsOp.dbGet.Close(); err != nil {
		return errors.Wrap(err, onClose+": can't .dbGet.Close()")
	}

	if personsOp.dbSet != personsOp.dbGet {
		if err = personsOp.dbSet.Close(); err != nil {
			return errors.Wrap(err, onClose+": can't .dbSet.Close()")
		}
	}

	return nil
}

// cleaner -----------------------------------------------------------------------------------------------------------------

var _ db.Cleaner = &personsPg{}

const onClean = "on personsPg.Clean()"

func (personsOp personsPg) Clean() error {
	if env := os.Getenv("ENV"); env != "test" {
		return fmt.Errorf("wrong ENV environment value (%s), must be 'test'", env)
	}

	if _, err := personsOp.stmClean.Exec(); err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, personsOp.sqlClean, nil)
	}

	return nil
}
