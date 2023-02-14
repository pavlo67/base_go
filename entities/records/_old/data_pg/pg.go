package data_pg

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/identity"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/libraries/sqllib"
	"github.com/pavlo67/common/common/libraries/sqllib/sqllib_pg"
	"github.com/pavlo67/common/common/libraries/strlib"
	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/common/common/selectors/logic"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
	"github.com/pkg/errors"
)

var fieldsToInsert = []string{"data_key", "url", "title", "summary", "embedded", "tags", "type_key", "content", "owner_key", "viewer_key", "history"}
var fieldsToInsertStr = strings.Join(fieldsToInsert, ", ")

var fieldsToUpdate = fieldsToInsert

var fieldsToRead = append(fieldsToUpdate, "updated_at", "created_at")
var fieldsToReadStr = strings.Join(fieldsToRead, ", ")

var fieldsToList = append([]string{"id"}, fieldsToRead...)

// var fieldsToListStr = strings.Join(fieldsToList, ", ")

var _ data.Operator = &dataPg{}

type dataPg struct {
	domain       identity.Domain
	interfaceKey joiner.InterfaceKey

	db              *sql.DB
	table           string
	fieldsToListStr string

	sqlInsert, sqlUpdate, sqlRead, sqlRemove, sqlClean string
	stmInsert, stmUpdate, stmRead, stmRemove           *sql.Stmt

	taggerOp      tagger.Operator
	taggerCleaner db.Cleaner
}

const onNew = "on dataPg.New(): "

