package persons_pg

//import (
//	"usersbase/sql"
//	"encoding/json"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/pavlo67/workshop/common/identity"
//
//	"github.com/pavlo67/workshop/common"
//	"github.com/pavlo67/workshop/common/config"
//	"github.com/pavlo67/workshop/common/crud"
//	"github.com/pavlo67/workshop/common/joiner"
//	"github.com/pavlo67/workshop/common/libraries/sqllib"
//	"github.com/pavlo67/workshop/common/libraries/sqllib/sqllib_pg"
//	"github.com/pavlo67/workshop/common/libraries/strlib"
//	"github.com/pavlo67/workshop/common/selectors"
//	"github.com/pavlo67/workshop/common/selectors/selectors_sql"
//	"github.com/pkg/errors"
//
//	"github.com/pavlo67/workshop/components/users"
//	"github.com/pavlo67/workshop/components/tagger"
//)
//
//var fieldsToInsert = []string{"users_key", "url", "title", "summary", "embedded", "tags", "type_key", "content", "history"}
//var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")
//
//var fieldsToUpdate = fieldsToInsert
//
//var fieldsToRead = append(fieldsToUpdate, "updated_at", "created_at")
//var fieldsToReadStr = strings.Join(fieldsToRead, ", ")
//
//var fieldsToList = append([]string{"id"}, fieldsToRead...)
//var fieldsToListStr = strings.Join(fieldsToList, ", ")
//
//var _ users.Operator = &usersPg{}
//
//type usersPg struct {
//	db    *sql.DB
//	table string
//
//	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlList, sqlClean string
//	stmInsert, stmUpdate, stmRead, stmRemove, stmList           *sql.Stmt
//
//	taggerOp      tagger.Operator
//	interfaceKey  joiner.InterfaceKey
//	taggerCleaner crud.Cleaner
//}
//
//const onNew = "on usersPg.New(): "
//
//func New(access config.Access, table string, interfaceKey joiner.InterfaceKey, taggerOp tagger.Operator, taggerCleaner crud.Cleaner) (users.Operator, crud.Cleaner, error) {
//	db, err := sqllib_pg.Connect(access)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, onNew)
//	}
//
//	if table == "" {
//		table = users.CollectionDefault
//	}
//
//	usersOp := usersPg{
//		db:    db,
//		table: table,
//
//		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",
//		sqlUpdate: "UPDATE " + table + " SET " + sqllib_pg.WildcardsForUpdate(fieldsToUpdate) + " WHERE id = $" + strconv.Itoa(len(fieldsToUpdate)+1),
//		sqlRemove: "DELETE FROM " + table + " where id = $1",
//
//		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1",
//		sqlList: sqllib.SQLList(table, fieldsToListStr, "", &crud.GetOptions{OrderBy: []string{"created_at DESC"}}),
//
//		sqlClean: "DELETE FROM " + table,
//
//		taggerOp:      taggerOp,
//		interfaceKey:  interfaceKey,
//		taggerCleaner: taggerCleaner,
//	}
//
//	sqlStmts := []sqllib.SqlStmt{
//		{&usersOp.stmInsert, usersOp.sqlInsert},
//		{&usersOp.stmUpdate, usersOp.sqlUpdate},
//		{&usersOp.stmRemove, usersOp.sqlRemove},
//
//		{&usersOp.stmRead, usersOp.sqlRead},
//		{&usersOp.stmList, usersOp.sqlList},
//	}
//
//	for _, sqlStmt := range sqlStmts {
//		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
//			return nil, nil, errors.Wrap(err, onNew)
//		}
//	}
//
//	return &usersOp, &usersOp, nil
//}
//
//const onSave = "on usersPg.Save(): "
//
//func (usersOp *usersPg) Save(item users.Item, options *crud.SaveOptions) (common.ID, error) {
//
//	var actor *identity.Key
//	if options != nil {
//		actor = options.ActorKey
//	}
//
//	item.History = append(item.History, crud.Action{
//		ActorKey:  actor,
//		Key:    crud.SavedAction,
//		DoneAt: time.Now(),
//	})
//
//	history, err := json.Marshal(item.History)
//	if err != nil {
//		return "", errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item)
//	}
//
//	var embedded, tags []byte
//
//	if len(item.Embedded) > 0 {
//		embedded, err = json.Marshal(item.Embedded)
//		if err != nil {
//			return "", errors.Wrapf(err, onSave+"can't marshal .Embedded(%#v)", item)
//		}
//	}
//
//	if len(item.Tags) > 0 {
//		tags, err = json.Marshal(item.Tags)
//		if err != nil {
//			return "", errors.Wrapf(err, onSave+"can't marshal .Tags(%#v)", item)
//		}
//	}
//
//	values := []interface{}{item.Key, item.URL, item.Title, item.Summary, embedded, tags, item.users.TypeKey, item.users.Content01, history}
//
//	var id common.ID
//
//	if item.ID == "" {
//
//		var lastInsertId uint64
//
//		err := usersOp.stmInsert.QueryRow(values...).Scan(&lastInsertId)
//		if err != nil {
//			return "", errors.Wrapf(err, onSave+sqllib.CantExec, usersOp.sqlInsert, strlib.Stringify(values))
//		}
//
//		id = common.ID(strconv.FormatUint(lastInsertId, 10))
//
//		if usersOp.taggerOp != nil && len(item.Tags) > 0 {
//			err = usersOp.taggerOp.AddTags(joiner.Link{usersOp.interfaceKey, id}, item.Tags, nil)
//			if err != nil {
//				return "", errors.Wrapf(err, onSave+": can't .AddTags(%#v)", item.Tags)
//			}
//		}
//
//	} else {
//		id = item.ID
//
//		values := append(values, item.ID)
//
//		_, err := usersOp.stmUpdate.Exec(values...)
//		if err != nil {
//			return "", errors.Wrapf(err, onSave+sqllib.CantExec, usersOp.sqlUpdate, strlib.Stringify(values))
//		}
//
//		if usersOp.taggerOp != nil {
//			err = usersOp.taggerOp.ReplaceTags(joiner.Link{usersOp.interfaceKey, item.ID}, item.Tags, nil)
//			if err != nil {
//				return "", errors.Wrapf(err, onSave+": can't .ReplaceTags(%#v)", item.Tags)
//			}
//		}
//
//	}
//
//	return id, nil
//}
//
//const onRead = "on usersPg.Read(): "
//
//func (usersOp *usersPg) Read(id common.ID, _ *crud.GetOptions) (*users.Item, error) {
//	if len(id) < 1 {
//		return nil, errors.New(onRead + "empty Key")
//	}
//
//	idNum, err := strconv.ParseUint(string(id), 10, 64)
//	if err != nil {
//		return nil, errors.Errorf(onRead+"wrong Key (%s)", id)
//	}
//
//	item := users.Item{ID: id}
//	var embedded, tags, history []byte
//	var createdAtStr string
//	var updatedAtPtr *string
//
//	err = usersOp.stmRead.QueryRow(idNum).Scan(
//		&item.Key, &item.URL, &item.Title, &item.Summary, &embedded, &tags, &item.users.TypeKey, &item.users.Content01, &history, &updatedAtPtr, &createdAtStr,
//	)
//
//	if err == sql.ErrNoRows {
//		return nil, common.ErrNotFound
//	}
//	if err != nil {
//		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, usersOp.sqlRead, idNum)
//	}
//
//	if len(tags) > 0 {
//		err = json.Unmarshal(tags, &item.Tags)
//		if err != nil {
//			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tags)
//		}
//	}
//
//	if len(embedded) > 0 {
//		err = json.Unmarshal(embedded, &item.Embedded)
//		if err != nil {
//			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Embedded (%s)", embedded)
//		}
//	}
//
//	if len(history) > 0 {
//		err = json.Unmarshal(history, &item.History)
//		if err != nil {
//			return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", history)
//		}
//	}
//
//	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
//	if err != nil {
//		// TODO??? return &item, errors.Wrapf(err, onRead+"can't parse .CreatedAt (%s)", createdAtStr)
//	} else {
//		item.History = item.History.SaveAction(crud.Action{Key: crud.CreatedAction, DoneAt: createdAt, Related: &joiner.Link{InterfaceKey: users.InterfaceKey, ID: id}})
//	}
//
//	if updatedAtPtr != nil {
//		updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
//		if err != nil {
//			// TODO??? return &item, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
//		}
//		item.History = item.History.SaveAction(crud.Action{Key: crud.UpdatedAction, DoneAt: updatedAt, Related: &joiner.Link{InterfaceKey: users.InterfaceKey, ID: id}})
//	}
//
//	return &item, nil
//}
//
//const onRemove = "on usersPg.Remove()"
//
//func (usersOp *usersPg) Remove(id common.ID, _ *crud.RemoveOptions) error {
//	if len(id) < 1 {
//		return errors.New(onRemove + "empty Key")
//	}
//
//	idNum, err := strconv.ParseUint(string(id), 10, 64)
//	if err != nil {
//		return errors.Errorf(onRemove+"wrong Key (%s)", id)
//	}
//
//	_, err = usersOp.stmRemove.Exec(idNum)
//	if err != nil {
//		return errors.Wrapf(err, onRemove+sqllib.CantExec, usersOp.sqlRemove, idNum)
//	}
//
//	if usersOp.taggerOp != nil {
//		err = usersOp.taggerOp.ReplaceTags(joiner.Link{usersOp.interfaceKey, id}, nil, nil)
//		if err != nil {
//			return errors.Wrapf(err, onRemove+": can't .ReplaceTags(%#v)", nil)
//		}
//	}
//
//	return nil
//}
//
//const onList = "on usersPg.List()"
//
//func (usersOp *usersPg) List(term *selectors.Term, options *crud.GetOptions) ([]users.Item, error) {
//	condition, values, err := selectors_sql.Use(term)
//	if err != nil {
//		return nil, errors.Errorf(onList+"wrong selector (%#v): %s", term, err)
//	}
//
//	query := usersOp.sqlList
//	stm := usersOp.stmList
//
//	if condition != "" || options != nil {
//		query = sqllib.SQLList(usersOp.table, fieldsToListStr, condition, options)
//		stm, err = usersOp.db.Prepare(query)
//		if err != nil {
//			return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
//		}
//	}
//
//	//l.Infof("%s / %#v\n%s", condition, values, query)
//
//	rows, err := stm.Query(values...)
//
//	if err == sql.ErrNoRows {
//		return nil, nil
//	} else if err != nil {
//		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, query, values)
//	}
//	defer rows.Close()
//
//	var items []users.Item
//
//	for rows.Next() {
//		var idNum int64
//		var item users.Item
//		var embedded, tags, history []byte
//		var createdAtStr string
//		var updatedAtPtr *string
//
//		err := rows.Scan(
//			&idNum, &item.Key, &item.URL, &item.Title, &item.Summary, &embedded, &tags, &item.users.TypeKey, &item.users.Content01, &history, &updatedAtPtr, &createdAtStr,
//		)
//
//		if err != nil {
//			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
//		}
//
//		if len(tags) > 0 {
//			if err = json.Unmarshal(tags, &item.Tags); err != nil {
//				return items, errors.Wrapf(err, onList+"can't unmarshal .Tags (%s)", tags)
//			}
//		}
//
//		if len(embedded) > 0 {
//			if err = json.Unmarshal(embedded, &item.Embedded); err != nil {
//				return items, errors.Wrapf(err, onList+"can't unmarshal .Embedded (%s)", embedded)
//			}
//		}
//
//		item.ID = common.ID(strconv.FormatInt(idNum, 10))
//
//		if len(history) > 0 {
//			err = json.Unmarshal(history, &item.History)
//			if err != nil {
//				return items, errors.Wrapf(err, onList+"can't unmarshal .History (%s)", history)
//			}
//		}
//
//		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
//		if err != nil {
//			// TODO??? return &item, errors.Wrapf(err, onList+"can't parse .CreatedAt (%s)", createdAtStr)
//		} else {
//			item.History = item.History.SaveAction(crud.Action{Key: crud.CreatedAction, DoneAt: createdAt, Related: &joiner.Link{InterfaceKey: users.InterfaceKey, ID: item.ID}})
//		}
//
//		if updatedAtPtr != nil {
//			updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
//			if err != nil {
//				// TODO??? return &item, errors.Wrapf(err, onList+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
//			}
//			item.History = item.History.SaveAction(crud.Action{Key: crud.UpdatedAction, DoneAt: updatedAt, Related: &joiner.Link{InterfaceKey: users.InterfaceKey, ID: item.ID}})
//		}
//
//		items = append(items, item)
//	}
//	err = rows.Err()
//	if err != nil {
//		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
//	}
//
//	return items, nil
//}
//
//const onCount = "on usersPg.Count(): "
//
//func (usersOp *usersPg) Count(term *selectors.Term, options *crud.GetOptions) (uint64, error) {
//	condition, values, err := selectors_sql.Use(term)
//	if err != nil {
//		termStr, _ := json.Marshal(term)
//		return 0, errors.Wrapf(err, onCount+": can't selectors_sql.Use(%s)", termStr)
//	}
//
//	query := sqllib_pg.CorrectWildcards(sqllib.SQLCount(usersOp.table, condition, options))
//	stm, err := usersOp.db.Prepare(query)
//	if err != nil {
//		return 0, errors.Wrapf(err, onCount+": can't db.Prepare(%s)", query)
//	}
//
//	var num uint64
//
//	err = stm.QueryRow(values...).Scan(&num)
//	if err != nil {
//		return 0, errors.Wrapf(err, onCount+sqllib.CantScanQueryRow, query, values)
//	}
//
//	return num, nil
//}
//
//func (usersOp *usersPg) Close() error {
//	return errors.Wrap(usersOp.db.Close(), "on usersPg.Close()")
//}
