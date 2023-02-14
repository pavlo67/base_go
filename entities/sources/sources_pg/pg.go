package sources_pg

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/pkg/errors"

// 	"github.com/pavlo67/common/common"
// 	"github.com/pavlo67/common/common/config"
// 	"github.com/pavlo67/common/common/crud"
// 	"github.com/pavlo67/common/common/joiner"
// 	"github.com/pavlo67/common/common/libraries/sqllib"
// 	"github.com/pavlo67/common/common/libraries/sqllib/sqllib_pg"
// 	"github.com/pavlo67/common/common/selectors"
// 	"github.com/pavlo67/common/common/selectors/selectors_sql"

// 	"github.com/pavlo67/workshop/components/packs"
// )

// var fieldsToInsert = []string{"key", "address_from", "address_to", "options", "type_key", "content", "history", "created_at"}
// var fieldsToInsertStr = strings.Join(fieldsToInsert, ",")

// var fieldsToRead = fieldsToInsert
// var fieldsToReadStr = strings.Join(fieldsToRead, ",")

// var fieldsToList = append([]string{"id"}, fieldsToRead...)
// var fieldsToListStr = strings.Join(fieldsToList, ",")

// var _ packs.ActorKey = &packsPg{}

// type packsPg struct {
// 	db    *sql.DB
// 	table string

// 	sqlInsert, sqlRead, sqlList, sqlReadToStart, sqlStart, sqlReadToAddHistory, sqlAddHistory, sqlClean string
// 	stmInsert, stmRead, stmList, stmReadToStart, stmStart, stmReadToAddHistory, stmAddHistory           *sql.Stmt

// 	interfaceKey joiner.InterfaceKey
// }

// const onNew = "on packsPg.New(): "

// func New(access config.Access, table string, interfaceKey joiner.InterfaceKey) (packs.ActorKey, db.Cleaner, error) {
// 	db, err := sqllib_pg.Connect(access)
// 	if err != nil {
// 		return nil, nil, errors.Wrap(err, onNew)
// 	}

// 	if table == "" {
// 		table = packs.CollectionDefault
// 	}

// 	packsOp := packsPg{
// 		db:    db,
// 		table: table,

// 		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",

// 		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1",
// 		sqlList: sqllib.SQLList(table, fieldsToListStr, "", &crud.GetOptions{OrderBy: []string{"created_at"}}),

// 		sqlReadToAddHistory: "SELECT history FROM " + table + " WHERE id = $1",
// 		sqlAddHistory:       "UPDATE " + table + " SET history = $1 WHERE id = $2",

// 		//sqlRemove: "DELETE FROM " + table + " where Key = $1",
// 		sqlClean: "DELETE FROM " + table,

// 		interfaceKey: interfaceKey,
// 	}

// 	sqlStmts := []sqllib.SqlStmt{
// 		{&packsOp.stmInsert, packsOp.sqlInsert},
// 		{&packsOp.stmRead, packsOp.sqlRead},
// 		{&packsOp.stmList, packsOp.sqlList},
// 		{&packsOp.stmReadToAddHistory, packsOp.sqlReadToAddHistory},
// 		{&packsOp.stmAddHistory, packsOp.sqlAddHistory},

// 		//{&packsOp.stmRemove, packsOp.sqlRemove},
// 	}

// 	for _, sqlStmt := range sqlStmts {
// 		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
// 			return nil, nil, errors.Wrap(err, onNew)
// 		}
// 	}

// 	return &packsOp, &packsOp, nil
// }

// const onSave = "on packsPg.Save(): "

// func (packsOp *packsPg) Save(pack *packs.Pack, _ *crud.SaveOptions) (common.Key, error) {
// 	if pack == nil {
// 		return "", errors.New(onSave + "nothing to save")
// 	}

// 	var toBytes []byte
// 	if len(pack.To) > 0 {
// 		var err error
// 		toBytes, err = json.Marshal(pack.To)
// 		if err != nil {
// 			return "", errors.Wrapf(err, onSave+"can't .Marshal(.To == %#v)", pack.To)
// 		}
// 	}

// 	var optionsBytes []byte
// 	if len(pack.Options) > 0 {
// 		var err error
// 		optionsBytes, err = json.Marshal(pack.Options)
// 		if err != nil {
// 			return "", errors.Wrapf(err, onSave+"can't .Marshal(.Options == %#v)", pack.Options)
// 		}
// 	}

