package flowmysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pavlo67/punctum/interfaces"
	"github.com/pavlo67/punctum/interfaces/confidenter"
	"github.com/pavlo67/punctum/interfaces/controller"
	"github.com/pavlo67/punctum/interfaces/controller/rights"
	"github.com/pavlo67/punctum/interfaces/flow"

	"encoding/json"
	"log"
	"strconv"

	"github.com/pavlo67/punctum/basis/clients"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/crud"
	"github.com/pavlo67/punctum/interfaces/selectors"
	"github.com/pkg/errors"
)

var _ flow.Operator = &FlowMySQL{}

const MaxVarcharLen = 255

type FlowMySQL struct {
	ctrl         controller.Operator
	dbh          *sql.DB
	itemTable    string
	stmtCreate   *sql.Stmt
	stmtRead     *sql.Stmt
	stmtUpdate   *sql.Stmt
	stmtIsNew    *sql.Stmt
	stmtDelete   *sql.Stmt
	stmtImportTo *sql.Stmt
	sqlReadAll   string
	dataManagers controller.Managers
	crudBuffer   *flow.Item

	sqlCreate, sqlRead, sqlUpdate, sqlDelete, sqlIsNew, sqlImportTo string
}

// NewFlowMySQL ...
func NewFlowMySQL(ctrlOp controller.Operator, mysqlConfig clients.MySQLConfig, itemTable string, dataManagers controller.Managers) (*FlowMySQL, error) {
	dbh, err := clients.ConnectToMysql(mysqlConfig)
	if err != nil {
		return nil, err
	}

	flowFields := "`r_view`, `r_owner`, `fount_is`, `original_id`, `fount_url`, `url`, `title`, `summary`, `content`, `original`, `media`, `created_at`, `imported_to`"

	f := FlowMySQL{
		ctrl:         ctrlOp,
		dbh:          dbh,
		itemTable:    itemTable,
		dataManagers: dataManagers,
		sqlReadAll:   "select SQL_CALC_FOUND_ROWS `id`, " + flowFields + " from `" + itemTable + "` ",
	}

	f.sqlCreate = "insert into `" + itemTable + "` (" + flowFields + ") values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), '')"
	f.sqlRead = "select `id`, " + flowFields + " from `" + itemTable + "` where `id`=?"
	f.sqlUpdate = "update `" + itemTable + "` set `r_view` = ?, `r_owner` = ?, `fount_is` = ?, `original_id` = ?, `fount_url` = ?, `url` = ?, `title` = ?, `summary` = ?, `content` = ?, `original` = ?, `media` = ?, `imported_to` = ?  where `id`=?"
	f.sqlDelete = "delete from `" + itemTable + "` where `id`=?"
	f.sqlIsNew = "select count(*) as cnt from `" + itemTable + "` where `r_owner` = ? and `fount_is` = ? and `original_id` = ?"
	//f.sqlImportTo = "update `" + itemTable + "` set `imported_to` = concat(imported_to, ?)  where `id`=? and r_view=?"
	f.sqlImportTo = "update `" + itemTable + "` set `imported_to` = ?  where `id`=? and r_view=?"

	sqlStmts := []clients.SqlStmt{
		{&f.stmtCreate, f.sqlCreate},
		{&f.stmtRead, f.sqlRead},
		{&f.stmtUpdate, f.sqlUpdate},
		{&f.stmtDelete, f.sqlDelete},
		{&f.stmtIsNew, f.sqlIsNew},
		{&f.stmtImportTo, f.sqlImportTo},
	}
	for _, sqlStmt := range sqlStmts {
		if err = clients.CreateStmt(dbh, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, err
		}
	}
	if len(f.dataManagers) == 0 {
		f.dataManagers = controller.Managers{rights.Create: controller.Anyone}
	}
	return &f, nil
}

