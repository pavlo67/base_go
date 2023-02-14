package links_mysql

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/pavlo67/partes/crud/selector"
	"github.com/pavlo67/partes/libs/mysqllib"
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"

	"github.com/pavlo67/workshop/confidenter/groups"
	"github.com/pavlo67/workshop/confidenter/rights"

	"os"

	"github.com/pavlo67/workshop/notebook/links"
)

// linksMySQL is a struct to implement Object interface using MySQL.
type linksMySQL struct {
	ctrlOp       groups.Operator
	dataManagers rights.Managers

	dbh       *sql.DB
	linkTable string

	stmtCreate *sql.Stmt
	stmtDelete *sql.Stmt

	stmtqueryByOwner *sql.Stmt
	stmtdeleteLinks  *sql.Stmt

	stmtQueryByTag      *sql.Stmt
	stmtQueryByObjectID *sql.Stmt

	tplDelete, tplCountLinked, tplCountLinkedFin string
	tplQueryTags, tplQueryTagsFin                string
	tplQuery                                     string

	tplQueryTagsByOwner, tplQueryTagsByOwnerFin string

	sqlqueryByOwner, sqldeleteLinks string

	sqlCreate, sqlDelete, sqlQueryTags, sqlQueryByTag, sqlQueryByObjectID string
}

var _ links.Operator = &linksMySQL{}

const onNew = "on links_mysql.New()"

var linkFieldsList = []string{
	links.FieldType, links.FieldTag, "object_id", "r_view", "r_owner",
}

// New ...
func New(mysqlConfig config.ServerAccess, linkTable string, ctrlOp groups.Operator, managers rights.Managers) (*linksMySQL, error) {
	if ctrlOp == nil {
		return nil, errors.New("no ctrlOp... WTF?")
	}

	dbh, err := mysqllib.ConnectToMysql(mysqlConfig)
	if err != nil {
		return nil, errors.Wrapf(err, onNew+": credentials = %#v", mysqlConfig)
	}

	l := new(linksMySQL)

	l.ctrlOp = ctrlOp
	l.dataManagers = managers
	if len(l.dataManagers) < 1 {
		l.dataManagers = rights.Managers{rights.Create: libs.Anyone}
	}

	l.dbh = dbh
	l.linkTable = linkTable

	linkFieldsToCreate := "linked_type, linked_id, " + strings.Join(linkFieldsList, ", ")

	//linkFields := "id, " + strings.Join123(linkFieldsList, ", ") // + ", created_at, updated_at"
	//linkFieldsToUpdate := strings.Join123(linkFieldsList, " = ?, ") + " = ?"

	l.tplDelete = "delete from `" + linkTable + "` where "
	l.tplCountLinked = "select object_id, count(*) from `" + linkTable + "`"
	l.tplCountLinkedFin = " group by object_id"

	l.tplQueryTags = "select `tag`, count(*) as `num` from `" + linkTable + "`"
	l.tplQueryTagsFin = "group by `tag` order by `num` desc"

	l.sqlqueryByOwner = "select object_id from `" + linkTable + "` where linked_type = ? and linked_id = ? and r_owner = ?"
	l.sqldeleteLinks = "delete from `" + linkTable + "` where linked_type = ? and linked_id = ? and id > 0"

	l.sqlCreate = "insert into `" + linkTable + "` (" + linkFieldsToCreate + ") values (?,?, ?,?,?,?,?)"
	l.sqlDelete = l.tplDelete + "id = ?"

	l.sqlQueryTags = "select tag, count(*) as `num` from `" + linkTable + "` where r_view = ? group by tag order by `num` desc"
	l.tplQueryTagsByOwner = "select `tag`, count(*) as `num` from `" + linkTable + "` "
	l.tplQueryTagsByOwnerFin = " group by `tag` order by `num` desc"

	l.tplQuery = "select linked_type, linked_id, `type`, tag, object_id from `" + linkTable + "`"

	l.sqlQueryByTag = "select linked_type, linked_id, `type`, tag, object_id from `" + linkTable + "` where tag = ?"
	l.sqlQueryByObjectID = "select linked_type, linked_id, `type`, tag, object_id from `" + linkTable + "` where object_id = ?"

	sqlStmts := []mysqllib.SqlStmt{
		{&l.stmtCreate, l.sqlCreate},
		{&l.stmtDelete, l.sqlDelete},

		{&l.stmtqueryByOwner, l.sqlqueryByOwner},
		{&l.stmtdeleteLinks, l.sqldeleteLinks},

		{&l.stmtQueryByTag, l.sqlQueryByTag},

		{&l.stmtQueryByObjectID, l.sqlQueryByObjectID},
	}

	for _, sqlStmt := range sqlStmts {
		if err = mysqllib.CreateStmt(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, errors.Wrap(err, onNew)
		}
	}

	return l, nil
}

const onClear = "on links_mysql.Clean()"

func (linksOp *linksMySQL) Clear(selector.Selector) error {
	if linksOp == nil {
		return nil
	}

	if _, ok := os.LookupEnv("ADMIN_MODE"); !ok {
		return errors.New(onClear + ": insufficient environment")
	}

	_, err := linksOp.dbh.Exec("truncate `" + linksOp.linkTable + "`")
	return err
}