// 	var contentBytes []byte
// 	if pack.Content != nil {
// 		var err error
// 		contentBytes, err = json.Marshal(pack.Content)
// 		if err != nil {
// 			return "", errors.Wrapf(err, onSave+"can't .Marshal(.Content == %#v)", pack.Content)
// 		}
// 	}

// 	var historyBytes []byte
// 	if len(pack.History) > 0 {
// 		var err error
// 		historyBytes, err = json.Marshal(pack.History)
// 		if err != nil {
// 			return "", errors.Wrapf(err, onSave+"can't .Marshal(.History == %#v)", pack.History)
// 		}
// 	}

// 	var createdAt time.Time
// 	if len(pack.History) > 0 {
// 		createdAt = pack.History[0].DoneAt
// 		if pack.History[0].Key == crud.CreatedAction {
// 			pack.History = pack.History[1:]
// 		}
// 	} else {
// 		createdAt = time.Now()
// 	}

// 	values := []interface{}{pack.Key, pack.From, toBytes, optionsBytes, pack.TypeKey, contentBytes, historyBytes, createdAt}

// 	var lastInsertId uint64
// 	err := packsOp.stmInsert.QueryRow(values...).Scan(&lastInsertId)
// 	if err != nil {
// 		return "", errors.Wrapf(err, onSave+sqllib.CantExec, packsOp.sqlInsert, values)
// 	}

// 	return common.Key(strconv.FormatUint(lastInsertId, 10)), nil
// }

// const onRead = "on packsPg.Read(): "

// func (packsOp *packsPg) Read(id common.Key, _ *crud.GetOptions) (*packs.Item, error) {
// 	if len(id) < 1 {
// 		return nil, errors.New(onRead + "empty Key")
// 	}

// 	idNum, err := strconv.ParseUint(string(id), 10, 64)
// 	if err != nil {
// 		return nil, errors.Errorf(onRead+"wrong Key (%s)", id)
// 	}

// 	item := packs.Item{Key: id}

// 	var toBytes, optionsBytes, historyBytes []byte
// 	var createdAtStr string

// 	err = packsOp.stmRead.QueryRow(idNum).Scan(
// 		&item.Key, &item.From, &toBytes, &optionsBytes, &item.TypeKey, &item.ContentRaw, &historyBytes, &createdAtStr,
// 	)
// 	if err == sql.ErrNoRows {
// 		return nil, common.ErrNotFound
// 	}
// 	if err != nil {
// 		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, packsOp.sqlRead, idNum)
// 	}

// 	if len(toBytes) > 0 {
// 		err = json.Unmarshal(toBytes, &item.To)
// 		if err != nil {
// 			return &item, errors.Wrapf(err, onRead+"can't unmarshal .To (%s)", toBytes)
// 		}
// 	}

// 	if len(optionsBytes) > 0 {
// 		err = json.Unmarshal(optionsBytes, &item.Options)
// 		if err != nil {
// 			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Options (%s)", optionsBytes)
// 		}
// 	}

// 	if len(historyBytes) > 0 {
// 		err = json.Unmarshal(historyBytes, &item.History)
// 		if err != nil {
// 			return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", historyBytes)
// 		}
// 	}

// 	if !(len(item.History) > 0 && item.History[0].Key == crud.CreatedAction) {
// 		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
// 		if err != nil {
// 			return &item, errors.Wrapf(err, onRead+"can't parse .At (%s)", createdAtStr)
// 		}

// 		item.History = append([]crud.ActionKey{{Key: crud.CreatedAction, DoneAt: createdAt}}, item.History...)
// 	}

// 	return &item, nil
// }

// const onRemove = "on packsPg.Remove()"

// func (packsOp *packsPg) Remove(common.Key, *crud.RemoveOptions) error {
// 	return common.ErrNotImplemented
// }

// const onList = "on packsPg.List()"

// func (packsOp *packsPg) List(term *selectors.Term, options *crud.GetOptions) ([]packs.Item, error) {
// 	condition, values, err := selectors_sql.Use(term)
// 	if err != nil {
// 		return nil, errors.Errorf(onList+"wrong selector (%#v): %s", term, err)
// 	}

// 	query := packsOp.sqlList
// 	stm := packsOp.stmList

// 	if condition != "" || options != nil {
// 		query = sqllib_pg.CorrectWildcards(sqllib.SQLList(packsOp.table, fieldsToListStr, condition, options))

