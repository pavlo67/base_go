package fountsmysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pavlo67/punctum/interfaces/confidenter"
	"github.com/pavlo67/punctum/interfaces/controller"
	"github.com/pavlo67/punctum/interfaces/controller/rights"
	"github.com/pavlo67/punctum/interfaces/founts"

	"strconv"
	"strings"

	"log"

	"github.com/pavlo67/punctum/basis/clients"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces"
	"github.com/pavlo67/punctum/interfaces/crud"
	"github.com/pavlo67/punctum/interfaces/selectors"
	"github.com/pkg/errors"
)

type MySQLFount struct {
	ctrl         controller.Operator
	dataManagers controller.Managers
	dbh          *sql.DB
	crudBuffer   *founts.Fount

	fountTable     string
	fountStatTable string
	scanStatTable  string

	stmtCreate *sql.Stmt
	stmtRead   *sql.Stmt
	stmtUpdate *sql.Stmt
	stmtDelete *sql.Stmt

	stmtSettings   *sql.Stmt
	stmtAddTag     *sql.Stmt
	stmtRemoveTags *sql.Stmt
	//stmtTagFounts  *sql.Stmt
	sqlReadAll  string
	sqlReadTags string

	stmtStatCreate *sql.Stmt

	stmtScanCreate *sql.Stmt
}

var sqlCreate, sqlRead, sqlUpdate, sqlDelete, sqlStatCreate, sqlScannerCreate, sqlAddTag, sqlRemoveTags, sqlExportSettings string