func (linksOp *linksMySQL) Close() error {
	if linksOp == nil {
		return nil
	}

	return errors.Wrap(linksOp.dbh.Close(), "on linksMySQL.dbh.Close()")
}

const onSetLinks = "on linksMySQL.setLinks()"

func (linksOp *linksMySQL) SetLinks(userIS common.ID, linkedType, linkedID string, linksListTmp []links.Item) ([]links.LinkedInfo, error) {
	if linksOp == nil {
		return nil, nil
	}

	is := userIS

	linksList := links.Correct(is, linksListTmp)

	values := []interface{}{linkedType, linkedID, string(is)}
	ids, err := mysqllib.QueryIDs(linksOp.stmtqueryByOwner, linksOp.sqlqueryByOwner, values...)
	if err != nil {
		return nil, errors.Wrap(err, onSetLinks)
	}

	var idStrings []string
	for _, id := range ids {
		if id > 0 {
			idStrings = append(idStrings, strconv.FormatUint(id, 10))
		}
	}

	values = []interface{}{linkedType, linkedID}
	_, err = linksOp.stmtdeleteLinks.Exec(values...)
	if err != nil {
		return nil, errors.Wrapf(err, onSetLinks+": "+libs.CantExecQuery, linksOp.sqldeleteLinks, values)
	}

	for _, n := range linksList {
		// !!! another owned tags.comp must be filtered on linksListTmp --> linksList
		if !groups.OneOf(userIS, linksOp.ctrlOp, n.RView) {
			n.RView = is
		}

		idString := strings.TrimSpace(n.To)
		id, _ := strconv.ParseUint(idString, 10, 64)

		values := []interface{}{linkedType, linkedID, n.Type, n.Name, id, string(n.RView), string(n.ROwner)}
		_, err = linksOp.stmtCreate.Exec(values...)
		if err != nil {
			return nil, errors.Wrapf(err, onSetLinks+"linksOp: "+libs.CantExecQuery, linksOp.sqlCreate, values)
		}

		if id > 0 {
			duplicateID := false
			for _, d := range ids {
				if d == id {
					duplicateID = true
					break
				}
			}
			if !duplicateID {
				idStrings = append(idStrings, idString)
			}
		}
	}

	if len(idStrings) < 1 {
		return nil, nil
	}

	var linkedInfo []links.LinkedInfo

	// TODO: QueryAccessible
	sql := linksOp.tplCountLinked + " where r_view = '" + string(libs.Anyone) + "' and object_id in (" + strings.Join(idStrings, ",") + ") " + linksOp.tplCountLinkedFin
	rows, err := linksOp.dbh.Query(sql)
	if err != nil {
		return nil, errors.Wrapf(err, onSetLinks+": "+libs.CantPrepareQuery, sql, nil)
	}
	defer rows.Close()

	for rows.Next() {
		var objectID string
		var countLinked uint
		if err := rows.Scan(&objectID, &countLinked); err != nil {
			return linkedInfo, errors.Wrapf(err, onSetLinks+": "+mysqllib.CantScanQueryRow, sql, nil)
		}
		linkedInfo = append(linkedInfo, links.LinkedInfo{ObjectID: objectID, CountLinked: countLinked})
	}
	err = rows.Err()
	if err != nil {
		return linkedInfo, errors.Wrapf(err, onSetLinks+": "+mysqllib.CantScanQueryRow, sql, nil)
	}

	return linkedInfo, nil
}

const onQuery = "on linksMySQL.Query"

func (linksOp *linksMySQL) Query(userIS common.ID, selector selector.Selector) (linked []links.Linked, err error) {
	if linksOp == nil {
		return nil, nil
	}

	condition, values, err := selector.Mysql(userIS, selector)
	if err != nil {
		return nil, errors.Wrapf(err, onQuery+": on selector.Mysql(%#v)", selector)
	}

	if strings.TrimSpace(condition) != "" {
		condition = "where " + condition
	}
	sqlQuery := linksOp.tplQuery + " " + condition

	rows, err := linksOp.dbh.Query(sqlQuery, values...)
	if err != nil {
		return nil, errors.Wrapf(err, onQuery+": "+libs.CantExecQuery, sqlQuery, values)
	}
	defer rows.Close()

	for rows.Next() {
		var li links.Linked
		if err := rows.Scan(&li.LinkedType, &li.LinkedID, &li.Type, &li.Tag, &li.ObjectID); err != nil {
			return linked, errors.Wrapf(err, onQuery+": "+mysqllib.CantScanQueryRow, sqlQuery, values)
		}
		linked = append(linked, li)
	}
	err = rows.Err()
	if err != nil {
		return linked, errors.Wrapf(err, onQuery+": "+mysqllib.CantScanQueryRow, sqlQuery, values)
	}

	return linked, nil
}

const onQueryByTag = "on linksMySQL.QueryByTag"

