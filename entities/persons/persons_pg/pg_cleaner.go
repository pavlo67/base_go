package persons_pg

//import (
//	"usersbase/sql"
//	"strings"
//
//	"github.com/pkg/errors"
//
//	"github.com/pavlo67/common/common"
//	"github.com/pavlo67/common/common/crud"
//	"github.com/pavlo67/common/common/libraries/sqllib"
//	"github.com/pavlo67/common/common/selectors"
//	"github.com/pavlo67/common/common/selectors/logic"
//	"github.com/pavlo67/common/common/selectors/selectors_sql"
//)
//
//var _ crud.Cleaner = &usersPg{}
//
//const onIDs = "on usersPg.IDs()"
//
//func (usersOp *usersPg) ids(condition string, values []interface{}) ([]interface{}, error) {
//	if strings.TrimSpace(condition) != "" {
//		condition = " WHERE " + condition
//	}
//
//	query := "SELECT id FROM " + usersOp.table + condition
//	stm, err := usersOp.db.Prepare(query)
//	if err != nil {
//		return nil, errors.Wrapf(err, onIDs+": can't db.Prepare(%s)", query)
//	}
//
//	rows, err := stm.Query(values...)
//	if err == sql.ErrNoRows {
//		return nil, nil
//	} else if err != nil {
//		return nil, errors.Wrapf(err, onIDs+sqllib.CantQuery, query, values)
//	}
//	defer rows.Close()
//
//	var ids []interface{}
//
//	for rows.Next() {
//		var id common.ID
//
//		err := rows.Scan(&id)
//		if err != nil {
//			return ids, errors.Wrapf(err, onIDs+sqllib.CantScanQueryRow, query, values)
//		}
//
//		ids = append(ids, id)
//	}
//	err = rows.Err()
//	if err != nil {
//		return ids, errors.Wrapf(err, onIDs+": "+sqllib.RowsError, query, values)
//	}
//
//	return ids, nil
//}
//
//const onClean = "on usersPg.Clean(): "
//
//func (usersOp *usersPg) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
//	var termTags *selectors.Term
//
//	condition, values, err := selectors_sql.Use(term)
//	if err != nil {
//		return errors.Errorf(onClean+"wrong selector (%#v): %s", term, err)
//	}
//
//	query := usersOp.sqlClean
//
//	if strings.TrimSpace(condition) != "" {
//		ids, err := usersOp.ids(condition, values)
//		if err != nil {
//			return errors.Wrap(err, onClean+"can't usersOp.ids(condition, values)")
//		}
//		termTags = logic.AND(selectors.In("key", usersOp.interfaceKey), selectors.In("id", ids...))
//
//		query += " WHERE " + condition
//
//	} else {
//		termTags = selectors.In("joiner_key", usersOp.interfaceKey) // TODO!!! correct field key
//
//	}
//
//	_, err = usersOp.db.Exec(query, values...)
//	if err != nil {
//		return errors.Wrapf(err, onClean+sqllib.CantExec, query, values)
//	}
//
//	if usersOp.taggerCleaner != nil {
//		err = usersOp.taggerCleaner.Clean(termTags, nil)
//		if err != nil {
//			return errors.Wrap(err, onClean)
//		}
//	}
//
//	return err
//}
//
//const onSelectToClean = "on usersPg.SelectToClean(): "
//
//func (usersOp *usersPg) SelectToClean(options *crud.RemoveOptions) (*selectors.Term, error) {
//	var limit uint64
//
//	if options != nil && options.Limit > 0 {
//		limit = options.Limit
//	} else {
//		return nil, errors.New(onSelectToClean + "no clean limit is defined")
//	}
//
//	queryMax := "SELECT MAX(id) from " + usersOp.table
//
//	var maxID uint64
//	row := usersOp.db.QueryRow(queryMax)
//
//	err := row.Scan(&maxID)
//	if err != nil {
//		return nil, errors.Errorf(onSelectToClean+": error on query (%s)", queryMax)
//	}
//
//	return selectors.Binary(selectors.Le, "id", selectors.Value{V: maxID - limit}), nil
//}