// Newfountsmysql ...
func NewMySQLFount(ctrl controller.Operator, mysqlConfig clients.MySQLConfig, fountTable, fountTagsTable, fountStatTable, scannerStatTable string, dataManagers controller.Managers) (*MySQLFount, error) {

	dbh, err := clients.ConnectToMysql(mysqlConfig)
	if err != nil {
		return nil, err
	}

	fountFields := "`url`, `title`, `import_type`, `to_flow`, `to_object`, `tags`, `r_owner`, `r_view`, `managers`, `import_details_type`, `import_details_params`, `created_at`"
	sqlCreate = "insert into `" + fountTable + "` (" + fountFields + ") values (?,?,?,?,?,?,?,?,?,?,?,NOW())"
	sqlRead = "select `id`, " + fountFields + ", `updated_at` from `" + fountTable + "` where `id`=?"
	sqlUpdate = "update `" + fountTable + "` set `url`=?, `title`=?, `import_type`=?, `to_flow`=?, `to_object`=?, `tags`=?, `r_owner`=?, `r_view`=?, `managers`=?, `import_details_type`=?, `import_details_params`=? where `id`=? "
	sqlDelete = "delete from `" + fountTable + "` where `id`=?"
	sqlExportSettings = "update `" + fountTable + "` set `import_type`=?, `import_details_params`=? where url=?"

	sqlAddTag = "insert into `" + fountTagsTable + "` (`fount_id`, `tag`, `r_view`) values (?,?,?)"
	sqlRemoveTags = "delete from `" + fountTagsTable + "` where fount_id = ?"
	//sqlTagFounts = "select `fount_id` from `" + fountTagsTable + "` where r_view=? and tag=?"

	fountStatFields := "`scanner_start`, `fount_id`, `start`, `duration`, `response_error`, `last_item_error`, `item_errors`, `items_taken`, `items_new`"
	sqlStatCreate = "insert into `" + fountStatTable + "` (" + fountStatFields + ") values (?,?,?,?,?,?,?,?,?)"

	scannerStatFields := "`start`, `duration`, `founts_num`, `errors_num`, `items_taken`, `items_new`"
	sqlScannerCreate = "insert into `" + scannerStatTable + "` (" + scannerStatFields + ") values (?,?,?,?,?,?)"

	f := MySQLFount{
		ctrl:           ctrl,
		dbh:            dbh,
		fountTable:     fountTable,
		fountStatTable: fountStatTable,
		scanStatTable:  scannerStatTable,
		dataManagers:   dataManagers,
		sqlReadAll:     "select SQL_CALC_FOUND_ROWS `id`, " + fountFields + ", `updated_at` from `" + fountTable + "`",
		sqlReadTags:    "select `id`, `fount_id`, `tag`, `r_view` from `" + fountTagsTable + "` ",
	}

	sqlStmts := []clients.SqlStmt{
		{&f.stmtCreate, sqlCreate},
		{&f.stmtRead, sqlRead},
		{&f.stmtUpdate, sqlUpdate},
		{&f.stmtDelete, sqlDelete},
		{&f.stmtStatCreate, sqlStatCreate},
		{&f.stmtScanCreate, sqlScannerCreate},
		{&f.stmtAddTag, sqlAddTag},
		{&f.stmtRemoveTags, sqlRemoveTags},
		{&f.stmtSettings, sqlExportSettings},
		//{&f.stmtTagFounts, sqlTagFounts},
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
func (fms *MySQLFount) Create(identity *confidenter.Identity, newFount founts.Fount) (*confidenter.Identity, error) {

	newFount.URL = strings.TrimSpace(newFount.URL)
	if newFount.URL == "" {
		return nil, errors.New("newFount.URL is empty")
	}

	rView, rOwner, _, managers, err := controller.SetRights(identity, fms.ctrl, fms.dataManagers, newFount.Managers)
	if err != nil {
		return nil, errors.Wrap(err, "can't .SetRights)")
	}

	sqlValues := []interface{}{newFount.URL, newFount.Title, string(newFount.ImportType), newFount.ToFlow, newFount.ToObject, newFount.Tags, string(rOwner), string(rView), managers, newFount.ImportDetailsType, newFount.ImportDetailsParams}
	res, err := fms.stmtCreate.Exec(sqlValues...)
	if err != nil {
		return nil, errors.Wrapf(err, "can't exec SQL: %s, %v", sqlCreate, sqlValues)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrapf(err, "can't get LastInsertId() SQL: %s, %v", sqlCreate, sqlValues)
	}
	for _, v := range strings.Split(newFount.Tags, ";") {
		v = strings.Trim(v, " ")
		if v != "" {
			_, err = fms.stmtAddTag.Exec(id, v, string(rView))
			if err != nil {
				log.Println("can't add fount's tag:", id, v, rView, err)
			}
		}
	}

	return &confidenter.Identity{program.Domain(), "fount", strconv.FormatInt(id, 10), ""}, nil
}

const onRead = "on MySQLFount.Read"

// Read
func (fms *MySQLFount) Read(identity *confidenter.Identity, toRead string) (*founts.Fount, error) {
	var f founts.Fount
	id, err := strconv.ParseUint(toRead, 10, 64)
	if err != nil {
		return nil, errors.Wrap(crud.ErrBadSelector, onRead+": "+toRead)
	}

	err = fms.stmtRead.QueryRow(id).Scan(&f.ID, &f.URL, &f.Title, &f.ImportType, &f.ToFlow, &f.ToObject, &f.Tags, &f.ROwner, &f.RView, &f.ManagersRaw, &f.ImportDetailsType, &f.ImportDetailsParams, &f.CreatedAt, &f.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, interfaces.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, "can't exec QueryRow: %s, id=%s", sqlRead, id)
	}

	err = controller.IsManager(identity, fms.ctrl, controller.Managers{rights.View: f.RView}, rights.View)
	if err != nil {
		return nil, errors.Wrapf(err, "can't confirm rights to read fount (identity = %+v) (rights.View=%v)", identity, f.RView)
	}
	return &f, nil
}

// Update ...
func (fms *MySQLFount) Update(identity *confidenter.Identity, fountNew founts.Fount) (crud.Result, error) {
	f, err := fms.Read(identity, strconv.FormatInt(fountNew.ID, 10))
	if err != nil {
		return crud.Result{}, errors.Wrap(err, "can't .Read()")
	}
	fountNew.URL = strings.TrimSpace(fountNew.URL)

	rView, rOwner, managers, err := controller.CheckAndUpdateRights(identity, fms.ctrl, string(f.RView), string(f.ROwner), f.ManagersRaw, fms.dataManagers, fountNew.Managers)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, "can't .CheckAndUpdateRights()")
	}

	values := []interface{}{fountNew.URL, fountNew.Title, string(fountNew.ImportType), fountNew.ToFlow, fountNew.ToObject, fountNew.Tags, string(rView), string(rOwner), managers, fountNew.ImportDetailsType, fountNew.ImportDetailsParams, f.ID}
	res, err := fms.stmtUpdate.Exec(values...)

	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't exec sql: %s (%v)", sqlUpdate, values)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't get RowsAffected(): %s (%v)", sqlUpdate, values)
	}

	//delete old tags
	_, err = fms.stmtRemoveTags.Exec(f.ID)
	if err != nil {
		log.Println("can't remove tags for fount:", f.ID, err)
	}
	// add new tags
	for _, v := range strings.Split(fountNew.Tags, ";") {
		v = strings.Trim(v, " ")
		if v != "" {
			_, err = fms.stmtAddTag.Exec(f.ID, v, string(rView))
			if err != nil {
				log.Println("can't add fount's tag:", f.ID, v, rView, err)
			}
		}
	}

	return crud.Result{NumOk: cnt}, nil
}

