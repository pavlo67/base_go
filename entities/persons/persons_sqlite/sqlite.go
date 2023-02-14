package persons_sqlite

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/data/components/ns"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/sqllib"
	"github.com/pavlo67/common/common/strlib"
	"github.com/pavlo67/data/components/selectors"

	"github.com/pavlo67/data/entities/persons"
)

var fields = []string{"nickname", "email", "roles", "info", "creds", "tags", "owner_nss", "viewer_nss", "history"}

var fieldsToInsert = append(fields, "urn")
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToUpdate = append(fields, "updated_at")
var fieldsToUpdateStr = strings.Join(fieldsToUpdate, " = ?, ") + " = ?"

var fieldsToRead = append(fields, "urn", "updated_at", "created_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append(fieldsToRead, "id")
var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ persons.Operator = &personsSQLite{}

type personsSQLite struct {
	db    *sql.DB
	table string

	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlStat, sqlClean string
	stmInsert, stmUpdate, stmRead, stmRemove, stmStat, stmClean *sql.Stmt
}

const onNew = "on personsSQLite.New(): "

func New(db *sql.DB, table string) (persons.Operator, db.Cleaner, error) {
	if table == "" {
		table = persons.CollectionDefault
	}

	personsOp := personsSQLite{
		db:    db,
		table: table,

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + strings.Repeat(",? ", len(fieldsToInsert))[1:] + ")",
		sqlUpdate: "UPDATE " + table + " SET " + fieldsToUpdateStr + " WHERE id = ?",
		sqlRemove: "DELETE FROM " + table + " where id = ?",
		sqlRead:   "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = ?",
		sqlStat:   "SELECT COUNT(*) FROM " + table,

		sqlClean: "DELETE FROM " + table,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&personsOp.stmInsert, personsOp.sqlInsert},
		{&personsOp.stmUpdate, personsOp.sqlUpdate},
		{&personsOp.stmRead, personsOp.sqlRead},
		{&personsOp.stmRemove, personsOp.sqlRemove},
		{&personsOp.stmStat, personsOp.sqlStat},
		{&personsOp.stmClean, personsOp.sqlClean},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &personsOp, &personsOp, nil
}

const onSave = "on personsSQLite.Add(): "

func (personsOp *personsSQLite) Save(item persons.Item, identity *auth.Identity) (auth.ID, error) {
	if identity == nil || (item.ID != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return "", errors.CommonError(common.NoRightsKey, common.Map{"on": onSave, "item": item})
	}

	// TODO!!! append to item.History
	rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, historyBytes, urnBytes, err := item.FoldIntoJSON()
	if err != nil {
		return "", errors.CommonError(err, onSave)
	}

	//l.Infof("URN BYTES: %#v", urnBytes)

	// "nickname", "email", "roles", "creds", "info",
	// "urn", "tags", "owner_nss", "viewer_nss", "history"

	if item.ID != "" {
		itemOld, err := personsOp.read(item.Identity.ID)
		if err != nil || itemOld == nil {
			errorStr := fmt.Sprintf("got %#v / %s", itemOld, err)
			if identity.HasRole(rbac.RoleAdmin) {
				return "", errors.CommonError(common.WrongIDKey, common.Map{"on": onSave, "item": item, "reason": errorStr})
			} else {
				l.Error(errorStr)
				return "", errors.CommonError(common.NoRightsKey, common.Map{"on": onSave, "item": item, "requestedRole": rbac.RoleAdmin})
			}
		}

		// ... "updated_at", "id"
		values := []interface{}{
			item.Nickname, emailBytes, rolesBytes, credsBytes, infoBytes,
			tagsBytes, item.OwnerNSS, item.ViewerNSS, historyBytes,
			time.Now(), item.ID,
		}

		if _, err = personsOp.stmUpdate.Exec(values...); err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantExec, personsOp.sqlUpdate, strlib.Stringify(values))
		}

	} else {
		values := []interface{}{
			item.Nickname, emailBytes, rolesBytes, credsBytes, infoBytes,
			tagsBytes, item.OwnerNSS, item.ViewerNSS, historyBytes,
			urnBytes,
		}

		res, err := personsOp.stmInsert.Exec(values...)
		if err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantExec, personsOp.sqlInsert, strlib.Stringify(values))
		}

		idSQLite, err := res.LastInsertId()
		if err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantGetLastInsertId, personsOp.sqlInsert, strlib.Stringify(values))
		}

		item.ID = auth.ID(strconv.FormatInt(idSQLite, 10))
	}

	return item.ID, nil
}

const onRemove = "on personsSQLite.Remove()"

func (personsOp *personsSQLite) Remove(id auth.ID, identity *auth.Identity) error {
	if identity == nil || (id != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return errors.CommonError(common.NoRightsKey, common.Map{"on": onRemove, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return fmt.Errorf(onRemove+"wrong id (%s)", id)
	}

	if _, err = personsOp.stmRemove.Exec(idNum); err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, personsOp.sqlRemove, idNum)
	}

	return nil
}

const onread = "on personsSQLite.read(): "