// 		stm, err = packsOp.db.Prepare(query)
// 		if err != nil {
// 			return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
// 		}
// 	}

// 	// l.Infof("%s / %#v\n%s", condition, values, query)

// 	rows, err := stm.Query(values...)
// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	} else if err != nil {
// 		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, query, values)
// 	}
// 	defer rows.Close()

// 	var items []packs.Item

// 	for rows.Next() {
// 		var idNum int64
// 		var item packs.Item

// 		var toBytes, optionsBytes, historyBytes []byte
// 		var createdAtStr string

// 		err := rows.Scan(
// 			&idNum, &item.Key, &item.From, &toBytes, &optionsBytes, &item.TypeKey, &item.ContentRaw, &historyBytes, &createdAtStr,
// 		)
// 		if err != nil {
// 			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
// 		}

// 		item.Key = common.Key(strconv.FormatInt(idNum, 10))

// 		if len(toBytes) > 0 {
// 			err = json.Unmarshal(toBytes, &item.To)
// 			if err != nil {
// 				return items, errors.Wrapf(err, onList+"can't unmarshal .To (%s)", toBytes)
// 			}
// 		}

// 		if len(optionsBytes) > 0 {
// 			err = json.Unmarshal(optionsBytes, &item.Options)
// 			if err != nil {
// 				return items, errors.Wrapf(err, onList+"can't unmarshal .Options (%s)", optionsBytes)
// 			}
// 		}

// 		if len(historyBytes) > 0 {
// 			err = json.Unmarshal(historyBytes, &item.History)
// 			if err != nil {
// 				return items, errors.Wrapf(err, onList+"can't unmarshal .History (%s)", historyBytes)
// 			}
// 		}

// 		if !(len(item.History) > 0 && item.History[0].Key == crud.CreatedAction) {
// 			createdAt, err := time.Parse(time.RFC3339, createdAtStr)
// 			if err != nil {
// 				return items, errors.Wrapf(err, onList+"can't parse .At (%s)", createdAtStr)
// 			}

// 			item.History = append([]crud.ActionKey{{Key: crud.CreatedAction, DoneAt: createdAt}}, item.History...)
// 		}

// 		items = append(items, item)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
// 	}

// 	return items, nil
// }

// const onAddHistory = "on packsPg.AddHistory(): "

// func (packsOp *packsPg) AddHistory(id common.Key, historyToAdd crud.History, _ *crud.SaveOptions) (crud.History, error) {

// 	// nothing to do

// 	if len(historyToAdd) < 1 {
// 		return nil, nil
// 	}

// 	// reading old .History

// 	if len(id) < 1 {
// 		return nil, errors.New(onAddHistory + "empty Key")
// 	}

// 	idNum, err := strconv.ParseUint(string(id), 10, 64)
// 	if err != nil {
// 		return nil, errors.Errorf(onAddHistory+"wrong Key (%s)", id)
// 	}

// 	var historyBytes []byte

// 	err = packsOp.stmReadToAddHistory.QueryRow(idNum).Scan(&historyBytes)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, onAddHistory+sqllib.CantScanQueryRow, packsOp.sqlReadToAddHistory, idNum)
// 	}

// 	// adding the result ------------------------------------------------------------------------------

// 	var history []crud.ActionKey

// 	if len(historyBytes) > 0 {
// 		err = json.Unmarshal(historyBytes, &history)
// 		if err != nil {
// 			return nil, errors.Wrapf(err, onAddHistory+"can't unmarshal .History (%s)", historyBytes)
// 		}
// 	}

// 	historyNew := append(history, historyToAdd...)
// 	historyBytesNew, err := json.Marshal(historyNew)
// 	if err != nil {
// 		return historyNew, errors.Wrapf(err, onAddHistory+"can't .Marshal(%#v)", historyNew)
// 	}

// 	// saving the updates -----------------------------------------------------------------------------

// 	values := []interface{}{historyBytesNew, idNum}
// 	_, err = packsOp.stmAddHistory.Exec(values...)
// 	if err != nil {
// 		return historyNew, errors.Wrapf(err, onAddHistory+sqllib.CantExec, packsOp.sqlAddHistory, values)
// 	}

// 	return historyNew, nil
// }

// func (packsOp *packsPg) Close() error {
// 	return errors.Wrap(packsOp.db.Close(), "on packsPg.Close()")
// }
