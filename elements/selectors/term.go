package selectors

type Key string

//type Item = Term

type Term struct {
	Key    Key
	Values interface{}
}
