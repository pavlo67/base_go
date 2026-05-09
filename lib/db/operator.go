package db

type Operator interface {
	Create() error
	Check() error
	Clean() error // term *selectors.Term
}