func New(access config.Access, domain identity.Domain, table string, interfaceKey joiner.InterfaceKey, taggerOp tagger.Operator, taggerCleaner db.Cleaner) (data.Operator, db.Cleaner, error) {

	domain = domain.Normalize()
	if domain == "" {
		return nil, nil, errors.New("no service name")
	}

	db, err := sqllib_pg.Connect(access)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	if table == "" {
		table = data.CollectionDefault
	}

	var fieldsToListTabled []string
	for _, f := range fieldsToList {
		fieldsToListTabled = append(fieldsToListTabled, table+"."+f)
	}

	dataOp := dataPg{
		domain: domain,

		db:              db,
		table:           table,
		fieldsToListStr: strings.Join(fieldsToListTabled, ","),

		sqlInsert: "INSERT INTO " + table + " (" + fieldsToInsertStr + ") VALUES (" + sqllib_pg.WildcardsForInsert(fieldsToInsert) + ") RETURNING id",
		sqlUpdate: "UPDATE " + table + " SET " + sqllib_pg.WildcardsForUpdate(fieldsToUpdate) +
			" WHERE id = $" + strconv.Itoa(len(fieldsToUpdate)+1) + " AND owner_key = $" + strconv.Itoa(len(fieldsToUpdate)+2),
		sqlRemove: "DELETE FROM " + table + " WHERE id = $1 AND owner_key = $2",

		sqlRead: "SELECT " + fieldsToReadStr + " FROM " + table + " WHERE id = $1 AND viewer_key IN ('', $2)",
		//sqlList: sqllib.SQLList(table, fieldsToListStr, "", &crud.GetOptions{OrderBy: []string{"created_at DESC"}}),

		sqlClean: "DELETE FROM " + table,

		taggerOp:      taggerOp,
		interfaceKey:  interfaceKey,
		taggerCleaner: taggerCleaner,
	}

	sqlStmts := []sqllib.SqlStmt{
		{&dataOp.stmInsert, dataOp.sqlInsert},
		{&dataOp.stmUpdate, dataOp.sqlUpdate},
		{&dataOp.stmRemove, dataOp.sqlRemove},

		{&dataOp.stmRead, dataOp.sqlRead},
		//{&dataOp.stmList, dataOp.sqlList},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.Prepare(db, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return &dataOp, &dataOp, nil
}

const onSave = "on dataPg.Add(): "

func (dataOp *dataPg) Save(item data.Item, options *crud.SaveOptions) (common.ID, error) {

	if options == nil || options.ActorKey == "" {
		return "", errors.Errorf(onSave + "no user")
	}

	optionsToGet := &crud.GetOptions{ActorKey: options.ActorKey}

	item.ID = item.ID.Normalize()
	itemIdent := item.Key.Identity()
	if itemIdent != nil && itemIdent.Domain == dataOp.domain && itemIdent.Path == string(dataOp.interfaceKey) {
		if item.ID == "" {
			return "", errors.Errorf(onSave+"can't insert new data with predefined local .Key (%s --> %#v)", item.Key, itemIdent)
		} else if item.ID != itemIdent.ID.Normalize() {
			return "", errors.Errorf(onSave+"can't update data (.ID = %s) with incorrect local .Key: (%s --> %#v)", item.ID, item.Key, itemIdent)
		}
	}

	var err error
	var embedded, tags []byte

	if len(item.Embedded) > 0 {
		embedded, err = json.Marshal(item.Embedded)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal .Additions(%#v)", item)
		}
	}

	if len(item.Tags) > 0 {
		tags, err = json.Marshal(item.Tags)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal .Tags(%#v)", item)
		}
	}

	item.History = append(item.History, crud.Action{
		ActorKey: options.ActorKey,
		Key:      crud.SavedAction,
		DoneAt:   time.Now(),
	})

	// TODO: do it more clever
	item.OwnerKey = options.ActorKey
	item.ViewerKey = options.ActorKey

	var id common.ID

	if item.ID == "" {

		history, err := json.Marshal(item.History)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item)
		}

		values := []interface{}{item.Key, item.URL, item.Title, item.Summary, embedded, tags, item.Data.TypeKey, item.Data.Content, item.OwnerKey, item.ViewerKey, history}

		var lastInsertId uint64

		err = dataOp.stmInsert.QueryRow(values...).Scan(&lastInsertId)
		if err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlInsert, strlib.Stringify(values))
		}

		id = common.ID(strconv.FormatUint(lastInsertId, 10))

		if dataOp.taggerOp != nil && len(item.Tags) > 0 {
			err = dataOp.taggerOp.AddTags(joiner.Link{dataOp.interfaceKey, id}, item.Tags, options)
			if err != nil {
				return "", errors.Wrapf(err, onSave+": can't .AddTags(%#v)", item.Tags)
			}
		}

	} else {
		id = item.ID

		itemOld, err := dataOp.Read(id, optionsToGet)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't read old item with id = %s", id)
		}
		if itemOld == nil {
			return "", errors.Errorf(onSave+"old item with id = %s is nil", id)
		}

		err = item.History.CheckOn(itemOld.History)
		if err != nil {
			return "", errors.Wrap(err, onSave)
		}

		history, err := json.Marshal(item.History)
		if err != nil {
			return "", errors.Wrapf(err, onSave+"can't marshal .History(%#v)", item)
		}

		values := []interface{}{
			item.Key, item.URL, item.Title, item.Summary, embedded, tags, item.Data.TypeKey, item.Data.Content, item.OwnerKey, item.ViewerKey, history,
			item.ID, item.OwnerKey,
		}

		_, err = dataOp.stmUpdate.Exec(values...)
		if err != nil {
			return "", errors.Wrapf(err, onSave+sqllib.CantExec, dataOp.sqlUpdate, strlib.Stringify(values))
		}

		if dataOp.taggerOp != nil {
			// TODO: use one common transaction

			linkToTagged := joiner.Link{dataOp.interfaceKey, item.ID}

			err = dataOp.taggerOp.RemoveTagsAll(linkToTagged, options)
			if err != nil {
				return "", errors.Wrapf(err, onSave+": can't .RemoveTagsAll(%#v, %#v)", linkToTagged, options)
			}

			err = dataOp.taggerOp.AddTags(linkToTagged, item.Tags, options)
			if err != nil {
				return "", errors.Wrapf(err, onSave+": can't .AddTags(%#v, %#v, %#v)", linkToTagged, item.Tags, options)
			}
		}

	}

	return id, nil
}

const onRead = "on dataPg.Read(): "

