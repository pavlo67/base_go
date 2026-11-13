package sqllib

import (
	"database/sql"
	"strings"

	"github.com/pavlo67/common/common/auth"

	"github.com/pkg/errors"
)

type CorrectWildcards func(query string) string

const CantPrepare = "can't .Prepare(%s)"
const CantQuery = "can't .Query('%s', %#v)"
const CantExec = "can't .Exec('%s', %#v)"

const CantGetLastInsertId = "can't .LastInsertId('%s', %#v)"
const CantGetRowsAffected = "can't .RowsAffected('%s', %#v)"
const NoRowOnQuery = "no row on query('%s', %#v)"
const CantScanQueryRow = "can't scan query row ('%s', %#v)"
const RowsError = "error on .Rows ('%s', %#v)"

var ErrNoTable = errors.New("table doesn't exist")

func SQLList(table, fields, condition string, identity *auth.Identity) string {
	if strings.TrimSpace(condition) != "" {
		condition = " WHERE " + condition
	}

	var limit, order string

	//ranges := options.GetRanges()
	//
	//if ranges != nil {
	//	if len(ranges.OrderBy) > 0 {
	//		order = strings.Join(ranges.OrderBy, ", ")
	//	}
	//
	//	if ranges.Offset+ranges.Limit > 0 {
	//		if ranges.Limit > 0 {
	//			limit += " LIMIT " + strconv.FormatUint(ranges.Limit, 10)
	//		}
	//
	//		if ranges.Offset > 0 {
	//			limit += " OFFSET " + strconv.FormatUint(ranges.Offset, 10)
	//		}
	//
	//		// TODO: sqlite & mysql version
	//	}
	// " ORDER BY " +
	//}

	return "SELECT " + fields + " FROM " + table + condition + order + limit
}

func SQLCount(table, condition string, _ *auth.Identity) string {
	query := "SELECT COUNT(*) FROM " + table

	if strings.TrimSpace(condition) != "" {
		return query + " WHERE " + condition
	}

	return query
}

//const defaultPageLengthStr = "200"
//
//func OrderAndLimit(sortBy []string, limits []uint64) string {
//	var sortStr, limitsStr string
//	if len(sortBy) > 0 {
//		for _, s := range sortBy {
//			if s == "" {
//				continue
//			}
//			desc := ""
//			if s[len(s)-1:] == "-" {
//				s = s[:len(s)-1]
//				desc = " DESC"
//			} else if s[len(s)-1:] == "+" {
//				s = s[:len(s)-1]
//			}
//			if sortStr != "" {
//				sortStr += ", "
//			}
//			sortStr += "`" + s + "`" + desc
//		}
//		if sortStr != "" {
//			sortStr = " ORDER BY " + sortStr
//		}
//	}
//	if len(limits) > 1 {
//		// limit[0] can be equal to 0
//		var pageLengthStr string
//		if limits[1] > 0 {
//			pageLengthStr = strconv.FormatUint(limits[1], 10)
//		} else {
//			pageLengthStr = defaultPageLengthStr
//		}
//		limitsStr = " LIMIT " + strconv.FormatUint(limits[0], 10) + ", " + pageLengthStr
//	} else if len(limits) > 0 {
//		if limits[0] > 0 {
//			limitsStr = " LIMIT " + strconv.FormatUint(limits[0], 10)
//		} else {
//			limitsStr = " LIMIT " + defaultPageLengthStr
//		}
//	}
//	return sortStr + limitsStr
//}

type SqlStmt struct {
	Stmt **sql.Stmt
	Sql  string
}

func Prepare(dbh *sql.DB, sqlQuery string, stmt **sql.Stmt) error {
	var err error

	*stmt, err = dbh.Prepare(sqlQuery)
	if err != nil {
		return errors.Wrapf(err, "can't dbh.Prepare(%s)", sqlQuery)
	}

	return nil
}

func Exec(dbh *sql.DB, sqlQuery string, values ...interface{}) (*sql.Result, error) {
	stmt, err := dbh.Prepare(sqlQuery)
	if err != nil {
		return nil, errors.Wrapf(err, CantPrepare, sqlQuery)
	}

	res, err := stmt.Exec(values...)
	if err != nil {
		return nil, errors.Wrapf(err, CantExec, sqlQuery, values)
	}

	return &res, nil
}

func Query(dbh *sql.DB, sqlQuery string, values ...interface{}) (*sql.Rows, error) {
	stmt, err := dbh.Prepare(sqlQuery)
	if err != nil {
		return nil, errors.Wrapf(err, CantPrepare, sqlQuery)
	}

	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, errors.Wrapf(err, CantExec, sqlQuery, values)
	}

	return rows, nil
}

func QueryStrings(stmt *sql.Stmt, sql string, values ...interface{}) (results []string, err error) {
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, errors.Wrapf(err, CantExec, sql, values)
	}
	defer rows.Close()

	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return results, errors.Wrapf(err, CantScanQueryRow, sql, values)
		}

		results = append(results, r)
	}

	err = rows.Err()
	if err != nil {
		return results, errors.Wrapf(err, CantScanQueryRow, sql, values)
	}

	return results, nil
}

//func QueryIDs(stmt *sql.Stmt, sql string, values ...interface{}) (ids []uint64, err error) {
//	rows, err := stmt.Query(values...)
//	if err != nil {
//		return nil, errors.Wrapf(err, basis.CantExecQuery, sql, values)
//	}
//	defer rows.Close()
//
//	for rows.Right() {
//		var id uint64
//		if err := rows.Scan(&id); err != nil {
//			return ids, errors.Wrapf(err, CantScanQueryRow, sql, values)
//		}
//
//		ids = append(ids, id)
//	}
//
//	err = rows.Err()
//	if err != nil {
//		return ids, errors.Wrapf(err, CantScanQueryRow, sql, values)
//	}
//
//	return ids, nil
//}