func (personsOp *personsSQLite) read(id auth.ID) (*persons.Item, error) {

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, fmt.Errorf(onRead+"wrong id (%s)", id)
	}

	var item persons.Item
	var emailBytes, rolesBytes, credsBytes, infoBytes, tagsBytes, urnBytes, historyBytes []byte

	// "nickname", "email", "roles", "creds", "info",
	// "tags", "owner_nss", "viewer_nss", "history"
	// "urn",
	// "updated_at", "created_at", "id"

	if err = personsOp.stmRead.QueryRow(idNum).Scan(
		&item.Nickname, &emailBytes, &rolesBytes, &credsBytes, &infoBytes,
		&tagsBytes, &item.OwnerNSS, &item.ViewerNSS, &historyBytes,
		&urnBytes,
		&item.UpdatedAt, &item.CreatedAt,
	); err == sql.ErrNoRows {
		return nil, errors.CommonError(common.ErrNotFound, onread)
	} else if err != nil {
		return nil, errors.Wrapf(err, onread+sqllib.CantScanQueryRow, personsOp.sqlRead, idNum)
	}

	if err := item.UnfoldFromJSON(id, rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, urnBytes, historyBytes); err != nil {
		return nil, errors.CommonError(err, onread)
	}

	return &item, nil
}

const onRead = "on personsSQLite.Read(): "

func (personsOp *personsSQLite) Read(id auth.ID, identity *auth.Identity) (*persons.Item, error) {
	if identity == nil || (id != identity.ID && !identity.HasRole(rbac.RoleAdmin)) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onRead, "id": id, "requestedRole": rbac.RoleAdmin})
	}

	return personsOp.read(id)
}

const onList = "on personsSQLite.List()"

func (personsOp *personsSQLite) List(selector *selectors.Term, identity *auth.Identity) ([]persons.Item, error) {
	if !identity.HasRole(rbac.RoleAdmin) {
		return nil, errors.CommonError(common.NoRightsKey, common.Map{"on": onList, "requestedRole": rbac.RoleAdmin})
	}

	var condition string
	var values []interface{}

	if selector != nil {

		var valuesStr []string
		switch v := selector.Values.(type) {
		case []string:
			valuesStr = v
		case string:
			valuesStr = []string{v}

		// TODO!!! remove the kostyl
		case ns.URN:
			valuesStr = []string{string(v)}
		default:
			return nil, fmt.Errorf(onList+": wrong selector.Values: %#v --> %T", selector, selector.Values)
		}

		switch selector.Key {
		case persons.HasEmail:
			if len(valuesStr) != 1 {
				return nil, fmt.Errorf(onList+": wrong values list in selector: %#v / %#v", selector, valuesStr)
			}
			condition = `email = ?`
			values = []interface{}{valuesStr[0]}

		case persons.HasNickname:
			if len(valuesStr) != 1 {
				return nil, fmt.Errorf(onList+": wrong values list in selector: %#v / %#v", selector, valuesStr)
			}
			condition = `nickname = ?`
			values = []interface{}{valuesStr[0]}

		default:
			return nil, fmt.Errorf(onList+": wrong selector.Key: %#v", selector)
		}
	}

	query := sqllib.SQLList(personsOp.table, fieldsToListStr, condition, identity)
	stm, err := personsOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
	}

	rows, err := stm.Query(values...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	var items []persons.Item

	for rows.Next() {
		var idNum int64
		var item persons.Item

		var emailBytes, rolesBytes, credsBytes, infoBytes, tagsBytes, urnBytes, historyBytes []byte

		// "nickname", "email", "roles", "creds", "info",
		// "tags", "owner_nss", "viewer_nss", "history"
		// "urn",
		// "updated_at", "created_at", "id"

		if err := rows.Scan(
			&item.Nickname, &emailBytes, &rolesBytes, &credsBytes, &infoBytes,
			&tagsBytes, &item.OwnerNSS, &item.ViewerNSS, &historyBytes,
			&urnBytes,
			&item.UpdatedAt, &item.CreatedAt,
			&idNum,
		); err != nil {
			return nil, errors.Wrapf(err, onList+": "+sqllib.CantScanQueryRow, query, values)
		}

		if err := item.UnfoldFromJSON(auth.ID(strconv.FormatInt(idNum, 10)), rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, urnBytes, historyBytes); err != nil {
			return nil, errors.CommonError(err, onList)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onStat = "on personsSQLite.Stat()"

func (personsOp *personsSQLite) Stat(*selectors.Term, *auth.Identity) (db.StatMap, error) {
	var num int
	if err := personsOp.stmStat.QueryRow().Scan(&num); err == sql.ErrNoRows {
		return nil, errors.CommonError(common.ErrNotFound, onStat)
	} else if err != nil {
		return nil, errors.Wrapf(err, onStat+sqllib.CantScanQueryRow, personsOp.sqlStat, nil)
	}

	return db.StatMap{"*": db.Stat{num}}, nil
}

func (personsOp *personsSQLite) Close() error {
	return errors.Wrap(personsOp.db.Close(), "on personsSQLite.Close()")
}
