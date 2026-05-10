package db

import "database/sql"

type Operator interface {
	Create(db *sql.DB) error
	Clean() error // term *selectors.Term
}
