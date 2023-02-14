package records_sqlite

import (
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/sqllib"
)

var _ db.Cleaner = &recordsSQLite{}

const onClean = "on recordsSQLite.Clean(): "

func (recordsOp *recordsSQLite) Clean() error {
	if _, err := recordsOp.stmClean.Exec(); err != nil {
		return errors.Wrapf(err, onClean+sqllib.CantExec, recordsOp.sqlClean, nil)
	}
	return nil

	//var termTags *selectors.Term
	//
	//condition, values, err := selectors_sql.Use(nil)
	//if err != nil {
	//	return fmt.Errorf(onClean+"wrong selector: %s", err)
	//}

	//if strings.TrimSpace(condition) != "" {
	//	ids, err := recordsOp.ids(condition, values)
	//	if err != nil {
	//		return errors.Wrap(err, onClean+"can't recordsOp.ids(condition, values)")
	//	}
	//	termTags = logic.AND(selectors.In("key", recordsOp.interfaceKey), selectors.In("id", ids...))
	//
	//	query += " WHERE " + condition
	//
	//} else {
	//	termTags = selectors.In("key", recordsOp.interfaceKey) // TODO!!! correct field key
	//
	//}

	//if recordsOp.taggerCleaner != nil {
	//	err = recordsOp.taggerCleaner.Clean(termTags, nil)
	//	if err != nil {
	//		return errors.Wrap(err, onClean)
	//	}
	//}

}

//const onIDs = "on recordsSQLite.IDs()"
//
//func (dataOp *recordsSQLite) ids(condition string, values []interface{}) ([]interface{}, error) {
//	if strings.TrimSpace(condition) != "" {
//		condition = " WHERE " + condition
//	}
//
//	query := "SELECT id FROM " + dataOp.table + condition
//	stm, err := dataOp.db.Prepare(query)
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