// Delete ...
func (fms *MySQLFount) Delete(identity *confidenter.Identity, toDelete string) (crud.Result, error) {
	f, err := fms.Read(identity, toDelete)
	if err != nil {
		return crud.Result{}, errors.Wrap(err, "can't .Read()")
	}
	if !controller.CanDelete(identity, fms.ctrl, f.ROwner, f.ManagersRaw) {
		return crud.Result{}, errors.New("no rights to .Delete()")
	}

	res, err := fms.stmtDelete.Exec(f.ID)

	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't exec SQL: %s, %s", sqlDelete, f.ID)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return crud.Result{}, errors.Wrapf(err, "can't get RowsAffected(): %s, %s", sqlDelete, f.ID)
	}
	_, err = fms.stmtRemoveTags.Exec(f.ID)
	if err != nil {
		log.Println("can't remove tags for fount:", f.ID, err)
	}
	return crud.Result{cnt}, nil
}

// ReadAll ...
func (fms *MySQLFount) ReadAll(userIdentity *confidenter.Identity, options *crud.ReadAllOptions, sel interfaces.Selector) ([]founts.Fount, int64, error) {
	condition, values, err := selectors.Mysql(userIdentity, sel)
	if err != nil {
		return nil, 0, errors.Wrapf(err, ": bad selector ('%v')", sel)
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
	}
	sqlQuery := fms.sqlReadAll + " where " + condition + " " + order
	valuesAll := append([]interface{}{}, values...)
	rector := program.Identity()
	if userIdentity.String() == rector.String() {
		// rector get all founts for scanner
		stmt, err := fms.dbh.Prepare(sqlQuery)
		if err != nil {
			return nil, 0, errors.Wrapf(err, "can't prepare sql:", sqlQuery)
		}
		defer stmt.Close()

		rows, err = stmt.Query(valuesAll...)
	} else {
		rows, err = fms.ctrl.QueryAccessible(fms.dbh, userIdentity.String(), fms.sqlReadAll, condition, order, valuesAll)
	}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, errors.New("founts not found")
		}
		return nil, 0, errors.Wrapf(err, ": can't get query (sql='%v', values='%v')", sqlQuery, valuesAll)
	}
	fountAll := []founts.Fount{}
	for rows.Next() {
		r := founts.Fount{}
		if err := rows.Scan(&r.ID, &r.URL, &r.Title, &r.ImportType, &r.ToFlow, &r.ToObject, &r.Tags, &r.ROwner, &r.RView, &r.ManagersRaw, &r.ImportDetailsType, &r.ImportDetailsParams, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return fountAll, 0, errors.Wrapf(err, ": can't scan queryRow (sql='%v', values='%v')", sqlQuery, valuesAll)
		}
		fountAll = append(fountAll, r)
	}
	var allCount int64
	stmtAll, err := fms.dbh.Prepare("SELECT FOUND_ROWS()")
	if err != nil {
		return nil, 0, errors.Wrapf(err, "can't prepare ('SELECT FOUND_ROWS()') for sql=", sqlQuery)
	}
	defer stmtAll.Close()

	err = stmtAll.QueryRow().Scan(&allCount)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "can't scan ('SELECT FOUND_ROWS()') for sql=", sqlQuery)
	}
	return fountAll, allCount, nil
}

func (fms *MySQLFount) Close() {
	fms.dbh.Close()
}

// Create Stat...
func (fms *MySQLFount) AddFountStat(identity *confidenter.Identity, stat founts.FountStat) error {
	if len(stat.LastItemError) > 255 {
		stat.LastItemError = stat.LastItemError[0:255]
	}
	//fountStatFields :=                 "`scanner_start`,   `fount_id`,    `start`,    `duration`,   `response_error`,  `last_item_error`,   `item_errors`,   `items_taken`,   `items_new`"
	sqlValues := []interface{}{stat.ScannerStart, stat.FountID, stat.Start, stat.Duration, stat.ResponseError, stat.LastItemError, stat.ItemErrors, stat.ItemsTaken, stat.ItemsNew}
	res, err := fms.stmtStatCreate.Exec(stat.ScannerStart, stat.FountID, stat.Start, stat.Duration, stat.ResponseError, stat.LastItemError, stat.ItemErrors, stat.ItemsTaken, stat.ItemsNew)
	if err != nil {
		return errors.Wrapf(err, "can't exec SQL: %s, %v", sqlStatCreate, sqlValues)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return errors.Wrapf(err, "can't get LastInsertId() SQL:  %s, %v", sqlStatCreate, sqlValues)
	}
	return nil
}