func (linksOp *linksMySQL) QueryByTag(userIS common.ID, tag string) (linked []links.Linked, err error) {
	if linksOp == nil {
		return nil, nil
	}

	rows, err := linksOp.stmtQueryByTag.Query(tag)
	if err != nil {
		return nil, errors.Wrapf(err, onQueryByTag+": "+libs.CantExecQuery, linksOp.sqlQueryByTag, tag)
	}
	defer rows.Close()

	for rows.Next() {
		var li links.Linked
		if err := rows.Scan(&li.LinkedType, &li.LinkedID, &li.Type, &li.Tag, &li.ObjectID); err != nil {
			return linked, errors.Wrapf(err, onQueryByTag+": "+mysqllib.CantScanQueryRow, linksOp.sqlQueryByTag, tag)
		}
		linked = append(linked, li)
	}
	err = rows.Err()
	if err != nil {
		return linked, errors.Wrapf(err, onQueryByTag+": "+mysqllib.CantScanQueryRow, linksOp.sqlQueryByTag, tag)
	}

	return linked, nil
}

const onQueryByObjectID = "on linksMySQL.QueryByObjectID"

func (linksOp *linksMySQL) QueryByObjectID(userIS common.ID, id string) (linked []links.Linked, err error) {
	if linksOp == nil {
		return nil, nil
	}

	rows, err := linksOp.stmtQueryByObjectID.Query(strings.TrimSpace(id))
	if err != nil {
		return nil, errors.Wrapf(err, onQueryByObjectID+": "+libs.CantExecQuery, linksOp.sqlQueryByObjectID, id)
	}
	defer rows.Close()

	for rows.Next() {
		var li links.Linked
		if err := rows.Scan(&li.LinkedType, &li.LinkedID, &li.Type, &li.Tag, &li.ObjectID); err != nil {
			return linked, errors.Wrapf(err, onQueryByObjectID+": "+mysqllib.CantScanQueryRow, linksOp.sqlQueryByObjectID, id)
		}
		linked = append(linked, li)
	}
	err = rows.Err()
	if err != nil {
		return linked, errors.Wrapf(err, onQueryByObjectID+": "+mysqllib.CantScanQueryRow, linksOp.sqlQueryByObjectID, id)
	}

	return linked, nil
}

const onQueryTags = "on linksMySQL.QueryTags"

// QueryTags selects all tags.comp with selector accordingly to user's rights.
func (linksOp *linksMySQL) QueryTags(userIS common.ID, selector selector.Selector) ([]links.TagInfo, error) {
	if linksOp == nil {
		return nil, nil
	}

	condition, values, err := selector.Mysql(userIS, selector)
	if err != nil {
		return nil, errors.Wrapf(err, "on selector.Mysql(%#v)", selector)
	}

	sqlQueryTags := linksOp.tplQueryTags + "::" + condition + "::" + linksOp.tplQueryTagsFin

	_, rows, err := groups.QueryAccessible(linksOp.ctrlOp, linksOp.dbh, userIS, linksOp.tplQueryTags, condition, linksOp.tplQueryTagsFin, values)

	if err != nil {
		return nil, errors.Wrapf(err, onQueryTags+": "+libs.CantExecQuery, sqlQueryTags)
	}
	defer rows.Close()

	tagsAll := []links.TagInfo{}
	for rows.Next() {
		r := links.TagInfo{}
		if err := rows.Scan(&r.Tag, &r.Count); err != nil {
			return tagsAll, errors.Wrapf(err, onQueryTags+": "+mysqllib.CantScanQueryRow, sqlQueryTags)
		}
		tagsAll = append(tagsAll, r)
	}
	err = rows.Err()
	if err != nil {
		return tagsAll, errors.Wrapf(err, onQueryTags+": "+mysqllib.CantScanQueryRow, sqlQueryTags)
	}

	return tagsAll, nil
}

const onQueryTagsByOwner = "on linksMySQL.QueryTagsByOwner"

func (linksOp *linksMySQL) QueryTagsByOwner(userIS common.ID, rOwner common.ID) ([]links.TagInfo, error) {
	if linksOp == nil {
		return nil, nil
	}

	_, rows, err := groups.QueryAccessible(linksOp.ctrlOp, linksOp.dbh, userIS, linksOp.tplQueryTagsByOwner, " r_owner = ? ", linksOp.tplQueryTagsByOwnerFin, []interface{}{string(rOwner)})

	sql := linksOp.tplQueryTagsByOwner + "::" + linksOp.tplQueryTagsByOwnerFin
	if err != nil {
		return nil, errors.Wrapf(err, onQueryTagsByOwner+": "+libs.CantExecQuery, sql, rOwner)
	}
	defer rows.Close()

	tagsAll := []links.TagInfo{}
	for rows.Next() {
		r := links.TagInfo{}
		if err := rows.Scan(&r.Tag, &r.Count); err != nil {
			return tagsAll, errors.Wrapf(err, onQueryTagsByOwner+": "+mysqllib.CantScanQueryRow, sql, rOwner)
		}
		tagsAll = append(tagsAll, r)
	}
	err = rows.Err()
	if err != nil {
		return tagsAll, errors.Wrapf(err, onQueryTagsByOwner+": "+mysqllib.CantScanQueryRow, sql, rOwner)
	}

	return tagsAll, nil
}
