package db

type Cleaner interface {
	Clean() error // term *selectors.Term
}
