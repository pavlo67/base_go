package contacts

type Type string

type Item struct {
	Type
	Value     string
	Connected []Item
}
