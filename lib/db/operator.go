package db

import "database/sql"

type Operator interface {
	Create(db *sql.DB) error
	Clean(opts interface{}) error // term *selectors.Term
}