// Create ...
func (f *FlowMySQL) Create(identity *confidenter.Identity, item flow.Item) (confidenter.Identity, error) {
	var original, media []byte
	var err error
	m := controller.Managers{rights.View: item.RView, rights.Owner: item.ROwner}
	rView, rOwner, _, _, err := controller.SetRights(identity, f.ctrl, f.dataManagers, m)
	if err != nil {
		return confidenter.Identity{}, errors.Wrap(err, "can't .SetRights)")
	}
	if item.Original != nil {
		original, err = json.Marshal(item.Original)
		if err != nil {
			return confidenter.Identity{}, errors.Wrapf(err, "can't marshal item.Original: %v in flow.Create", item.Original)
		}
	}
	if item.Media != nil {
		media, err = json.Marshal(item.Media)
		if err != nil {
			return confidenter.Identity{}, errors.Wrapf(err, "can't marshal item.Media: %v in flow.Create", item.Media)
		}
	}

	if len([]rune(item.Title)) > MaxVarcharLen {
		item.Title = string([]rune(item.Title)[:MaxVarcharLen])
	}
	if len([]rune(item.Summary)) > MaxVarcharLen {
		item.Summary = string([]rune(item.Summary)[:MaxVarcharLen])
	}

	sqlValues := []interface{}{string(rView), string(rOwner), string(item.FountIS), item.OriginalID, item.FountURL, item.URL, item.Title, item.Summary, item.Content, string(original), string(media)}
	res, err := f.stmtCreate.Exec(sqlValues...)
	if err != nil {
		return confidenter.Identity{}, errors.Wrapf(err, "can't exec SQL: %s, %v", f.sqlCreate, sqlValues)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return confidenter.Identity{}, errors.Wrapf(err, "can't get LastInsertId() SQL: %s, %v", f.sqlCreate, sqlValues)
	}
	return confidenter.Identity{
		Domain: program.Domain(),
		Path:   "flow",
		ID:     strconv.FormatInt(id, 10),
	}, nil
}

func (f *FlowMySQL) IsNew(item flow.Item) (bool, error) {
	// TODO: do it correctly!!!
	var is int
	err := f.stmtIsNew.QueryRow(string(item.ROwner), string(item.FountIS), item.OriginalID).Scan(&is)
	if err != nil {
		return false, errors.Wrapf(err, "can't exec QueryRow: %s, values=%v, %v, %v", f.sqlIsNew, item.ROwner, item.FountIS, item.OriginalID)
	}
	if is > 0 {
		return false, nil
	}
	return true, nil
}

// Read
func (f *FlowMySQL) Read(identity *confidenter.Identity, itemID confidenter.IdentityString) (*flow.Item, error) {
	//itemIdentity := itemIS.Identity()
	var item flow.Item
	var media []byte
	err := f.stmtRead.QueryRow(string(itemID)).Scan(&item.ID, &item.RView, &item.ROwner, &item.FountIS, &item.OriginalID, &item.FountURL, &item.URL, &item.Title, &item.Summary, &item.Content, &item.Original, &media, &item.CreatedAt, &item.ImportedTo)
	if err == sql.ErrNoRows {
		return nil, errors.New("item not found")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "can't exec QueryRow: %s, is=%s", f.sqlRead, itemID)
	}
	if string(media) != "" {
		err = json.Unmarshal(media, &item.Media)
		if err != nil {
			return nil, errors.Wrapf(err, "can't unmarshal flow.media: %s, is=%s", media, itemID)
		}
	}
	err = controller.IsManager(identity, f.ctrl, controller.Managers{rights.View: item.RView}, rights.View)
	if err != nil {
		return nil, errors.Wrapf(err, "can't confirm rights to read item (identity = %+v)", identity)
	}
	return &item, nil
}

// Update ...
func (f *FlowMySQL) Update(identity *confidenter.Identity, itemIS confidenter.IdentityString, item flow.Item) (crud.Result, error) {

	var i *flow.Item
	var err error
	i, err = f.Read(identity, itemIS)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, "can't .Read()")
	}
	rView, rOwner, _, err := controller.CheckAndUpdateRights(identity, f.ctrl, string(i.RView), string(i.ROwner), "", f.dataManagers, nil)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, "can't .CheckAndUpdateRights()")
	}
	var original, media []byte
	if item.Original != nil {
		original, err = json.Marshal(item.Original)
		if err != nil {
			log.Println("can't marshal item.Original:", item.Original, "in flow.Update")
		}
	}
	if item.Media != nil {
		media, err = json.Marshal(item.Media)
		if err != nil {
			return crud.Result{}, errors.Wrapf(err, "can't marshal item.Media: %v in flow.Update", item.Media)
		}
	}

	if len([]rune(item.Title)) > MaxVarcharLen {
		item.Title = string([]rune(item.Title)[:MaxVarcharLen])
	}
	if len([]rune(item.Summary)) > MaxVarcharLen {
		item.Summary = string([]rune(item.Summary)[:MaxVarcharLen])
	}

	values := []interface{}{string(rView), string(rOwner), string(item.FountIS), item.OriginalID, item.FountURL, item.URL, item.Title, item.Summary, item.Content, string(original), string(media), item.ImportedTo, i.ID}
	res, err := f.stmtUpdate.Exec(values...)

	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't exec sql: %s (%v)", f.sqlUpdate, values)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't get RowsAffected(): %s (%v)", f.sqlUpdate, values)
	}
	return crud.Result{NumOk: cnt}, nil
}

