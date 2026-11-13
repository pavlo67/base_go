package sqllib_sqlite

//import (
//	"database/sql"
//	"strings"
//
//	_ "github.com/mattn/go-sqlite3"
//
//	"github.com/pavlo67/common/common/config"
//	"github.com/pavlo67/common/common/errors"
//)
//
//func Connect(cfg config.Access) (*sql.DB, error) {
//	if strings.TrimSpace(cfg.Path) == "" {
//		return nil, errors.New("no path to SQLite database is defined")
//	}
//
//	db, err := sql.Open("sqlite3", cfg.Path)
//	if err != nil {
//		return nil, errors.Wrapf(err, "wrong db connect (cfg = %#v)", cfg)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		return nil, errors.Wrapf(err, "wrong .Ping on db connect (cfg = %#v)", cfg)
//	}
//
//	return db, nil
//}