func (dataOp *dataPg) Read(id common.ID, options *crud.GetOptions) (*data.Item, error) {
	if len(id) < 1 {
		return nil, errors.New(onRead + "empty Key")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return nil, errors.Errorf(onRead+"wrong Key (%s)", id)
	}

	var viewerKey identity.Key
	if options != nil {
		viewerKey = options.ActorKey

		// TODO!!! check if options.Term == nil

	}

	// TODO: check viewer_key for groups

	values := []interface{}{idNum, viewerKey}

	item := data.Item{ID: id}
	var embedded, tags, history []byte
	// var createdAtStr string
	var createdAt time.Time
	var updatedAtPtr *string

	err = dataOp.stmRead.QueryRow(values...).Scan(
		&item.Key, &item.URL, &item.Title, &item.Summary, &embedded, &tags, &item.Data.TypeKey, &item.Data.Content, &item.OwnerKey, &item.ViewerKey, &history, &updatedAtPtr,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, common.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, onRead+sqllib.CantScanQueryRow, dataOp.sqlRead, values)
	}

	if item.Key.Normalize() == "" {
		item.Key = (&identity.Item{Domain: dataOp.domain, Path: strings.TrimSpace(string(dataOp.interfaceKey)), ID: id}).Key()
	}

	if len(tags) > 0 {
		err = json.Unmarshal(tags, &item.Tags)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Tags (%s)", tags)
		}
	}

	if len(embedded) > 0 {
		err = json.Unmarshal(embedded, &item.Embedded)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .Additions (%s)", embedded)
		}
	}

	if len(history) > 0 {
		err = json.Unmarshal(history, &item.History)
		if err != nil {
			return &item, errors.Wrapf(err, onRead+"can't unmarshal .History (%s)", history)
		}
	}

	l.Info(createdAt)

	//createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	//if err != nil {
	//	// TODO??? return &item, errors.Wrapf(err, onRead+"can't parse .At (%s)", createdAtStr)
	//} else {
	//	item.History = item.History.SaveAction(crud.Action{Key: crud.CreatedAction, DoneAt: createdAt, Related: &joiner.Link{InterfaceKey: data.InterfaceKey, ID: id}})
	//}

	item.History = item.History.SaveAction(crud.Action{Key: crud.CreatedAction, DoneAt: createdAt, Related: &joiner.Link{InterfaceKey: data.InterfaceKey, ID: id}})

	if updatedAtPtr != nil {
		updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
		if err != nil {
			// TODO??? return &item, errors.Wrapf(err, onRead+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
		}
		item.History = item.History.SaveAction(crud.Action{Key: crud.UpdatedAction, DoneAt: updatedAt, Related: &joiner.Link{InterfaceKey: data.InterfaceKey, ID: id}})
	}

	return &item, nil
}

const onRemove = "on dataPg.Remove()"

func (dataOp *dataPg) Remove(id common.ID, options *crud.RemoveOptions) error {
	if options == nil || options.ActorKey == "" {
		return errors.Errorf(onRemove + "no user")
	}

	if len(id) < 1 {
		return errors.New(onRemove + "empty Key")
	}

	idNum, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return errors.Errorf(onRemove+"wrong Key (%s)", id)
	}

	// TODO: check owner_key for groups
	// TODO: deny the action if owner_key is empty

	values := []interface{}{idNum, options.ActorKey}

	_, err = dataOp.stmRemove.Exec(values...)
	if err != nil {
		return errors.Wrapf(err, onRemove+sqllib.CantExec, dataOp.sqlRemove, values)
	}

	if dataOp.taggerOp != nil {
		err = dataOp.taggerOp.RemoveTagsAll(joiner.Link{dataOp.interfaceKey, id}, &crud.SaveOptions{ActorKey: options.ActorKey})
		if err != nil {
			return errors.Wrapf(err, onRemove+": can't .ReplaceTags(%#v)", nil)
		}
	}

	return nil
}

const onList = "on dataPg.List()"

func (dataOp *dataPg) List(options *crud.GetOptions) ([]data.Item, error) {

	// TODO!!! use default key's value searching data with empty .Key

	query, values, err := sqllib.SQLList(dataOp.table, dataOp.fieldsToListStr, options, sqllib_pg.CorrectWildcards)
	if err != nil {
		return nil, errors.Errorf(onList+"can't sqllib.SQLList(%s, %s, %#v, sqllib_pg.CorrectWildcards): %s", dataOp.table, dataOp.fieldsToListStr, options, err)
	}

	// l.Infof("%s / %#v\n%s", condition, values, query)

	stm, err := dataOp.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, onList+": can't db.Prepare(%s)", query)
	}

	rows, err := stm.Query(values...)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrapf(err, onList+sqllib.CantQuery, query, values)
	}
	defer rows.Close()

	items := []data.Item{} // to return [] to front-end instead null

	for rows.Next() {
		var idNum int64
		var item data.Item
		var embedded, tags, history []byte
		var createdAtStr string
		var updatedAtPtr *string

		err := rows.Scan(
			&idNum, &item.Key, &item.URL, &item.Title, &item.Summary, &embedded, &tags, &item.Data.TypeKey, &item.Data.Content, &item.OwnerKey, &item.ViewerKey, &history,
			&updatedAtPtr, &createdAtStr,
		)

		if err != nil {
			return items, errors.Wrapf(err, onList+sqllib.CantScanQueryRow, query, values)
		}

		item.ID = common.ID(strconv.FormatInt(idNum, 10))

		if item.Key.Normalize() == "" {
			item.Key = (&identity.Item{Domain: dataOp.domain, Path: strings.TrimSpace(string(dataOp.interfaceKey)), ID: item.ID}).Key()
		}

		if len(tags) > 0 {
			if err = json.Unmarshal(tags, &item.Tags); err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Tags (%s)", tags)
			}
		}

		if len(embedded) > 0 {
			if err = json.Unmarshal(embedded, &item.Embedded); err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .Additions (%s)", embedded)
			}
		}

		if len(history) > 0 {
			err = json.Unmarshal(history, &item.History)
			if err != nil {
				return items, errors.Wrapf(err, onList+"can't unmarshal .History (%s)", history)
			}
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			// TODO??? return &item, errors.Wrapf(err, onList+"can't parse .At (%s)", createdAtStr)
		} else {
			item.History = item.History.SaveAction(crud.Action{Key: crud.CreatedAction, DoneAt: createdAt, Related: &joiner.Link{InterfaceKey: data.InterfaceKey, ID: item.ID}})
		}

		if updatedAtPtr != nil {
			updatedAt, err := time.Parse(time.RFC3339, *updatedAtPtr)
			if err != nil {
				// TODO??? return &item, errors.Wrapf(err, onList+"can't parse .UpdatedAt (%s)", *updatedAtPtr)
			}
			item.History = item.History.SaveAction(crud.Action{Key: crud.UpdatedAction, DoneAt: updatedAt, Related: &joiner.Link{InterfaceKey: data.InterfaceKey, ID: item.ID}})
		}

		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		return items, errors.Wrapf(err, onList+": "+sqllib.RowsError, query, values)
	}

	return items, nil
}