// Delete ...
func (f *FlowMySQL) Delete(identity *confidenter.Identity, itemIS confidenter.IdentityString) (crud.Result, error) {

	i, err := f.Read(identity, itemIS)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, "can't .Read()")
	}
	if !controller.CanDelete(identity, f.ctrl, i.ROwner, "") {
		return crud.Result{}, errors.New("no rights to .Delete()")
	}

	res, err := f.stmtDelete.Exec(i.ID)
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't exec SQL: %s, %s", f.sqlDelete, i.ID)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrap(err, "can't get RowsAffected()")
	}
	if cnt < 1 {
		return crud.Result{}, errors.Errorf("can't delete item (id, cnt): %v, %d", i.ID, cnt)
	}

	return crud.Result{NumOk: cnt}, nil
}

// ReadAll ...
func (f *FlowMySQL) ReadAll(userIdentity *confidenter.Identity, options *crud.ReadAllOptions, selector interfaces.Selector) ([]flow.Item, int64, error) {

	condition, values, err := selectors.Mysql(userIdentity, selector)
	if err != nil {
		return nil, 0, errors.Wrapf(err, ": bad selector ('%v')", selector)
	}
	for k, v := range values {
		if value, ok := v.(confidenter.IdentityString); ok {
			values[k] = string(value)
		}
	}
	var order string
	var rows *sql.Rows
	if options != nil {
		order = clients.GetOrderAndLimit(options.SortBy, options.Limits)
	} else {
		order = " order by created_at desc "
	}
	sqlQuery := f.sqlReadAll + " where " + condition + " " + order
	valuesAll := append([]interface{}{}, values...)
	rows, err = f.ctrl.QueryAccessible(f.dbh, userIdentity.String(), f.sqlReadAll, condition, order, valuesAll)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, errors.New("items not found")
		}
		return nil, 0, errors.Wrapf(err, ": can't get query (sql='%v', values='%v')", sqlQuery, valuesAll)
	}
	itemAll := []flow.Item{}
	var allCount int64
	for rows.Next() {
		r := flow.Item{}
		var media []byte
		if err := rows.Scan(&r.ID, &r.RView, &r.ROwner, &r.FountIS, &r.OriginalID, &r.FountURL, &r.URL, &r.Title, &r.Summary, &r.Content, &r.Original, &media, &r.CreatedAt, &r.ImportedTo); err != nil {
			return itemAll, 0, errors.Wrapf(err, ": can't scan queryRow (sql='%v', values='%v')", f.sqlReadAll+condition, valuesAll)
		}
		if string(media) != "" {
			err = json.Unmarshal(media, &r.Media)
			if err != nil {
				return itemAll, 0, errors.Wrapf(err, "can't unmarshal flow.media: %s, id=%s", media, r.ID)
			}
		}

		itemAll = append(itemAll, r)
	}
	stmtAll, err := f.dbh.Prepare("SELECT FOUND_ROWS()")
	if err != nil {
		return nil, 0, errors.Wrapf(err, "can't prepare ('SELECT FOUND_ROWS()') for sql=", sqlQuery)
	}
	defer stmtAll.Close()

	err = stmtAll.QueryRow().Scan(&allCount)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "can't scan ('SELECT FOUND_ROWS()') for sql=", sqlQuery)
	}
	return itemAll, allCount, nil
}

func (f *FlowMySQL) Close() {
	f.dbh.Close()
}

func (f *FlowMySQL) ImportTo(identity *confidenter.Identity, id int64, importIS string) error {

	_, err := f.stmtImportTo.Exec(importIS, id, string(identity.String()))
	if err != nil {
		return errors.Wrapf(err, "can't exec sql:%s, values: %v, %v, %v", f.sqlImportTo, importIS, id, identity.String())
	}
	return nil
}
