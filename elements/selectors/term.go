package selectors

type Key string

// DEPRECATED
type Item = Term

type Term struct {
	Key    Key
	Values interface{}
}