const onCount = "on dataPg.Count(): "

func (dataOp *dataPg) Count(options *crud.GetOptions) (uint64, error) {

	// TODO: check viewer_key

	query, values, err := sqllib.SQLCount(dataOp.table, options, sqllib_pg.CorrectWildcards)
	if err != nil {
		return 0, errors.Wrapf(err, onCount+": can't sqllib.SQLCount(%s, %#v, sqllib_pg.CorrectWildcards)", dataOp.table, options)
	}

	stm, err := dataOp.db.Prepare(query)
	if err != nil {
		return 0, errors.Wrapf(err, onCount+": can't db.Prepare(%s)", query)
	}

	var num uint64

	err = stm.QueryRow(values...).Scan(&num)
	if err != nil {
		return 0, errors.Wrapf(err, onCount+sqllib.CantScanQueryRow, query, values)
	}

	return num, nil
}

func (dataOp *dataPg) Tagger() tagger.Operator {
	return dataOp.taggerOp
}

const onListTagged = "on dataPg.ListTagged(): "

func (dataOp *dataPg) ListTagged(tagLabel string, options *crud.GetOptions) ([]data.Item, error) {
	if dataOp.taggerOp == nil {
		return nil, errors.New(onListTagged + ": no tagger.Operator")
	}

	// select storage.* from storage join tagged on storage.id::text = tagged.id and tag = 'www';
	if options == nil {
		options = &crud.GetOptions{}
	}
	options.JoinTo = crud.JoinTo{
		Clause: "join tagged on " + dataOp.table + ".id::text = tagged.id and tag = ?",
		Values: []interface{}{tagLabel},
	}

	return dataOp.List(options)

	//return data.ListTagged(dataOp, dataOp.taggerOp, &dataOp.interfaceKey, tagLabel, options)
}

const onListUntagged = "on dataPg.ListTagged(): "

func (dataOp *dataPg) ListUntagged(tagLabel string, options *crud.GetOptions) ([]data.Item, error) {
	if dataOp.taggerOp == nil {
		return nil, errors.New(onListTagged + ": no tagger.Operator")
	}

	// select storage.* from storage join tagged on storage.id::text = tagged.id and tag = 'www';
	if options == nil {
		options = &crud.GetOptions{}
	}
	options.JoinTo = crud.JoinTo{
		Clause: "LEFT JOIN tagged ON " + dataOp.table + ".id::text = tagged.id AND tag = ?",
		Values: []interface{}{tagLabel},
	}
	if options.Term == nil {
		options.Term = selectors.String("tag IS NULL")
	} else {
		options.Term = logic.AND(options.Term, selectors.String("tag IS NULL"))
	}

	return dataOp.List(options)

	//return data.ListTagged(dataOp, dataOp.taggerOp, &dataOp.interfaceKey, tagLabel, options)
}

func (dataOp *dataPg) Close() error {
	return errors.Wrap(dataOp.db.Close(), "on dataPg.Close()")
}