func (fms *MySQLFount) ReadAllFountStat(identity *confidenter.Identity, options *crud.ReadAllOptions, selector interfaces.Selector) ([]founts.FountStat, int64, error) {
	return nil, 0, nil
}

func (fms *MySQLFount) DeleteAllFountStat(identity *confidenter.Identity, selector interfaces.Selector) (crud.Result, error) {
	return crud.Result{}, nil
}

func (fms *MySQLFount) AddScannerStat(identity *confidenter.Identity, stat founts.ScannerStat) error {
	//scannerStatFields := "             `start`,    `duration`,   `founts_num`,   `errors_num`,   `items_taken`,   `items_new`"
	sqlValues := []interface{}{stat.Start, stat.Duration, stat.FountsNum, stat.ErrorsNum, stat.ItemsTaken, stat.ItemsNew}
	_, err := fms.stmtScanCreate.Exec(stat.Start, stat.Duration, stat.FountsNum, stat.ErrorsNum, stat.ItemsTaken, stat.ItemsNew)
	if err != nil {
		return errors.Wrapf(err, "can't exec SQL: %s, %v", sqlScannerCreate, sqlValues)
	}
	return nil
}

func (fms *MySQLFount) ReadAllScannerStat(identity *confidenter.Identity, options *crud.ReadAllOptions, selector interfaces.Selector) ([]founts.ScannerStat, int64, error) {
	return nil, 0, nil
}

func (fms *MySQLFount) DeleteAllScannerStat(identity *confidenter.Identity, selector interfaces.Selector) (crud.Result, error) {
	return crud.Result{}, nil
}

//func (fms *MySQLFount) GetFountsForTag(identity *confidenter.Identity, tag string) ([]int64, error) {
//
//	var fountsID []int64
//	rows, err := fms.stmtTagFounts.Query(string(identity.String()), tag)
//	if err != nil {
//		return nil, errors.Wrapf(err, "can't query sql:%v; values:%v, %v", sqlTagFounts, identity.String(), tag)
//	}
//	if rows != nil {
//		defer rows.Close()
//	}
//	for rows.Next() {
//		var id int64
//		err = rows.Scan(&id)
//		if err != nil {
//			return nil, errors.Wrapf(err, "can't scan sql:%v; values:%v, %v", sqlTagFounts, identity.String(), tag)
//		}
//		fountsID = append(fountsID, id)
//	}
//	return fountsID, nil
//}

func (fms *MySQLFount) ReadTags(userIdentity *confidenter.Identity, sel interfaces.Selector) ([]founts.FountTag, error) {

	var tags []founts.FountTag
	condition, values, err := selectors.Mysql(userIdentity, sel)
	if err != nil {
		return nil, errors.Wrapf(err, ": bad selector ('%v')", sel)
	}
	for k, v := range values {
		if value, ok := v.(confidenter.IdentityString); ok {
			values[k] = string(value)
		}
	}
	sqlQuery := fms.sqlReadAll + " where " + condition
	valuesAll := append([]interface{}{}, values...)
	rows, err := fms.ctrl.QueryAccessible(fms.dbh, userIdentity.String(), fms.sqlReadTags, condition, " order by tag ", valuesAll)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("fount tags not found")
		}
		return nil, errors.Wrapf(err, ": can't get query (sql='%v', values='%v')", sqlQuery, valuesAll)
	}
	for rows.Next() {
		r := founts.FountTag{}
		var rView confidenter.IdentityString
		if err := rows.Scan(&r.ID, &r.FountID, &r.Tag, &rView); err != nil {
			return tags, errors.Wrapf(err, ": can't scan queryRow (sql='%v', values='%v')", sqlQuery, valuesAll)
		}
		r.RView = rView.Identity()
		tags = append(tags, r)
	}
	return tags, nil
}

func (fms *MySQLFount) ExportSettings(url, fountType, importParams string) (int64, error) {

	res, err := fms.stmtSettings.Exec(fountType, importParams, url)
	if err != nil {
		return 0, errors.Wrapf(err, "can't exec sql=%v values: %v, %v, %v", sqlExportSettings, fountType, importParams, url)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "can't get RowsAffected() sql=%v values: %v, %v, %v", sqlExportSettings, fountType, importParams, url)
	}
	return cnt, nil
}
